gauge
=====

About
-----
A gauge is used to keep track of discrete values over time,
and export them to /private/v as a list of numbers. Thresholds
can be defined to be evaluated over the values, and plugged into
the healthcheck framework.

Example
-------

Create a gauge.IntGauger, with a name for the varexp name, and a history of 30 values.

```go
numNodes, err = gauge.Ints("my-gauge-name", 30)
```

Call .Gauge to emit a new value.

```go
for range time.Tick(30 * time.Second) {
	nodes, _ := fetchLiveNodes(e.maxNodeAge, e.session)
	numNodes.Gauge(len(nodes))
}
```

Define a Healthcheck using various gauge.Threshold.

```go
// DependencySet creates the set of healthchecks for BongoViewer.
func DependencySet(numNodes gauge.IntGauger) libhealth.DependencySet {
	numNodes.Set(gauge.MinIntThreshold{
		Threshold:   0,
		LastN:       1,
		Severity:    health.MINOR,
		Description: "lost connectivity once, stay calm",
	}).Set(gauge.MinIntThreshold{
		Threshold:   0,
		AnyN:        3,
		Severity:    health.MAJOR,
		Description: "lost connectivity a few times, network spotty?",
	}).Set(gauge.MinIntThreshold{
		Threshold:   0,
		LastN:       5,
		Severity:    health.OUTAGE,
		Description: "lost connectivity five times in a row, network is bad",
	})
 
	return libhealth.NewBasicDependencySet(
		libhealth.NewMonitor(
			"number-of-live-nodes-found",
			"if this goes to zero we have lost connectivity",
			"https://example.com/TODO",
			libhealth.REQUIRED,
			b.numNodes.Health,
		),
	)
}
```
