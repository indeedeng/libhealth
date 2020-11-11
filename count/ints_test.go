package count

import (
	"testing"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/libhealth"
)

func makeTicker() (BucketPeriod, *int) {
	tick := 0
	return func() int { return tick }, &tick
}

func check(
	t *testing.T,
	h libhealth.Health,
	expState libhealth.Status,
	expMessage string) {
	require.Equal(t, expState, h.Status, "exp state %s got %s", expState, h.Status)
	require.Equal(t, libhealth.Message(expMessage), h.Message)
}

func Test_MaxIntThreshold(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Ints("test-ints-max", period, 3)
	require.NoError(t, err)
	*tick = 0

	// 0, 0, 0 (no threshold)
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Set(MaxIntThreshold{
		Threshold:   5,
		Description: "max5",
		Severity:    libhealth.MAJOR,
	})

	counter.Increment(0) // 0, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(2) // 2, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(5) // 7, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	*tick = 1 // 0, 7, 0
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	counter.Increment(4) // 4, 7, 0
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	counter.Increment(2) // 6, 7, 0
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	*tick = 2 // 0, 6, 7
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	counter.Increment(3) // 3, 6, 7
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	*tick = 3 // 0, 3, 6
	counter.Increment(0)
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(5) // 5, 3, 6
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Set(MaxIntThreshold{
		Threshold:   3,
		Description: "max3",
		Severity:    libhealth.MINOR,
	})

	check(t, counter.Health(), libhealth.MINOR, "max3")

	*tick = 4 // 0, 5, 3
	counter.Increment(0)
	check(t, counter.Health(), libhealth.MAJOR, "max5")

	*tick = 5 // 0, 0, 5
	check(t, counter.Health(), libhealth.OK, "ok")
}

func Test_MinIntThreshold(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Ints("test-ints-min", period, 3)
	require.NoError(t, err)
	*tick = 0

	// 0, 0, 0 (no threshold)
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Set(MinIntThreshold{
		Threshold:   5,
		Description: "min5",
		Severity:    libhealth.MAJOR,
	})

	counter.Increment(0) // 0, 0, 0
	check(t, counter.Health(), libhealth.MAJOR, "min5")

	counter.Increment(7) // 7, 0, 0
	check(t, counter.Health(), libhealth.MAJOR, "min5")

	*tick = 1
	counter.Increment(3) // 3, 7, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	*tick = 2
	counter.Increment(8) // 8, 3, 7
	check(t, counter.Health(), libhealth.MAJOR, "min5")

	counter.Set(MinIntThreshold{
		Threshold:   9,
		Description: "min9",
		Severity:    libhealth.MINOR,
	})

	check(t, counter.Health(), libhealth.MAJOR, "min5")

	*tick = 3
	counter.Increment(5) // 5, 8, 3
	check(t, counter.Health(), libhealth.MINOR, "min9")

	*tick = 4
	counter.Increment(11) // 11, 5, 8
	check(t, counter.Health(), libhealth.MAJOR, "min5")

	*tick = 5
	counter.Increment(1) // 1, 11, 5
	check(t, counter.Health(), libhealth.OK, "ok")
}

func Test_MaxSumIntThreshold(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Ints("test-ints-max-sum", period, 3)
	require.NoError(t, err)
	counter.Set(MaxSumIntThreshold{
		Threshold:   8,
		Description: "max sum8",
		Severity:    libhealth.MAJOR,
	})

	*tick = 0 // 0, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	*tick = 1
	counter.Increment(4) // 4, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	*tick = 2
	counter.Increment(5) // 5, 4, 0
	check(t, counter.Health(), libhealth.MAJOR, "max sum8")

	*tick = 4
	counter.Increment(2) // 2, 0, 4
	check(t, counter.Health(), libhealth.OK, "ok")
}

func Test_MinSumIntThreshold(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Ints("test-ints-min-sum", period, 3)
	require.NoError(t, err)
	counter.Set(MinSumIntThreshold{
		Threshold:   3,
		Description: "min sum3",
		Severity:    libhealth.MAJOR,
	})

	*tick = 0 // 0, 0, 0
	check(t, counter.Health(), libhealth.MAJOR, "min sum3")

	counter.Increment(2) // 2, 0, 0
	check(t, counter.Health(), libhealth.MAJOR, "min sum3")

	*tick = 2
	counter.Increment(3) // 3, 0, 2
	check(t, counter.Health(), libhealth.OK, "ok")
}

func Test_IntMix(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Ints("test-ints-mix", period, 5)
	require.NoError(t, err)
	*tick = 0 // 0, 0, 0, 0, 0

	counter.Set(MaxIntThreshold{
		Threshold:   10,
		Description: "max10",
		Severity:    libhealth.OUTAGE,
	}).Set(MaxIntThreshold{
		Threshold:   8,
		Description: "max8",
		Severity:    libhealth.MAJOR,
	}).Set(MinIntThreshold{
		Threshold:   3,
		Description: "min3",
		Severity:    libhealth.MINOR,
	}).Set(MinIntThreshold{
		Threshold:   0,
		Description: "min0",
		Severity:    libhealth.OUTAGE,
	}).Set(MinSumIntThreshold{
		Threshold:   -1,
		Description: "min sum-1",
		Severity:    libhealth.MAJOR,
	}).Set(MaxSumIntThreshold{
		Threshold:   100,
		Description: "max sum100",
		Severity:    libhealth.OUTAGE,
	})

	check(t, counter.Health(), libhealth.OUTAGE, "min0")

	counter.Increment(2)
	*tick = 1 // 0, 2, 0, 0, 0
	check(t, counter.Health(), libhealth.MINOR, "min3")

	counter.Increment(4)
	*tick = 2 // 0, 4, 2, 0, 0
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(9)
	*tick = 3 // 0, 9, 4, 2, 0
	check(t, counter.Health(), libhealth.MAJOR, "max8")

	counter.Increment(11)
	*tick = 4 // 0, 11, 9, 4, 2
	check(t, counter.Health(), libhealth.OUTAGE, "max10")

	counter.Increment(4)
	*tick = 5 // 4, 11, 9, 4, 2
	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(-50) // -46, 11, 9, 4, 2
	check(t, counter.Health(), libhealth.MAJOR, "min sum-1")

	counter.Increment(1000) // 954, 11, 9, 4, 2
	check(t, counter.Health(), libhealth.OUTAGE, "max sum100")
}

func Test_Ints_Zero(t *testing.T) {
	period, _ := makeTicker()
	_, err := Ints("test-ints-zero", period, 0)
	require.Error(t, err)
}

func Test_Ints_One(t *testing.T) {
	period, tick := makeTicker()
	counter, err := Ints("test-ints-single", period, 1)
	require.NoError(t, err)
	*tick = 0

	counter.Set(MaxIntThreshold{
		Threshold:   5,
		Description: "max5",
		Severity:    libhealth.MAJOR,
	})

	check(t, counter.Health(), libhealth.OK, "ok")

	counter.Increment(4)
	check(t, counter.Health(), libhealth.OK, "ok")
	*tick = 1

	counter.Increment(11)
	check(t, counter.Health(), libhealth.MAJOR, "max5")
}
