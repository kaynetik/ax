package main

import (
	"fmt"
	"os"

	"github.com/kaynetik/ax"
)

const (
	oneInt = int(1)

	flagCompareHelp    = "-help"
	flagCompareArchive = "-archive"
)

func main() {
	args := os.Args

	flags := parseAllFlags()

	if compareFirstArg(flagCompareHelp, args) {
		printHelp()
	}

	if compareFirstArg(flagCompareArchive, args) {
		conf := prepareConfigForArchiving(flags)

		err := archive(conf)
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

func prepareConfigForArchiving(scannedFlags *cmdScan) *ax.Config {
	zc := ax.NewDefaultZipConfig()

	return &ax.Config{
		PathToArchive:  scannedFlags.pathToArchive,
		OutputPath:     "",
		NewArchiveName: "",
		ZipConfig:      &zc,
	}
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
