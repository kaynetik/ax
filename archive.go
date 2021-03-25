package ax

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	cmd7z                = "7z"
	archiveType          = "7z"
	defaultArchiveOutput = "tmp_archive"
)

var (
	// ErrNotDir - path is not a directory.
	// Currently only archiving dirs has been thoroughly tested.
	ErrNotDir = errors.New("path is not a directory")

	// ErrPathEmpty - path can not be empty.
	ErrPathEmpty = errors.New("path can not be empty")
)

// PathConfig - path config required for the archiving process.
type PathConfig struct {
	// PathToArchive - path which points to the directory which should be archived.
	// Should be relative to the current working dir.
	PathToArchive string

	// OutputPath - path where the temporary archive files should be placed.
	OutputPath string

	// NewArchiveName - if set it will represent the base for the name of output archive(s).
	NewArchiveName string
}

// ArchiveConfig - represents the config required for the archiving process.
//
// If it's not provided, then default one will be used. Note that default setting is 'Ultra' (highest compression).
// This config will be just append-only, to keep API consistency.
// Preferred over params in order to avoid breaking API with future updates.
type ArchiveConfig struct {
	// PathConfig - embeds required path information for archiving process.
	PathConfig

	// Password - if set, it will be used to encrypt the archive.
	Password []byte

	// ArchiveType - default setting '-t7z'.
	ArchiveType string

	// BlockSize - default setting 'm' [BlockSizeMB].
	BlockSize BlockSize // b,k,m,g - size representation

	// VolumeSize - default setting '-v9m' - representing volumes of 9Megabytes each.
	VolumeSize uint64

	// FastBytes - default setting '-mfb=64', where the number set represents the number of Fast bytes for LZMA.
	FastBytes uint16

	// DictSize - default setting '-md=64m.
	DictSize uint16

	// ApplyPassword - if true, then password flag for 7zip cmd will be used.
	ApplyPassword bool

	// HeadersEncryption - default setting '-mhe=on'. Enables Headers encryption.
	HeadersEncryption bool // -he=on

	// Compression - default setting '-mx=9', where 9 represents Ultra and 0 would be none compression at all.
	Compression uint8

	// SolidArchive - default setting '-ms=on'
	SolidArchive bool
}

// Archive - used to create archive zip volume(s) from a chosen directory.
func Archive(conf *ArchiveConfig) error {
	err := validatePathToArchive(conf)
	if err != nil {
		return fmt.Errorf("path validation issue: %w", err)
	}

	err = executeCommand(cmd7z, cmdArgsArchive(conf))
	if err != nil {
		return fmt.Errorf("failed executing 7zip: %w", err)
	}

	return nil
}

func validatePathToArchive(conf *ArchiveConfig) error {
	if conf.PathToArchive == "" {
		return ErrPathEmpty
	}

	stat, err := os.Stat(conf.PathToArchive)
	if err != nil {
		return fmt.Errorf("failed getting path stat: %w", err)
	}

	if !stat.IsDir() {
		return ErrNotDir
	}

	return nil
}

// BlockSize - represents size of the volume blocks in b|k|m|g.
type BlockSize string

func (v *BlockSize) String() string {
	return string(*v)
}

// DetermineBlockSize - returns BlockSize casted from an input string.
func DetermineBlockSize(bs string) BlockSize {
	switch bs {
	case letterB:
		return BlockSizeByte
	case letterK:
		return BlockSizeKB
	case letterG:
		return BlockSizeGB
	default:
		return defaultBlockSize
	}
}

const (
	letterB = "b"
	letterK = "k"
	letterM = "m"
	letterG = "g"

	// BlockSizeByte - Byte representative character.
	BlockSizeByte = BlockSize(letterB)

	// BlockSizeKB - KB representative character.
	BlockSizeKB = BlockSize(letterK)

	// BlockSizeMB - MB representative character.
	BlockSizeMB = BlockSize(letterM)

	// BlockSizeGB - GB representative character.
	BlockSizeGB = BlockSize(letterG)
)

const (
	defaultArchiveType      = archiveType
	defaultBlockSize        = BlockSizeMB
	defaultVolumeSize       = uint64(90)
	defaultFastBytes        = uint16(64)
	defaultDictSize         = uint16(64)
	defaultHeadersEnc       = true
	defaultCompressionLevel = uint8(9)
	defaultSolidArchive     = true
	defaultPasswordStr      = ""
)

