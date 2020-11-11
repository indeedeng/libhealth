package libhealth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Overall(t *testing.T) {
	results := []Result{
		{Health: Health{Status: OUTAGE, Urgency: WEAK}},
		{Health: Health{Status: OK, Urgency: REQUIRED}},
	}

	now := time.Date(2016, 4, 26, 16, 52, 0, 0, time.UTC)

	summary := NewSummary(now, results)

	require.Equal(t, MINOR, summary.Overall())
	require.Equal(t, now, summary.Executed())
}

func Test_SummaryState(t *testing.T) {
	results := []Result{
		{name: "foo1", Health: Health{Status: OUTAGE, Urgency: WEAK}},
		{name: "foo2", Health: Health{Status: MAJOR, Urgency: REQUIRED}},
		{name: "foo3", Health: Health{Status: OK, Urgency: REQUIRED}},
	}

	now := time.Date(2016, 4, 26, 16, 52, 0, 0, time.UTC)

	summary := NewSummary(now, results)

	require.Equal(t, OK, summary.Status())
	require.Equal(t, OK, summary.Status("foo3"))
	require.Equal(t, MAJOR, summary.Status("foo2"))
	require.Equal(t, OUTAGE, summary.Status("foo1"))
	require.Equal(t, MAJOR, summary.Status("foo2", "foo3"))
	require.Equal(t, MAJOR, summary.Status("foo3", "foo2"))
	require.Equal(t, OUTAGE, summary.Status("foo1", "foo2", "foo3"))
	require.Equal(t, OUTAGE, summary.Status("foo3", "foo2", "foo1"))
}

func Test_SummaryStateWithUrgency(t *testing.T) {
	results := []Result{
		{name: "foo1", Health: Health{Status: OUTAGE, Urgency: WEAK}},
		{name: "foo2", Health: Health{Status: MAJOR, Urgency: REQUIRED}},
		{name: "foo3", Health: Health{Status: OK, Urgency: REQUIRED}},
	}

	now := time.Date(2016, 4, 26, 16, 52, 0, 0, time.UTC)

	summary := NewSummary(now, results)

	require.Equal(t, OK, summary.StatusWithUrgency())
	require.Equal(t, OK, summary.StatusWithUrgency("foo3"))
	require.Equal(t, MAJOR, summary.StatusWithUrgency("foo2"))
	require.Equal(t, MINOR, summary.StatusWithUrgency("foo1"))
	require.Equal(t, MINOR, summary.StatusWithUrgency("foo1", "foo3"))
	require.Equal(t, MAJOR, summary.StatusWithUrgency("foo2", "foo3"))
	require.Equal(t, MAJOR, summary.StatusWithUrgency("foo3", "foo2"))
	require.Equal(t, MAJOR, summary.StatusWithUrgency("foo1", "foo2", "foo3"))
	require.Equal(t, MAJOR, summary.StatusWithUrgency("foo3", "foo2", "foo1"))
}
