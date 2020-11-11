package libhealth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HealthSuite struct {
	suite.Suite
}

func TestHealthSuite(t *testing.T) {
	suite.Run(t, new(HealthSuite))
}

func (s *HealthSuite) TestString() {
	h := Health{
		Status:   MAJOR,
		Urgency:  STRONG,
		Time:     time.Date(2015, 12, 14, 11, 19, 0, 0, time.UTC),
		Message:  "this is a test",
		Duration: 21 * time.Millisecond,
	}
	exp := "MAJOR STRONG at 2015-12-14 11:19:00 +0000 UTC, this is a test"
	require.Equal(s.T(), exp, h.String())
}
