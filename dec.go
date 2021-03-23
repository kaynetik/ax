package ax

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

// FileDecryption - decrypt a file.
func FileDecryption(key []byte, encFileName, decFileName string) {
	inFile, err := os.Open(encFileName)
	if err != nil {
		panic(err)
	}

	defer inFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewOFB(block, iv)

	outFile, err := os.OpenFile(decFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, encFilePerm)
	if err != nil {
		panic(err)
	}

	defer outFile.Close()

	reader := &cipher.StreamReader{S: stream, R: inFile}

	if _, err = io.Copy(outFile, reader); err != nil {
		panic(err)
	}
}

// DefaultFileDecryption -- represents basic usage of the FileDecryption func.
func DefaultFileDecryption(passwd []byte, fileList []string) error {
	key := sha256.Sum256(passwd)

	for _, file := range fileList {
		fileNameSlc := strings.Split(file, ".")
		decryptedFileName := strings.Join(fileNameSlc[:len(fileNameSlc)-2], ".")

		FileDecryption(key[:], file, decryptedFileName)

		err := os.Remove(file)
		if err != nil {
			return fmt.Errorf("failed removing file at path [%s]: %w", file, err)
		}
	}

	fmt.Println("Archives Decrypted!")

	return nil
}
