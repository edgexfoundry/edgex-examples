.PHONY: build clean

GO=CGO_ENABLED=1 GO111MODULE=on go

APP_SERVICES=custom/*

.PHONY: build $(APP_SERVICES)

build: $(APP_SERVICES)

$(APP_SERVICES):
	cd $@; \
	make build

clean:
	rm -f custom/*/app-service