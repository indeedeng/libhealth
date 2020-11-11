package libhealth

import (
	"os"
)

// return the machine hostname or unknown
func hostname() string {
	h, e := os.Hostname()
	if e != nil {
		return "unknown"
	}
	return h
}
