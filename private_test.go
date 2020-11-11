package libhealth

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/libtime"
)

func Test_copyComponents(t *testing.T) {
	exTime := time.Date(2017, 11, 28, 11, 38, 1, 0, time.UTC)
	r1Time := time.Date(2017, 11, 28, 11, 39, 0, 0, time.UTC)
	r1LastGood := time.Date(2017, 11, 28, 11, 37, 0, 0, time.UTC)
	r2Time := time.Date(2017, 11, 28, 11, 37, 0, 0, time.UTC)
	r2LastGood := time.Date(2017, 11, 28, 11, 37, 0, 0, time.UTC)
	summary := Summary{
		executed: exTime,
		results: []Result{
			{
				Health: Health{
					Status:   OUTAGE,
					Urgency:  WEAK,
					Time:     r1Time,
					Message:  "the thing is broken",
					Duration: 2 * time.Second,
				},
				docurl:   "https://example.com",
				desc:     "description1",
				lastGood: r1LastGood,
				period:   12 * time.Second,
				name:     "check1",
			},
			{
				Health: Health{
					Status:   MINOR,
					Urgency:  REQUIRED,
					Time:     r2Time,
					Message:  "the thing is sort of broken",
					Duration: 10 * time.Second,
				},
				docurl:   "https://example.com",
				desc:     "description1",
				lastGood: r2LastGood,
				period:   12 * time.Second,
				name:     "check1",
			},
		},
	}

	components := copyComponents(summary)
	require.Equal(t, 2, len(components))
	require.Equal(t, libtime.ToMilliseconds(r1Time), components[0].Timestamp)
	require.Equal(t, libtime.ToMilliseconds(r1LastGood), components[0].LastGood)
	require.Equal(t, libtime.ToMilliseconds(r2Time), components[1].Timestamp)
	require.Equal(t, libtime.ToMilliseconds(r2LastGood), components[1].LastGood)
	require.Equal(t, "the thing is broken", components[0].Message)
	require.Equal(t, "the thing is sort of broken", components[1].Message)
}

func Test_times(t *testing.T) {
	// the "system" time is in Japan (a fixed timezone) for testing
	japan, err := time.LoadLocation("Japan")
	require.NoError(t, err)

	now := time.Date(2017, 11, 14, 11, 51, 33, 0, time.UTC).In(japan)
	started := time.Date(2017, 10, 31, 2, 3, 4, 0, time.UTC).In(japan)
	appTimes := times(now, started)

	require.Equal(t, "2017-11-14T20:51:33.000+0900", appTimes.AppStartDateSystem)
	require.Equal(t, "2017-11-14T11:51:33.000+0000", appTimes.AppStartDateUTC)
	require.Equal(t, "1509415384000", appTimes.AppStartUnixTimestamp)
	require.Equal(t, "345h48m29s", appTimes.AppUpTimeReadable)
	require.Equal(t, "1244909", appTimes.AppUpTimeSeconds)
}

// -- test env reporting --

func Test_Private_start_times(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				OK,
				"doing nothing is healthy",
			)
		}, nil))
	deps.waitUntilInitialRun()

	priv := NewPrivate("test_live", deps)
	raw, status := priv.generate(true, "test_live")
	result := string(raw)

	require.Equal(t, 200, status)
	require.Contains(t, result, "appStartDateSystem")
	require.Contains(t, result, "appStartDateUTC")
	require.Contains(t, result, "appStartUnixTimestamp")
	require.Contains(t, result, "appUpTimeReadable")
	require.Contains(t, result, "appUpTimeSeconds")
	require.Contains(t, result, `"results":{"OK":[{"timestamp":`)
}

// -- ok --

func Test_Private_live(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				OK,
				"doing nothing is healthy",
			)
		}, nil))
	deps.waitUntilInitialRun()

	priv := NewPrivate("test_live", deps)
	raw, status := priv.generate(true, "test_live")
	result := string(raw)

	require.Equal(t, 200, status)
	require.Contains(t, result, `"condition":"OK"`)
}

