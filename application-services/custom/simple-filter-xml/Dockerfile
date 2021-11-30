#
# Copyright (c) 2021 Intel Corporation
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

FROM golang:1.16-alpine3.14 AS builder

LABEL license='SPDX-License-Identifier: Apache-2.0' \
  copyright='Copyright (c) 2019: Intel'

# add git for go modules
RUN apk update && apk add --no-cache make git gcc libc-dev libsodium-dev zeromq-dev
WORKDIR /simple-filter-xml

COPY go.mod .

RUN go mod download

COPY . .
RUN apk info -a zeromq-dev

RUN make build

# Next image - Copy built Go binary into new workspace
FROM alpine

LABEL license='SPDX-License-Identifier: Apache-2.0' \
  copyright='Copyright (c) 2019: Intel'

RUN apk --no-cache add zeromq

# Turn off secure mode for examples. Not recommended for production
ENV EDGEX_SECURITY_SECRET_STORE=false

COPY --from=builder /simple-filter-xml/res /res
COPY --from=builder /simple-filter-xml/app-service /simple-filter-xml

CMD [ "/simple-filter-xml", "-cp=consul.http://edgex-core-consul:8500", "--registry","--confdir=/res"]