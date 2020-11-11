package count

import (
	"testing"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/libhealth"
)

func Test_MaxFloatThreshold(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Floats("test-floats-max", period, 3)
	require.NoError(t, err)

	*tick = 0

	// 0, 0, 0 (no threshold)
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Set(MaxFloatThreshold{
		Threshold:   5.0,
		Description: "max5",
		Severity:    libhealth.MAJOR,
	})

	counter.Increment(0) // 0, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(2.1) // 2.1, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(5.1) // 7.2, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	*tick = 1 // 0, 7.2, 0
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	counter.Increment(4.1) // 4.1, 7.2, 0
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	counter.Increment(2.1) // 6.2, 7.2, 0
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	*tick = 2 // 0, 6.2, 7.2
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	counter.Increment(3.1) // 3.1, 6.2, 7.2
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	*tick = 3 // 0, 3.1, 6.2
	counter.Increment(0)
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(5.1) // 5.1, 3.1, 6.2
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Set(MaxFloatThreshold{
		Threshold:   6,
		Description: "max6",
		Severity:    libhealth.OUTAGE,
	})

	*tick = 4 // 7.1, 5.1, 3.1
	counter.Increment(7.1)
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	*tick = 5
	counter.Increment(1.1) // 1.1, 7.1, 5.1
	check(t, counter.Health(), libhealth.OUTAGE, "max6")
}

func Test_MinFloatThreshold(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Floats("test-floats-min", period, 3)
	require.NoError(t, err)

	*tick = 0

	// 0, 0, 0 (no threshold)
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Set(MinFloatThreshold{
		Threshold:   5.2,
		Description: "min5",
		Severity:    libhealth.MAJOR,
	})

	counter.Increment(0) // 0, 0, 0
	check(t, counter.Health(), libhealth.MAJOR, "min5")

	counter.Increment(7.1) // 7.1, 0, 0
	check(t, counter.Health(), libhealth.MAJOR, "min5")

	*tick = 1
	counter.Increment(3.1) // 3.1, 7.1, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	*tick = 2
	counter.Increment(8.1) // 8.1, 3.1, 7.1
	check(t, counter.Health(), libhealth.MAJOR, "min5")

	counter.Set(MinFloatThreshold{
		Threshold:   4,
		Description: "min4",
		Severity:    libhealth.OUTAGE,
	})

	check(t, counter.Health(), libhealth.OUTAGE, "min4")

	*tick = 3
	counter.Increment(3.1) // 3.1, 8.1, 3.1
	check(t, counter.Health(), libhealth.OK, "ok")

	*tick = 4
	counter.Increment(5.1) // 5.1, 3.1, 8.1
	check(t, counter.Health(), libhealth.OUTAGE, "min4")
}

func Test_MaxSumFloatThreshold(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Floats("test-floats-max-sum", period, 3)
	require.NoError(t, err)

	counter.Set(MaxSumFloatThreshold{
		Threshold:   8.8,
		Severity:    libhealth.MAJOR,
		Description: "max sum8.8",
	})

	*tick = 0 // 0, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	*tick = 1
	counter.Increment(4.5) // 4.5, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	*tick = 2
	counter.Increment(4.5) // 4.5, 4.5, 0
	check(t, counter.Health(), libhealth.MAJOR, "max sum8.8")

	*tick = 4
	counter.Increment(2.5) // 2.5, 0, 4.5
	check(t, counter.Health(), libhealth.OK, "ok")
}

func Test_MinSumFloatThreshold(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Floats("test-floats-min-sum", period, 3)
	require.NoError(t, err)

	counter.Set(MinSumFloatThreshold{
		Threshold:   3.0,
		Description: "min sum3",
		Severity:    libhealth.MAJOR,
	})

	*tick = 0 // 0, 0, 0
	check(t, counter.Health(), libhealth.MAJOR, "min sum3")

	counter.Increment(2.1) // 2.1, 0, 0
	check(t, counter.Health(), libhealth.MAJOR, "min sum3")

	*tick = 2
	counter.Increment(3.1) // 3.1, 0, 2.1
	check(t, counter.Health(), libhealth.OK, "ok")
}

func Test_FloatMix(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Floats("test-floats-mix", period, 5)
	require.NoError(t, err)

	*tick = 0 // 0, 0, 0, 0, 0

	counter.Set(MaxFloatThreshold{
		Threshold:   10.5,
		Description: "max10",
		Severity:    libhealth.OUTAGE,
	}).Set(MaxFloatThreshold{
		Threshold:   8.5,
		Description: "max8",
		Severity:    libhealth.MAJOR,
	}).Set(MinFloatThreshold{
		Threshold:   3.5,
		Description: "min3",
		Severity:    libhealth.MINOR,
	}).Set(MinFloatThreshold{
		Threshold:   0,
		Description: "min0",
		Severity:    libhealth.OUTAGE,
	}).Set(MinSumFloatThreshold{
		Threshold:   -1.5,
		Description: "min sum-1.5",
		Severity:    libhealth.MAJOR,
	}).Set(MaxSumFloatThreshold{
		Threshold:   100.5,
		Description: "max sum100",
		Severity:    libhealth.OUTAGE,
	})

	check(t, counter.Health(), libhealth.OUTAGE, "min0")

	counter.Increment(1.1)
	*tick = 1 // 0, 1.1, 0, 0, 0
	check(t, counter.Health(), libhealth.MINOR, "min3")

	counter.Increment(3.6)
	*tick = 2 // 0, 3.6, 1.1, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(9)
	*tick = 3 // 0, 9, 3.6, 1.1, 0
	check(t, counter.Health(), libhealth.MAJOR, "max8")

	counter.Increment(11)
	*tick = 4 // 0, 11, 9, 3,6, 1.1
	check(t, counter.Health(), libhealth.OUTAGE, "max10")

	counter.Increment(3.7)
	*tick = 5 // 3.7, 11, 9, 3.6, 1.1
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(-50) // -46.3, 11, 9, 3.6, 1.1
	check(t, counter.Health(), libhealth.MAJOR, "min sum-1.5")

	counter.Increment(1000) // 953.7, 11, 9, 3.6, 1.1
	check(t, counter.Health(), libhealth.OUTAGE, "max sum100")
}
