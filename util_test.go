package libhealth

import (
	"testing"
)

func Test_hostname(t *testing.T) {
	if hostname() == "unknown" {
		t.Error("unable to get os hostname")
	}
}
