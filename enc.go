package ax

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

const (
	encFilePerm = 0o600
)

// FileEncryption - encrypt a file.
func FileEncryption(bytKey []byte, inFileName, encFileName string) {
	inFile, err := os.Open(inFileName)
	if err != nil {
		panic(err)
	}

	defer inFile.Close()

	block, err := aes.NewCipher(bytKey)
	if err != nil {
		panic(err)
	}

	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewOFB(block, iv)

	outFile, err := os.OpenFile(encFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, encFilePerm)
	if err != nil {
		panic(err)
	}

	defer outFile.Close()

	writer := &cipher.StreamWriter{S: stream, W: outFile}

	_, err = io.Copy(writer, inFile)
	if err != nil {
		panic(err)
	}
}

// DefaultFileEncryption - represents basic usage of the FileEncryption func.
func DefaultFileEncryption(passwd []byte, fileList []string) error {
	key := sha256.Sum256(passwd)

	for i, file := range fileList {
		FileEncryption(key[:], file, fmt.Sprintf("%s.enc.%d", file, i))

		err := os.Remove(file)
		if err != nil {
			return fmt.Errorf("failed removing previous file: %w", err)
		}
	}

	fmt.Println("\nArchive(s) encrypted!")

	return nil
}
