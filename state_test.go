package libhealth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseState(t *testing.T) {
	testcases := []struct {
		inputs []string
		exp    Status
	}{
		{
			inputs: []string{"OUTAGE", "Outage", "outage"},
			exp:    OUTAGE,
		},
		{
			inputs: []string{"MAJOR", "Major", "major"},
			exp:    MAJOR,
		},
		{
			inputs: []string{"MINOR", "Minor", "minor"},
			exp:    MINOR,
		},
		{
			inputs: []string{"OK", "Ok", "ok"},
			exp:    OK,
		},
		{
			inputs: []string{"", " ", "\t", "abc123"},
			exp:    OUTAGE,
		},
	}

	for _, test := range testcases {
		for _, input := range test.inputs {
			parsed := ParseStatus(input)
			require.Equal(t, test.exp, parsed)
		}
	}
}

func Test_State(t *testing.T) {
	require.True(t, MINOR.WorseThan(OK))
	require.False(t, MINOR.WorseThan(MAJOR))
	require.False(t, MINOR.WorseThan(MINOR))
	require.True(t, MINOR.SameOrWorseThan(OK))
	require.True(t, MINOR.SameOrWorseThan(MINOR))
	require.False(t, MINOR.SameOrWorseThan(MAJOR))
	require.True(t, MINOR.SameOrBetterThan(MAJOR))
	require.True(t, MINOR.SameOrBetterThan(MINOR))
	require.False(t, MINOR.SameOrBetterThan(OK))
	require.True(t, MINOR.BetterThan(MAJOR))
	require.False(t, MINOR.BetterThan(MINOR))
	require.False(t, MINOR.BetterThan(OK))
	require.True(t, MINOR.SameAs(MINOR))
	require.False(t, MINOR.SameAs(MAJOR))
}

func Test_BestState(t *testing.T) {
	require.Equal(t, OK, BestState(MINOR, OK))
	require.Equal(t, OK, BestState(OK, MINOR))
	require.Equal(t, MINOR, BestState(MINOR, MINOR))
}

func Test_WorstState(t *testing.T) {
	require.Equal(t, MINOR, WorstState(MINOR, OK))
	require.Equal(t, MINOR, WorstState(OK, MINOR))
	require.Equal(t, MINOR, WorstState(MINOR, MINOR))
}
