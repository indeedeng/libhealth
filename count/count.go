// Package count enables tracking accumulating values.
package count

import (
	"fmt"
	"strings"
	"sync"

	"oss.indeed.com/go/libhealth"
	"oss.indeed.com/go/libhealth/count/internal/data"
)

var (
	_ Threshold = (*MaxIntThreshold)(nil)
	_ Threshold = (*MinIntThreshold)(nil)
	_ Threshold = (*MaxSumIntThreshold)(nil)
	_ Threshold = (*MinSumIntThreshold)(nil)

	_ Threshold = (*MaxFloatThreshold)(nil)
	_ Threshold = (*MinFloatThreshold)(nil)
	_ Threshold = (*MaxSumFloatThreshold)(nil)
	_ Threshold = (*MinSumFloatThreshold)(nil)
)

const (
	// OkMessage is used to override a health.Status failure
	// message when things are in the health.OK state.
	OkMessage = "ok"
)

// A Countable represents something which can apply thresholds to
// accumulating values over time, producing resultant health.Health
// which represents the current state of the thing being counted.
type Countable interface {
	Health() libhealth.Health
	Set(threshold Threshold) Countable
}

// An IntCounter represents Countable int values.
type IntCounter interface {
	Countable
	Increment(int)
}

// A FloatCounter represents Countable float values.
type FloatCounter interface {
	Countable
	Increment(float64)
}

// A Threshold represents some predicate which can be applied to
// an ordered list of values and return either degraded or an
// OK health.Status, along with a description.
type Threshold interface {
	Apply(bucket *data.Buckets) (libhealth.Status, string)
}

// A BucketPeriod represents how much time is alloted to each
// bucket of values over time.
type BucketPeriod data.Ticker

// Most common BucketPeriod values provided for convenience.
var (
	SizeOneSecond      = BucketPeriod(data.OneSecond)
	SizeOneMinute      = BucketPeriod(data.OneMinute)
	SizeFiveMinutes    = BucketPeriod(data.FiveMinutes)
	SizeFifteenMinutes = BucketPeriod(data.FifteenMinutes)
	SizeOneHour        = BucketPeriod(data.OneHour)
)

type container struct {
	lock       sync.RWMutex
	varexp     string
	length     int
	thresholds []Threshold
	buckets    *data.Buckets
}

func newContainer(varname string, size BucketPeriod, length int, zero data.Value) (*container, error) {
	if err := checkLength(length); err != nil {
		return nil, err
	}

	c := &container{
		varexp:     varname,
		length:     length,
		thresholds: make([]Threshold, 0),
		buckets:    data.NewBuckets(data.Ticker(size), length, zero),
	}

	return c, nil
}

func (c *container) health() libhealth.Health {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Keep track of the messages associated with the worst state
	// so that we can display only those. If the worst state is OK,
	// we return only "ok" so as to avoid confusion.
	messages := make([]string, 0, len(c.thresholds))
	worst := libhealth.OK
	for _, threshold := range c.thresholds {
		state, description := threshold.Apply(c.buckets)
		switch {
		case state.WorseThan(worst):
			worst = state // downgrade to new worst
			messages = []string{description}
		case state.SameAs(worst):
			messages = append(messages, description)
		}
	}

	message := strings.Join(messages, ", ")
	if worst == libhealth.OK {
		// hide error messages if state is OK
		message = "ok"
	}

	return libhealth.NewHealth(worst, message)
}

func (c *container) set(threshold Threshold) *container {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.thresholds = append(c.thresholds, threshold)
	return c
}

func checkLength(length int) error {
	if length <= 0 {
		return fmt.Errorf("a counter must keep track of at least one value, len: %d", length)
	}
	return nil
}
