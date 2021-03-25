package ax

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"
)

const (
	nameOfArchiveGoSrcFile = "archive.go"
	invalidPath            = "./invalid\\on_each\\/os"
)

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
						PathToArchive: invalidPath,
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
				const randomDefaultVal = "defaultVal"

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
						randomDefaultVal,
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

func (s *Suite) TestUnitBlockSizeStringReturn() {
	testCases := []TestCase{
		{
			Name: "success stringer interface on BlockSize",
			Assert: func() {
				cases := []struct {
					bsExpectedStr string
					bs            BlockSize
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
				}

				for _, c := range cases {
					assert.EqualValues(s.T(), c.bsExpectedStr, c.bs.String())
				}
			},
		},
	}

	RunTestCases(s, testCases)
}

func (s *Suite) TestUnitGetDefaultPasswordByte() {
	testCases := []TestCase{
		{
			Name: "success get default password",
			Assert: func() {
				got := getDefaultPassword()

				assert.EqualValues(s.T(), []byte(defaultPasswordStr), got)
			},
		},
	}

	RunTestCases(s, testCases)
}

func (s *Suite) TestUnitDefaultConfig() {
	testCases := []TestCase{
		{
			Name: "err get default config",
			PreRequisites: func() {
				s.testArchiveConfig = &ArchiveConfig{}
			},
			Assert: func() {
				gotAC := NewDefaultArchiveConfig()

				assert.NotEqualValues(s.T(), s.testArchiveConfig, &gotAC)
			},
		},
		{
			Name: "success get default config",
			PreRequisites: func() {
				s.testArchiveConfig = &ArchiveConfig{
					Password:          getDefaultPassword(),
					ArchiveType:       defaultArchiveType,
					BlockSize:         defaultBlockSize,
					VolumeSize:        defaultVolumeSize,
					FastBytes:         defaultFastBytes,
					DictSize:          defaultDictSize,
					HeadersEncryption: defaultHeadersEnc,
					Compression:       defaultCompressionLevel,
					SolidArchive:      defaultSolidArchive,
				}
			},
			Assert: func() {
				gotAC := NewDefaultArchiveConfig()

				assert.EqualValues(s.T(), s.testArchiveConfig, &gotAC)
			},
		},
	}

	RunTestCases(s, testCases)
}

func (s *Suite) TestUnitListFiles() {
	testCases := []TestCase{
		{
			Name: "err walking path",
			PreRequisites: func() {
				s.wfBuilder = func(fl *[]string) filepath.WalkFunc {
					return func(path string, f os.FileInfo, err error) error {
						return errors.New("triggered filepath.WalkFunc err")
					}
				}
			},
			Assert: func() {
				dotPath := fmt.Sprintf(".%c", os.PathSeparator)
				gotFileList, err := ListFiles(dotPath, s.wfBuilder)

				assert.NotNil(s.T(), err)
				assert.Equal(s.T(), 0, len(gotFileList))
			},
		},
		{
			// Note: Here we are testing just the handling within ListFiles func.
			// The actual success filepath.WalkFunc case is covered in the TestUnitDefaultPathWalkerFunc.
			Name: "success walking path",
			PreRequisites: func() {
				s.wfBuilder = func(fl *[]string) filepath.WalkFunc {
					return func(path string, f os.FileInfo, err error) error {
						return nil
					}
				}
			},
			Assert: func() {
				dotPath := fmt.Sprintf(".%c", os.PathSeparator)
				gotFileList, err := ListFiles(dotPath, s.wfBuilder)

				assert.Nil(s.T(), err)
				assert.Equal(s.T(), 0, len(gotFileList))
			},
		},
	}

	RunTestCases(s, testCases)
}

func (s *Suite) TestUnitDefaultPathWalkerFunc() {
	testCases := []TestCase{
		{
			Name: "success walking path",
			PreRequisites: func() {
				s.wfBuilder = DefaultPathWalkerFunc
			},
			Assert: func() {
				dotPath := fmt.Sprintf(".%c", os.PathSeparator)
				fileList := make([]string, 0)

				err := filepath.Walk(dotPath, s.wfBuilder(&fileList))

				assert.Nil(s.T(), err)
				assert.GreaterOrEqual(s.T(), len(fileList), 10)
			},
		},
		{
			Name: "err walking path - invalid path",
			PreRequisites: func() {
				s.wfBuilder = DefaultPathWalkerFunc

			},
			Assert: func() {
				fileList := make([]string, 0)

				err := filepath.Walk(invalidPath, s.wfBuilder(&fileList))

				assert.NotNil(s.T(), err)
				assert.Equal(s.T(), 0, len(fileList))
			},
		},
		{
			Name: "err walking path - does-not-exist path",
			PreRequisites: func() {
				s.wfBuilder = DefaultPathWalkerFunc

			},
			Assert: func() {
				fileList := make([]string, 0)

				err := filepath.Walk("does-not-exist", s.wfBuilder(&fileList))

				assert.NotNil(s.T(), err)
				assert.Equal(s.T(), 0, len(fileList))
			},
		},
	}

	RunTestCases(s, testCases)
}
