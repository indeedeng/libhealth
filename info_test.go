package libhealth

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- ok --

func Test_Info_live(t *testing.T) {
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
		},
		nil))
	deps.waitUntilInitialRun()

	info := NewInfo(deps)
	raw, status := info.generate(true, "test_live")
	result := string(raw)

	if status != 200 {
		t.Error("status not 200")
	}

	if !strings.Contains(result, `"condition":"OK"`) {
		t.Error("condition not OK")
	}
}

func Test_Info_background(t *testing.T) {
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
		},
		nil))
	deps.waitUntilInitialRun()

	info := NewInfo(deps)
	raw, status := info.generate(false, "test_background")
	result := string(raw)

	if status != 200 {
		t.Error("status not 200")
	}

	if !strings.Contains(result, `"condition":"OK"`) {
		t.Error("condition not OK")
	}
}

// -- not ok --

func Test_Info_live_bad(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				MAJOR,
				"doing nothing is healthy",
			)
		},
		nil))
	deps.waitUntilInitialRun()

	info := NewInfo(deps)
	raw, status := info.generate(true, "test_live")
	result := string(raw)

	if status != 200 {
		t.Error("status not 200")
	}

	if !strings.Contains(result, `"condition":"MAJOR"`) {
		t.Error("condition not MAJOR")
	}
}

func Test_Info_background_bad(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				MAJOR,
				"i am unhealthy",
			)
		},
		nil))
	deps.waitUntilInitialRun()

	info := NewInfo(deps)
	raw, status := info.generate(false, "test_background")
	result := string(raw)

	if status != 200 {
		t.Error("status not 200")
	}

	if !strings.Contains(result, `"condition":"MAJOR"`) {
		t.Error("condition not MAJOR")
	}
}

// -- outage --

func Test_Info_live_outage_avoided(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				OUTAGE,
				"i am dead",
			)
		},
		nil))
	deps.waitUntilInitialRun()

	info := NewInfo(deps)
	raw, status := info.generate(true, "test_live")
	result := string(raw)

	assert.Equal(t, http.StatusOK, status, "STRONG dep at OUTAGE for /info/ should not fail")
	assert.Contains(t, result, `"condition":"MAJOR"`)
}

func Test_Info_live_outage(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		REQUIRED,
		func(ctx context.Context) Health {
			return NewHealth(
				OUTAGE,
				"i am dead",
			)
		},
		nil))
	deps.waitUntilInitialRun()

	info := NewInfo(deps)
	raw, status := info.generate(true, "test_live")
	result := string(raw)

	assert.Equal(t, http.StatusInternalServerError, status, "REQUIRED dep at OUTAGE for /info/ should fail")
	assert.Contains(t, result, `"condition":"OUTAGE"`)
}

func Test_Info_background_outage_avoided(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		STRONG,
		func(ctx context.Context) Health {
			return NewHealth(
				OUTAGE,
				"i am dead",
			)
		},
		nil))
	deps.waitUntilInitialRun()

	info := NewInfo(deps)
	raw, status := info.generate(false, "test_background")
	result := string(raw)

	assert.Equal(t, http.StatusOK, status, "STRONG dep at OUTAGE for /info/ should not fail")
	assert.Contains(t, result, `"condition":"MAJOR"`)
}

func Test_Info_background_outage(t *testing.T) {
	deps := NewBasicDependencySet(NewMonitor(
		"example-daemon-dependency-check",
		"does not really do anything",
		"http://example.com/wiki/ExampleDaemon",
		REQUIRED,
		func(ctx context.Context) Health {
			return NewHealth(
				OUTAGE,
				"i am dead",
			)
		},
		nil))
	deps.waitUntilInitialRun()

	info := NewInfo(deps)
	raw, status := info.generate(false, "test_background")
	result := string(raw)

	assert.Equal(t, http.StatusInternalServerError, status, "REQUIRED dep at OUTAGE for /info/ should fail")
	assert.Contains(t, result, `"condition":"OUTAGE"`)
}
