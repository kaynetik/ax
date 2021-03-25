package ax

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	testArchiveConfig *ArchiveConfig
	wfBuilder         walkFuncBuilder
}

func TestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(Suite))
}
