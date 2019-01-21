# Nest Exporter

[![Build Status](https://travis-ci.com/scnewma/nest_exporter.svg?branch=master)](https://travis-ci.com/scnewma/nest_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/scnewma/nest_exporter)](https://goreportcard.com/report/github.com/scnewma/nest_exporter)

Prometheus exporter for Nest thermostat devices.

## Building and Running

Prerequisites:

* Go
* [Nest developer API token](https://developers.nest.com/guides/api/how-to-auth) (use PIN-based authorization)

Building:

```
go get github.com/scnewma/nest_exporter
cd ${GOPATH}/src/github.com/scnewma/nest_exporter
make
./nest_exporter --nest.token=[TOKEN]
```

## Running Tests

```
make test
```

## Using Docker

```
# put Nest token in .token file then
make docker-run

# OR

docker run -p 9264:9264 -d scnewma/nest_exporter --nest.token=[TOKEN]
```

## Credit

The idea for how to structure/build the project and accept flags for running the exporter came from reviewing [node_exporter](https://github.com/prometheus/node_exporter).

There is already a [nest exporter](https://github.com/jcollie/nest_exporter) available. I borrowed the metrics port from that exporter. I wrote a new exporter for learning purposes.
