.PHONY: build clean

GO=CGO_ENABLED=1 go

APPVERSION=2.0.0

# This pulls the version of the SDK from the go.mod file. If the SDK is the only required module,
# it must first remove the word 'required' so the offset of $2 is the same if there are multiple required modules
SDKVERSION=$(shell cat ./go.mod | grep 'github.com/edgexfoundry/app-functions-sdk-go v' | sed 's/require//g' | awk '{print $$2}')

MICROSERVICE=app-service
GOFLAGS=-ldflags "-X github.com/edgexfoundry/app-functions-sdk-go/internal.SDKVersion=$(SDKVERSION) -X github.com/edgexfoundry/app-functions-sdk-go/internal.ApplicationVersion=$(APPVERSION)"

GIT_SHA=$(shell git rev-parse HEAD)

build:
	go mod tidy
	$(GO) build $(GOFLAGS) -o $(MICROSERVICE)

clean:
	rm -f $(MICROSERVICE)

