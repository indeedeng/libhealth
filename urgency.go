package libhealth

import "strings"

// Urgency is a level of requirement for a service to be operational. A REQUIRED
// service would cause major service disruption if it is not healthy. Likewise a
// WEAK service can fail without major issues.
type Urgency int

const (
	REQUIRED Urgency = iota
	STRONG
	WEAK
	NONE
	UNKNOWN
)

// ParseUrgency converts the given string into an Urgency.
// If the string is malformed, UNKNOWN is returned.
func ParseUrgency(urgency string) Urgency {
	switch strings.ToUpper(urgency) {
	case "REQUIRED":
		return REQUIRED
	case "STRONG":
		return STRONG
	case "WEAK":
		return WEAK
	case "NONE":
		return NONE
	default:
		return UNKNOWN
	}
}

// String provides an obvious representation of the Urgency level
func (u Urgency) String() string {
	switch u {
	case REQUIRED:
		return "REQUIRED"
	case STRONG:
		return "STRONG"
	case WEAK:
		return "WEAK"
	case NONE:
		return "NONE"
	case UNKNOWN:
		return "UNKNOWN"
	}
	return "UNKNOWN"
}

// Detail provides a detailed representation of the Urgency level
func (u Urgency) Detail() string {
	switch u {
	case REQUIRED:
		return "Required: Failure of this dependency would result in complete system outage"
	case STRONG:
		return "Strong: Failure of this dependency would result in major functional degradation"
	case WEAK:
		return "Weak: Failure of this dependency would result in minor functionality loss"
	case NONE:
		return "None: Failure of this dependency would result in no loss of functionality"
	case UNKNOWN:
		return "Unknown"
	}
	return "Invalid"
}

// DowngradeWith returns the downgraded Outage state according to HCv3 math.
func (u Urgency) DowngradeWith(systemState, newState Status) Status {
	switch u {
	case REQUIRED:
		return WorstState(systemState, newState)
	case STRONG:
		bounded := BestState(newState, MAJOR)
		return WorstState(systemState, bounded)
	case WEAK:
		bounded := BestState(newState, MINOR)
		return WorstState(systemState, bounded)
	case NONE:
		return systemState
	default: // UNKNOWN
		return WorstState(systemState, newState)
	}
}
