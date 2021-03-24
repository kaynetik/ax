package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/kaynetik/ax/cmd/cli/flags"

	"github.com/kaynetik/ax"
)

const (
	oneInt = int(1)

	flagCompareHelp           = "-help"
	flagCompareArchiveIn      = "-arc-in"
	flagCompareArchiveExtract = "-arc-extract"
	flagCompareEncryptIn      = "-enc-in"
	flagCompareDecryptIn      = "-dec-in"
	flagCompareGitRepo        = "-git-repo"
)

func main() {
	args := os.Args
	cmdScan := &flags.CmdScan{}

	if len(args) <= 1 {
		printInteractiveModeHelp()

		getScannerFn := func() *bufio.Scanner {
			return bufio.NewScanner(os.Stdin)
		}

		cmdScan.ReadUserInput(getScannerFn)

		archiveEncryptAndPushToGit(cmdScan)

		return
	}

	cmdScan = flags.ParseAllFlags()

	switch args[oneInt] {
	case flagCompareHelp:
		printHelp()
	case flagCompareArchiveIn:
		conf := prepareConfigForArchiving(cmdScan)

		err := archive(conf)
		if err != nil {
			panic(err)
		}

		return
	case flagCompareArchiveExtract:
		conf := prepareConfigForExtracting(cmdScan)

		err := extract(conf)
		if err != nil {
			panic(err)
		}

		return
	case flagCompareEncryptIn:
		fileList, err := ax.ListFiles(cmdScan.EncryptPath)
		if err != nil {
			panic(err)
		}

		err = ax.DefaultFileEncryption(cmdScan.EncryptPassword, fileList)
		if err != nil {
			panic(err)
		}

		return
	case flagCompareDecryptIn:
		fileList, err := ax.ListFiles(cmdScan.DecryptPath)
		if err != nil {
			panic(err)
		}

		err = ax.DefaultFileDecryption(cmdScan.DecryptPassword, fileList)
		if err != nil {
			panic(err)
		}

		return
	case flagCompareGitRepo:
		archiveEncryptAndPushToGit(cmdScan)
	default:
		panic(errors.New("unknown flag provided"))
	}
}

func archiveEncryptAndPushToGit(cs *flags.CmdScan) {
	// Cleanup
	err := os.RemoveAll(cs.ArchiveOutPath)
	if err != nil {
		panic(err)
	}

	// Archive
	arcConf := prepareConfigForArchiving(cs)

	err = archive(arcConf)
	if err != nil {
		panic(err)
	}

	// Encrypt
	fileList, err := ax.ListFiles(cs.EncryptPath)
	if err != nil {
		panic(err)
	}

	err = ax.DefaultFileEncryption(cs.EncryptPassword, fileList)
	if err != nil {
		panic(err)
	}

	// Push to GIT Repository
	err = os.Chdir(cs.ArchiveOutPath)
	if err != nil {
		panic(err)
	}

	err = ax.PushToGIT(cs.GitRepo)
	if err != nil {
		panic(err)
	}

	fmt.Println("Pushed to GIT! Your Archive(s) have been backed up!")
}

func archive(conf *ax.ArchiveConfig) error {
	err := ax.Archive(conf)
	if err != nil {
		return fmt.Errorf("an issue occurred while archiving: %w", err)
	}

	fmt.Println("Finished Archiving!")

	return nil
}

func prepareConfigForArchiving(scannedFlags *flags.CmdScan) *ax.ArchiveConfig {
	ac := ax.NewDefaultArchiveConfig()
	ac.Password = scannedFlags.PasswordByte
	ac.ApplyPassword = scannedFlags.ProtectArchiveWithPasswd
	ac.PathToArchive = scannedFlags.PathToArchive
	ac.OutputPath = scannedFlags.ArchiveOutPath
	ac.NewArchiveName = scannedFlags.NewArchiveName

	return &ac
}

func extract(conf *ax.ExtractConfig) error {
	err := ax.Extract(conf)
	if err != nil {
		return fmt.Errorf("an issue occurred while etxtacting archive(s): %w", err)
	}

	fmt.Println("Finished Extracting!")

	return nil
}

func prepareConfigForExtracting(scannedFlags *flags.CmdScan) *ax.ExtractConfig {
	ec := ax.ExtractConfig{
		Password:    scannedFlags.PasswordByte,
		ExtractPath: scannedFlags.ArchiveExtract,
	}

	return &ec
}

func printHelp() {
	fmt.Printf("TODO: Additional info about the crypto CLI implementation.\n\n")
	fmt.Printf("When called without any arguments/flags, interactive mode will be initated.\n")
	fmt.Printf("When called with arguments/flags, those that are left out will assume their default required values.\n\n")
}

func printInteractiveModeHelp() {
	fmt.Printf("You are using AX interactive mode. Prompts that are left blank will use default values.\n")
	fmt.Printf("Currently we support only GIT Push via ssh: git@github.com:{USER}/{REPOSITORY}.git\n")
	fmt.Println("\nChoose which action do you want to preform:")
	fmt.Println("\n1. Archive, Encrypt & Push to GIT Repo")
	fmt.Println("\n2. .... TBD - This is Work in Progress")
	fmt.Println("\n\nNote: Defaulting to option number 1, as others aren't supported yet!")
}
