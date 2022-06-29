# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# This makefile contains the base configuration settings all Makefiles will share
# To utilize it, add the following to the top of your Makefile (substitute actual path):
# 	include ../make/base.mk
#

# Guard to prevent it from being included more than once
ifndef BASE_MK_DEFINED
BASE_MK_DEFINED := 1

# The directory this specific makefile resides
make_base_dir := $(dir $(lastword $(MAKEFILE_LIST)))

SHELL := /bin/bash

# If the users passes TRACE=1, enable shell command tracing
TRACE ?= 0
ifeq ($(TRACE),1)
.SHELLFLAGS = -xc
endif

.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables --no-builtin-rules --no-print-directory

include $(make_base_dir)/colors.mk

include $(make_base_dir)/helpers.mk

endif # BASE_MK_DEFINED
