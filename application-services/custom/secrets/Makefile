.PHONY: build docker clean

GO=CGO_ENABLED=1 GO111MODULE=on go

build:
	go mod tidy
	$(GO) build -o app-service

clean:
	rm -f app-service
