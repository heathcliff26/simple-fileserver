SHELL := bash

REPOSITORY ?= localhost
CONTAINER_NAME ?= simple-fileserver
TAG ?= latest

GO_BUILD_FLAGS ?= -ldflags="-w -s"

default: build

build:
	go build $(GO_BUILD_FLAGS) -o bin/simple-fileserver ./cmd/

build-image:
	podman build -t $(REPOSITORY)/$(CONTAINER_NAME):$(TAG) .

test:
	go test -v ./...

.PHONY: \
	default \
	build \
	build-image \
	test \
	$(NULL)
