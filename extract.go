package ax

import (
	"fmt"
	"os"
)

const (
	archiveWildcard001 = "*.001"
)

// ExtractConfig - represents configuration which is required for 7zip extraction process.
type ExtractConfig struct {
	// Password - if set, it will be used to decrypt the archive.
	Password []byte

	// ExtractPath - path which points to the directory to archive(s) location (for extraction).
	ExtractPath string
}

// Extract - used to extract the archive(s).
func Extract(conf *ExtractConfig) error {
	err := executeCommand(cmd7z, cmdArgsArchiveExtract(conf))
	if err != nil {
		return fmt.Errorf("failed executing 7zip: %w", err)
	}

	return nil
}

// cmdArgsArchiveExtract - used to build command arguments for Archive Extraction process.
// Returned string will be transformed into arguments slice which is later used for Output() func on [exec.Command].
func cmdArgsArchiveExtract(ec *ExtractConfig) string {
	cmdString := "x"

	// Append password, if defined.
	if ec.Password != nil || string(ec.Password) != "" {
		appendToString(&cmdString, fmt.Sprintf("-p%s", ec.Password))
	}

	// TODO: Here a check for volumes should be applied. If the archive wasn't previously split into volumes, then exact
	// Append path to the archive which has to be extracted.
	appendToString(&cmdString, fmt.Sprintf("%s%c%s", ec.ExtractPath, os.PathSeparator, archiveWildcard001))

	// Append path where we want files to be extracted.
	appendToString(&cmdString, fmt.Sprintf("-o%s", ec.ExtractPath))

	return cmdString
}
