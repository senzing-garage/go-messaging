# Makefile for Go project

# Detect the operating system and architecture.

include makefiles/osdetect.mk

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

# "Simple expanded" variables (':=')

# PROGRAM_NAME is the name of the GIT repository.
PROGRAM_NAME := $(shell basename `git rev-parse --show-toplevel`)
MAKEFILE_PATH := $(abspath $(firstword $(MAKEFILE_LIST)))
MAKEFILE_DIRECTORY := $(shell dirname $(MAKEFILE_PATH))
TARGET_DIRECTORY := $(MAKEFILE_DIRECTORY)/target
DIST_DIRECTORY := $(MAKEFILE_DIRECTORY)/dist
BUILD_TAG := $(shell git describe --always --tags --abbrev=0  | sed 's/v//')
BUILD_ITERATION := $(shell git log $(BUILD_TAG)..HEAD --oneline | wc -l | sed 's/^ *//')
BUILD_VERSION := $(shell git describe --always --tags --abbrev=0 --dirty  | sed 's/v//')
DOCKER_CONTAINER_NAME := $(PROGRAM_NAME)
DOCKER_IMAGE_NAME := senzing/$(PROGRAM_NAME)
DOCKER_BUILD_IMAGE_NAME := $(DOCKER_IMAGE_NAME)-build
GIT_REMOTE_URL := $(shell git config --get remote.origin.url)
GIT_REPOSITORY_NAME := $(shell basename `git rev-parse --show-toplevel`)
GIT_VERSION := $(shell git describe --always --tags --long --dirty | sed -e 's/\-0//' -e 's/\-g.......//')
GO_PACKAGE_NAME := $(shell echo $(GIT_REMOTE_URL) | sed -e 's|^git@github.com:|github.com/|' -e 's|\.git$$||' -e 's|Senzing|senzing|')
PATH := $(MAKEFILE_DIRECTORY)/bin:$(PATH)

# Recursive assignment ('=')

GO_OSARCH = $(subst /, ,$@)
GO_OS = $(word 1, $(GO_OSARCH))
GO_ARCH = $(word 2, $(GO_OSARCH))

# Conditional assignment. ('?=')
# Can be overridden with "export"
# Example: "export LD_LIBRARY_PATH=/path/to/my/senzing/er/lib"

GOBIN ?= $(shell go env GOPATH)/bin
LD_LIBRARY_PATH ?= /opt/senzing/er/lib

# Export environment variables.

.EXPORT_ALL_VARIABLES:

# -----------------------------------------------------------------------------
# The first "make" target runs as default.
# -----------------------------------------------------------------------------

.PHONY: default
default: help

# -----------------------------------------------------------------------------
# Operating System / Architecture targets
# -----------------------------------------------------------------------------

-include makefiles/$(OSTYPE).mk
-include makefiles/$(OSTYPE)_$(OSARCH).mk


.PHONY: hello-world
hello-world: hello-world-osarch-specific

# -----------------------------------------------------------------------------
# Dependency management
# -----------------------------------------------------------------------------

.PHONY: dependencies-for-development
dependencies-for-development: dependencies-for-development-osarch-specific
	@go install github.com/daixiang0/gci@latest
	@go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
	@go install github.com/vladopajic/go-test-coverage/v2@latest
	@go install golang.org/x/tools/cmd/godoc@latest


.PHONY: dependencies
dependencies:
	@go get -u ./...
	@go get -t -u ./...
	@go mod tidy

# -----------------------------------------------------------------------------
# Setup
# -----------------------------------------------------------------------------

.PHONY: setup
setup: setup-osarch-specific

# -----------------------------------------------------------------------------
# Lint
# -----------------------------------------------------------------------------

.PHONY: lint
lint: golangci-lint

# -----------------------------------------------------------------------------
# Build
# -----------------------------------------------------------------------------

PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64
$(PLATFORMS):
	$(info Building $(TARGET_DIRECTORY)/$(GO_OS)-$(GO_ARCH)/$(PROGRAM_NAME))
	@GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build -o $(TARGET_DIRECTORY)/$(GO_OS)-$(GO_ARCH)/$(PROGRAM_NAME)


.PHONY: build
build: build-osarch-specific

# -----------------------------------------------------------------------------
# Run
# -----------------------------------------------------------------------------

.PHONY: run
run: run-osarch-specific

# -----------------------------------------------------------------------------
# Test
# -----------------------------------------------------------------------------

.PHONY: test
test: test-osarch-specific

