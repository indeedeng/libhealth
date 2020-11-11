package libhealth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

// Paths to info healthchecks.
const (
	InfoHealthCheck     = `/info/healthcheck`
	InfoHealthCheckLive = `/info/healthcheck/live`
	infoBad             = `{"condition":"info healthcheck error"}`
)

// Info is an http.Handler for
//      /info/healthcheck
//      /info/healthcheck/live.
type Info struct {
	deps DependencySet
}

// NewInfo will create a new Info handler for a given DependencySet set.
func NewInfo(d DependencySet) *Info {
	return &Info{d}
}

// InfoResult represents the body of an info healthcheck.
type InfoResult struct {
	Condition string `json:"condition"`
	Hostname  string `json:"hostname"`
	Duration  int64  `json:"duration"`
}

func (i *Info) generate(live bool, hostname string) ([]byte, int) {
	var s Summary
	if live {
		s = i.deps.Live()
	} else {
		s = i.deps.Background()
	}

	r := InfoResult{
		Condition: s.Overall().String(),
		Hostname:  hostname,
		Duration:  int64(s.Duration()) / 1000000,
	}

	bytes, err := json.Marshal(r)
	if err != nil {
		return []byte(infoBad), http.StatusInternalServerError
	}
	return bytes, ComputeStatusCode(true, s)
}

// ServeHTTP is intended to be used by a net/http.ServeMux for serving formatted json.
func (i *Info) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	live := strings.HasSuffix(r.URL.String(), "/live")
	healthcheckJSON, code := i.generate(live, hostname())
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, healthcheckJSON, "", "  "); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	if _, err := w.Write(prettyJSON.Bytes()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
