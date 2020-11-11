package count

import (
	"oss.indeed.com/go/libhealth"
	"oss.indeed.com/go/libhealth/count/internal/data"
)

// Floats will create a FloatCounter.
func Floats(varname string, size BucketPeriod, length int) (FloatCounter, error) {
	c, err := newContainer(varname, size, length, data.NewFloat(0))
	return (*floats)(c), err
}

type floats container

func (f *floats) Increment(delta float64) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.buckets.Increment(data.NewFloat(delta))
}

func (f *floats) Health() libhealth.Health {
	return (*container)(f).health()
}

func (f *floats) Set(threshold Threshold) Countable {
	return (*floats)((*container)(f).set(threshold))
}

type floatThreshold struct {
	// For Min/Max, Threshold represents the limit of the
	// value of the current bucket.
	// For SumMin/SumMax, Threshold represents the limit of
	// the values of all buckets summed together.
	Threshold   float64
	Description string
	Severity    libhealth.Status
}

// MaxFloatThreshold represents a maximum value that will trigger
// an unhealthy health.Status if at least n Infringements are
// exceeded.
type MaxFloatThreshold floatThreshold

// MinFloatThreshold represents a maximum value that will trigger
// an unhealthy health.Status if at least n Infringements are
// exceeded.
type MinFloatThreshold floatThreshold

func (f floatThreshold) state(violated bool) (libhealth.Status, string) {
	if violated {
		return f.Severity, f.Description
	}
	return libhealth.OK, OkMessage
}

// Apply the limits of t over buckets.
//
// Caller is responsible for protecting buckets with a lock.
func (t MaxFloatThreshold) Apply(buckets *data.Buckets) (libhealth.Status, string) {
	exceeded := buckets.Compare(data.GreaterEq, data.NewFloat(t.Threshold))
	return floatThreshold(t).state(exceeded)
}

// Apply the limits of t over buckets.
//
// Caller is responsible for protecting buckets with a lock.
func (t MinFloatThreshold) Apply(buckets *data.Buckets) (libhealth.Status, string) {
	exceeded := buckets.Compare(data.LessEq, data.NewFloat(t.Threshold))
	return floatThreshold(t).state(exceeded)
}

// MaxSumFloatThreshold represents a maximum sum of all buckets
// before triggering an unhealthy state.Status.
type MaxSumFloatThreshold floatThreshold

// MinSumFloatThreshold represents a minimum sum of all buckets
// before triggering an unhealthy state.Status.
type MinSumFloatThreshold floatThreshold

// Apply the limits of t over buckets.
//
// Caller is responsible for protecting buckets with a lock.
func (t MaxSumFloatThreshold) Apply(buckets *data.Buckets) (libhealth.Status, string) {
	sum := buckets.Sum()
	crossed := data.NewFloat(t.Threshold).Less(sum)
	return floatThreshold(t).state(crossed)
}

// Apply the limits of t over buckets.
//
// Caller is responsible for protecting buckets with a lock.
func (t MinSumFloatThreshold) Apply(buckets *data.Buckets) (libhealth.Status, string) {
	sum := buckets.Sum()
	crossed := sum.Less(data.NewFloat(t.Threshold))
	return floatThreshold(t).state(crossed)
}
