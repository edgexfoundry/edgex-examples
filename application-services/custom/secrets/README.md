# Secrets example

This example demonstrates storing a secret to the secret store (Vault) and retrieving those secrets.

When running in secure mode, the secrets are stored and retrieved from Vault based on the `SecretStore` section of the configuration file.

Please refer to the [Application Functions SDK documentation](https://docs.edgexfoundry.org/latest/microservices/application/AdvancedTopics/#secrets)  for more details on storing and getting secrets using the SDK.

## Build and Run the app with Secure Edgex services

**Steps:**

1. From the `secrets` folder under app-services/custom run:

   ```console
   make build
   ```

2. Down load default EdgeX secure compose file 

   Copy to https://github.com/edgexfoundry/edgex-compose/blob/ireland/docker-compose.yml to local folder.

3. Modify compose file so this `Secrets`example service has access to Vault for it's SecretStore 

   Add the service's service-key, `app-secrets`, to the `secretstore-setup` service's `ADD_SECRETSTORE_TOKENS` environment variable as shown below:

   ```yaml
     secretstore-setup:
       container_name: edgex-security-secretstore-setup
       depends_on:
       - security-bootstrapper
       - vault
       environment:
         ADD_SECRETSTORE_TOKENS: 'app-secrets'
   ```

   This creates a Vault back SecretStore for our example and populates it with then known `redisdb` secret. In addition creating the SecretStore, service's can request known secrets and to be added to the API Gateway. See the [Configuring Add-on Service](https://docs.edgexfoundry.org/latest/security/Ch-Configuring-Add-On-Services/) security documentation for complete details.

4. Run the Edgex services in Docker using the modified compose file from above.

   Run the following command from the same folder the compose file resides.

   ```console
   docker-compose -p edgex up -d
   ```

   Now all the EdgeX service will be running. This can be verified by running the following command:

   ```console
   docker-compose -p edgex ps
   ```

   Which will output the following:

   ```console
   edgex-app-rules-engine             /edgex-init/ready_to_run_w ...   Up       48095/tcp, 127.0.0.1:59701->59701/tcp
   edgex-core-command                 /edgex-init/ready_to_run_w ...   Up       127.0.0.1:59882->59882/tcp
   edgex-core-consul                  /edgex-init/consul_wait_in ...   Up       8300/tcp, 8301/tcp, 8301/udp, 8302/tcp, 8302/udp, 127.0.0.1:8500->8500/tcp, 8600/tcp, 8600/udp
   edgex-core-data                    /edgex-init/ready_to_run_w ...   Up       127.0.0.1:5563->5563/tcp, 127.0.0.1:59880->59880/tcp
   edgex-core-metadata                /edgex-init/ready_to_run_w ...   Up       127.0.0.1:59881->59881/tcp
   edgex-device-rest                  /edgex-init/ready_to_run_w ...   Up       127.0.0.1:59986->59986/tcp
   edgex-device-virtual               /edgex-init/ready_to_run_w ...   Up       127.0.0.1:59900->59900/tcp
   edgex-kong                         /edgex-init/kong_wait_inst ...   Up       0.0.0.0:8000->8000/tcp,:::8000->8000/tcp, 8001/tcp, 127.0.0.1:8100->8100/tcp, 0.0.0.0:8443->8443/tcp,:::8443->8443/tcp, 8444/tcp
   edgex-kong-db                      /edgex-init/postgres_wait_ ...   Up       127.0.0.1:5432->5432/tcp
   edgex-kuiper                       /edgex-init/kuiper_wait_in ...   Up       20498/tcp, 127.0.0.1:59720->59720/tcp, 9081/tcp
   edgex-redis                        /edgex-init/redis_wait_ins ...   Up       127.0.0.1:6379->6379/tcp
   edgex-security-bootstrapper        /entrypoint.sh gate              Up
   edgex-security-proxy-setup         /edgex-init/proxy_setup_wa ...   Exit 0
   edgex-security-secretstore-setup   entrypoint.sh                    Up
   edgex-support-notifications        /edgex-init/ready_to_run_w ...   Up       127.0.0.1:59860->59860/tcp
   edgex-support-scheduler            /edgex-init/ready_to_run_w ...   Up       127.0.0.1:59861->59861/tcp
   edgex-sys-mgmt-agent               /edgex-init/ready_to_run_w ...   Up       127.0.0.1:58890->58890/tcp
   edgex-vault                        /edgex-init/vault_wait_ins ...   Up       127.0.0.1:8200->8200/tcp
   ```

5. Run the `secrets` example service as root

   The service must run as root so that it can access it's SecretStore token. Also the environment variable `EDGEX_SECURITY_SECRET_STORE` must not be set to `false`. Either not set at all or set to `true`

   ```console
   sudo EDGEX_SECURITY_SECRET_STORE=true ./app-service
   ```

6. Follow the instructions in [Storing and Getting Secrets](#storing-and-getting-secrets) in order to test storing and retrieving secrets from the secret store.

## Storing and Getting Secrets

These tests use a collection of Postman requests, in *SecretsExample.postman_collection.json*, to store and retrieve secrets from the secret store.

1. Import the collection *SecretsExample.postman_collection.json* into Postman.

2. Execute the `Store Secrets` request in the Postman collection to push secrets to Vault. This is going through the App Service REST API, not directly to Vault. As such, the secret is exclusively for that app service instance.

3. Execute the `Get Secrets with App Service HTTP` request in the Postman collection.

   This request triggers an EdgeX event to the application service which causes execution of the pipeline function that calls the GetSecrets API.  As a result, the app service will get the exclusive secrets that were just pushed.

4. View the service's logs to verify that the secrets were retrieved. We'll view the secrets in the application's console (in production, NEVER log your application's secrets. This is done in the example service to demonstrate the functionality).

   ```console
   level=INFO ts=2021-07-19T21:29:44.2351952Z app=app-secrets source=getsecrets.go:52 msg="--- Get secrets at location /mqtt, keys: []  ---"
   level=INFO ts=2021-07-19T21:29:44.2398679Z app=app-secrets source=getsecrets.go:59 msg="key:username, value:app-user"
   level=INFO ts=2021-07-19T21:29:44.2399432Z app=app-secrets source=getsecrets.go:59 msg="key:password, value:SuperDuperSecretPassword"
   level=INFO ts=2021-07-19T21:29:44.2399976Z app=app-secrets source=getsecrets.go:52 msg="--- Get secrets at location /mqtt, keys: [password]  ---"
   level=INFO ts=2021-07-19T21:29:44.2400476Z app=app-secrets source=getsecrets.go:59 msg="key:password, value:SuperDuperSecretPassword"
   ```

   
