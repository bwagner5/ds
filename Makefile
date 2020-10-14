BUILD_DIR ?= $(dir $(realpath -s $(firstword $(MAKEFILE_LIST))))/build
VERSION ?= $(shell git describe --tags --always --dirty)
GOOS ?= $(uname | tr '[:upper:]' '[:lower:]')
GOARCH ?= $([[ uname -m = "x86_64" ]] && amd64 || arm64 )
GOPROXY ?= "https://proxy.golang.org,direct"

$(shell mkdir -p ${BUILD_DIR})

all: verify test build

build:
	go build -a -ldflags="-s -w -X main.versionID=${VERSION}" -o ${BUILD_DIR}/ds-${GOOS}-${GOARCH} ${BUILD_DIR}/../cmd/main.go

test:
	go test -bench=. ${BUILD_DIR}/../... -v -coverprofile=coverage.out -covermode=atomic -outputdir=${BUILD_DIR}

verify:
	go mod tidy
	go mod download
	go vet ./...
	go fmt ./...

version:
	@echo ${VERSION}

help:
	@grep -E '^[a-zA-Z_-]+:.*$$' $(MAKEFILE_LIST) | sort

.PHONY: all build test verify help