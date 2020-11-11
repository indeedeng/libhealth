libhealth
=========
[![Go Report Card](https://goreportcard.com/badge/oss.indeed.com/go/libhealth)](https://goreportcard.com/report/oss.indeed.com/go/libhealth)
[![Build Status](https://travis-ci.com/indeedeng/libhealth.svg?branch=master)](https://travis-ci.com/indeedeng/libhealth)
[![GoDoc](https://godoc.org/oss.indeed.com/go/libhealth?status.svg)](https://godoc.org/oss.indeed.com/go/libhealth)
[![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/indeedeng/libhealth.svg)](OSSMETADATA)
[![GitHub](https://img.shields.io/github/license/indeedeng/libhealth.svg)](LICENSE)

libhealth is a Golang library that provides flexible components to report the current state of external systems
that an application depends on, as well as the current health of any internal aspects of the application.

## components
### dependency set
Most applications have a set of dependencies whose health state needs to be tracked. A dependency set can
compute the health state of an application from health monitors and their associated urgency levels.

The final health state of an application is computed from their status and urgency. A weakly coupled component
in a complete outage results in a MINOR healthcheck outage state. A more strongly coupled or required
component in an OUTAGE results in correspondingly higher outage states. The matrix below illustrates
the relationship between the "urgency" of a component and computed status.

| Status \ Urgency | REQUIRED | STRONG | WEAK   | NONE |
| ---------------- | -------- | ------ | ------ | ---- |
| OUTAGE           | OUTAGE   | MAJOR  | MINOR  | OK   |
| MAJOR            | MAJOR    | MAJOR  | MINOR  | OK   |
| MINOR            | MINOR    | MINOR  | MINOR  | OK   |
| OK               | OK       | OK     | OK     | OK   |


While applications can implement a dependency set of their own, a basic dependency set is provided which fits
most use cases. Applications typically need one basic dependency set, and libhealth will update monitors
and track their state in background goroutines.

Example:
```go
import 	"oss.indeed.com/go/libhealth"

func setupHealth() {
	deps := libhealth.NewBasicDependencySet()
	deps.Register(libhealth.NewMonitor(
                                "health-monitor-name",
                                "monitor description",
                                "https://docs/to/your/monitor",
                                libhealth.WEAK,
                                func(ctx context.Context) libhealth.Health {
                                    // calculate monitor health here
                                    return libhealth.NewHealth(libhealth.OK, "everything is fine")
                                }))
}
```

### healthcheck endpoints
Typical applications expose several healthcheck endpoints to an HTTP server for tracking their state.
libhealth provides two classes of endpoints: public "info" and private healthcheck endpoints. The
public endpoints are typically consumed by other software (e.g loadbalancers, HAProxy, nginx, etc).
They return a 200 or 500 status code and very simple json payload indicating the source of the response.

An example response to the /info endpoints is shown below:
```json
{
  "condition" : "OK",
  "duration" : 0,
  "hostname" : "aus-worker11"
}
```

The private healthcheck endpoints expose significantly more information about the runtime and environment
of the process to aid debugging outages, however these should only be exposed to whitelisted ips.
If this is not possible, consider only exposing endpoints for the less verbose public endpoints.

A common pattern is to expose the following endpoints:
```
/info/healthcheck
/info/healthcheck/live
/private/healthcheck
/private/healthcheck/live
```

HTTP handlers are provided by libhealth for serving each of these routes:
```go
import "oss.indeed.com/go/libhealth"

func healthRouter(d libhealth.DependencySet) *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/info/healthcheck", libhealth.NewInfo(d))
	router.Handle("/info/healthcheck/live", libhealth.NewInfo(d))
	router.Handle("/private/healthcheck", libhealth.NewPrivate("my-app-name", d))
	router.Handle("/private/healthcheck/live", libhealth.NewPrivate("my-app-name", d))
	return router
}
```

Alternatively, you can use the helper function `WrapServeMux`, which will register all these handlers for you:
```go
import "oss.indeed.com/go/libhealth"

...
router := http.NewServeMux()
libhealth.WrapServeMux(router, "my-app-name", dependencies)
```

# Contributing

We welcome contributions! Feel free to help make `libhealth` better.

### Process

- Open an issue and describe the desired feature / bug fix before making
changes. It's useful to get a second pair of eyes before investing development
effort.
- Make the change. If adding a new feature, remember to provide tests that
demonstrate the new feature works, including any error paths. If contributing
a bug fix, add tests that demonstrate the erroneous behavior is fixed.
- Open a pull request. Automated CI tests will run. If the tests fail, please
make changes to fix the behavior, and repeat until the tests pass.
- Once everything looks good, one of the indeedeng members will review the
PR and provide feedback.

# Maintainers

The `oss.indeed.com/go/libhealth` project is maintained by Indeed Engineering.

While we are always busy helping people get jobs, we will try to respond to
GitHub issues, pull requests, and questions within a couple of business days.

# Code of Conduct

`oss.indeed.com/go/libhealth` is governed by the [Contributer Covenant v1.4.1](CODE_OF_CONDUCT.md)

For more information please contact opensource@indeed.com.

## License

The `oss.indeed.com/go/libhealth` project is open source under the [Apache 2.0](LICENSE) license.
