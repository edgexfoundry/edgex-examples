# Copyright (C) 2022 Intel Corporation
# SPDX-License-Identifier: Apache-2.0

include ../make/base.mk

.PHONY: run download-models install patch-pipelines sniff-events stop

GITHUB_URL = https://github.com/intel/edge-video-analytics-microservice.git
VERSION ?= v0.7.0
SRC = src_$(VERSION)
MODELS_DIR = models
PIPELINES_LINK = pipelines

DOCKER_COMPOSE = docker-compose
COMPOSE_ARGS += -p edge-video-analytics

$(SRC):
	git clone --depth=1 --branch $(VERSION) $(GITHUB_URL) $@

$(PIPELINES_LINK): $(SRC)
	ln -sf $(SRC)/$@ $@

$(MODELS_DIR):
	mkdir -p $@

download-models: $(SRC) | $(MODELS_DIR)
	chmod +x $(SRC)/tools/model_downloader/model_downloader.sh $(SRC)/docker/run.sh
	$(SRC)/tools/model_downloader/model_downloader.sh --model-list $(SRC)/models_list/models.list.yml

# patch pipeline to accept rtsp video
patch-pipelines: $(PIPELINES_LINK)
	sed -Ei "s|\{auto_source\} ! decodebin|{auto_source} ! application/x-rtp,media=video ! decodebin|g" $(PIPELINES_LINK)/*/*/pipeline.json

install: $(SRC) $(PIPELINES_LINK) download-models patch-pipelines | $(MODELS_DIR)

run: $(PIPELINES_LINK) | $(MODELS_DIR)
	$(DOCKER_COMPOSE) $(COMPOSE_ARGS) up

stop:
	$(DOCKER_COMPOSE) $(COMPOSE_ARGS) down

sniff-events:
	@docker run --rm --net=host hivemq/mqtt-cli sub -t "incoming/data/edge-video-analytics/inference-event" | jq .
