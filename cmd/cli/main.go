package main

import (
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

	if len(args) <= 1 {
		fmt.Println("This will parse the interactive mode! WIP")

		return
	}

	cmdScan := flags.ParseAllFlags()

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
		// Cleanup
		err := os.RemoveAll(cmdScan.ArchiveOutPath)
		if err != nil {
			panic(err)
		}

		// Archive
		arcConf := prepareConfigForArchiving(cmdScan)

		err = archive(arcConf)
		if err != nil {
			panic(err)
		}

		// Encrypt
		fileList, err := ax.ListFiles(cmdScan.EncryptPath)
		if err != nil {
			panic(err)
		}

		err = ax.DefaultFileEncryption(cmdScan.EncryptPassword, fileList)
		if err != nil {
			panic(err)
		}

		// Push to GIT Repository
		err = os.Chdir(cmdScan.ArchiveOutPath)
		if err != nil {
			panic(err)
		}

		err = ax.PushToGIT(cmdScan.GitRepo)
		if err != nil {
			panic(err)
		}
	default:
		panic(errors.New("unknown flag provided"))
	}
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
