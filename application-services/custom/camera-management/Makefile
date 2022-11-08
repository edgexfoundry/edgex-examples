# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

.PHONY: build-app test clean docker tidy run-app run-tests

GO=CGO_ENABLED=1 go

# VERSION file is not needed for local development, In the CI/CD pipeline, a temporary VERSION file is written
# if you need a specific version, just override below
APPVERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)

# This pulls the version of the SDK from the go.mod file. If the SDK is the only required module,
# it must first remove the word 'required' so the offset of $2 is the same if there are multiple required modules
SDKVERSION=$(shell cat ./go.mod | grep 'github.com/edgexfoundry/app-functions-sdk-go/v2 v' | sed 's/require//g' | awk '{print $$2}')

MICROSERVICE=app-camera-management
GOFLAGS=-ldflags "-X github.com/edgexfoundry/app-functions-sdk-go/v2/internal.SDKVersion=$(SDKVERSION) -X github.com/edgexfoundry/app-functions-sdk-go/v2/internal.ApplicationVersion=$(APPVERSION)"

GIT_SHA=$(shell git rev-parse HEAD)

build: build-app
build-app: tidy
	$(GO) build $(GOFLAGS) -o $(MICROSERVICE)

tidy:
	go mod tidy

run-app:
	EDGEX_SECURITY_SECRET_STORE=false ./$(MICROSERVICE)

# TODO: Change the registries (edgexfoundry, nexus3.edgexfoundry.org:10004) below as needed.
#       Leave them as is if service is to be upstreamed to EdgeX Foundry
# NOTE: This is only used for local development. Jenkins CI does not use this make target
docker:
	docker build \
	    --build-arg http_proxy \
	    --build-arg https_proxy \
		-f Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/$(MICROSERVICE):$(GIT_SHA) \
		-t edgexfoundry/$(MICROSERVICE):$(APPVERSION)-dev \
		-t nexus3.edgexfoundry.org:10004/$(MICROSERVICE):$(APPVERSION)-dev \
		.

run-tests:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) vet ./...
	gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")
	[ "`gofmt -l $$(find . -type f -name '*.go'| grep -v "/vendor/")`" = "" ]

clean:
	rm -f $(MICROSERVICE)

vendor:
	$(GO) mod vendor


MAKEFLAGS += --warn-undefined-variables --no-builtin-rules --no-print-directory

EVAM_CMDS = run install stop sniff-events

.PHONY: web-ui serve-ui $(EVAM_CMDS:%=%-edge-video-analytics)

$(EVAM_CMDS:%=%-edge-video-analytics): %-edge-video-analytics :
	$(MAKE) -C edge-video-analytics $*

web-ui:
	$(MAKE) -C web-ui prod

serve-ui:
	$(MAKE) -C web-ui serve
