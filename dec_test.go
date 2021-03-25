package ax

import (
	"crypto/sha256"
	"github.com/stretchr/testify/assert"
	"os"
)

const decFileNameOut = testArchiveNewName + test7z001

func (s *Suite) TestUnitFileDecrypt() {
	testCases := []TestCase{
		{
			Name: "full file decryption",
			Assert: func() {
				// Generate encrypted file.
				outDir := "./tests/lorem_data_out/"
				fileIn := outDir + testArchiveNewName + test7z001
				pwdKey := sha256.Sum256([]byte("defaultPwdKey"))

				FileEncryption(pwdKey[:], fileIn, outDir+encFileNameOut)

				// Execute the func being tested.
				FileDecryption(pwdKey[:], outDir+encFileNameOut, outDir+decFileNameOut)

				// Assert.
				assert.FileExists(s.T(), outDir+encFileNameOut)
				_ = os.Remove(outDir + encFileNameOut)
			},
		},
	}

	RunTestCases(s, testCases)
}

func (s *Suite) TestUnitDefaultDecryption() {
	testCases := []TestCase{
		{
			Name: "default decryption - with cleanup",
			Assert: func() {
				pwdKey := []byte("defaultPwdKey")
				// Setup dir structure.
				inFilePath := "./tests/lorem_data_in/lorem.md"
				outFilePath := "./tests/lorem_enc"
				outFile := outFilePath + "/lorem.md"
				_ = os.Mkdir(outFilePath, os.ModePerm)
				copyFileToEnc(s.T(), inFilePath, outFile)

				// Get listing of the temporary out path.
				fl, err := ListFiles(outFilePath, DefaultPathWalkerFunc)
				assert.Nil(s.T(), err)
				err = DefaultFileEncryption(pwdKey, fl)
				assert.Nil(s.T(), err)

				// Get listing of the temporary out path.
				fl, err = ListFiles(outFilePath, DefaultPathWalkerFunc)
				assert.Nil(s.T(), err)

				// Execute the func being tested.
				err = DefaultFileDecryption(pwdKey, fl)

				// Assert.
				assert.Nil(s.T(), err)
				assert.GreaterOrEqual(s.T(), len(fl), 1)
				assert.FileExists(s.T(), outFile)

				// Cleanup.
				err = os.RemoveAll(outFilePath)
				if err != nil {
					s.T().Fatal(err)
				}
			},
		},
	}

	RunTestCases(s, testCases)
}
