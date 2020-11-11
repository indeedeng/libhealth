package gauge

import (
	"container/list"

	"oss.indeed.com/go/libhealth"
)

// Floats will create a FloatGauger.
//
// The value of varname is exported by varexp, and so it should
// follow the convention of being all lowercase with dashes,
// preferably prefixed with the name of the executable.
//
// The specified length must be greater than zero.
//
// Operations made available by the underlying implementation
// are threadsafe.
func Floats(varname string, length int) (FloatGauger, error) {
	c, err := newContainer(varname, length)
	return (*floats)(c), err
}

type floats container

func (f *floats) Gauge(value float64) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.values.PushFront(value)
	for f.values.Len() > f.length {
		f.values.Remove(f.values.Back())
	}
}

func (f *floats) Health() libhealth.Health {
	return (*container)(f).health()
}

func (f *floats) Set(threshold Threshold) Gaugeable {
	return (*floats)((*container)(f).set(threshold))
}

type floatThreshold struct {
	Threshold   float64
	LastN       int
	AnyN        int
	Description string
	Severity    libhealth.Status
}

// MaxFloatThreshold represents a maximum threshold over a sequence of float64.
type MaxFloatThreshold floatThreshold

// Apply the limits of t over values.
func (t MaxFloatThreshold) Apply(values *list.List) (libhealth.Status, string) {
	i, numViolations := 1, 0
	// Walk forwards starting at the front of the list, because that's
	// where the most recent elements are. The most recent elements are
	// what are being checked for the LastN threshold.
	for elem := values.Front(); elem != nil; elem = elem.Next() {
		value := elem.Value.(float64)
		if value >= t.Threshold {
			numViolations++
		}

		// Check to see if the LastN elements were all over Threshold.
		if (i == t.LastN) && (numViolations == t.LastN) {
			return t.Severity, t.Description
		}
		i++
	}

	// Check to see if at least AnyN elements were over Threshold.
	if (t.AnyN > 0) && (numViolations >= t.AnyN) {
		return t.Severity, t.Description
	}

	// Otherwise everything is fine.
	return libhealth.OK, OkMessage
}

// MinFloatThreshold represents a minimum threshold over a sequence of float64.
type MinFloatThreshold floatThreshold

// Apply the limits of t over values.
//
// Caller is responsible for protecting values with a lock.
func (t MinFloatThreshold) Apply(values *list.List) (libhealth.Status, string) {
	i, numViolations := 1, 0
	// Walk forwards starting at the front of the list, because that's
	// where the most recent elements are. The most recent elements are
	// what are being checked for the LastN threshold.
	for elem := values.Front(); elem != nil; elem = elem.Next() {
		value := elem.Value.(float64)
		if value <= t.Threshold {
			numViolations++
		}

		// Check to see if the LastN elements were all under Threshold.
		if (i == t.LastN) && (numViolations == t.LastN) {
			return t.Severity, t.Description
		}
		i++
	}

	// Check to see if at least AnyN elements were over t.Threshold.
	if (t.AnyN > 0) && (numViolations >= t.AnyN) {
		return t.Severity, t.Description
	}

	// Otherwise everything is fine.
	return libhealth.OK, OkMessage
}
