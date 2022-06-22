# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# Helper functions for use in Makefiles to reduce boiler-plate logic
#

# Guard to prevent it from being included more than once
ifndef HELPERS_MK_DEFINED
HELPERS_MK_DEFINED := 1

# Remove/delete all docker images that match given TAG
#
# Usage: $(call remove_docker_images,TAG)
define remove_docker_images =
@echo -e "$(cyan)Deleting docker images for $1...$(clear)"
@if [ ! -z "$$(docker image ls -q $1)" ]; then \
  docker image ls -q $1 | xargs docker rmi -f; \
else \
	echo -e "$(yellow)No images found for $1$(clear)"; \
fi
endef # remove_docker_images

# Checks for existing docker images with the given TAG, and if they are missing,
# removes the SENTINEL_FILENAME file if present.
#
# Note 1: the + in front of $(MAKE) ensures make passes the jobserver
# Note 2: the body of the if statement does not have line continuations  " ... \ " on purpose. This
# 		  ensures the result is multiple commands (and also allows +$(MAKE) to work)
#
# Usage: $(call run_image_check,SENTINEL_FILENAME,TAG)
define run_image_check =
$(if $(and $(wildcard $1),$(if $(shell docker images -q $2),,missing)),\
@echo -e "$(red)$(bold)$2 image not found$(comma) flagging for re-creation.$(clear)"
rm -f $1
+$(MAKE) --always-make $1,\
)
$(NOOP)
endef # run_image_check

# Adds targets that build and clean Docker images.
#
# Usage: $(call docker_targets,$(IMAGE_NAMES),$(IMAGE_DEPENDENCIES))
#
# You'll get a $(name).imagebuilt and .PHONY $(name) target for each image name,
# and they can be used (or depended upon) to build the image and output a sentinel.
# You'll also get a .PHONY clean-$(name) target to remove the image and sentinel.
# Finally, if not previously defined, it'll add a $(name)_dockerfile variable,
# which it lists as a dependency of the build target and passes to `docker build`.
#
# If you need to use a different Dockerfile for you build,
# you can define $(name)_dockerfile yourself
# or, if its name ends in the suffix 'Dockerfile', just list it as a dependency.
# As a result, all of these are equivalent:
#
#   Inferred automatically:
#     $(call docker_targets,myimage,foo bar)
#   Listed explicitly:
#     $(call docker_targets,myimage,Dockerfile foo bar)
#   Defined ahead of time:
#     myimage_dockerfile = Dockerfile
#     $(call docker_targets,myimage,foo bar)
#
# If you're building from multiple Dockerfiles, use one of the alternatives.
# You can, of course, call the function multiple times to list different dependencies.
#
# The "magic" of this function comes from using define/foreach/eval, as described here:
# https://www.gnu.org/software/make/manual/html_node/Eval-Function.html#Eval-Function
docker_targets = $(foreach img,$(1),$(eval $(call docker_image_tmpl,$(img),$(2))))
DOCKER_BUILD ?= docker build
define docker_image_tmpl =
# This sets the Dockerfile variable and ensures there's only one.
$(1)_dockerfile ?= $$(or $(filter %Dockerfile,$(2)),Dockerfile)
$$(if $$(word 2,$$($(1)_dockerfile)), \
	$$(error "multiple Dockerfiles dependencies for $(1): $$($(1)_dockerfile)"))

# This target builds the image
$(1).imagebuilt: $$($(1)_dockerfile) $(2) $$(wildcard \.dockerignore) | $(1).image-check
	$$(DOCKER_BUILD) -f $$< -t $$(REPOSITORY:%=%/)$(1) .
	touch $(1).imagebuilt

$(1).image-check:
	$$(call run_image_check,$(1).imagebuilt,$$(REPOSITORY:%=%/)$(1))

# This target aligns the sentinel with the image, then builds it if needed.
$(1): $(1).imagebuilt

# This target removes the image and sentinel file.
clean-$(1):
	$$(call remove_docker_images,$$(REPOSITORY:%=%/)$(1))
	-rm -rf $(1).imagebuilt

show-deps-$(1): $$($(1)_dockerfile) $(2) $$(wildcard \.dockerignore)
	@echo -e "$(bold)$(green)$(1) Dependencies:$(clear)" $$(addprefix "\n - ",$$(^))

.PHONY: $(1) clean-$(1) $(1).image-check show-deps-$(1)
endef

# $(call json-arr,$(some-var)) converts words to a JSON array.
#
# Examples: $(call json-arr,a b c d)    --> ["a","b","c","d"]
#           $(call json-arr,a   b c d ) --> ["a","b","c","d"]
#           $(call json-arr,a)          --> ["a"]
#           $(call json-arr,)           --> []
#
# Note that it doesn't handle any sort of actual JSON processing/validation,
# so it's only useful on simple word lists:
#            # `make var='"a,b", c' some-target`
#            $(call json-arr,$(var)) --> [""a,b",","c"]
empty :=
space := $(empty) $(empty)
comma := ,
json-arr = [$(subst $(space),$(comma),$(foreach w,$(1),"$(strip $(w))"))]

# gofmt_check checks to ensure all go files are properly formatted, and if not, will
# print the offending files and exit with an error
define gofmt_check =
$(eval ufiles := $(shell gofmt -l .))
$(if $(ufiles),@printf "$(red)$(bold)ERROR: The following $(words $(ufiles)) Go file(s) are not formatted:$(addprefix \n⨯ $(CURDIR)/,$(ufiles))\n$(clear)"; exit 1)
endef

# Add $(NOOP) at the end of certain dynamic target to suppress messages such as
# "make[1]: Nothing to be done for 'target-name'."
NOOP = @:

# use this before a long-running command to trap ctrl-c and print out a notice that services are
# still running in the background:
# -$(tail_trap) tail -f log.txt
tail_trap = trap 'printf "\e[1A\e[K\n$(cyan)[NOTE]$(clear) Service(s) will continue running in the background\n"; exit 0' INT;

endif # HELPERS_MK_DEFINED
