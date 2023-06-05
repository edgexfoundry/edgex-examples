#  * Copyright 2022 Intel Corporation.
#  *
#  * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
#  * in compliance with the License. You may obtain a copy of the License at
#  *
#  * http://www.apache.org/licenses/LICENSE-2.0
#  *
#  * Unless required by applicable law or agreed to in writing, software distributed under the License
#  * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
#  * or implied. See the License for the specific language governing permissions and limitations under
#  * the License.
#  *******************************************************************************/

networks:
  edgex-network:
    driver: bridge
services:
  sshd-remote:
    image: edgex-sshd-remote:latest
    build:
      context: sshd-remote
    container_name: edgex-sshd-remote
    hostname: edgex-sshd-remote
    ports:
    - "2223:22"
    read_only: true
    restart: always
    security_opt:
    - no-new-privileges:true
    networks:
      edgex-network:
        aliases:
        - edgex-core-consul
        - edgex-core-data
        - edgex-core-metadata
        - edgex-redis
        - edgex-security-spire-server
        - edgex-security-spiffe-token-provider
        - edgex-vault
    tmpfs:
    - /run
    volumes:
    - spire-remote-agent:/srv/spiffe/remote-agent:z
    - /tmp/edgex/secrets/spiffe:/tmp/edgex/secrets/spiffe:z
  remote-spire-agent:
    build:
      context: remote-spire-agent
    command: docker-entrypoint.sh
    container_name: edgex-remote-spire-agent
    depends_on:
    - sshd-remote
    hostname: edgex-security-spire-agent
    image: nexus3.edgexfoundry.org:10004/security-spire-agent:latest
    networks:
      edgex-network: {}
    pid: host
    privileged: true
    read_only: true
    restart: always
    security_opt:
    - no-new-privileges:true
    tmpfs:
    - /run
    user: root:root
    volumes:
    - spire-remote-agent:/srv/spiffe/remote-agent:z
    - /tmp/edgex/secrets/spiffe:/tmp/edgex/secrets/spiffe:z
    - /var/run/docker.sock:/var/run/docker.sock:rw
  device-virtual:
    command: /device-virtual -cp=consul.http://edgex-core-consul:8500 --registry --confdir=/res
    container_name: edgex-device-virtual
    depends_on:
    - remote-spire-agent
    environment:
      CLIENTS_CORE_COMMAND_HOST: edgex-core-command
      CLIENTS_CORE_DATA_HOST: edgex-core-data
      CLIENTS_CORE_METADATA_HOST: edgex-core-metadata
      EDGEX_SECURITY_SECRET_STORE: "true"
      MESSAGEQUEUE_HOST: edgex-redis
      REGISTRY_HOST: edgex-core-consul
      SECRETSTORE_HOST: edgex-vault
      SECRETSTORE_PORT: '8200'
      SECRETSTORE_RUNTIMETOKENPROVIDER_ENABLED: "true"
      SECRETSTORE_RUNTIMETOKENPROVIDER_ENDPOINTSOCKET: /tmp/edgex/secrets/spiffe/public/api.sock
      SECRETSTORE_RUNTIMETOKENPROVIDER_HOST: edgex-security-spiffe-token-provider
      SECRETSTORE_RUNTIMETOKENPROVIDER_PORT: 59841
      SECRETSTORE_RUNTIMETOKENPROVIDER_PROTOCOL: https
      SECRETSTORE_RUNTIMETOKENPROVIDER_REQUIREDSECRETS: redisdb
      SECRETSTORE_RUNTIMETOKENPROVIDER_TRUSTDOMAIN: edgexfoundry.org
      SERVICE_HOST: edgex-device-virtual
    hostname: edgex-device-virtual
    image: nexus3.edgexfoundry.org:10004/device-virtual:latest
    networks:
      edgex-network: {}
    ports:
    - 127.0.0.1:59900:59900/tcp
    read_only: true
    restart: always
    security_opt:
    - no-new-privileges:true
    user: 2002:2001
    volumes:
    - /tmp/edgex/secrets/device-virtual:/tmp/edgex/secrets/device-virtual:ro,z
    - /tmp/edgex/secrets/spiffe/public:/tmp/edgex/secrets/spiffe/public:ro,z
version: '3.7'
volumes:
  spire-remote-agent: {}
