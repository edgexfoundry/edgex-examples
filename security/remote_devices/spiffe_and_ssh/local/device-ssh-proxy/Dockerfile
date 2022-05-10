# Copyright (C) 2020-2022 Intel Corporation
# 
# SPDX-License-Identifier: Apache-2.0

FROM alpine:latest

RUN apk add --no-cache --update dumb-init openssl openssh-client && rm -rf /var/cache/apk/*

COPY openssl.conf /usr/local/etc/
COPY docker-entrypoint.sh /usr/local/bin/

RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT [ "docker-entrypoint.sh" ]
