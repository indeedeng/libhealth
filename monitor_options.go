package libhealth

import "time"

type MonitorOption func(monitor *Monitor)

// WithStatusChan configures the optional subscription channel to notify consumers of health changes.
func WithStatusChan(statusChan chan HealthStatus) MonitorOption {
	return func(monitor *Monitor) {
		monitor.statusChan = statusChan
	}
}

// WithTimeout configures an optional timeout for the monitor. If not provided, the default is used.
func WithTimeout(timeout time.Duration) MonitorOption {
	return func(monitor *Monitor) {
		monitor.timeout = timeout
	}
}

// WithPeriod configures an optional interval for the montor to be run on. If not provided, the default is used.
func WithPeriod(period time.Duration) MonitorOption {
	return func(monitor *Monitor) {
		monitor.period = period
	}
}
