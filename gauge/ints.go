package gauge

import (
	"container/list"

	"oss.indeed.com/go/libhealth"
)

// Ints will create an IntGauger.
//
// The value of varname is exported by varexp, and so it should
// follow the convention of being all lowercase with dashes, and
// preferably prefixed with the name of the executable.
//
// The specified length must be greater than zero.
//
// Operations made available by the underlying implementation
// are threadsafe.
func Ints(varname string, length int) (IntGauger, error) {
	c, err := newContainer(varname, length)
	return (*ints)(c), err
}

type ints container

// Gauge value. Old historical values will be removed from i
// as necessary.
func (i *ints) Gauge(value int) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.values.PushFront(value)
	for i.values.Len() > i.length {
		i.values.Remove(i.values.Back())
	}
}

// Health computes the current health of i by applying the defined
// thresholds over the current history of gauged values.
func (i *ints) Health() libhealth.Health {
	return (*container)(i).health()
}

// Set a new threshold on i.
func (i *ints) Set(threshold Threshold) Gaugeable {
	// wtf we started programming in C?
	return (*ints)((*container)(i).set(threshold))
}

type intThreshold struct {
	Threshold   int
	LastN       int
	AnyN        int
	Description string
	Severity    libhealth.Status
}

// MaxIntThreshold represents a maximum threshold over a sequence of int.
type MaxIntThreshold intThreshold

// Apply the limits of t over values.
//
// Caller is responsible for protecting values with a lock.
func (t MaxIntThreshold) Apply(values *list.List) (libhealth.Status, string) {
	i, numViolations := 1, 0
	// Walk forwards starting at the front of the list, because that's
	// where the most recent elements are. The most recent elements are
	// what are being checked for the LastN threshold.
	for elem := values.Front(); elem != nil; elem = elem.Next() {
		value := elem.Value.(int)
		if value >= t.Threshold {
			numViolations++
		}

		// Check to see if the LastN libhealth were all over Threshold.
		if (i == t.LastN) && (numViolations == t.LastN) {
			return t.Severity, t.Description
		}
		i++
	}

	// Check to see if at least AnyN libhealth were over Threshold.
	if (t.AnyN > 0) && (numViolations >= t.AnyN) {
		return t.Severity, t.Description
	}

	// Otherwise everything is fine.
	return libhealth.OK, OkMessage
}

// MinIntThreshold represents a minimum threshold over a sequence of integers.
type MinIntThreshold intThreshold

// Apply the limits of t over values.
//
// Caller is responsible for protecting values with a lock.
func (t MinIntThreshold) Apply(values *list.List) (libhealth.Status, string) {
	i, numViolations := 1, 0
	// Walk forwards starting at the front of the list, because that's
	// where the most recent elements are. The most recent elements are
	// what are being checked for the LastN threshold.
	for elem := values.Front(); elem != nil; elem = elem.Next() {
		value := elem.Value.(int)
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
