package gauge

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/libhealth"
)

func floatList(floats ...float64) *list.List {
	l := list.New()
	for _, val := range floats {
		l.PushFront(val)
	}
	return l
}

func Test_MaxFloatThreshold_Apply(t *testing.T) {
	tests := []struct {
		values    *list.List
		threshold MaxFloatThreshold
		expState  libhealth.Status
	}{
		// only AnyN
		{
			floatList(1.1),
			MaxFloatThreshold{Threshold: 5.5, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(7.7),
			MaxFloatThreshold{Threshold: 5.4, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(1.1, 2.2, 3.3, 4.4),
			MaxFloatThreshold{Threshold: 5.5, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(1.1, 2.2, 3.3, 4.4, 5.5),
			MaxFloatThreshold{Threshold: 5.5, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(1.1, 2.2, 3.3, 4.4, 5.5),
			MaxFloatThreshold{Threshold: 5.5, AnyN: 2, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(1.1, 3.3, 5.5, 7.7, 9.9, 2.3, 3.3),
			MaxFloatThreshold{Threshold: 6.6, AnyN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(1.1, 3.3, 5.5, 7.7, 9.9, 2.2, 3.3, 4.4, 5.5),
			MaxFloatThreshold{Threshold: 6.6, AnyN: 3, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		// only LastN
		{
			floatList(1.1),
			MaxFloatThreshold{Threshold: 5.5, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(5.5),
			MaxFloatThreshold{Threshold: 5.5, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(3.3),
			MaxFloatThreshold{Threshold: 5.5, LastN: 9, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(4.4),
			MaxFloatThreshold{Threshold: 3.3, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 2.2),
			MaxFloatThreshold{Threshold: 5.5, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 2.2),
			MaxFloatThreshold{Threshold: 5.5, LastN: 3, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7),
			MaxFloatThreshold{Threshold: 5.5, LastN: 3, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(1, 2, 3, 4, 5, 6, 7, 1, 8),
			MaxFloatThreshold{Threshold: 5, LastN: 3, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		// combined AnyN and LastN
		{
			floatList(3.3),
			MaxFloatThreshold{Threshold: 5.5, AnyN: 1, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(3.3),
			MaxFloatThreshold{Threshold: 3.3, AnyN: 1, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(2.2, 3.3, 4.4, 5.5),
			MaxFloatThreshold{Threshold: 4.4, AnyN: 3, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(9.9, 7.7, 3.3, 2.2, 8.8, 1.1, 4.4),
			MaxFloatThreshold{Threshold: 4.4, AnyN: 3, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(4.4, 5.5, 6.6, 12.12, 3.3, 7.7, 8.7),
			MaxFloatThreshold{Threshold: 8.8, AnyN: 3, LastN: 3, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
	}

	for _, test := range tests {
		state, _ := test.threshold.Apply(test.values)
		require.Equal(t, test.expState, state, "expected state %s got %s", test.expState, state)
	}
}

func Test_MinFloatThreshold_Apply(t *testing.T) {
	tests := []struct {
		values    *list.List
		threshold MinFloatThreshold
		expState  libhealth.Status
	}{
		// only AnyN
		{
			floatList(1.1),
			MinFloatThreshold{Threshold: 3.3, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(5.5),
			MinFloatThreshold{Threshold: 3.3, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(4.4, 5.5, 6.6, 7.7, 8.8),
			MinFloatThreshold{Threshold: 4.4, AnyN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(4.4, 5.5, 6.6, 7.7, 8.8),
			MinFloatThreshold{Threshold: 4.4, AnyN: 2, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(3.3, 4.4, 7.7, 8.8, 1.1, 2.2, 1.1),
			MinFloatThreshold{Threshold: 3.3, AnyN: 4, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		// only LastN
		{
			floatList(5.5),
			MinFloatThreshold{Threshold: 3.3, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(1.1),
			MinFloatThreshold{Threshold: 3.3, LastN: 1, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(2.2, 5.5, 6.6, 1.1, 4.4, 4.4, 2.2),
			MinFloatThreshold{Threshold: 3.3, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(2.2, 5.5, 6.6, 1.1, 4.4, 4.4, 2.2),
			MinFloatThreshold{Threshold: 4.4, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		// combined AnyN and LastN
		{
			floatList(1.1),
			MinFloatThreshold{Threshold: 5.5, AnyN: 4, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.OK,
		},
		{
			floatList(1.1, 8.8, 2.2, 1.1, 1.1, 3.3, 7.7),
			MinFloatThreshold{Threshold: 5.5, AnyN: 4, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
		{
			floatList(9.9, 8.8, 9.9, 1.1, 7.7, 3.3, 5.5),
			MinFloatThreshold{Threshold: 5.5, AnyN: 4, LastN: 2, Severity: libhealth.MAJOR},
			libhealth.MAJOR,
		},
	}

	for _, test := range tests {
		state, _ := test.threshold.Apply(test.values)
		require.Equal(t, test.expState, state, "expected state %s got %s", test.expState, state)
	}
}

func gaugeFloats(g FloatGauger, vals ...float64) {
	for _, val := range vals {
		g.Gauge(val)
	}
}

func Test_Floats_AnyN(t *testing.T) {
	tests := []struct {
		gauge    FloatGauger
		expState libhealth.Status
	}{
		{
			gauge: func() FloatGauger {
				g, err := Floats("test", 1)
				require.NoError(t, err)
				return g
			}(), expState: libhealth.OK}, // default to OK if no values yet
		{
			gauge: func() FloatGauger {
				g, err := Floats("test", 1)
				require.NoError(t, err)
				g.Set(MaxFloatThreshold{
					Threshold: 5.5,
					AnyN:      1,
					Severity:  libhealth.MAJOR,
				})
				gaugeFloats(g, 1.1)
				return g
			}(),
			expState: libhealth.OK,
		},
		{
			gauge: func() FloatGauger {
				g, err := Floats("test", 1)
				require.NoError(t, err)
				g.Set(MaxFloatThreshold{
					Threshold: 5.5,
					AnyN:      1,
					Severity:  libhealth.MAJOR,
				})
				gaugeFloats(g, 7.7)
				return g
			}(),
			expState: libhealth.MAJOR,
		},
		{
			gauge: func() FloatGauger {
				g, err := Floats("test", 10)
				require.NoError(t, err)
				g.Set(MaxFloatThreshold{
					Threshold: 5.5,
					AnyN:      3,
					Severity:  libhealth.MAJOR,
				})
				gaugeFloats(g, 9.9, 1.1, 5.5, 1.1, 3.3, 2.2, 0.0, 1.1, 2.2)
				return g
			}(),
			expState: libhealth.OK,
		},
		{
			gauge: func() FloatGauger {
				g, err := Floats("test", 10)
				require.NoError(t, err)
				g.Set(MaxFloatThreshold{
					Threshold: 5.5,
					AnyN:      3,
					Severity:  libhealth.MAJOR,
				})
				gaugeFloats(g,
					9.9, 5.5, 1.1, 8.8, 3.3, 6.6, 1.1,
					3.3, 6.6, 1.1, 5.5, 7.7, 4.4, 1.1, 4.4,
				)
				return g
			}(),
			expState: libhealth.MAJOR,
		},
		{
			gauge: func() FloatGauger {
				g, err := Floats("test", 10)
				require.NoError(t, err)
				g.Set(MaxFloatThreshold{
					Threshold: 5.5,
					AnyN:      3,
					Severity:  libhealth.MINOR,
				}).Set(MaxFloatThreshold{
					Threshold: 8.8,
					AnyN:      2,
					Severity:  libhealth.OUTAGE,
				})
				gaugeFloats(g, 1.1, 4.4, 2.2, 9.9, 4.4, 2.2, 1.1, 4.4, 8.8)
				return g
			}(),
			expState: libhealth.OUTAGE,
		},
		{
			gauge: func() FloatGauger {
				g, err := Floats("test", 10)
				require.NoError(t, err)
				g.Set(MaxFloatThreshold{
					Threshold: 4.4,
					AnyN:      2,
					Severity:  libhealth.MINOR,
				}).Set(MaxFloatThreshold{
					Threshold: 5.5,
					AnyN:      2,
					Severity:  libhealth.MAJOR,
				}).Set(MaxFloatThreshold{
					Threshold: 6.6,
					AnyN:      3,
					Severity:  libhealth.OUTAGE,
				})
				gaugeFloats(g, 2.2, 9.9, 4.4, 7.7, 1.1, 4.4, 4.4)
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

func Test_Floats_Mixed(t *testing.T) {
	tests := []struct {
		gauge    FloatGauger
		expState libhealth.Status
		expMsg   string
	}{
		{
			gauge: func() FloatGauger {
				g, err := Floats("test", 5)
				require.NoError(t, err)
				g.Set(MinFloatThreshold{
					Threshold:   3.3,
					LastN:       2,
					Severity:    libhealth.MINOR,
					Description: "min1",
				}).Set(MinFloatThreshold{
					Threshold:   1.1,
					LastN:       2,
					Severity:    libhealth.MAJOR,
					Description: "min2",
				}).Set(MaxFloatThreshold{
					Threshold:   8.8,
					LastN:       2,
					Severity:    libhealth.MINOR,
					Description: "max1",
				}).Set(MaxFloatThreshold{
					Threshold:   9.9,
					LastN:       2,
					Severity:    libhealth.MAJOR,
					Description: "max2",
				})
				gaugeFloats(g,
					5.5, 6.6, 4.4, 8.8, 9.9, 7.7, 1.1, 5.5, 2.2, 1.1,
				)
				return g
			}(),
			expState: libhealth.MINOR,
			expMsg:   "min1",
		},
		{
			gauge: func() FloatGauger {
				g, err := Floats("test", 10)
				require.NoError(t, err)
				g.Set(MinFloatThreshold{
					Threshold:   3.4,
					AnyN:        2,
					Severity:    libhealth.MINOR,
					Description: "min1",
				}).Set(MinFloatThreshold{
					Threshold:   2.2,
					AnyN:        2,
					Severity:    libhealth.MAJOR,
					Description: "min2",
				}).Set(MinFloatThreshold{
					Threshold:   1.1,
					AnyN:        1,
					Severity:    libhealth.OUTAGE,
					Description: "min3",
				}).Set(MaxFloatThreshold{
					Threshold:   6.6,
					AnyN:        2,
					Severity:    libhealth.MINOR,
					Description: "max1",
				}).Set(MaxFloatThreshold{
					Threshold:   7.7,
					AnyN:        2,
					Severity:    libhealth.MAJOR,
					Description: "max2",
				}).Set(MaxFloatThreshold{
					Threshold:   8.8,
					AnyN:        3,
					Severity:    libhealth.OUTAGE,
					Description: "max3",
				})
				gaugeFloats(g,
					2.2, 4.4, 2.2, 1.1, 8.8, 4.4,
					7.7, 6.6, 9.9, 12.12, 2.2,
				)
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
