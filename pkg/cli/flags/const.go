package flags

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
	promptAnswerNo                          = "no"
)
