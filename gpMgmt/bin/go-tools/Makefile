# Default target
.DEFAULT_GOAL := all

# Environment settings
LINUX_ENV := env GOOS=linux GOARCH=amd64
MAC_ENV := env GOOS=darwin GOARCH=amd64

# Directories
DIRS := gpservice gpctl

# Phony targets
.PHONY: all unit integration test-coverage depend-dev lint build build_linux build_mac proto cert install format pipeline clean uninstall

# Default target
all: build

# Unit tests
unit:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) unit;)

# Integration tests
integration: PACKAGE = ''
integration: UTILITY = '""'
integration:
ifeq ($(PACKAGE), gpservice)
	@$(MAKE) -C gpservice integration UTILITY=$(UTILITY) FILE=$(FILE)
else ifeq ($(PACKAGE), gpctl)
	@$(MAKE) -C gpctl integration UTILITY=$(UTILITY) FILE=$(FILE)
else
	@$(MAKE) -C gpservice integration UTILITY=$(UTILITY) FILE=$(FILE)
	@$(MAKE) -C gpctl integration UTILITY=$(UTILITY) FILE=$(FILE)
endif

# Test coverage
test-coverage:
	@echo "Coverage for gpservice"
	@$(MAKE) -C gpservice test-coverage
	@echo "\n\nCoverage for gpctl"
	@$(MAKE) -C gpctl test-coverage

# Dependencies for development
depend-dev: export GOBIN := $(CURDIR)/dev-bin
depend-dev: export GOFLAGS := -mod=readonly # do not update dependencies during installation
depend-dev:
	mkdir -p $(GOBIN)
	go install github.com/golang/protobuf/protoc-gen-go@v1.3.2
	go install github.com/golang/mock/mockgen@v1.6.0
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.52.2

# Linting
lint:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) lint;)

# Building
build:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) build;)

build_linux:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) build_linux;)

build_mac:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) build_mac;)

# Proto generation
proto:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) proto;)

# Certificate generation
cert:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) cert;)

# Installation
install:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) install;)

# Formatting
format:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) format;)

# Pipeline
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
GIT_REMOTE ?= $(shell git ls-remote --get-url)
PIPELINE_NAME ?= gpdb-dev-$(shell git rev-parse --abbrev-ref HEAD | tr '/' '-')
pipeline:
	fly -t dev set-pipeline -p $(PIPELINE_NAME) \
		-c ci/pipelines/pipeline.yml \
		-l ../../../concourse/vars/common_prod.yml \
		-l ../../../concourse/vars/common_dev.yml \
		-v gpdb-git-remote=$(GIT_REMOTE) \
		-v gpdb-git-branch=$(GIT_BRANCH) \
		-v pipeline-name=$(PIPELINE_NAME)

# Cleaning
clean:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) clean;)

# Uninstallation
uninstall:
	@$(foreach dir, $(DIRS), $(MAKE) -C $(dir) uninstall;)
