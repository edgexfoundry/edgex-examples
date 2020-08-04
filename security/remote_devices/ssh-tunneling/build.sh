#!/bin/bash
#  ----------------------------------------------------------------------------------
#  Copyright 2020 Intel Corp.
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
#  SPDX-License-Identifier: Apache-2.0'
#  ----------------------------------------------------------------------------------

set -e

DEFAULT_SSH_KEY="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDKvsFf5HocBOBWXdVJKfQzkhf0K8lSLjZn9PX84VdhHyP8n1mzfpZywA4vsz8+A3OsGHAr2xpkyzOS0YkwD7nrI3q1x0A0+ANhQNOaKbnfQRepTAES3FPm5n0AbNVfgOre3RR2NLOt6M5m3mA/MERNer1fEp6BM96sdU0o3KjqwFGkPufoQrVkpz2691MZ6/ACDc+lk7uQrinsB4YxM7ctiLNl4I1A3TJgVv0jkJImUCHaThYj3XoaqUqUjQFTS7SlFfkXuk13EjNfRzqPwKFnVvGTUaYzaBV5S4wt5XCxhLfs497M2k5zmNx3HFY/GEyeoroCpjsiXkm+HcgdIYb7 root"

# override the new SSH public key from the environment variable if any; otherwise use default
SSH_PUBLIC_KEY=${SSH_PUBLIC_KEY-"$DEFAULT_SSH_KEY"}
export SSH_PUBLIC_KEY

echo "SSH_PUBLIC_KEY to be injected:" $SSH_PUBLIC_KEY

docker build --build-arg SSH_PUBLIC_KEY="$SSH_PUBLIC_KEY" -t eg_sshd .
