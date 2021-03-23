package flags

import (
	"flag"
	"fmt"
	"syscall"

	"golang.org/x/term"
)

const (
	flagNameArchiveIn      = "arc-in"
	flagNamePass           = "arc-pass"
	flagNameArchiveOutPath = "arc-out"
	flagNameNewArchiveName = "arc-name"

	flagValArchiveIn      = "../tmp_to_archive"
	flagValPass           = "on"
	flagValArchiveOutPath = "../tmp_archive_out"
	flagValNewArchiveName = "new_archive"

	flagUsageArchiveIn      = "Select the path which you wish to Archive"
	flagUsagePass           = "If you want to be prompted for a password, or not (default on)"
	flagUsageArchiveOutPath = "Select the path where you want to store temporary Archive(s)"
	flagUsageNewArchiveName = "Choose the name of new (temporary) Archive(s)"

	promptEnterPasswordForArchiveEncryption = "Enter Password for Archive(s) Encryption"
)

// CmdScan - represents scanned flags from the stdin.
type CmdScan struct {
	ProtectArchiveWithPasswd bool
	PasswordByte             []byte
	PathToArchive            string
	ArchiveOutPath           string
	NewArchiveName           string
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

	flag.Parse()

	if flagPass == flagValPass {
		bytePassword, err = protectedScan(promptEnterPasswordForArchiveEncryption)
		if err != nil {
			panic(err)
		}

		cs.ProtectArchiveWithPasswd = true
		cs.PasswordByte = bytePassword
	} else {
		cs.ProtectArchiveWithPasswd = false
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