func Test_Private_background(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				OK,
				"doing nothing is healthy",
			)
		}, nil))
	deps.waitUntilInitialRun()

	priv := NewPrivate("test_background", deps)
	raw, status := priv.generate(false, "test_background")
	result := string(raw)

	require.Equal(t, 200, status)
	require.Contains(t, result, `"condition":"OK"`)
}

// -- bad --

func Test_Private_live_STRONG_MAJOR(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				MAJOR,
				"doing nothing is major",
			)
		}, nil))
	deps.waitUntilInitialRun()

	priv := NewPrivate("test_live", deps)
	raw, status := priv.generate(true, "test_live")
	j, err := decodeRaw(raw)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, status)

	// the overall check is capped by the urgency
	require.Equal(t, "MAJOR", j.Condition, "overall healthcheck is capped")
	// the individual results are also capped by urgency
	require.EqualValues(t, []string{"MAJOR"}, j.Statuses())
	require.Len(t, j.Results["MAJOR"], 1)
	require.Equal(t, "example-daemon-dependency-check", j.Results["MAJOR"][0].ID)
}

func Test_Private_live_STRONG_OUTAGE(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				OUTAGE,
				"doing nothing is outage",
			)
		}, nil))
	deps.waitUntilInitialRun()

	priv := NewPrivate("test_live", deps)
	raw, status := priv.generate(true, "test_live")
	j, err := decodeRaw(raw)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, status)

	// the overall check is capped by the urgency
	require.Equal(t, "MAJOR", j.Condition, "overall healthcheck is capped")
	// the individual results are also capped by urgency
	require.EqualValues(t, []string{"MAJOR"}, j.Statuses())
	require.Len(t, j.Results["MAJOR"], 1)
	require.Equal(t, "example-daemon-dependency-check", j.Results["MAJOR"][0].ID)
}

func Test_Private_background_STRONG_MAJOR(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				MAJOR,
				"doing nothing is major",
			)
		}, nil))
	deps.waitUntilInitialRun()

	priv := NewPrivate("test_background", deps)
	raw, status := priv.generate(false, "test_background")
	j, err := decodeRaw(raw)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, status)

	// the overall check is capped by the urgency
	require.Equal(t, "MAJOR", j.Condition, "overall healthcheck is capped")
	// the individual results are also capped by urgency
	require.EqualValues(t, []string{"MAJOR"}, j.Statuses())
	require.Len(t, j.Results["MAJOR"], 1)
	require.Equal(t, "example-daemon-dependency-check", j.Results["MAJOR"][0].ID)
}

func Test_Private_background_STRONG_OUTAGE(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				OUTAGE,
				"doing nothing is outage",
			)
		}, nil))
	deps.waitUntilInitialRun()

	priv := NewPrivate("test_background", deps)
	raw, status := priv.generate(false, "test_background")
	j, err := decodeRaw(raw)
	require.NoError(t, err)

	require.Equal(t, http.StatusInternalServerError, status)

	// the overall check is capped by the urgency
	require.Equal(t, "MAJOR", j.Condition, "overall healthcheck is capped")
	// the individual results are also capped by urgency
	require.EqualValues(t, []string{"MAJOR"}, j.Statuses())
	require.Len(t, j.Results["MAJOR"], 1)
	require.Equal(t, "example-daemon-dependency-check", j.Results["MAJOR"][0].ID)
}

func decodeRaw(raw []byte) (*privateResponse, error) {
	var resp privateResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type privateResponse struct {
	Condition string `json:"condition"`
	Results   map[string][]*privateDep
}

func (p *privateResponse) Statuses() []string {
	ret := make([]string, 0, len(p.Results))
	for k := range p.Results {
		ret = append(ret, k)
	}
	return ret
}

type privateDep struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Urgency string `json:"urgency"`
}
