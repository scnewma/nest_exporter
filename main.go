package main

import (
	"os"
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
		app           = kingpin.New(
			"nest_exporter",
			"A prometheus exporter for Nest.",
		)
		listenAddress = app.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default(":9264").Envar("LISTEN_ADDRESS").String()
		metricsPath = app.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").Envar("METRICS_PATH").String()
		logLevel = app.Flag(
			"log.level",
			"Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]",
		).Default(log.InfoLevel.String()).Envar("LOG_LEVEL").String()
		token = app.Flag(
			"nest.token",
			"Nest authorization token that has access to developer API.",
		).Required().Envar("NEST_TOKEN").String()
	)
	kingpin.Version(version.Print())
	kingpin.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	setLogLevel(*logLevel)

	log.Infof("Starting %s %s", app.Name, version.Info())

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
