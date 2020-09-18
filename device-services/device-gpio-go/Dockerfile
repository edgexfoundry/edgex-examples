#
# Copyright (c) 2020 Intel
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

ARG BASE=golang:1.15-alpine
FROM ${BASE} AS builder

ARG MAKE='make build'

WORKDIR $GOPATH/src/github.com/edgexfoundry/device-gpio-go

LABEL license='SPDX-License-Identifier: Apache-2.0' \
  copyright='Copyright (c) 2020: Intel'

RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories

# add git for go modules
RUN apk add --update --no-cache make git

COPY . .

RUN ${MAKE}

# Next image - Copy built Go binary into new workspace
FROM scratch

LABEL license='SPDX-License-Identifier: Apache-2.0' \
  copyright='Copyright (c) 2020: Intel'

ENV APP_PORT=49960
#expose command data port
EXPOSE $APP_PORT

WORKDIR /
COPY --from=builder /go/src/github.com/edgexfoundry/device-gpio-go/cmd/device-gpio-go /usr/local/bin/
COPY --from=builder /go/src/github.com/edgexfoundry/device-gpio-go/cmd/res/configuration.toml /res/
COPY --from=builder /go/src/github.com/edgexfoundry/device-gpio-go/cmd/res/device-gpio-profile.yaml /res/

ENTRYPOINT ["/usr/local/bin/device-gpio-go"]
CMD ["-cp=consul.http://172.17.0.1:8500", "--registry", "--confdir=/res"]
