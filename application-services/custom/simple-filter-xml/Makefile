.PHONY: build clean

build:
	go mod tidy
	go build -o app-service

docker:
	docker build \
	--build-arg http_proxy \
	--build-arg https_proxy \
	-f Dockerfile \
	-t edgexfoundry/docker-simple-filter-xml:dev \
	.

clean:
	rm -f app-service
