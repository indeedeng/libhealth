count
=====

About
-----
A count is used to keep track of accumulating values over time intervals,
and export them to /private/v as a list of numbers. Thresholds can be
defined to be evaluated over the values, and plugged into the healthcheck
framework.

Example
-------

Create a count.IntCounter, and set a name for the varexp name,
with 20 buckets representing 5 minutes of time each.

```go
	counter, err := count.Ints("my-counter", libhealth.SizeFiveMinutes, 20)
```

Call .Increment to increment current value.

```go
	for range time.Tick(30 * time.Second) {
		nodes, err := fetchLiveNodes(e.maxNodeAge, e.session)
		if err != nil {
			counter.Increment(1)
			log.Warnf("session fetchLiveNodes error: %v", err)
		}
		...
	}
```

Define a Healthcheck using count.Threshold.

```go
// DependencySet creates the set of healthchecks for background fetch task
func DependencySet(counter counts.IntCounter) libhealth.DependencySet {
	// thresholds for rate of fetch errors (1 fetch per 30 seconds, buckets of 5 minutes)
	counter.Set(count.MaxIntThreshold{
		Threshold:   1,
		Severity:	health.MINOR,
		Description: "failed to fetch data once in the last 5 minutes",
	}).Set(count.MaxSumIntThreshold{
		Threshold:   5,
		Severity:	health.MAJOR,
		Description: "failed to fetch data>= 5 times in 5 minutes",
	}).Set(count.MaxSumIntThreshold{
		Threshold:   10,
		Severity:	health.OUTAGE,
		Description: "failed to fetch data >= 10 times in 5 minutes",
	})

	return libhealth.NewBasicDependencySet(
		libhealth.NewMonitor(
			"rate-of-data-fetch-errors",
			"if this is non-zero we cannot talk to the service",
			"https://example.com/TODO",
			libhealth.STRONG,
			counter.Health,
		),
	)
}
```
