package libhealth

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	ok = `
{
  "hostname" : "tst-log3",
  "duration" : 0,
  "condition" : "OK",
  "dcStatus" : "OK"
}`

	outage = `
{
  "hostname" : "tst-log3",
  "duration" : 0,
  "condition" : "OUTAGE",
  "dcStatus" : "OUTAGE"
}`
)

type TransitiveSuite struct {
	suite.Suite
}

func Test_TransitiveSuite(t *testing.T) {
	suite.Run(t, new(TransitiveSuite))
}

func (s *TransitiveSuite) Test_ok() {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, ok)
		}))
	defer ts.Close()

	tm := TransitiveMonitor(
		ts.URL,
		"test-hc-ok",
		"example description",
		"www.example.org",
		REQUIRED,
		nil,
	)

	result := tm.Check(context.TODO())

	require.Equal(s.T(), REQUIRED, result.Urgency)
	require.Equal(s.T(), OK, result.Status)
	require.Contains(s.T(), result.Message, "200 OK")
}

func (s *TransitiveSuite) Test_outage() {
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, outage)
		}))
	defer ts.Close()

	tm := TransitiveMonitor(
		ts.URL,
		"test-hc-outage",
		"example description",
		"www.example.org",
		REQUIRED,
		nil,
	)

	result := tm.Check(context.TODO())

	require.Equal(s.T(), REQUIRED, result.Urgency)
	require.Equal(s.T(), OUTAGE, result.Status)
	require.Contains(s.T(), result.Message, "500 Internal Server Error")
}

func (s *TransitiveSuite) Test_no_connection() {
	tm := TransitiveMonitor(
		"http://localhost:0",
		"test-hc-no-connection",
		"example description",
		"www.example.org",
		REQUIRED,
		nil,
	)

	result := tm.Check(context.TODO())

	require.Equal(s.T(), REQUIRED, result.Urgency)
	require.Equal(s.T(), OUTAGE, result.Status)
	require.Contains(s.T(), result.Message, "error checking transitive monitor")
}
