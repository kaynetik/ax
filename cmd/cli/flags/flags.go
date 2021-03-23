package flags

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

const (
	flagNameArchiveIn      = "arc-in"
	flagNamePass           = "arc-pass"
	flagNameArchiveOutPath = "arc-out"
	flagNameNewArchiveName = "arc-name"
	flagNameArchiveExtract = "arc-extract"
	flagNameGitRepo        = "git-repo"

	flagNameEncryptIn = "enc-in"
	flagNameDecryptIn = "dec-in"

	flagValArchiveIn      = "../tmp_to_archive"
	flagValPass           = "on"
	flagValArchiveOutPath = "../tmp_archive_out"
	flagValNewArchiveName = "new_archive"
	flagValArchiveExtract = "../tmp_archive_out"
	flagValGitRepo        = "git@github.com:USER/REPOSITORY.git"

	flagValEncryptIn = "../tmp_archive_out"
	flagValDecryptIn = "../tmp_archive_out"

	flagUsageArchiveIn      = "Select the path which you wish to Archive"
	flagUsagePass           = "If you want to be prompted for a password, or not (default on)"
	flagUsageArchiveOutPath = "Select the path where you want to store temporary Archive(s)"
	flagUsageNewArchiveName = "Choose the name of new (temporary) Archive(s)"
	flagUsageArchiveExtract = "Choose the path of Archive(s) location, which should be Extracted"
	flagUsageGitRepo        = "Enter the remote GIT Repository where you wish to persist your backup"

	flagUsageEncryptIn = "Select the path in which files for Encryption are located"
	flagUsageDecryptIn = "Select the path in which files for Decryption are located " +
		"\nIf the password isn't correct, no warning will be provided, " +
		"to disable possibility of brute forcing the correct one"

	promptEnterPasswordForArchiveEncryption = "Enter Password for to protect Archive(s)"
	promptEnterPasswordForEncryption        = "Enter Password for Archive(s) Encryption"
	promptEnterPasswordForDecryption        = "Enter Password for Archive(s) Decryption"
)

// CmdScan - represents scanned flags from the stdin.
type CmdScan struct {
	ProtectArchiveWithPasswd bool
	PasswordByte             []byte
	PathToArchive            string
	ArchiveOutPath           string
	NewArchiveName           string
	ArchiveExtract           string
	GitRepo                  string
	EncryptPath              string
	DecryptPath              string
	EncryptPassword          []byte
	DecryptPassword          []byte
}

// ParseAllFlags - parses flags from the tty and applies validation for that input.
func ParseAllFlags() *CmdScan {
	var (
		err          error
		cs           CmdScan
		bytePassword []byte
		flagPass     string
	)

	flag.StringVar(&cs.PathToArchive, flagNameArchiveIn, flagValArchiveIn, flagUsageArchiveIn)
	flag.StringVar(&flagPass, flagNamePass, flagValPass, flagUsagePass)
	flag.StringVar(&cs.ArchiveOutPath, flagNameArchiveOutPath, flagValArchiveOutPath, flagUsageArchiveOutPath)
	flag.StringVar(&cs.NewArchiveName, flagNameNewArchiveName, flagValNewArchiveName, flagUsageNewArchiveName)
	flag.StringVar(&cs.ArchiveExtract, flagNameArchiveExtract, flagValArchiveExtract, flagUsageArchiveExtract)
	flag.StringVar(&cs.GitRepo, flagNameGitRepo, flagValGitRepo, flagUsageGitRepo)
	flag.StringVar(&cs.EncryptPath, flagNameEncryptIn, flagValEncryptIn, flagUsageEncryptIn)
	flag.StringVar(&cs.DecryptPath, flagNameDecryptIn, flagValDecryptIn, flagUsageDecryptIn)

	flag.Parse()

	var (
		argsStr          string
		archiveCalled    bool
		encryptionCalled bool
		decryptionCalled bool
		pushCalled       bool
	)

	argsStr = strings.Join(os.Args, " ")
	encryptionCalled = strings.Contains(argsStr, flagNameEncryptIn)
	decryptionCalled = strings.Contains(argsStr, flagNameDecryptIn)
	archiveCalled = flagPass == flagValPass && !encryptionCalled && !decryptionCalled
	pushCalled = encryptionCalled && cs.GitRepo != flagValGitRepo && flagPass == flagValPass

	if archiveCalled || pushCalled {
		bytePassword, err = protectedScan(promptEnterPasswordForArchiveEncryption)
		if err != nil {
			panic(err)
		}

		cs.ProtectArchiveWithPasswd = true
		cs.PasswordByte = bytePassword
	} else {
		cs.ProtectArchiveWithPasswd = false
	}

	if encryptionCalled || pushCalled {
		bytePassword, err = protectedScan(promptEnterPasswordForEncryption)
		if err != nil {
			panic(err)
		}

		cs.EncryptPassword = bytePassword
	}

	if decryptionCalled {
		bytePassword, err = protectedScan(promptEnterPasswordForDecryption)
		if err != nil {
			panic(err)
		}

		cs.DecryptPassword = bytePassword
	}

	return &cs
}

func protectedScan(prompt string) ([]byte, error) {
	fmt.Printf("\n%s: ", prompt)

	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return nil, fmt.Errorf("protected scan failed: %w", err)
	}

	return bytePassword, nil
}
