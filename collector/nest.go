package collector

import (
	"github.com/jtsiros/nest/device"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	nestUpDesc = prometheus.NewDesc(
		"nest_up",
		"Whether the scrape succeeded.",
		nil, nil)

	currentTempFDesc = prometheus.NewDesc(
		"nest_current_temperature_fahrenheit",
		"Current temperature monitored by Nest thermostat.",
		[]string{"thermostat"}, nil)

	currentTempCDesc = prometheus.NewDesc(
		"nest_current_temperature_celcius",
		"Current temperature monitored by Nest thermostat.",
		[]string{"thermostat"}, nil)

	targetTempCDesc = prometheus.NewDesc(
		"nest_target_temperature_celcius",
		"Target temperature of the Nest thermostat.",
		[]string{"thermostat"}, nil)

	targetTempFDesc = prometheus.NewDesc(
		"nest_target_temperature_fahrenheit",
		"Target temperature of the Nest thermostat.",
		[]string{"thermostat"}, nil)

	currentHumidityDesc = prometheus.NewDesc(
		"nest_current_humidity",
		"Current humidity monitored by Nest thermostat.",
		[]string{"thermostat"}, nil)
)

// NestClient is a client that can perform operations
// to retrieve device data from Nest products.
type NestClient interface {
	Devices() (*device.Devices, error)
}

// NewNestCollector creates a nestCollector given a nest client.
// This function will fatal if no client is provided.
func NewNestCollector(client NestClient) nestCollector {
	if client == nil {
		log.Fatal("client must be provided")
	}

	return nestCollector{client}
}

type nestCollector struct {
	client NestClient
}

func (c nestCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c nestCollector) Collect(ch chan<- prometheus.Metric) {
	devices, err := c.client.Devices()
	if err != nil {
		log.WithError(err).Warn("failed to get devices")
		ch <- newUpMetric(false)
		return
	}

	ch <- newUpMetric(true)

	for _, t := range devices.Thermostats {
		ch <- newCurrentTempFMetric(t.Name, t.AmbientTemperatureF)
		ch <- newCurrentTempCMetric(t.Name, t.AmbientTemperatureC)
		ch <- newTargetTempFMetric(t.Name, t.TargetTemperatureF)
		ch <- newTargetTempCMetric(t.Name, t.TargetTemperatureC)
		ch <- newCurrentHumidityMetric(t.Name, t.Humidity)
	}
}

func newUpMetric(up bool) prometheus.Metric {
	val := 0

	if up {
		val = 1
	}

	return prometheus.MustNewConstMetric(
		nestUpDesc,
		prometheus.GaugeValue,
		float64(val),
	)
}

func newCurrentTempFMetric(thermostatName string, temp int) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		currentTempFDesc,
		prometheus.GaugeValue,
		float64(temp),
		thermostatName,
	)
}

func newCurrentTempCMetric(thermostatName string, temp float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		currentTempCDesc,
		prometheus.GaugeValue,
		temp,
		thermostatName,
	)
}

func newTargetTempFMetric(thermostatName string, temp int) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		targetTempFDesc,
		prometheus.GaugeValue,
		float64(temp),
		thermostatName,
	)
}

func newTargetTempCMetric(thermostatName string, temp float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		targetTempCDesc,
		prometheus.GaugeValue,
		temp,
		thermostatName,
	)
}

func newCurrentHumidityMetric(thermostatName string, humidity int) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		currentHumidityDesc,
		prometheus.GaugeValue,
		float64(humidity),
		thermostatName,
	)
}
