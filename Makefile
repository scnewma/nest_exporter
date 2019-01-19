PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test
test: lint
	go test $(PKGS)

BIN_DIR := $(shell go env GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter --errors ./... --vendor

BINARY := nest_exporter
VERSION ?= vlatest
BASE_PKG := github.com/scnewma/nest_exporter
LD_FLAGS := -ldflags "-X $(BASE_PKG)/version.Version=$(VERSION)"

PLATFORMS := linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build $(LD_FLAGS) -o release/$(BINARY)-$(VERSION)-$(os)-amd64

.PHONY: release
release: linux darwin

.PHONY: clean
clean:
	rm -rf release/
