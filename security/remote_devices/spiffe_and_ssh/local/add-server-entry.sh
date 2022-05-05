
docker exec -ti edgex-security-spire-config spire-server entry create -socketPath /tmp/edgex/secrets/spiffe/private/api.sock  -parentID spiffe://edgexfoundry.org/spire/agent/x509pop/cn/remote-agent -dns "edgex-device-virtual" -spiffeID  spiffe://edgexfoundry.org/service/device-virtual -selector "docker:label:com.docker.compose.service:device-virtual" 

#Entry ID         : f62bfec6-b19c-43ea-94b8-975f7e9a258e
#SPIFFE ID        : spiffe://edgexfoundry.org/service/device-virtual
#Parent ID        : spiffe://edgexfoundry.org/spire/agent/x509pop/cn/remote-agent
#Revision         : 0
#TTL              : default
#Selector         : docker:label:com.docker.compose.service:device-virtual
#DNS name         : edgex-device-virtual
