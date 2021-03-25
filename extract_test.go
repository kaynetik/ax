package ax

import (
	"fmt"
	"os"

	"github.com/stretchr/testify/assert"
)

const (
	testDirLoremIn     = "lorem_data_in"
	testLoremFileName  = "lorem.md"
	testPathToArchive  = "." + string(os.PathSeparator) + "tests" + string(os.PathSeparator) + testDirLoremIn
	testPathArchiveOut = "./tests/lorem_data_out"
	testArchiveNewName = "test_new_name"
)

func (s *Suite) TestUnitExtract() {
	testCases := []TestCase{
		{
			Name: "err extracting",
			PreRequisites: func() {
				s.ec = &ExtractConfig{}
			},
			Assert: func() {
				err := Extract(s.ec)

				assert.NotNil(s.T(), err)
			},
		},
		{
			Name: "success extracting",
			PreRequisites: func() {
				ac := NewDefaultArchiveConfig()
				ac.PathToArchive = testPathToArchive
				ac.OutputPath = testPathArchiveOut
				ac.NewArchiveName = testArchiveNewName
				s.ac = &ac

				_ = os.RemoveAll(s.ac.OutputPath)

				err := Archive(s.ac)
				if err != nil {
					s.T().Fatalf("prerequisit: Archive, failed with err: %q", err)
				}

				s.ec = &ExtractConfig{
					Password:    getDefaultPassword(),
					ExtractPath: s.ac.OutputPath,
				}
			},
			Assert: func() {
				err := Extract(s.ec)
				expectedFileExistence := fmt.Sprintf("%s%c%s%c%s",
					testPathArchiveOut, os.PathSeparator, testDirLoremIn, os.PathSeparator, testLoremFileName)

				assert.Nil(s.T(), err)
				assert.FileExists(s.T(), expectedFileExistence)
			},
		},
	}

	RunTestCases(s, testCases)
}
