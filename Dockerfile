FROM quay.io/prometheus/busybox:glibc

COPY nest_exporter /bin/nest_exporter

EXPOSE 9264
USER nobody
ENTRYPOINT [ "/bin/nest_exporter" ]
