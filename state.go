package libhealth

import "strings"

// Status is a representation of the health of the component.
type Status int

const (
	OUTAGE Status = iota
	MAJOR
	MINOR
	OK
)

// ParseStatus parses the given string into a Status.
// If the string is malformed, OUTAGE is returned.
func ParseStatus(state string) Status {
	switch strings.ToUpper(state) {
	case "OUTAGE":
		return OUTAGE
	case "MAJOR":
		return MAJOR
	case "MINOR":
		return MINOR
	case "OK":
		return OK
	default:
		return OUTAGE
	}
}

// String provides a regular string representation of a Status.
func (s Status) String() string {
	switch s {
	case OUTAGE:
		return "OUTAGE"
	case MAJOR:
		return "MAJOR"
	case MINOR:
		return "MINOR"
	case OK:
		return "OK"
	}
	return "INVALID"
}

// WorseThan compares a Status to another Status.
func (s Status) WorseThan(level Status) bool {
	return s < level
}

// SameOrWorseThan compares a Status to another Status.
func (s Status) SameOrWorseThan(level Status) bool {
	return s <= level
}

// SameOrBetterThan compares a Status to another Status.
func (s Status) SameOrBetterThan(level Status) bool {
	return s >= level
}

// BetterThan compares a Status to another Status.
func (s Status) BetterThan(level Status) bool {
	return s > level
}

// SameAs compares a Status to another Status.
func (s Status) SameAs(level Status) bool {
	return s == level
}

// BestState returns the more cheerful of two states.
func BestState(left, right Status) Status {
	if left.SameOrBetterThan(right) {
		return left
	}
	return right
}

// WorstState returns the less positive of two states.
func WorstState(left, right Status) Status {
	if left.SameOrWorseThan(right) {
		return left
	}
	return right
}
