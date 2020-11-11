package libhealth

// DependencySet is an interface which can be used to represent a set of
// HealthChecker instances. These HealthChecker instances are used to determine the
// healthiness of a service.
type DependencySet interface {
	Register(monitors ...HealthMonitor)
	Background() Summary
	Live() Summary
}

// HealthTracker is the interface a framework that takes care of managing
// a HealthServer should satisfy so that it can register dependencies defined
// by the thing the framework is running.
type HealthTracker interface {
	DependencySet() DependencySet       // implemented by your app
	RegisterDependencies(DependencySet) // implemented by libhealth
}
