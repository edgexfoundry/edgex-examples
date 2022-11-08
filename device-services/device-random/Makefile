.PHONY: build test clean prepare update docker

GO = CGO_ENABLED=0 GO111MODULE=on go
GOCGO = CGO_ENABLED=1 GO111MODULE=on go

MICROSERVICES=cmd/device-random

.PHONY: $(MICROSERVICES)

DOCKERS=device_random
.PHONY: $(DOCKERS)

VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)
GIT_SHA=$(shell git rev-parse HEAD)
GOFLAGS=-ldflags "-X github.com/edgexfoundry/device-random.Version=$(VERSION)"

build: $(MICROSERVICES)

tidy:
	go mod tidy

cmd/device-random:
	$(GOCGO) build $(GOFLAGS) -o $@ ./cmd

test:
	go mod tidy
	$(GOCGO) test ./... -coverprofile=coverage.out
	$(GOCGO) vet ./...
	gofmt -l .
	[ "`gofmt -l .`" = "" ]

clean:
	rm -f $(MICROSERVICES)

docker: $(DOCKERS)

device_random:
	docker build \
		--label "git_sha=$(GIT_SHA)" \
		-t edgexfoundry/device-random:$(GIT_SHA) \
		-t edgexfoundry/device-random:$(VERSION)-dev \
		.
