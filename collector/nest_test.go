package collector

import (
	"sync"
	"testing"

	"github.com/jtsiros/nest/device"
	"github.com/prometheus/client_golang/prometheus"
)

func TestDescribe(t *testing.T) {
	collector := nestCollector{mockClient()}
	resultsCh := make(chan *prometheus.Desc)
	var results []*prometheus.Desc

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for m := range resultsCh {
			results = append(results, m)
		}
		wg.Done()
	}()

	collector.Describe(resultsCh)
	close(resultsCh)

	wg.Wait()

	expectedCount := 6
	if len(results) != expectedCount {
		t.Errorf("expected %d metrics, got %d\n", expectedCount, len(results))
	}
}

func TestCollect(t *testing.T) {
	collector := nestCollector{mockClient()}
	resultsCh := make(chan prometheus.Metric)
	var results []prometheus.Metric

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for m := range resultsCh {
			results = append(results, m)
		}
		wg.Done()
	}()

	collector.Collect(resultsCh)
	close(resultsCh)

	wg.Wait()

	expectedCount := 6
	if len(results) != expectedCount {
		t.Errorf("expected %d metrics, got %d\n", expectedCount, len(results))
	}
}

func mockClient() mockNestClient {
	return mockNestClient{
		DevicesFn: func() (*device.Devices, error) {
			return &device.Devices{
				Thermostats: map[string]*device.Thermostat{
					"x": {
						Name:                "Thermy",
						Humidity:            50,
						TargetTemperatureC:  20.2,
						TargetTemperatureF:  70,
						AmbientTemperatureC: 20.4,
						AmbientTemperatureF: 72,
					},
				},
			}, nil
		},
	}
}

type mockNestClient struct {
	DevicesFn func() (*device.Devices, error)
}

func (c mockNestClient) Devices() (*device.Devices, error) {
	return c.DevicesFn()
}
