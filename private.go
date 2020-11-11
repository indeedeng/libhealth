package libhealth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"oss.indeed.com/go/libtime"
)

// The exposed endpoints for private healthchecks.
const (
	PrivateHealthCheck     = `/private/healthcheck`
	PrivateHealthCheckLive = `/private/healthcheck/live`
)

const (
	privateBad = `{"condition":"private healthcheck error"}`
	timeFormat = "2006-01-02T15:04:05.000-0700" // this REFERENCE time must be MST
)

var allowlistEnv = []string{"HOME", "LANG", "PATH", "PWD", "TMPDIR", "SHELL", "USER"}

// Private is an http.Handler for
//      /private/healthcheck
//      /private/healthcheck/live
type Private struct {
	dependencies DependencySet
	appName      string
	startTime    time.Time
}

// NewPrivate creates a new Private so that service appName can
// be register its DependencySet.
func NewPrivate(appName string, set DependencySet) *Private {
	return &Private{
		dependencies: set,
		appName:      appName,
		startTime:    time.Now(),
	}
}

// Components of a particular Status. We list them explicitly
// so that during json encoding they are ordered.
type Components struct {
	Outage []Component `json:"OUTAGE,omitempty"`
	Major  []Component `json:"MAJOR,omitempty"`
	Minor  []Component `json:"MINOR,omitempty"`
	Ok     []Component `json:"OK,omitempty"`
}

// A Component is the healthcheck status of one
// component in a /private/healthcheck result.
type Component struct {
	Timestamp   int64  `json:"timestamp"`
	DocURL      string `json:"documentationUrl"`
	Urgency     string `json:"urgency"`
	Description string `json:"description"`
	State       string `json:"status"`
	Message     string `json:"errorMessage"`
	Duration    int64  `json:"duration"`
	LastGood    int64  `json:"lastKnownGoodTimestamp"`
	Period      int64  `json:"period"`
	ID          string `json:"id"`
	Date        string `json:"date"`
}

// A PrivateResult is the struct (and JSON) definition of what
// the response to a private healthcheck endpoint should be
// composed of.
type PrivateResult struct {
	AppName                   string            `json:"appName"`
	Condition                 string            `json:"condition"`
	Duration                  int64             `json:"duration"`
	Hostname                  string            `json:"hostname"`
	Environment               map[string]string `json:"environment"`
	CWD                       string            `json:"cwd"`
	AppStartDateSystem        string            `json:"appStartDateSystem"`
	AppStartDateUTC           string            `json:"appStartDateUTC"`
	AppStartUnixTimestamp     string            `json:"appStartUnixTimestamp"`
	AppUpTimeReadable         string            `json:"appUpTimeReadable"`
	AppUpTimeSeconds          string            `json:"appUpTimeSeconds"`
	LeastRecentlyExecutedDate string            `json:"leastRecentlyExecutedDate"`
	LeastRecentlyExecutedTime int64             `json:"leastRecentlyExecutedTimestamp"`
	Results                   Components        `json:"results"`
}

func categorize(results []Component) Components {
	c := Components{}
	for _, result := range results {
		switch ParseStatus(result.State) {
		case OUTAGE:
			c.Outage = append(c.Outage, result)
		case MAJOR:
			c.Major = append(c.Major, result)
		case MINOR:
			c.Minor = append(c.Minor, result)
		case OK:
			c.Ok = append(c.Ok, result)
		}
	}
	return c
}

func copyComponents(s Summary) []Component {
	components := make([]Component, 0, len(s.results))
	for _, result := range s.results {
		components = append(components, Component{
			Timestamp:   libtime.ToMilliseconds(result.Time),
			DocURL:      result.docurl,
			Urgency:     result.Urgency.Detail(),
			Description: result.desc,
			State:       result.Status.String(),
			Message:     string(result.Message),
			Duration:    result.Duration.Nanoseconds() / 1000,
			LastGood:    libtime.ToMilliseconds(result.lastGood),
			Period:      result.period.Nanoseconds() / 1000000000,
			ID:          result.name,
			Date:        result.Time.Format(timeFormat),
		})
	}
	return components
}

func (p *Private) generate(live bool, hostname string) (hc []byte, code int) {
	var summary Summary
	if live {
		summary = p.dependencies.Live()
	} else {
		summary = p.dependencies.Background()
	}

	components := copyComponents(summary)
	categorized := categorize(components)

	appTimes := times(time.Now(), p.startTime)
	cwd, _ := os.Getwd()
	env := os.Environ()
	envMap := make(map[string]string, len(env))
	for _, envVar := range allowlistEnv {
		if val, found := os.LookupEnv(envVar); found {
			envMap[envVar] = val
		}
	}

	result := PrivateResult{
		AppName:                   p.appName,
		Condition:                 summary.Overall().String(),
		Hostname:                  hostname,
		Environment:               envMap,
		CWD:                       cwd,
		AppStartDateSystem:        appTimes.AppStartDateSystem,
		AppStartDateUTC:           appTimes.AppStartDateUTC,
		AppStartUnixTimestamp:     appTimes.AppStartUnixTimestamp,
		AppUpTimeReadable:         appTimes.AppUpTimeReadable,
		AppUpTimeSeconds:          appTimes.AppUpTimeSeconds,
		LeastRecentlyExecutedDate: summary.Executed().Format(timeFormat),
		LeastRecentlyExecutedTime: libtime.ToMilliseconds(summary.Executed()),
		Results:                   categorized,
	}

	hc, err := json.Marshal(result)
	if err != nil {
		return []byte(privateBad), http.StatusInternalServerError
	}
	return hc, ComputeStatusCode(false, summary)
}

type timeCollection struct {
	AppStartDateSystem    string
	AppStartDateUTC       string
	AppStartUnixTimestamp string
	AppUpTimeReadable     string
	AppUpTimeSeconds      string
}

func times(now, started time.Time) timeCollection {
	return timeCollection{
		AppStartDateSystem:    now.Format(timeFormat),
		AppStartDateUTC:       now.In(time.UTC).Format(timeFormat),
		AppStartUnixTimestamp: strconv.FormatInt(libtime.ToMilliseconds(started), 10),
		AppUpTimeReadable:     now.Sub(started).String(),
		AppUpTimeSeconds:      strconv.FormatInt(int64(now.Sub(started).Seconds()), 10),
	}
}

func (p *Private) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	live := strings.HasSuffix(r.URL.String(), "/live")
	healthcheck, code := p.generate(live, hostname())
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, healthcheck, "", "  "); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	if _, err := w.Write(prettyJSON.Bytes()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
