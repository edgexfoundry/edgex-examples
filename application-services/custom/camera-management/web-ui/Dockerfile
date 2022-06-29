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

# using alpine based image to address vulnerabilities found by Snyk scans
FROM node:18-alpine3.16

# add node_modules binaries to $PATH (eg. `ng`)
ENV PATH /app/node_modules/.bin:$PATH

# install the latest stable chromium for automated testing
RUN apk add --no-cache chromium
ENV CHROME_BIN=/usr/bin/chromium-browser

# update npm to latest
RUN npm install -g npm

ARG USER=1000
# delete the old 'node' user and create new one with same UID as local user's UID
RUN deluser --remove-home node; \
    adduser -u $USER -D -s /bin/sh node
# use as local user for file permission purposes when mounting
USER $USER:$USER

# set working directory (will be volume mounted)
WORKDIR /app
VOLUME ["/app"]

# NOTE: We do not copy any files or install any package.json deps
# as this images is meant to be used with '/app' bind-mounted to
# the local repository root.
