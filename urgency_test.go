package libhealth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseUrgency(t *testing.T) {
	testcases := []struct {
		urgencies []string
		exp       Urgency
	}{
		{
			urgencies: []string{"REQUIRED", "Required", "required"},
			exp:       REQUIRED,
		},
		{
			urgencies: []string{"STRONG", "Strong", "strong"},
			exp:       STRONG,
		},
		{
			urgencies: []string{"WEAK", "Weak", "weak"},
			exp:       WEAK,
		},
		{
			urgencies: []string{"NONE", "None", "none"},
			exp:       NONE,
		},
		{
			urgencies: []string{"UNKNOWN", "Unknown", "unknown", "", "\t", "abc123"},
			exp:       UNKNOWN,
		},
	}

	for _, test := range testcases {
		for _, urgency := range test.urgencies {
			parsed := ParseUrgency(urgency)
			require.Equal(t, test.exp, parsed)
		}
	}
}

func Test_DowngradeWith_required(t *testing.T) {
	subject := REQUIRED
	require.Equal(t, OK, subject.DowngradeWith(OK, OK))
	require.Equal(t, MINOR, subject.DowngradeWith(OK, MINOR))
	require.Equal(t, MAJOR, subject.DowngradeWith(OK, MAJOR))
	require.Equal(t, OUTAGE, subject.DowngradeWith(OK, OUTAGE))
}

func Test_DowngradeWith_strong(t *testing.T) {
	subject := STRONG
	require.Equal(t, OK, subject.DowngradeWith(OK, OK))
	require.Equal(t, MINOR, subject.DowngradeWith(OK, MINOR))
	require.Equal(t, MAJOR, subject.DowngradeWith(OK, MAJOR))
	require.Equal(t, MAJOR, subject.DowngradeWith(OK, OUTAGE))
}

func Test_DowngradeWith_weak(t *testing.T) {
	subject := WEAK
	require.Equal(t, OK, subject.DowngradeWith(OK, OK))
	require.Equal(t, MINOR, subject.DowngradeWith(OK, MINOR))
	require.Equal(t, MINOR, subject.DowngradeWith(OK, MAJOR))
	require.Equal(t, MINOR, subject.DowngradeWith(OK, OUTAGE))
}

func Test_DowngradeWith_none(t *testing.T) {
	subject := NONE
	require.Equal(t, OK, subject.DowngradeWith(OK, OK))
	require.Equal(t, OK, subject.DowngradeWith(OK, MINOR))
	require.Equal(t, OK, subject.DowngradeWith(OK, MAJOR))
	require.Equal(t, OK, subject.DowngradeWith(OK, OUTAGE))
}

func Test_DowngradeWith_unknown(t *testing.T) {
	subject := UNKNOWN
	require.Equal(t, OK, subject.DowngradeWith(OK, OK))
	require.Equal(t, MINOR, subject.DowngradeWith(OK, MINOR))
	require.Equal(t, MAJOR, subject.DowngradeWith(OK, MAJOR))
	require.Equal(t, OUTAGE, subject.DowngradeWith(OK, OUTAGE))
}
