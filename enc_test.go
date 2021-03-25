package ax

import (
	"crypto/sha256"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	test7z001      = ".7z.001"
	encFileNameOut = testArchiveNewName + test7z001 + ".1"
)

func (s *Suite) TestUnitEncrypt() {
	testCases := []TestCase{
		{
			Name: "full file encryption",
			Assert: func() {
				outDir := "./tests/lorem_data_out/"
				fileIn := outDir + testArchiveNewName + test7z001
				pwdKey := sha256.Sum256([]byte("defaultPwdKey"))

				FileEncryption(pwdKey[:], fileIn, outDir+encFileNameOut)

				assert.FileExists(s.T(), outDir+encFileNameOut)
			},
		},
	}

	RunTestCases(s, testCases)
}

func (s *Suite) TestUnitDefaultEncryption() {
	testCases := []TestCase{
		{
			Name: "default encryption - with cleanup",
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

				// Execute the func being tested.
				err = DefaultFileEncryption(pwdKey, fl)

				// Assert.
				assert.Nil(s.T(), err)
				assert.GreaterOrEqual(s.T(), len(fl), 1)
				assert.FileExists(s.T(), outFile+".enc.0")

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

func copyFileToEnc(t *testing.T, inFilePath, outFilePath string) {
	t.Helper()

	from, err := os.Open(inFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		t.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		t.Fatal(err)
	}
}
