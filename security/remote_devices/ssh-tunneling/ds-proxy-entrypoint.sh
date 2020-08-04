#!/usr/bin/dumb-init /bin/sh
#  ----------------------------------------------------------------------------------
#  Copyright (c) 2020 Intel Corporation
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

# Use dumb-init as PID 1 in order to reap zombie processes and forward system signals to 
# all processes in its session. This can alleviate the chance of leaking zombies, 
# thus more graceful termination of all sub-processes if any.

# ssh key generated and put under /root/.ssh/ in the container
rm -rf /root/.ssh && mkdir /root/.ssh \
&& cp -R /root/ssh/* /root/.ssh/ \
&& chmod -R 700 /root/.ssh/* \
&& chmod -R 600 /root/.ssh/id_rsa.* \
&& ls -al /root/.ssh/* \
&& cat /root/.ssh/id_rsa.pub

# ssh tunneling for both ways
sshTunneling="ssh -vv -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null \
  -N $TUNNEL_HOST \
  -L *:$SERVICE_PORT:$SERVICE_HOST:$SERVICE_PORT \
  -R 0.0.0.0:48080:edgex-core-data:48080 \
  -R 0.0.0.0:5563:edgex-core-data:5563 \
  -R 0.0.0.0:48081:edgex-core-metadata:48081 \
  -R 0.0.0.0:8500:edgex-core-consul:8500 \
  -p $TUNNEL_SSH_PORT && while true; do sleep 60; done"

echo "Executing $@"
"$@"

#sleep for some time to wait for creating authorized_keys on remote side
sleep 3

echo "Executing hook=$sshTunneling"
eval $sshTunneling
