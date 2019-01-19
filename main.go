package main

import (
	"context"
	"net/http"

	"github.com/jtsiros/nest"
	"github.com/jtsiros/nest/auth"
	"github.com/jtsiros/nest/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/scnewma/nest_exporter/collector"
	"github.com/scnewma/nest_exporter/version"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(":9264").String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		logLevel = kingpin.Flag(
			"log.level",
			"Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]",
		).Default(log.InfoLevel.String()).String()
		token = kingpin.Flag(
			"nest.token",
			"Nest authorization token that has access to developer API.",
		).Required().String()
	)
	kingpin.Version(version.Print())
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	setLogLevel(*logLevel)

	log.Infoln("Starting nest_exporter", version.Info())

	reg := prometheus.NewPedanticRegistry()

	nestClient := newNestClient(*token)

	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
		collector.NewNestCollector(nestClient),
	)

	http.Handle(*metricsPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Nest Exporter</title></head>
			<body>
			<h1>Nest Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Infoln("Listening on", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatal(err)
	}
}

func newNestClient(token string) *nest.Client {
	appConfig := config.Config{
		APIURL: config.APIURL,
	}
	conf := auth.NewConfig(appConfig)

	tok, err := auth.NewConfigWithToken(token).Token()
	if err != nil {
		log.Fatalf("Failed to get config from token, reason: %s\n", err.Error())
	}
	client := conf.Client(context.Background(), tok)

	n, err := nest.NewClient(appConfig, client)
	if err != nil {
		log.Fatalf("Failed to craete nest client, reason: %s\n", err.Error())
	}

	return n
}

func setLogLevel(logLevel string) {
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal("Unable to parse log level. Valid levels: [debug, info, warn, error, fatal]")
	}

	log.SetLevel(lvl)
}
