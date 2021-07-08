.PHONY: build clean

GO=CGO_ENABLED=1 GO111MODULE=on go

build:
	go mod tidy
	$(GO) build -o app-service

docker:
	docker build \
	--build-arg http_proxy \
	--build-arg https_proxy \
	-f Dockerfile \
	-t edgexfoundry/docker-simple-filter-xml:dev \
	.

clean:
	rm -f app-service
