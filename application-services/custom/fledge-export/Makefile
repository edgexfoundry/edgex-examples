.PHONY: build test clean docker

GO=CGO_ENABLED=1 go

# VERSION file is not needed for local development, In the CI/CD pipeline, a temporary VERSION file is written
# if you need a specific version, just override below
APPVERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)

# This pulls the version of the SDK from the go.mod file. If the SDK is the only required module,
# it must first remove the word 'required' so the offset of $2 is the same if there are multiple required modules
SDKVERSION=$(shell cat ./go.mod | grep 'github.com/edgexfoundry/app-functions-sdk-go v' | sed 's/require//g' | awk '{print $$2}')

MICROSERVICE=fledge-export
GOFLAGS=-ldflags "-X github.com/edgexfoundry/app-functions-sdk-go/internal.SDKVersion=$(SDKVERSION) -X github.com/edgexfoundry/app-functions-sdk-go/internal.ApplicationVersion=$(APPVERSION)"

GIT_SHA=$(shell git rev-parse HEAD)

build:
	$(GO) build $(GOFLAGS) -o $(MICROSERVICE)

# NOTE: This is only used for local development. Jenkins CI does not use this make target
docker:
	docker build \
	    --build-arg http_proxy \
	    --build-arg https_proxy \
		-f Dockerfile \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/docker-fledge-export:$(GIT_SHA) \
		-t edgexfoundry/docker-fledge-export:master-dev \
		-t nexus3.edgexfoundry.org:10004/docker-fledge-export:master-dev \
		.

test:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]
	./bin/test-go-mod-tidy.sh
	./bin/test-attribution-txt.sh

clean:
	rm -f $(MICROSERVICE)