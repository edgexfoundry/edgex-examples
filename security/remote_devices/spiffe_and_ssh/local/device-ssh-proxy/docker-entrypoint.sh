#!/usr/bin/dumb-init /bin/sh

# Copyright (C) 2020-2022 Intel Corporation
# 
# SPDX-License-Identifier: Apache-2.0

set -ex


# Wait for agent CA creation

while test ! -f "/srv/spiffe/ca/public/agent-ca.crt"; do
    echo "Waiting for /srv/spiffe/ca/public/agent-ca.crt"
    sleep 1
done

# Pre-create remote agent certificate

if test ! -f "/srv/spiffe/remote-agent/agent.crt"; then
    openssl ecparam -genkey -name secp521r1 -noout -out "/srv/spiffe/remote-agent/agent.key"
    SAN="" openssl req -subj "/CN=remote-agent" -config "/usr/local/etc/openssl.conf" -key "/srv/spiffe/remote-agent/agent.key" -sha512 -new -out "/run/agent.req.$$"
    SAN="" openssl x509 -sha512 -extfile /usr/local/etc/openssl.conf -extensions agent_ext -CA "/srv/spiffe/ca/public/agent-ca.crt" -CAkey "/srv/spiffe/ca/private/agent-ca.key" -CAcreateserial -req -in "/run/agent.req.$$" -days 3650 -out "/srv/spiffe/remote-agent/agent.crt"
    rm -f "/run/agent.req.$$"
fi


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
    /srv/spiffe/remote-agent/agent.key $TUNNEL_HOST:/srv/spiffe/remote-agent/agent.key
  scp -p \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -P $TUNNEL_SSH_PORT \
    /srv/spiffe/remote-agent/agent.crt $TUNNEL_HOST:/srv/spiffe/remote-agent/agent.crt
  scp -p \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -P $TUNNEL_SSH_PORT \
    /tmp/edgex/secrets/spiffe/trust/bundle $TUNNEL_HOST:/tmp/edgex/secrets/spiffe/trust/bundle    
  ssh \
    -o StrictHostKeyChecking=no \
    -o UserKnownHostsFile=/dev/null \
    -p $TUNNEL_SSH_PORT \
    $TUNNEL_HOST -- \
    chown -Rh 2002:2001 /tmp/edgex/secrets/spiffe
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
    -R 0.0.0.0:$SECURITY_SPIRE_SERVER_PORT:$SECURITY_SPIRE_SERVER_HOST:$SECURITY_SPIRE_SERVER_PORT \
    -R 0.0.0.0:$SECRETSTORE_RUNTIMETOKENPROVIDER_PORT:$SECRETSTORE_RUNTIMETOKENPROVIDER_HOST:$SECRETSTORE_RUNTIMETOKENPROVIDER_PORT \
    -p $TUNNEL_SSH_PORT \
    $TUNNEL_HOST 
  sleep 1
done
