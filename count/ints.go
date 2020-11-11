package count

import (
	"oss.indeed.com/go/libhealth"
	"oss.indeed.com/go/libhealth/count/internal/data"
)

// Ints will create an IntCounter.
func Ints(varname string, size BucketPeriod, length int) (IntCounter, error) {
	c, err := newContainer(varname, size, length, data.NewInt(0))
	return (*ints)(c), err
}

type ints container

func (i *ints) Increment(delta int) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.buckets.Increment(data.NewInt(delta))
}

func (i *ints) Health() libhealth.Health {
	return (*container)(i).health()
}

func (i *ints) Set(threshold Threshold) Countable {
	return (*ints)((*container)(i).set(threshold))
}

type intThreshold struct {
	// For Min/Max, Threshold represents the limit of the
	// value of the current bucket.
	// For SumMin/SumMax, Threshold represents the limit of
	// the values of all buckets summed together.
	Threshold   int
	Description string
	Severity    libhealth.Status
}

// MaxIntThreshold represents a maximum value that will trigger
// an unhealthy health.Status if at least n Infringements are
// exceeded.
type MaxIntThreshold intThreshold

// MinIntThreshold represents a mininum value that will trigger
// an unhealthy health.Status if at least n Infringements are
// exceeded.
type MinIntThreshold intThreshold

func (t intThreshold) state(violated bool) (libhealth.Status, string) {
	if violated {
		return t.Severity, t.Description
	}
	return libhealth.OK, OkMessage
}

// Apply the limits of t over buckets.
//
// Caller is responsible for protecting buckets with a lock.
func (t MaxIntThreshold) Apply(buckets *data.Buckets) (libhealth.Status, string) {
	exceeded := buckets.Compare(data.GreaterEq, data.NewInt(t.Threshold))
	return intThreshold(t).state(exceeded)
}

// Apply the limits of t over buckets.
//
// Caller is responsible for protecting buckets with a lock.
func (t MinIntThreshold) Apply(buckets *data.Buckets) (libhealth.Status, string) {
	exceeded := buckets.Compare(data.LessEq, data.NewInt(t.Threshold))
	return intThreshold(t).state(exceeded)
}

// MaxSumIntThreshold represents a maximum sum of all buckets
// before triggering an unhealthy state.Status.
type MaxSumIntThreshold intThreshold

// MinSumIntThreshold represents a minimum sum of all buckets
// before triggering an unhealthy state.Status.
type MinSumIntThreshold intThreshold

// Apply the limits of t over buckets.
//
// Caller is responsible for protecting buckets with a lock.
func (t MaxSumIntThreshold) Apply(buckets *data.Buckets) (libhealth.Status, string) {
	sum := buckets.Sum()
	crossed := data.NewInt(t.Threshold).Less(sum)
	return intThreshold(t).state(crossed)
}

// Apply the limits of t over buckets.
//
// Caller is responsible for protecting buckets with a lock.
func (t MinSumIntThreshold) Apply(buckets *data.Buckets) (libhealth.Status, string) {
	sum := buckets.Sum()
	crossed := sum.Less(data.NewInt(t.Threshold))
	return intThreshold(t).state(crossed)
}
