package ax

import (
	"errors"
	"os"
	"strings"

	"github.com/stretchr/testify/assert"
)

const (
	gitTestRepo = "git@github.com:kaynetik/test-bk.git"
)

func (s *Suite) TestUnitPushToGIT() {
	testCases := []TestCase{
		{
			Name: "err git add remote",
			PreRequisites: func() {
				outFilePath := "./tests/git_test"
				outFile := outFilePath + "/lorem.md"
				_ = os.Mkdir(outFilePath, os.ModePerm)
				copyFileToEnc(s.T(), testLoremInFile, outFile)
			},
			Assert: func() {
				err := PushToGIT(gitTestRepo)

				assert.NotNil(s.T(), err)
			},
		},
	}

	RunTestCases(s, testCases)
}

func (s *Suite) TestUnitCmdErrWrapper() {
	testCases := []TestCase{
		{
			Name: "success build err",
			Assert: func() {
				cmd := "non-existent"
				cmdArgs := "-b few -c fake -d args"
				wrapErrStr := "wrap-this"
				toWrapErr := errors.New(wrapErrStr)
				err := ErrCmdWrapFn(cmd, cmdArgs, toWrapErr)

				assert.NotNil(s.T(), err)

				assert.True(s.T(), strings.Contains(err.Error(), wrapErrStr))
				assert.True(s.T(), strings.Contains(err.Error(), cmdArgs))
				assert.True(s.T(), strings.Contains(err.Error(), cmd))
			},
		},
	}

	RunTestCases(s, testCases)
}
