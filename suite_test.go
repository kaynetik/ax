package ax

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite

	testArchiveConfig *ArchiveConfig
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
