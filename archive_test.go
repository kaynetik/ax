package ax

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
)

const nameOfArchiveGoSrcFile = "archive.go"

func (s *Suite) TestUnitValidatePath() {
	testCases := []TestCase{
		{
			Name: "err empty path",
			PreRequisites: func() {
				s.testArchiveConfig = &ArchiveConfig{}
			},
			Assert: func() {
				err := validatePathToArchive(s.testArchiveConfig)

				assert.NotNil(s.T(), err)
				assert.EqualValues(s.T(), ErrPathEmpty, err)
			},
		},
		{
			Name: "err getting filepath stat",
			PreRequisites: func() {
				s.testArchiveConfig = &ArchiveConfig{
					PathConfig: PathConfig{
						PathToArchive: "./invalid\\on_each\\/os",
					},
				}
			},
			Assert: func() {

				err := validatePathToArchive(s.testArchiveConfig)

				assert.NotNil(s.T(), err)
			},
		},
		{
			Name: "err file is not of dir type",
			PreRequisites: func() {
				s.testArchiveConfig = &ArchiveConfig{
					PathConfig: PathConfig{
						PathToArchive: fmt.Sprintf(".%c%s", os.PathSeparator, nameOfArchiveGoSrcFile),
					},
				}
			},
			Assert: func() {

				err := validatePathToArchive(s.testArchiveConfig)

				assert.NotNil(s.T(), err)
				assert.EqualValues(s.T(), ErrNotDir, err)
			},
		},
		{
			Name: "successful validation of path",
			PreRequisites: func() {
				s.testArchiveConfig = &ArchiveConfig{
					PathConfig: PathConfig{
						PathToArchive: fmt.Sprintf(".%ccmd", os.PathSeparator),
					},
				}
			},
			Assert: func() {

				err := validatePathToArchive(s.testArchiveConfig)

				assert.Nil(s.T(), err)
			},
		},
	}

	RunTestCases(s, testCases)
}

func (s *Suite) TestUnitArchive() {
	testCases := []TestCase{
		{
			Name: "err path validation",
			PreRequisites: func() {
				// Just triggering one of the cases, given that path validation has been covered
				//  with TestUnitValidatePath test.
				s.testArchiveConfig = &ArchiveConfig{}
			},
			Assert: func() {
				err := Archive(s.testArchiveConfig)

				assert.NotNil(s.T(), err)
			},
		},
		{
			Name: "successful command execution",
			PreRequisites: func() {
				s.testArchiveConfig = &ArchiveConfig{
					PathConfig: PathConfig{
						PathToArchive: fmt.Sprintf(".%ccmd", os.PathSeparator),
					},
				}
			},
			Assert: func() {
				err := Archive(s.testArchiveConfig)

				assert.Nil(s.T(), err)
			},
		},
	}

	RunTestCases(s, testCases)
}

func (s *Suite) TestUnitDetermineBlockSize() {
	testCases := []TestCase{
		{
			Name: "success convert all cases",

			Assert: func() {
				// Exception to the general rule - it is acceptable that this assert iterates over few cases,
				// given that they don't have an returning error, but are just casting letters to the BlockSize type.
				cases := []struct {
					blockSizeStr string
					expectedBS   BlockSize
				}{
					{
						letterB,
						BlockSizeByte,
					},
					{
						letterK,
						BlockSizeKB,
					},
					{
						letterM,
						BlockSizeMB,
					},
					{
						letterG,
						BlockSizeGB,
					},
					{
						"defaultValMB",
						defaultBlockSize,
					},
				}

				for _, c := range cases {
					gotBS := DetermineBlockSize(c.blockSizeStr)

					assert.EqualValues(s.T(), c.expectedBS, gotBS)
				}

			},
		},
	}

	RunTestCases(s, testCases)
}