func getDefaultPassword() []byte { return []byte(defaultPasswordStr) }

// NewDefaultArchiveConfig - returns ArchiveConfig with default values pre-set.
//
// Note that default setting is 'Ultra' (highest compression).
func NewDefaultArchiveConfig() ArchiveConfig {
	return ArchiveConfig{
		Password:          getDefaultPassword(),
		ArchiveType:       defaultArchiveType,
		BlockSize:         defaultBlockSize,
		VolumeSize:        defaultVolumeSize,
		FastBytes:         defaultFastBytes,
		DictSize:          defaultDictSize,
		HeadersEncryption: defaultHeadersEnc,
		Compression:       defaultCompressionLevel,
		SolidArchive:      defaultSolidArchive,
	}
}

// cmdArgsArchive - used to build command arguments for Archive Compression process.
// Returned string will be transformed into arguments slice which is later used for Output() func on [exec.Command].
func cmdArgsArchive(ac *ArchiveConfig) string {
	var (
		exportArchive struct{ name, typ string }
		cmdStr        = "a"
	)

	if ac.HeadersEncryption {
		appendToString(&cmdStr, "-mhe=on")
	}

	if ac.Password != nil && ac.ApplyPassword {
		appendToString(&cmdStr, fmt.Sprintf("-p%s", ac.Password))
	}

	if ac.ArchiveType != "" {
		exportArchive.typ = ac.ArchiveType
		appendToString(&cmdStr, fmt.Sprintf("-t%s", ac.ArchiveType))
	} else {
		exportArchive.typ = "7z"
	}

	if ac.Compression != 0 {
		appendToString(&cmdStr, fmt.Sprintf("-mx=%d", ac.Compression))
	}

	if ac.FastBytes != 0 {
		appendToString(&cmdStr, fmt.Sprintf("-mfb=%d", ac.FastBytes))
	}

	if ac.DictSize != 0 {
		appendToString(&cmdStr, fmt.Sprintf("-md=%dm", ac.DictSize))
	}

	if ac.VolumeSize != 0 {
		if ac.BlockSize == "" {
			ac.BlockSize = BlockSizeMB
		}

		appendToString(&cmdStr, fmt.Sprintf("-v%d%s", ac.VolumeSize, ac.BlockSize))
	}

	if ac.SolidArchive {
		appendToString(&cmdStr, "-ms=on")
	}

	if ac.NewArchiveName != "" {
		exportArchive.name = ac.NewArchiveName
	} else {
		exportArchive.name = "archive"
	}

	if ac.OutputPath == "" {
		ac.OutputPath = defaultArchiveOutput
	}

	if exportArchive.name != "" && exportArchive.typ != "" {
		appendToString(
			&cmdStr,
			fmt.Sprintf("%s%c%s.%s", ac.OutputPath, os.PathSeparator, exportArchive.name, exportArchive.typ),
		)
	}

	appendToString(&cmdStr, ac.PathToArchive)

	return cmdStr
}

func appendToString(command *string, flag string) {
	*command = fmt.Sprintf("%s %s", *command, flag)
}

// ListFiles - used to list files, without directories in a chosen path.
func ListFiles(pathToWalk string, walkFuncBuilder walkFuncBuilder) ([]string, error) {
	fileList := make([]string, 0)

	pathWalkerFn := walkFuncBuilder(&fileList)

	err := filepath.Walk(pathToWalk, pathWalkerFn)
	if err != nil {
		return fileList, fmt.Errorf("failed walking path: %s with error: %w", pathToWalk, err)
	}

	return fileList, nil
}

type walkFuncBuilder func(fileList *[]string) filepath.WalkFunc

// DefaultPathWalkerFunc - returns default implementation of filepath.WalkFunc.
//
// This approach enables the flexibility to override the filepath.WalkFunc used by our ListFiles func.
func DefaultPathWalkerFunc(fileList *[]string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk initiated with an error: %w", err)
		}

		s, err := os.Stat(path) //TODO: is it worth injecting this func, just for test coverage?
		if err != nil {
			return fmt.Errorf("failed reading path: %s: %w", path, err)
		}

		if !s.IsDir() && !strings.Contains(path, ".git/") {
			*fileList = append(*fileList, path)
		}

		return nil
	}
}
