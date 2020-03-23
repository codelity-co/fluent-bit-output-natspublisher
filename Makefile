GO111MODULE=on
GOCMD=go
GOLINT=golangci-lint run
GOBUILD=$(GOCMD) build

CGO_ENABLED=1
BINARY_NAME=fluentbit-plugin-natspublisher.so

ROOT := $$(git rev-parse --show-toplevel)

lint: 
		$(GOLINT)

linux: 
		GO111MODULE=$(GO111MODULE) CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GOBUILD) -buildmode=c-shared -o $(ROOT)/dist/linux/$(BINARY_NAME) -v

darwin:
		GO111MODULE=$(GO111MODULE) CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GOBUILD) -buildmode=c-shared -o $(ROOT)/dist/darwin/$(BINARY_NAME) -v