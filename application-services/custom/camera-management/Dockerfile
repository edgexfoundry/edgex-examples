#
# Copyright (c) 2022 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

#build stage
ARG BASE=golang:1.17-alpine3.15
FROM ${BASE} AS builder

ARG ALPINE_PKG_BASE="make git gcc libc-dev libsodium-dev zeromq-dev"
ARG ALPINE_PKG_EXTRA=""

LABEL license='SPDX-License-Identifier: Apache-2.0' \
    copyright='Copyright (c) 2022: Intel'
#disabled for now due to issues with the nl.alpinelinux.org mirror, Edgex is working on a fix for this.
#RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories
RUN apk add --update --no-cache ${ALPINE_PKG_BASE} ${ALPINE_PKG_EXTRA}
WORKDIR /app

COPY go.mod ./
RUN go mod download all

COPY . .

ARG MAKE="make build-app"
RUN $MAKE

#final stage
FROM alpine:3.15
LABEL license='SPDX-License-Identifier: Apache-2.0' \
  copyright='Copyright (c) 2022: Intel'
LABEL Name=app-camera-management Version=${VERSION}

# dumb-init is required as security-bootstrapper uses it in the entrypoint script
RUN apk add --update --no-cache ca-certificates zeromq dumb-init

COPY --from=builder /app/LICENSE /LICENSE
COPY --from=builder /app/res/ /res/
COPY --from=builder /app/app-camera-management /app-camera-management
COPY --from=builder /app/web-ui/dist /web-ui/dist

EXPOSE 59750

ENTRYPOINT ["/app-camera-management"]
CMD ["-cp=consul.http://edgex-core-consul:8500", "--registry", "--confdir=/res"]
