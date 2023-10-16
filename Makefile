# Makefile for go-messaging.

# Detect the operating system and architecture.

include Makefile.osdetect

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

# "Simple expanded" variables (':=')

# PROGRAM_NAME is the name of the GIT repository.
PROGRAM_NAME := $(shell basename `git rev-parse --show-toplevel`)
MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MAKEFILE_DIRECTORY := $(shell dirname $(MAKEFILE_PATH))
TARGET_DIRECTORY := $(MAKEFILE_DIRECTORY)/target
BUILD_VERSION := $(shell git describe --always --tags --abbrev=0 --dirty  | sed 's/v//')
BUILD_TAG := $(shell git describe --always --tags --abbrev=0  | sed 's/v//')
BUILD_ITERATION := $(shell git log $(BUILD_TAG)..HEAD --oneline | wc -l | sed 's/^ *//')
GIT_REMOTE_URL := $(shell git config --get remote.origin.url)
GO_PACKAGE_NAME := $(shell echo $(GIT_REMOTE_URL) | sed -e 's|^git@github.com:|github.com/|' -e 's|\.git$$||' -e 's|Senzing|senzing|')

# Recursive assignment ('=')

GO_OSARCH = $(subst /, ,$@)
GO_OS = $(word 1, $(GO_OSARCH))
GO_ARCH = $(word 2, $(GO_OSARCH))

# Conditional assignment. ('?=')
# Can be overridden with "export"
# Example: "export LD_LIBRARY_PATH=/path/to/my/senzing/g2/lib"

LD_LIBRARY_PATH ?= /opt/senzing/g2/lib

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

-include Makefile.$(OSTYPE)
-include Makefile.$(OSTYPE)_$(OSARCH)


.PHONY: hello-world
hello-world: hello-world-osarch-specific

# -----------------------------------------------------------------------------
# Generate code
# -----------------------------------------------------------------------------

.PHONY: generate-code
generate-code: generate-csharp generate-go generate-java generate-python generate-ruby generate-rust generate-typescript


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
# Dependency management
# -----------------------------------------------------------------------------

.PHONY: dependencies
dependencies:
	@go get -u ./...
	@go get -t -u ./...
	@go mod tidy

# -----------------------------------------------------------------------------
# Build
#  - docker-build: https://docs.docker.com/engine/reference/commandline/build/
# -----------------------------------------------------------------------------

PLATFORMS := darwin/amd64 linux/amd64 windows/amd64
$(PLATFORMS):
	@echo Building $(TARGET_DIRECTORY)/$(GO_OS)-$(GO_ARCH)/$(PROGRAM_NAME)
	@GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build -o $(TARGET_DIRECTORY)/$(GO_OS)-$(GO_ARCH)/$(PROGRAM_NAME)


.PHONY: build
build: build-osarch-specific

# -----------------------------------------------------------------------------
# Test
# -----------------------------------------------------------------------------

.PHONY: test
test: test-osarch-specific

# -----------------------------------------------------------------------------
# Run
# -----------------------------------------------------------------------------

.PHONY: run
run: run-osarch-specific

# -----------------------------------------------------------------------------
# Clean
# -----------------------------------------------------------------------------

.PHONY: clean-csharp
clean-csharp:
	@rm $(MAKEFILE_DIRECTORY)/csharp/* || true


.PHONY: clean-go
clean-go:
	@go clean -cache
	@go clean -testcache
	@rm -f $(GOPATH)/bin/$(PROGRAM_NAME) || true
	@rm $(MAKEFILE_DIRECTORY)/go/typedef/typedef.go || true


.PHONY: clean-java
clean-java:
	@rm $(MAKEFILE_DIRECTORY)/java/* || true


.PHONY: clean-python
clean-python:
	@rm $(MAKEFILE_DIRECTORY)/python/typedef/* || true


.PHONY: clean-ruby
clean-ruby:
	@rm $(MAKEFILE_DIRECTORY)/ruby/* || true


.PHONY: clean-rust
clean-rust:
	@rm $(MAKEFILE_DIRECTORY)/rust/* || true


.PHONY: clean-typescript
clean-typescript:
	@rm $(MAKEFILE_DIRECTORY)/typescript/* || true

# -----------------------------------------------------------------------------
# Utility targets
# -----------------------------------------------------------------------------

.PHONY: clean
clean: clean-osarch-specific
	@go clean -cache
	@go clean -testcache


.PHONY: clean-generated code
clean-generated code: clean-go clean-java clean-python clean-ruby clean-rust clean-typescript

	
.PHONY: help
help:
	@echo "Build $(PROGRAM_NAME) version $(BUILD_VERSION)-$(BUILD_ITERATION)".
	@echo "Makefile targets:"
	@$(MAKE) -pRrq -f $(firstword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs


.PHONY: print-make-variables
print-make-variables:
	@$(foreach V,$(sort $(.VARIABLES)), \
		$(if $(filter-out environment% default automatic, \
		$(origin $V)),$(warning $V=$($V) ($(value $V)))))


.PHONY: setup
setup: setup-osarch-specific


.PHONY: update-pkg-cache
update-pkg-cache:
	@GOPROXY=https://proxy.golang.org GO111MODULE=on \
		go get $(GO_PACKAGE_NAME)@$(BUILD_TAG)
