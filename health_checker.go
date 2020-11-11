package libhealth

import (
	"context"
	"time"
)

// HealthChecker is a func that you provide which checks the status of something's health.
type HealthChecker func(ctx context.Context) Health

// HealthMonitor is an interface that should be implemented for every custom type of
// health check you need. These can then be added to a DependencySet
// and used in Live and Background version of Info and Private healthchecks.
type HealthMonitor interface {
	// Name provides a unique name for the healthcheck.
	Name() string
	// Check is the function representing the work of a health check.
	Check(ctx context.Context) Health
	// Timeout is how long Check gets to run before defaulting to OUTAGE.
	Timeout() time.Duration
	// Period is how often Check should run.
	Period() time.Duration
	// Description is an informative string about what the healthcheck does.
	Description() string
	// Documentation is a url link that documents additional information about the healthcheck.
	Documentation() string
	// Urgency is how important this service.
	Urgency() Urgency
	// LastOk is the time when this healthcheck was last OK.
	LastOk() time.Time
	// Failed is the number of consecutive healthcheck failures. Returns zero when the healthcheck is OK.
	Failed() int
}
