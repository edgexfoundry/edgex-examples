# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

include ../make/base.mk

REPOSITORY ?= camera-management
SERVICE_NAME := web-ui

TEST_PORT ?= 9876
port ?= 4200

.PHONY: build image lint debug docs e2e test prod \
				nuke clean clean-images audit audit-fix default view-coverage \
				force-install install view-docs $(SERVICE_NAME).image-check \
				serve clear-cache

# This allows the user to pass in a variable amount of arguments to the command
# by setting the `args` variable: `make install args="--force"`, etc.
args ?=

# This allows the user to pass in additional arguments to `docker run ...`
# by setting the `run_args` variable.
run_args ?=

# This allows the user to pass in additional arguments to `docker build ...`
# by setting the `build_args` variable.
build_args ?=

# number of lines to tail when using `make tail` or `make logs`
n ?= 100

SRC_FILES = $(shell find src -type f)

# snippet to call docker run with provided arguments.
# usage: $(call DOCKER_RUN_EX,--arguments --here) command here
DOCKER_RUN_EX = docker run \
		-it \
		--rm \
		$1 \
		$(run_args) \
		-v "$(CURDIR):/app" \
		$(REPOSITORY)/$(SERVICE_NAME)

# shorthand to call DOCKER_RUN_EX without arguments
# usage: $(DOCKER_RUN) command here
DOCKER_RUN = $(call DOCKER_RUN_EX,)

# default docker-compose command and arguments
DOCKER_COMPOSE = docker-compose -f docker-compose.yml -p camera-management-web-ui

OPEN = xdg-open

# The default target is dynamically created based on the way make was called. By default it will
# execute a `make prod`, unless `ng`, `npm`, `exec`, or `gen` variables are defined.
# To run `npm` commands within the confines of the volume mounted dev container, set the `npm` variable
# 	to your argument list like so: `make npm="audit fix"`.
# To run ng commands within the confines of the volume mounted dev container, set the `ng` variable
# 	to your argument list like so: `make ng="generate component Test"`.
# To run any arbitrary command within the confines of the volume mounted dev container, set the `exec`
#		variable to your command and argument list like so: `make exec="node --version"`
# To generate angular components, services, fixtures, etc, set the `gen`
#		variable to the argument list of `ng generate ...` args like so: `make gen="component Test"`
default: image
ifdef npm
	$(DOCKER_RUN) npm $(npm)
else ifdef ng
	$(DOCKER_RUN) ng $(ng)
else ifdef exec
	$(DOCKER_RUN) /bin/ash -c '$(exec)'
else ifdef gen
	$(DOCKER_RUN) ng generate $(gen)
else
	$(MAKE) prod
endif

build: image

image: $(SERVICE_NAME).imagebuilt

$(SERVICE_NAME).imagebuilt: Dockerfile | $(SERVICE_NAME).image-check
	docker build \
		-f Dockerfile \
		--build-arg USER=$(shell id -u) \
		$(build_args) \
		-t $(REPOSITORY)/$(SERVICE_NAME) \
		.
	touch $@

# This target checks if images are present, if not deletes the .imagebuilt sentinel file
$(SERVICE_NAME).image-check:
	$(call run_image_check,$(SERVICE_NAME).imagebuilt,$(REPOSITORY)/$(SERVICE_NAME))

# this target downloads and installs new node modules only when needed and is used as a pre-req
# for other targets to make sure the dependencies are up-to-date
node_modules: package.json package-lock.json | image
	$(DOCKER_RUN) npm install $(args) || $(MAKE) force-install args=$(args)
	touch $@

prod: dist/index.html

dist/index.html: $(SRC_FILES) node_modules | image
	$(DOCKER_RUN) ng build --configuration production

test: image
	$(call DOCKER_RUN_EX,-p "$(TEST_PORT):$(TEST_PORT)" --shm-size=1g) ng test --code-coverage $(args)

view-coverage:
	@( $(OPEN) coverage/report-html/index.html & ) || printf "\n\e[31mPlease run 'make test' first.\nIf still not working, your OS may be incompatible with '$(OPEN)'.\e[0m\n"
	@sleep 0.5; echo

e2e: image
	$(call DOCKER_RUN_EX, --shm-size=1g) ng e2e $(args)

lint: image
	$(DOCKER_RUN) ng lint $(args)

docs: image
	$(DOCKER_RUN) typedoc --options typedoc.json src

view-docs:
	@( $(OPEN) docs/index.html & ) || printf "\n\e[31mPlease run 'make docs' first.\nIf still not working, your OS may be incompatible with '$(OPEN)'.\e[0m\n"
	@sleep 0.5; echo

install: node_modules

force-install:
	$(DOCKER_RUN) npm install --force $(args)
	touch node_modules

audit: image
	$(DOCKER_RUN) npm audit $(args)

audit-fix: image
	$(DOCKER_RUN) npm audit fix $(args)

update-angular: image
	$(DOCKER_RUN) ng update @angular/core @angular/cli rxjs || $(MAKE) force-install
	$(MAKE) force-install

# cleanup everything including npm packages
nuke: clean clean-images
	rm -rf node_modules

# cleanup cache + artifacts + images
clean:
	rm -rf dist coverage node_modules/.cache *.imagebuilt

# delete docker images and image build layers
clean-images:
	$(call remove_docker_images,$(REPOSITORY)/$(SERVICE_NAME))


# Launches an instance of the dev container with the current code mounted
# at /app. Useful for running various npm and ng commands.
debug: image
	$(call DOCKER_RUN_EX, \
		-e DISPLAY="unix${DISPLAY}" \
		-v "/tmp/.X11-unix:/tmp/.X11-unix" \
		-v "/var/run/dbus:/var/run/dbus" \
		--privileged \
		--shm-size=1g) \
		/bin/sh $(args)

serve: image node_modules clear-cache
	$(call DOCKER_RUN_EX,-p "$(port):$(port)") ng serve --host=0.0.0.0 --port=$(port)

clear-cache:
	rm -rf .angular/cache
