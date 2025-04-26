SHELL := bash

REPOSITORY ?= localhost
CONTAINER_NAME ?= simple-fileserver
TAG ?= latest

# Build the binary
build:
	hack/build.sh

# Build the container image
build-image:
	podman build -t $(REPOSITORY)/$(CONTAINER_NAME):$(TAG) .

# Run unit-tests
test:
	go test -v -coverprofile=coverprofile.out ./...

# Update dependencies
update-deps:
	hack/update-deps.sh

# Run linter
lint:
	golangci-lint run -v

# Format code
fmt:
	gofmt -s -w ./cmd ./pkg

# Validate that all generated files are up to date
validate:
	hack/validate.sh

# Generate cover profile
coverprofile:
	hack/coverprofile.sh

# Scan code for vulnerabilities using gosec
gosec:
	gosec ./...

# Clean build artifacts
clean:
	rm -rf bin coverprofiles coverprofile.out

# Show this help message
help:
	@echo "Available targets:"
	@echo ""
	@awk '/^#/{c=substr($$0,3);next}c&&/^[[:alpha:]][[:alnum:]_-]+:/{print substr($$1,1,index($$1,":")),c}1{c=0}' $(MAKEFILE_LIST) | column -s: -t
	@echo ""
	@echo "Run 'make <target>' to execute a specific target."

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
	gosec \
	clean \
	help \
	$(NULL)
