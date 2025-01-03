SHELL := bash

REPOSITORY ?= localhost
CONTAINER_NAME ?= simple-fileserver
TAG ?= latest

default: build

build:
	hack/build.sh

build-image:
	podman build -t $(REPOSITORY)/$(CONTAINER_NAME):$(TAG) .

test:
	go test -v -covermode=atomic -coverprofile=coverprofile.out ./...

update-deps:
	hack/update-deps.sh

lint:
	golangci-lint run -v

fmt:
	gofmt -s -w ./cmd ./pkg

validate:
	hack/validate.sh

coverprofile:
	hack/coverprofile.sh

.PHONY: \
	default \
	build \
	build-image \
	test \
	update-deps \
	lint \
	fmt \
	validate \
	coverprofile \
	$(NULL)
