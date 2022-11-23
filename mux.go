package libhealth

import (
	"net/http"
)

// WrapServeMux will wrap mux with Handlers for
//   - /private/healthcheck
//   - /private/healthcheck/live
//   - /info/healthcheck
//   - /info/healthcheck/live
func WrapServeMux(
	mux *http.ServeMux,
	appname string,
	provided DependencySet,
	additional ...HealthMonitor,
) {
	provided.Register(additional...)

	infoHandler := NewInfo(provided)
	privHandler := NewPrivate(appname, provided)

	mux.Handle(InfoHealthCheck, infoHandler)
	mux.Handle(InfoHealthCheckLive, infoHandler)
	mux.Handle(PrivateHealthCheck, privHandler)
	mux.Handle(PrivateHealthCheckLive, privHandler)
}
