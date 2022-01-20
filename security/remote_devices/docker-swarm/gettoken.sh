docker exec -ti `docker ps  --filter name=consul | awk '{ print $1 }' | tail -n 1` cat /tmp/edgex/secrets/device-virtual/secrets-token.json