# -----------------------------------------------------------------------------
# Coverage
# -----------------------------------------------------------------------------

.PHONY: coverage
coverage: coverage-osarch-specific


.PHONY: check-coverage
check-coverage: export SENZING_LOG_LEVEL=TRACE
check-coverage:
	@go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...
	@${GOBIN}/go-test-coverage --config=.github/coverage/testcoverage.yaml

# -----------------------------------------------------------------------------
# Documentation
# -----------------------------------------------------------------------------

.PHONY: documentation
documentation: documentation-osarch-specific

# -----------------------------------------------------------------------------
# Generate code
# -----------------------------------------------------------------------------

.PHONY: generate
generate: generate-csharp generate-go generate-java generate-python generate-ruby generate-rust generate-typescript


.PHONY: generate-csharp
generate-csharp:
	jtd-codegen \
		--csharp-system-text-namespace Senzing \
		--csharp-system-text-out ./csharp \
		--root-name SenzingMessage \
		message-RFC8927.json


.PHONY: generate-go
generate-go:
	jtd-codegen \
		--go-out ./go/typedef \
		--go-package typedef \
		--root-name SenzingMessage \
		message-RFC8927.json


.PHONY: generate-java
generate-java:
	jtd-codegen \
		--java-jackson-out ./java \
		--java-jackson-package com.senzing.schema \
		--root-name SenzingMessage \
		message-RFC8927.json


.PHONY: generate-python
generate-python:
	jtd-codegen \
		--python-out ./python/typedef \
		--root-name SenzingMessage \
		message-RFC8927.json


.PHONY: generate-ruby
generate-ruby:
	jtd-codegen \
		--root-name SenzingMessage \
		--ruby-module SenzingTypeDef \
		--ruby-out ./ruby \
		--ruby-sig-module SenzingSig \
		message-RFC8927.json


.PHONY: generate-rust
generate-rust:
	jtd-codegen \
		--root-name SenzingMessage \
		--rust-out ./rust \
		message-RFC8927.json


.PHONY: generate-typescript
generate-typescript:
	jtd-codegen \
		--root-name SenzingMessage \
		--typescript-out ./typescript \
		message-RFC8927.json

# -----------------------------------------------------------------------------
# Clean
# -----------------------------------------------------------------------------

.PHONY: clean
clean: clean-osarch-specific
	@go clean -cache
	@go clean -testcache


.PHONY: clean-csharp
clean-csharp:
	@rm $(MAKEFILE_DIRECTORY)/csharp/* || true


.PHONY: clean-go
clean-go:
	@go clean -cache
	@go clean -testcache
	@rm -f $(GOPATH)/bin/$(PROGRAM_NAME) || true
	@rm -f $(MAKEFILE_DIRECTORY)/go/typedef/typedef.go || true


.PHONY: clean-java
clean-java:
	@rm -fr $(MAKEFILE_DIRECTORY)/java/* || true


.PHONY: clean-python
clean-python:
	@rm -fr $(MAKEFILE_DIRECTORY)/python/typedef/* || true


.PHONY: clean-ruby
clean-ruby:
	@rm -fr $(MAKEFILE_DIRECTORY)/ruby/* || true


.PHONY: clean-rust
clean-rust:
	@rm -fr $(MAKEFILE_DIRECTORY)/rust/* || true


.PHONY: clean-typescript
clean-typescript:
	@rm -fr $(MAKEFILE_DIRECTORY)/typescript/* || true


.PHONY: clean-generated-code
clean-generated-code: clean-go clean-java clean-python clean-ruby clean-rust clean-typescript

# -----------------------------------------------------------------------------
# Utility targets
# -----------------------------------------------------------------------------

.PHONY: help
help:
	$(info Build $(PROGRAM_NAME) version $(BUILD_VERSION)-$(BUILD_ITERATION))
	$(info Makefile targets:)
	@$(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs


.PHONY: print-make-variables
print-make-variables:
	@$(foreach V,$(sort $(.VARIABLES)), \
		$(if $(filter-out environment% default automatic, \
		$(origin $V)),$(info $V=$($V) ($(value $V)))))


.PHONY: update-pkg-cache
update-pkg-cache:
	@GOPROXY=https://proxy.golang.org GO111MODULE=on \
		go get $(GO_PACKAGE_NAME)@$(BUILD_TAG)

# -----------------------------------------------------------------------------
# Specific programs
# -----------------------------------------------------------------------------

.PHONY: golangci-lint
golangci-lint:
	@${GOBIN}/golangci-lint run --config=.github/linters/.golangci.yaml
