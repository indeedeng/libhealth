package gauge

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/libhealth"
)

func intList(ints ...int) *list.List {
	l := list.New()
	for _, val := range ints {
		l.PushFront(val)
	}
	return l
}

func Test_MaxIntThreshold_Apply(t *testing.T) {
	tests := []struct {
		values    *list.List
		threshold MaxIntThreshold
		expState  libhealth.Status
	}{
		// only AnyN
		{
			intList(1),
			MaxIntThreshold{Threshold: 5, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(7),
			MaxIntThreshold{Threshold: 5, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(1, 2, 3, 4),
			MaxIntThreshold{Threshold: 5, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(1, 2, 3, 4, 5),
			MaxIntThreshold{Threshold: 5, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(1, 2, 3, 4, 5),
			MaxIntThreshold{Threshold: 5, AnyN: 2, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(1, 3, 5, 7, 9, 2, 3),
			MaxIntThreshold{Threshold: 6, AnyN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(1, 3, 5, 7, 9, 2, 3, 4, 5),
			MaxIntThreshold{Threshold: 6, AnyN: 3, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		// only LastN
		{
			intList(1),
			MaxIntThreshold{Threshold: 5, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(5),
			MaxIntThreshold{Threshold: 5, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(3),
			MaxIntThreshold{Threshold: 5, LastN: 9, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(4),
			MaxIntThreshold{Threshold: 3, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(1, 2, 3, 4, 5, 6, 2),
			MaxIntThreshold{Threshold: 5, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(1, 2, 3, 4, 5, 6, 2),
			MaxIntThreshold{Threshold: 5, LastN: 3, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(1, 2, 3, 4, 5, 6, 7),
			MaxIntThreshold{Threshold: 5, LastN: 3, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(1, 2, 3, 4, 5, 6, 7, 1, 8),
			MaxIntThreshold{Threshold: 5, LastN: 3, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		// combined AnyN and LastN
		{
			intList(3),
			MaxIntThreshold{Threshold: 5, AnyN: 1, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(3),
			MaxIntThreshold{Threshold: 3, AnyN: 1, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(2, 3, 4, 5),
			MaxIntThreshold{Threshold: 4, AnyN: 3, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(9, 7, 3, 2, 8, 1, 4),
			MaxIntThreshold{Threshold: 4, AnyN: 3, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(4, 5, 6, 12, 3, 7, 8),
			MaxIntThreshold{Threshold: 8, AnyN: 3, LastN: 3, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
	}

	for _, test := range tests {
		state, _ := test.threshold.Apply(test.values)
		require.Equal(t, test.expState, state, "expected state %s got %s", test.expState, state)
	}
}

func Test_MinIntThreshold_Apply(t *testing.T) {
	tests := []struct {
		values    *list.List
		threshold MinIntThreshold
		expState  libhealth.Status
	}{
		// only AnyN
		{
			intList(1),
			MinIntThreshold{Threshold: 3, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(5),
			MinIntThreshold{Threshold: 3, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(4, 5, 6, 7, 8),
			MinIntThreshold{Threshold: 4, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(4, 5, 6, 7, 8),
			MinIntThreshold{Threshold: 4, AnyN: 2, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(3, 4, 7, 8, 1, 2, 1),
			MinIntThreshold{Threshold: 3, AnyN: 4, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		// only LastN
		{
			intList(5),
			MinIntThreshold{Threshold: 3, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(1),
			MinIntThreshold{Threshold: 3, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(2, 5, 6, 1, 4, 4, 2),
			MinIntThreshold{Threshold: 3, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(2, 5, 6, 1, 4, 4, 2),
			MinIntThreshold{Threshold: 4, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		// combined AnyN and LastN
		{
			intList(1),
			MinIntThreshold{Threshold: 5, AnyN: 4, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			intList(1, 8, 2, 1, 1, 3, 7),
			MinIntThreshold{Threshold: 5, AnyN: 4, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			intList(9, 8, 9, 1, 7, 3, 5),
			MinIntThreshold{Threshold: 5, AnyN: 4, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
	}

	for _, test := range tests {
		state, _ := test.threshold.Apply(test.values)
		require.Equal(t, test.expState, state, "expected state %s got %s", test.expState, state)
	}
}

func gaugeInts(g IntGauger, vals ...int) {
	for _, val := range vals {
		g.Gauge(val)
	}
}

func Test_Ints_AnyN(t *testing.T) {
	tests := []struct {
		gauge    IntGauger
		expState libhealth.Status
	}{
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 1)
				require.NoError(t, err)
				return g
			}(),
			expState: libhealth.OK,
		},
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 1)
				require.NoError(t, err)
				g.Set(MaxIntThreshold{
					Threshold: 5,
					AnyN:      1,
					Severity:  libhealth.MAJOR,
				})
				gaugeInts(g, 1)
				return g
			}(),
			expState: libhealth.OK,
		},
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 1)
				require.NoError(t, err)
				g.Set(MaxIntThreshold{
					Threshold: 5,
					AnyN:      1,
					Severity:  libhealth.MAJOR,
				})
				gaugeInts(g, 7)
				return g
			}(),
			expState: libhealth.MAJOR,
		},
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 10)
				require.NoError(t, err)
				g.Set(MaxIntThreshold{
					Threshold: 5,
					AnyN:      3,
					Severity:  libhealth.MAJOR,
				})
				gaugeInts(g, 9, 1, 5, 1, 3, 2, 0, 1, 2)
				return g
			}(),
			expState: libhealth.OK,
		},
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 10)
				require.NoError(t, err)
				g.Set(MaxIntThreshold{
					Threshold: 5,
					AnyN:      3,
					Severity:  libhealth.MAJOR,
				})
				gaugeInts(g, 9, 5, 1, 8, 3, 6, 1, 3, 6, 1, 5, 7, 4, 1, 4)
				return g
			}(),
			expState: libhealth.MAJOR,
		},
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 10)
				require.NoError(t, err)
				g.Set(MaxIntThreshold{
					Threshold: 5,
					AnyN:      3,
					Severity:  libhealth.MINOR,
				}).Set(MaxIntThreshold{
					Threshold: 8,
					AnyN:      2,
					Severity:  libhealth.OUTAGE,
				})
				gaugeInts(g, 1, 4, 2, 9, 4, 2, 1, 4, 8)
				return g
			}(),
			expState: libhealth.OUTAGE,
		},
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 10)
				require.NoError(t, err)
				g.Set(MaxIntThreshold{
					Threshold: 4,
					AnyN:      2,
					Severity:  libhealth.MINOR,
				}).Set(MaxIntThreshold{
					Threshold: 5,
					AnyN:      2,
					Severity:  libhealth.MAJOR,
				}).Set(MaxIntThreshold{
					Threshold: 6,
					AnyN:      3,
					Severity:  libhealth.OUTAGE,
				})
				gaugeInts(g, 2, 9, 4, 7, 1, 4, 4)
				return g
			}(),
			expState: libhealth.MAJOR,
		},
	}

	for _, test := range tests {
		state := test.gauge.Health().Status
		require.Equal(t, test.expState, state, "expected %s got %s", test.expState, state)
	}
}

func Test_Ints_LastN(t *testing.T) {
	tests := []struct {
		gauge    IntGauger
		expState libhealth.Status
	}{
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 1)
				require.NoError(t, err)
				return g
			}(),
			expState: libhealth.OK,
		},
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 1)
				require.NoError(t, err)
				g.Set(MaxIntThreshold{
					Threshold: 4,
					LastN:     1,
					Severity:  libhealth.MAJOR,
				})
				gaugeInts(g, 4, 1, 4, 5, 7)
				return g
			}(),
			expState: libhealth.MAJOR,
		},
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 1)
				require.NoError(t, err)
				g.Set(MaxIntThreshold{
					Threshold: 4,
					LastN:     1,
					Severity:  libhealth.MAJOR,
				})
				gaugeInts(g, 4, 1, 4, 5, 7)
				return g
			}(),
			expState: libhealth.MAJOR,
		},
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 10)
				require.NoError(t, err)
				g.Set(MaxIntThreshold{
					Threshold: 2,
					LastN:     2,
					Severity:  libhealth.MINOR,
				}).Set(MaxIntThreshold{
					Threshold: 3,
					LastN:     3,
					Severity:  libhealth.MAJOR,
				}).Set(MaxIntThreshold{
					Threshold: 4,
					LastN:     4,
					Severity:  libhealth.OUTAGE,
				})
				gaugeInts(g, 1, 1, 3, 3, 9)
				return g
			}(),
			expState: libhealth.MAJOR,
		},
	}

	for _, test := range tests {
		state := test.gauge.Health().Status
		require.Equal(t, test.expState, state, "expected %s got %s", test.expState, state)
	}
}

func Test_Ints_Mixed(t *testing.T) {
	tests := []struct {
		gauge    IntGauger
		expState libhealth.Status
		expMsg   string
	}{
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 5)
				require.NoError(t, err)
				g.Set(MinIntThreshold{
					Threshold:   3,
					LastN:       2,
					Severity:    libhealth.MINOR,
					Description: "min1",
				}).Set(MinIntThreshold{
					Threshold:   1,
					LastN:       2,
					Severity:    libhealth.MAJOR,
					Description: "min2",
				}).Set(MaxIntThreshold{
					Threshold:   8,
					LastN:       2,
					Severity:    libhealth.MINOR,
					Description: "max1",
				}).Set(MaxIntThreshold{
					Threshold:   9,
					LastN:       2,
					Severity:    libhealth.MAJOR,
					Description: "max2",
				})
				gaugeInts(g, 5, 6, 4, 8, 9, 7, 1, 5, 2, 1)
				return g
			}(),
			expState: libhealth.MINOR,
			expMsg:   "min1",
		},
		{
			gauge: func() IntGauger {
				g, err := Ints("test", 10)
				require.NoError(t, err)
				g.Set(MinIntThreshold{
					Threshold:   3,
					AnyN:        2,
					Severity:    libhealth.MINOR,
					Description: "min1",
				}).Set(MinIntThreshold{
					Threshold:   2,
					AnyN:        2,
					Severity:    libhealth.MAJOR,
					Description: "min2",
				}).Set(MinIntThreshold{
					Threshold:   1,
					AnyN:        1,
					Severity:    libhealth.OUTAGE,
					Description: "min3",
				}).Set(MaxIntThreshold{
					Threshold:   6,
					AnyN:        2,
					Severity:    libhealth.MINOR,
					Description: "max1",
				}).Set(MaxIntThreshold{
					Threshold:   7,
					AnyN:        2,
					Severity:    libhealth.MAJOR,
					Description: "max2",
				}).Set(MaxIntThreshold{
					Threshold:   8,
					AnyN:        3,
					Severity:    libhealth.OUTAGE,
					Description: "max3",
				})
				gaugeInts(g, 2, 4, 2, 1, 8, 4, 7, 6, 9, 12, 2)
				return g
			}(),
			expState: libhealth.OUTAGE,
			expMsg:   "min3, max3",
		},
	}

	for _, test := range tests {
		health := test.gauge.Health()
		state := health.Status
		msg := string(health.Message)
		require.Equal(t, test.expState, state, "expected %s got %s", test.expState, state)
		require.Equal(t, test.expMsg, msg)
	}
}

func Test_Ints_Zero(t *testing.T) {
	_, err := Ints("test-ints-zero", 0)
	require.Error(t, err)
}
