GO111MODULE=on
GOARCH=amd64
GOOS=linux
GOCMD=go
GOLINT=golangci-lint run
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

CGO_ENABLED=1
BINARY_NAME=fluentbit-plugin-natspublisher.so
VET_REPORT = vet.report
TEST_REPORT = tests.xml
SHELL=/bin/bash
DOCKERCMD=docker

ROOT := $$(git rev-parse --show-toplevel)

all: lint build
.PHONY: all

.PHONY: lint
lint: 
		$(GOLINT)

.PHONY: build
build: 
		GO111MODULE=$(GO111MODULE) CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) $(GOBUILD) -buildmode=c-shared -o $(ROOT)/docker-compose/fluent-bit/plugins/$(BINARY_NAME) -v

