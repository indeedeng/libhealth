package libhealth

import (
	"context"
	"sync"
	"time"
)

var (
	epoch          = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	defaultTimeout = 60 * time.Second
	defaultPeriod  = 30 * time.Second
)

// Monitor is the standard implementation of HealthChecker.
type Monitor struct {
	name        string
	timeout     time.Duration
	period      time.Duration
	description string
	docURL      string
	urgency     Urgency
	checker     HealthChecker
	statusChan  chan HealthStatus

	previous Health
	lastOk   time.Time
	failed   int
	lock     sync.RWMutex // locks above data
}

var _ HealthMonitor = (*Monitor)(nil)

// NewMonitorWithOptions constructs a new monitor from the provided components and configures optional ones (such as
// a status channel, timeout, and interval) based on the provided in options.
func NewMonitorWithOptions(
	name,
	description,
	docURL string,
	urgency Urgency,
	check HealthChecker,
	options ...MonitorOption,
) *Monitor {
	monitor := &Monitor{
		name:        name,
		timeout:     defaultTimeout,
		period:      defaultPeriod,
		description: description,
		docURL:      docURL,
		urgency:     urgency,
		checker:     check,
		statusChan:  nil,

		previous: NewHealth(OK, "starting up"),
		lastOk:   epoch,
		failed:   0,
	}

	for _, option := range options {
		option(monitor)
	}

	return monitor
}

// NewMonitor will create a new monitor and set some defaults.
// If desired, pass in a channel and it will be published to when the health state of the monitor changes.
func NewMonitor(
	name,
	description,
	docURL string,
	urgency Urgency,
	check HealthChecker,
	statusChan chan HealthStatus,
) *Monitor {
	return NewMonitorWithOptions(
		name, description, docURL, urgency, check,
		WithStatusChan(statusChan),
	)
}

type HealthStatus struct {
	Monitor HealthMonitor

	Prev Status
	Next Health
}

// Check will execute the HealthChecker associated with the monitor.
func (m *Monitor) Check(ctx context.Context) Health {
	prev, next := m.checkOnce(ctx)

	m.record(next, prev.Status)
	return next
}

func (m *Monitor) checkOnce(ctx context.Context) (prev, next Health) {
	startTime := time.Now()
	next = m.checker(ctx)
	endTime := time.Now()

	next.Urgency = m.urgency
	next.Time = startTime
	next.Duration = endTime.Sub(startTime)

	m.lock.Lock()
	{
		if next.SameAs(OK) {
			m.lastOk = endTime
			m.failed = 0
		} else {
			m.failed++
		}

		prev = m.previous
		m.previous = next
	}
	m.lock.Unlock()

	return prev, next
}

func (m *Monitor) record(next Health, prev Status) {
	if m.statusChan == nil {
		return
	}
	status := HealthStatus{
		Monitor: m,
		Prev:    prev,
		Next:    next,
	}
	select {
	case m.statusChan <- status:
		return
	default:
		return
	}
}

func (m *Monitor) Name() string {
	return m.name
}

func (m *Monitor) Timeout() time.Duration {
	return m.timeout
}

func (m *Monitor) Period() time.Duration {
	return m.period
}

func (m *Monitor) Description() string {
	return m.description
}

func (m *Monitor) Documentation() string {
	return m.docURL
}

func (m *Monitor) Urgency() Urgency {
	return m.urgency
}

func (m *Monitor) LastOk() time.Time {
	m.lock.RLock()
	defer m.lock.RUnlock() //nolint:gocritic

	return m.lastOk
}

func (m *Monitor) Failed() int {
	m.lock.RLock()
	defer m.lock.RUnlock() //nolint:gocritic

	return m.failed
}
