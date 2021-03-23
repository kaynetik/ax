package ax

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	cmd7z = "7z"

	archiveType = "7z"
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

	// BlockSize - default setting 'm' [BlockSizeMByte].
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

// BlockSize - represents size of the volume blocks in b|k|m|g.
type BlockSize string

func (v *BlockSize) String() string {
	return string(*v)
}

// DetermineBlockSize - returns BlockSize casted from an input string.
func DetermineBlockSize(bs string) BlockSize {
	switch bs {
	case "b":
		return BlockSizeByte
	case "k":
		return BlockSizeKByte
	case "g":
		return BlockSizeGByte
	default:
		return BlockSizeMByte
	}
}

const (
	// BlockSizeByte - Byte representative character.
	BlockSizeByte = BlockSize("b")

	// BlockSizeKByte - KiloByte representative character.
	BlockSizeKByte = BlockSize("k")

	// BlockSizeMByte - MegaByte representative character.
	BlockSizeMByte = BlockSize("m")

	// BlockSizeGByte - GigaByte representative character.
	BlockSizeGByte = BlockSize("g")
)

// NewDefaultArchiveConfig - returns ArchiveConfig with default values pre-set.
//
// Note that default setting is 'Ultra' (highest compression).
func NewDefaultArchiveConfig() ArchiveConfig {
	return ArchiveConfig{
		Password:          []byte(""),
		ArchiveType:       archiveType,
		BlockSize:         BlockSizeByte, // TODO: Change to mb before release
		VolumeSize:        90,
		FastBytes:         64,
		DictSize:          64,
		HeadersEncryption: true,
		Compression:       9,
		SolidArchive:      true,
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
			ac.BlockSize = BlockSizeMByte
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
		ac.OutputPath = "tmp_archive"
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
func ListFiles(pathToWalk string) ([]string, error) {
	var fileList []string

	err := filepath.Walk(pathToWalk, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk initiated with an error: %w", err)
		}

		s, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed reading path: %s: %w", path, err)
		}

		if !s.IsDir() && !strings.Contains(path, "/.git/") {
			fileList = append(fileList, path)
		}

		return nil
	})
	if err != nil {
		return fileList, fmt.Errorf("failed walking path: %s with error: %w", pathToWalk, err)
	}

	return fileList, nil
}
