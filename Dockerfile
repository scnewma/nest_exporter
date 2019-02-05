FROM golang:1.11.1-alpine3.8 as builder

ARG VERSION
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

RUN adduser -D -g '' appuser

ADD . ${GOPATH}/src/app/
WORKDIR ${GOPATH}/src/app

RUN go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s -X github.com/scnewma/nest_exporter/version.Version=${VERSION}" -o /go/bin/nest_exporter

FROM scratch
ARG VCS_REF
ARG BUILD_DATE

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/scnewma/nest_exporter" \
      org.label-schema.docker.cmd="docker run -e NEST_TOKEN=token -p 9264:9264 scnewma/nest_exporter" \
      org.label-schema.docker.params="NEST_TOKEN=token to access Nest API,LISTEN_ADDRESS=Address to listen on. Default: 0.0.0.0:9264,METRICS_PATH=Path under which to expose metrics. Default: /metrics,LOG_LEVEL=Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]" \
      org.label-schema.schema-version="1.0" \
      maintainer="tom@whi.tw"

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /go/bin/nest_exporter /go/bin/nest_exporter

EXPOSE 9264

USER appuser

ENTRYPOINT [ "/go/bin/nest_exporter" ]
