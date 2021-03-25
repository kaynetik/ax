package ax

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	// ac - reference to ArchiveConfig.
	ac *ArchiveConfig
	// ec - reference to ArchiveConfig.
	ec *ExtractConfig

	// wfBuilder - instance of walkFuncBuilder.
	wfBuilder walkFuncBuilder
}

func TestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(Suite))
}
