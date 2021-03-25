package ax

import (
	"os"

	"github.com/stretchr/testify/assert"
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
				ac.PathToArchive = "../tmp_to_archive"
				ac.OutputPath = "../tmp_archive_out"
				ac.NewArchiveName = "test_archive"
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

				assert.Nil(s.T(), err)
			},
		},
	}

	RunTestCases(s, testCases)
}
