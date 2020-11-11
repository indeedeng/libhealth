package libhealth

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMonitor_Check(t *testing.T) {
	monitor := NewMonitor("test", "", "", REQUIRED, func(ctx context.Context) Health {
		return NewHealth(OK, "okay")
	}, nil)

	health := monitor.Check(context.Background())
	require.Equal(t, health.Status, OK)
}

func TestMonitor_Check_Channel(t *testing.T) {
	statusChan := make(chan HealthStatus)
	monitor := NewMonitor("test", "", "", REQUIRED, func(ctx context.Context) Health {
		return NewHealth(OUTAGE, "outage")
	}, statusChan)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		health := monitor.Check(context.Background())
		require.Equal(t, health.Status, OUTAGE)
		wg.Done()
	}()

	status := <-statusChan
	require.Equal(t, status.Monitor, monitor)
	require.Equal(t, status.Prev, OK)
	require.Equal(t, status.Next.Status, OUTAGE)
	wg.Wait()
}
