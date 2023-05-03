.PHONY: build clean

APP_SERVICES=custom/*

.PHONY: build $(APP_SERVICES)

build: $(APP_SERVICES)

$(APP_SERVICES):
	cd $@; \
	make build

clean:
	rm -f custom/*/app-service