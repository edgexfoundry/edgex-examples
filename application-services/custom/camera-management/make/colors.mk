# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

# This makefile contains ANSI color definitions
#

# Guard to prevent it from being included more than once
ifndef COLOR_MK_DEFINED
COLOR_MK_DEFINED := 1

# clear resets color and style back to defaults
clear := \e[0m

# no_color clears the foreground color
no_color := \e[39m

# normal colors
red := \e[31m
green := \e[32m
yellow := \e[33m
blue := \e[34m
magenta := \e[35m
cyan := \e[36m
white := \e[37m

# normal clears dim, bold, and underline
normal := \e[21m\e[22m\e[24m

# styles
dim := \e[2m
bold := \e[1m
uline := \e[4m

endif # COLOR_MK_DEFINED
