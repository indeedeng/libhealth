package libhealth

import (
	"context"
	"sync"
	"time"
)

// BasicDependencySet is the standard implementation of DependancySet that
// behaves the way you would expect. HealthMonitors are registered to a
// DependencySet instance, and then you can call Live() or Background() on
// that instance. Live() will cause all of the HealthMonitors to call their
// Check() methods, whereas Background() will retrieve the cached Health
// of that health check.
type BasicDependencySet struct {
	monitors map[string]HealthMonitor
	cached   map[string]Result
	lock     sync.RWMutex // locks the map structure, but not the values

	ctx context.Context

	initialRunWg sync.WaitGroup
}

// NewBasicDependencySet will create a new BasicDependencySet instance and
// register all of the provided HealthMonitor instances.
func NewBasicDependencySet(monitors ...HealthMonitor) *BasicDependencySet {
	return NewBasicDependencySetWithContext(context.Background(), monitors...)
}

// NewBasicDependencySet will create a new BasicDependencySet instance using
// a specific context and register all of the provided HealthMonitor instances.
func NewBasicDependencySetWithContext(ctx context.Context, monitors ...HealthMonitor) *BasicDependencySet {
	deps := &BasicDependencySet{
		monitors: make(map[string]HealthMonitor),
		cached:   make(map[string]Result),
		ctx:      ctx,
	}
	deps.Register(monitors...)
	return deps
}

// Register will register all of the provided HealthMonitor instances.
// HealthMonitors which have a 0 value Timeout will NOT be executed on
// Register(). They will only be executed on calls to Live(). This is important,
// because it means such a checker will be set to OUTAGE if only
// Background() is ever called.
func (d *BasicDependencySet) Register(monitors ...HealthMonitor) {
	for _, monitor := range monitors {
		notyetrun := fresh(monitor)

		// set status not-run-yet
		d.update(monitor, &notyetrun)

		d.initialRunWg.Add(1)
		firstTime := func(monitor HealthMonitor) {
			d.run(monitor, time.Now())
			d.initialRunWg.Done()
		}

		// then run the check every period
		if monitor.Period() > 0 {
			go func(monitor HealthMonitor) {
				// immediately run a check, then schedule more
				firstTime(monitor)

				// each healthcheck ticks and updates its associated health
				// if the check times out, the health is set to outage
				ticker := time.NewTicker(monitor.Period())
				defer ticker.Stop()

				for now := range ticker.C {
					d.run(monitor, now)
				}
			}(monitor)
		} else {
			// then immediately kick off a check
			go firstTime(monitor)
		}
	}
}

func (d *BasicDependencySet) waitUntilInitialRun() {
	d.initialRunWg.Wait()
}

func (d *BasicDependencySet) run(monitor HealthMonitor, now time.Time) Result {
	result := performCheck(d.ctx, monitor, now)
	d.update(monitor, &result)
	return result
}

func performCheck(ctx context.Context, monitor HealthMonitor, startTime time.Time) Result {
	ctx, cancelFunc := context.WithDeadline(ctx, startTime.Add(monitor.Timeout()))
	defer cancelFunc()

	select {
	case health := <-asyncCheck(ctx, monitor):
		return wrap(monitor, health)
	case <-ctx.Done():
		return timeout(monitor, startTime)
	}
}

func asyncCheck(ctx context.Context, monitor HealthMonitor) <-chan Health {
	healthc := make(chan Health, 1)
	go func() {
		healthc <- monitor.Check(ctx)
	}()

	return healthc
}

func fresh(m HealthMonitor) Result {
	h := NewHealth(OUTAGE, "healthcheck has not run yet")
	h.Urgency = m.Urgency()
	h.Time = time.Now()
	return wrap(m, h)
}

func timeout(m HealthMonitor, t time.Time) Result {
	h := NewHealth(OUTAGE, "healthcheck timed out")
	h.Urgency = m.Urgency()
	h.Time = t
	return wrap(m, h)
}

func wrap(m HealthMonitor, h Health) Result {
	// We don't care about the real state, just the downgraded one.
	h.Status = h.Urgency.DowngradeWith(OK, h.Status)
	return Result{
		h,
		m.Documentation(),
		m.Description(),
		m.LastOk(),
		m.Period(),
		m.Name(),
	}
}

// Background will retrieve the cached Health for each of the registered
// HealthChecker instances.
func (d *BasicDependencySet) Background() Summary {
	return NewSummary(time.Now(), d.snapshotResults())
}

// Live will force all of the HealthChecker instances to execute their
// Check methods, and will update all cached Health as well.
func (d *BasicDependencySet) Live() Summary {
	monitors := d.snapshotMonitors()

	checkResults := make(chan Result)
	start := time.Now()
	for _, monitor := range monitors {
		go func(m HealthMonitor) {
			checkResults <- d.run(m, start)
		}(monitor)
	}

	results := make([]Result, 0, len(monitors))
	for range monitors {
		results = append(results, <-checkResults)
	}
	return NewSummary(time.Now(), results)
}

func (d *BasicDependencySet) update(monitor HealthMonitor, result *Result) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.monitors[monitor.Name()] = monitor
	d.cached[monitor.Name()] = *result
}

func (d *BasicDependencySet) snapshotMonitors() []HealthMonitor {
	d.lock.RLock()
	defer d.lock.RUnlock()

	monitors := make([]HealthMonitor, 0, len(d.monitors))
	for _, monitor := range d.monitors {
		monitors = append(monitors, monitor)
	}

	return monitors
}

func (d *BasicDependencySet) snapshotResults() []Result {
	d.lock.RLock()
	defer d.lock.RUnlock()

	results := make([]Result, 0, len(d.cached))
	for _, result := range d.cached {
		results = append(results, result)
	}

	return results
}
