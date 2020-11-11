// Package gauge provides mechanisms for gauging values.
package gauge

import (
	"container/list"
	"fmt"
	"strings"
	"sync"

	"oss.indeed.com/go/libhealth"
)

var (
	_ Gaugeable = (*ints)(nil)
	_ Threshold = (*MaxIntThreshold)(nil)
	_ Threshold = (*MinIntThreshold)(nil)

	_ Gaugeable = (*floats)(nil)
	_ Threshold = (*MaxFloatThreshold)(nil)
	_ Threshold = (*MinFloatThreshold)(nil)
)

const (
	// OkMessage is the default health.Message returned
	// for things that are not currently on fire.
	OkMessage = "ok"
)

// Gaugeable represents something which can apply thresholds to
// values collected over time, producing a resultant health.Health
// which represents the current state of the thing being gauged.
type Gaugeable interface {
	Health() libhealth.Health
	Set(threshold Threshold) Gaugeable
}

// An IntGauger represents Gaugable int values.
type IntGauger interface {
	Gaugeable
	Gauge(int)
}

// A FloatGauger represents Gaugable float values.
type FloatGauger interface {
	Gaugeable
	Gauge(float64)
}

// A Threshold represents some predicate which can be applied
// to an ordered list of values and return an either degraded
// or OK health.Status, along with a description.
type Threshold interface {
	Apply(*list.List) (libhealth.Status, string)
}

type container struct {
	lock       sync.RWMutex
	varexp     string
	length     int
	thresholds []Threshold
	values     *list.List
}

func newContainer(varname string, length int) (*container, error) {
	if err := checkLength(length); err != nil {
		return nil, err
	}

	c := &container{
		varexp:     varname,
		length:     length,
		thresholds: make([]Threshold, 0),
		values:     list.New(),
	}

	return c, nil
}

func (c *container) set(threshold Threshold) *container {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.thresholds = append(c.thresholds, threshold)
	return c
}

func (c *container) health() libhealth.Health {
	c.lock.RLock()
	defer c.lock.RUnlock()

	// Keep track of only the messages associated with the worst
	// state, and display those. If the worst state is OK, we
	// will only display "ok" so as to not confuse the operator.
	messages := []string{}
	worst := libhealth.OK

	for _, threshold := range c.thresholds {
		state, description := threshold.Apply(c.values)

		switch {
		case state.WorseThan(worst):
			messages = []string{description}
			worst = state // downgrade to new worst
		case state.SameAs(worst):
			messages = append(messages, description)
		}
	}

	message := strings.Join(messages, ", ")
	if worst == libhealth.OK {
		// just a single "ok" if things are fine
		message = OkMessage
	}
	return libhealth.NewHealth(worst, message)
}

func checkLength(length int) error {
	if length <= 0 {
		return fmt.Errorf("a gauge must keep track of at least one value, len: %d", length)
	}
	return nil
}
