package libhealth

import (
	"fmt"
	"time"
)

// Message is a string with some useful information regarding a Health.
type Message string

// Health is a representation of the health of a service at a moment in time.
// It is composed of a Status, an Urgency, a Time, and a Message. Once created
// it should not be modified.
type Health struct {
	Status
	Urgency
	time.Time
	Message
	time.Duration
}

// NewHealth creates a Health for a fixed moment in time.
func NewHealth(
	state Status,
	message string,
) Health {
	return Health{
		Status:   state,
		Urgency:  UNKNOWN, // set by the owning Monitor
		Time:     time.Now(),
		Message:  Message(message),
		Duration: 0,
	}
}

// String returns a human readable summary
func (h Health) String() string {
	return fmt.Sprintf(
		"%s %s at %s, %s",
		h.Status,
		h.Urgency,
		h.Time,
		h.Message,
	)
}
