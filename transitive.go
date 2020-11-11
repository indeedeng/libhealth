package libhealth

import (
	"context"
	"net/http"
	"time"
)

var (
	sharedClient = http.Client{
		Timeout: 10 * time.Second,
	}
)

// TransitiveMonitor creates a Monitor that is a dependency on another service that responds to a healthcheck
func TransitiveMonitor(
	url,
	name,
	description,
	wikipage string,
	urgency Urgency,
	statusChan chan HealthStatus,
) *Monitor {
	errorHealth := func(err error, start time.Time) Health {
		msg := "error checking transitive monitor: " + err.Error()
		return Health{
			Status:   OUTAGE,
			Urgency:  urgency,
			Time:     start,
			Message:  Message(msg),
			Duration: time.Since(start),
		}
	}

	return NewMonitor(
		name,
		description,
		wikipage,
		urgency,

		func(ctx context.Context) Health {
			start := time.Now()

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return errorHealth(err, start)
			}
			resp, err := sharedClient.Do(req)
			if err != nil {
				return errorHealth(err, start)
			}
			defer resp.Body.Close()

			state := OK
			if resp.StatusCode != http.StatusOK {
				state = OUTAGE
			}

			return Health{
				Status:   state,
				Urgency:  urgency,
				Time:     start,
				Message:  Message(resp.Status),
				Duration: time.Since(start),
			}
		},
		statusChan,
	)
}
