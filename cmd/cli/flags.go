package main

import (
	"flag"
	"fmt"
	"syscall"

	"golang.org/x/term"
)

const (
	flagNameArchive = "archive"
	flagNamePass    = "pass"

	flagValArchive = "../tmp_to_archive"
	flagValPass    = "on"

	flagUsageArchive = "Select the path which you wish to Archive"
	flagUsagePass    = "If you want to be prompted for a password, or not (default on)"

	promptEnterPasswordForArchiveEncryption = "Enter Password for Archive(s) Encryption"
)

type cmdScan struct {
	passwd        []byte
	pathToArchive string
}

func parseAllFlags() *cmdScan {
	var (
		err               error
		cs                cmdScan
		bytePassword      []byte
		flagPathToArchive string
		flagPass          string
	)

	flag.StringVar(&flagPathToArchive, flagNameArchive, flagValArchive, flagUsageArchive)
	flag.StringVar(&flagPass, flagNamePass, flagValPass, flagUsagePass)

	flag.Parse()

	if flagPass == flagValPass {
		bytePassword, err = protectedScan(promptEnterPasswordForArchiveEncryption)
		if err != nil {
			panic(err)
		}
	}

	cs.passwd = bytePassword
	cs.pathToArchive = flagPathToArchive

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
