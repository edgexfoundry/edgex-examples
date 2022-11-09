#
# Copyright (C) 2022 Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0
#

.PHONY: all checkov test

all: test

checkov-install:
	pip3 install checkov

trivy-install:
	sudo apt-get install wget apt-transport-https gnupg lsb-release
	wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | sudo apt-key add -
	echo "deb https://aquasecurity.github.io/trivy-repo/deb $$(lsb_release -sc) main" | sudo tee -a /etc/apt/sources.list.d/trivy.list
	sudo apt-get update
	sudo apt-get install trivy

checkov:
	which checkov > /dev/null 2>&1 || $(MAKE) checkov-install
	checkov --framework helm -d .

snyk:
	helm template . --output dir out
	snyk iac test out/
	@rm -fr out/

trivy:
	which trivy > /dev/null 2>&1 || $(MAKE) trivy-install
	trivy config -s CRITICAL,HIGH,MEDIUM .

test: 
	echo "Please run $(MAKE) checkov|snyk|trivy"

