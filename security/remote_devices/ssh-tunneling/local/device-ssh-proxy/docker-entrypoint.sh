#!/usr/bin/dumb-init /bin/sh

# Copyright (C) 2020-2022 Intel Corporation
# 
# SPDX-License-Identifier: Apache-2.0

set -ex

if [ "`stat -c '%u %g %a' /root/.ssh`" != "0 0 700" ]; then
  chown 0:0 /root/.ssh
  chmod 700 /root/.ssh
fi
if [ "`stat -c '%u %g %a' /root/.ssh/id_rsa`" != "0 0 600" ]; then
  chown 0:0 /root/.ssh/id_rsa
  chmod 600 /root/.ssh/id_rsa
fi
if [ "`stat -c '%u %g %a' /root/.ssh/id_rsa.pub`" != "0 0 600" ]; then
  chown 0:0 /root/.ssh/id_rsa.pub
  chmod 600 /root/.ssh/id_rsa.pub
fi

while true; do
  scp -p \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -P $TUNNEL_SSH_PORT \
    /tmp/edgex/secrets/device-virtual/secrets-token.json $TUNNEL_HOST:/tmp/edgex/secrets/device-virtual/secrets-token.json
  ssh \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -p $TUNNEL_SSH_PORT \
    $TUNNEL_HOST -- \
    chown -Rh 2002:2001 /tmp/edgex/secrets/device-virtual
  ssh -N \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -L *:$SERVICE_PORT:$SERVICE_HOST:$SERVICE_PORT \
    -R 0.0.0.0:$SECRETSTORE_PORT:$SECRETSTORE_HOST:$SECRETSTORE_PORT \
    -R 0.0.0.0:6379:$MESSAGEQUEUE_HOST:6379 \
    -R 0.0.0.0:8500:$REGISTRY_HOST:8500 \
    -R 0.0.0.0:5563:$CLIENTS_CORE_DATA_HOST:5563 \
    -R 0.0.0.0:59880:$CLIENTS_CORE_DATA_HOST:59880 \
    -R 0.0.0.0:59881:$CLIENTS_CORE_METADATA_HOST:59881 \
    -p $TUNNEL_SSH_PORT \
    $TUNNEL_HOST 
  sleep 1
done
