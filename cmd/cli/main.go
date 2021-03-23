package main

import (
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
)

func main() {
	args := os.Args

	cmdScan := flags.ParseAllFlags()

	if compareFirstArg(flagCompareHelp, args) {
		printHelp()
	}

	if compareFirstArg(flagCompareArchiveIn, args) {
		conf := prepareConfigForArchiving(cmdScan)

		err := archive(conf)
		if err != nil {
			panic(err)
		}

		return
	}

	if compareFirstArg(flagCompareArchiveExtract, args) {
		conf := prepareConfigForExtracting(cmdScan)

		err := extract(conf)
		if err != nil {
			panic(err)
		}

		return
	}
}

func archive(conf *ax.Config) error {
	err := ax.Archive(conf)
	if err != nil {
		return fmt.Errorf("an issue occurred while archiving: %w", err)
	}

	fmt.Println("Finished archiving!")

	return nil
}

func prepareConfigForArchiving(scannedFlags *flags.CmdScan) *ax.Config {
	zc := ax.NewDefaultZipConfig()
	zc.Password = scannedFlags.PasswordByte
	zc.ApplyPassword = scannedFlags.ProtectArchiveWithPasswd

	return &ax.Config{
		PathToArchive:  scannedFlags.PathToArchive,
		OutputPath:     scannedFlags.ArchiveOutPath,
		NewArchiveName: scannedFlags.NewArchiveName,
		ZipConfig:      &zc,
	}
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

func compareFirstArg(flagForComparison string, args []string) bool {
	argsLen := len(args)

	if argsLen > oneInt && args[oneInt] == flagForComparison {
		return true
	}

	return false
}

func printHelp() {
	fmt.Printf("TODO: Additional info about the crypto CLI implementation.\n\n")
	fmt.Printf("When called without any arguments/flags, interactive mode will be initated.\n")
	fmt.Printf("When called with arguments/flags, those that are left out will assume their default required values.\n\n")
}
