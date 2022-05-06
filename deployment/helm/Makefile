#
# Copyright (C) 2022 Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0
#

.PHONY: all checkov test

all: test

checkov:
	which checkov > /dev/null 2>&1 || pip3 install checkov

test: checkov
	checkov --framework helm -d .

