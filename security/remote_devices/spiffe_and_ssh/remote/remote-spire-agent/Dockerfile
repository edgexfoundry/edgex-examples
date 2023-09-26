#  ----------------------------------------------------------------------------------
#  Copyright 2023 Intel Corporation
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
# 
#  ----------------------------------------------------------------------------------

# Deployment image
FROM alpine:3.17

LABEL license='SPDX-License-Identifier: Apache-2.0' \
      copyright='Copyright (c) 2023 Intel Corporation'

RUN apk update && apk --no-cache --update add dumb-init openssl gcompat

COPY --from=ghcr.io/spiffe/spire-agent:1.6.3 /opt/spire/bin/spire-agent /usr/local/bin
COPY --from=ghcr.io/spiffe/spire-server:1.6.3 /opt/spire/bin/spire-server /usr/local/bin

COPY docker-entrypoint.sh /usr/local/bin/
COPY agent.conf /usr/local/etc/spire/agent.conf.tpl

ENTRYPOINT [ "/usr/bin/dumb-init" ]
CMD [ "--verbose", "docker-entrypoint.sh" ]
