package flags

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
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

// protectedScan - used to read password from stdin. Input is being hidden while typing.
func protectedScan(prompt string) ([]byte, error) {
	fmt.Printf("\n%s: ", prompt)

	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("protected scan failed: %w", err)
	}

	return bytePassword, nil
}

// getScannerFn - signature of the function we expect for the ReadUserInput func.
type getScannerFn func() *bufio.Scanner

// scanner - internal struct that is used solely to wrap around bufio.Scanner.
type scanner struct {
	ioScanner *bufio.Scanner
}

// newScanner - returns instance of scanner, wrapper for bufio.Scanner.
func newScanner(io *bufio.Scanner) scanner {
	return scanner{ioScanner: io}
}

// scanWithMsg - prompts the user for input, with the provided message (prompt).
func (s *scanner) scanWithMsg(prompt string) string {
	fmt.Printf("\n%s: ", prompt)
	s.ioScanner.Scan()

	return s.ioScanner.Text()
}

// ReadUserInput - is used for interactive CLI mode.
//
// User will be prompted line by line for the required information. If something is left out, it will be ignored
//	if not required. But if it's required then a default value will be used instead.
// But for cases like Password for protecting the Archive - it can be left blank, but it's advised not to do so.
func (cs *CmdScan) ReadUserInput(fn getScannerFn) {
	// Get new bufio scanner instance [*bufio.Scanner].
	s := newScanner(fn())

	// Scan for Path which should be Archived.
	pathToArchive := s.scanWithMsg("Path to Archive")
	if pathToArchive != "" {
		cs.PathToArchive = pathToArchive
	} else {
		cs.PathToArchive = flagValArchiveIn
	}

	// Scan for Archive Temporary Output path.
	archOutPath := s.scanWithMsg("Path for Archive(s) Output (default '../tmp_archive_out')")
	if archOutPath != "" {
		cs.ArchiveOutPath = archOutPath
	} else {
		cs.ArchiveOutPath = flagValArchiveOutPath
	}

	// Scan for GIT Repository.
	gitRepo := s.scanWithMsg("GIT Repository")
	if gitRepo != "" {
		cs.GitRepo = gitRepo
	}

	// Scan for Passwords
	fmt.Printf("\n\nWARN: You will now be prompted for passwords. Those will be hidden, so just keep on typing.\n")

	// Scan for Archive Password.
	archivePasswd, err := protectedScan("Password for Archive protection")
	if err != nil {
		panic(err)
	}

	cs.PasswordByte = archivePasswd
	cs.ProtectArchiveWithPasswd = true

	// Scan if user wants two unique layers of encryption
	differentPasswd := s.scanWithMsg(
		"\nDo you want to chose different password for Encrypting generated Archive Volume(s)? (yes/no, default yes)",
	)
	if differentPasswd == promptAnswerNo {
		cs.EncryptPassword = archivePasswd
	} else {
		// Scan for Encryption Password.
		encryptionPassword, errPwd := protectedScan("Password for Archive Volume(s) Encryption")
		if errPwd != nil {
			panic(errPwd)
		}

		cs.EncryptPassword = encryptionPassword
	}

	cs.EncryptPath = cs.ArchiveOutPath // TODO: We might want this flexible?
}
