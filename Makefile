APP_NAME := nest_exporter
VERSION ?= $(shell git rev-parse --abbrev-ref HEAD) 
DOCKER_REPO ?= scnewma
DOCKER_TAG ?= $(VERSION) 

.DEFAULT_GOAL := clean-build

.PHONY: release
release: docker docker-publish

.PHONY: docker
docker: linux
	@docker build -t $(DOCKER_REPO)/$(APP_NAME):$(DOCKER_TAG) .

.PHONY: docker-run
docker-run:
	@docker run -p 9264:9264 $(DOCKER_REPO)/$(APP_NAME):$(DOCKER_TAG) --nest.token=$(shell cat .token)

.PHONY: docker-publish
docker-publish:
	@docker push $(DOCKER_REPO)/$(APP_NAME):$(DOCKER_TAG)

PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test
test: 
	go test $(PKGS)

BINARY := $(APP_NAME)
BASE_PKG := github.com/scnewma/nest_exporter
LD_FLAGS := -ldflags "-X $(BASE_PKG)/version.Version=$(VERSION)"

$(BINARY):
	mkdir -p bin
	go build $(LD_FLAGS) -o $(BINARY)

build: $(BINARY)

linux: clean
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build $(LD_FLAGS) -o $(BINARY)

.PHONY: clean
clean:
	rm -f $(BINARY)

.PHONY: clean-build
clean-build: clean build
