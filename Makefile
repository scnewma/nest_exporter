BINARY=nest_exporter

VERSION=`git describe --tags`

LD_FLAGS=-ldflags "-X version.Version=${VERSION}"

build: clean
	go build ${LD_FLAGS} -o bin/${BINARY}

clean:
	if [ -f bin/${BINARY} ]; then rm bin/${BINARY} ; fi

.PHONY: clean build
