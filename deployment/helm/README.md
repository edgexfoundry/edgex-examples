# EdgeX Foundry on Kubernetes

A [Helm](https://helm.sh/) chart to easily deploy the EdgeX IoT project on Kubernetes.
Based on EdgeX [Minnesota](https://github.com/edgexfoundry/edgex-compose/tree/minnesota) version.

The helm chart is based on the secure and non-secure versions of the
[EdgeX Docker deployment scripts](https://github.com/edgexfoundry/edgex-compose/).
This helm chart can be used as a starting point for your own EdgeX deployment.


## Prerequisites

- Kubernetes cluster 1.24+
- [Helm](https://helm.sh/) 3.7.0+

Before starting, make sure you have curl and openssl installed locally.
These tools are needed to generate TLS assets and test the configuration.

These instructions also assume that your Kubernetes cluster has an installed
ingress controller and an installed load balancer.
The examples were built with a stock Ngnix ingress controller,
and a MetalLB load balancer (with an assigned IP of 192.168.122.200).
Adaptations may be needed if your cluster uses something different.

To install the Nginx ingress controller,
follow the instructions at
<https://kubernetes.github.io/ingress-nginx/deploy/#quick-start>

To install MetalLB,
follow the instructions at
<https://metallb.universe.tf/installation/#installation-with-helm>

The example also arbitrarily configures the ingress route to
respond to the hostname `edgex` 
(which should be passed in using the TLS Server Name Identification (SNI) feature).
This can also be customized as well using the `edgex.security.tlsHost` setting.


## Quick-start Installation

Install the EdgeX helm chart with a release name edgex-minnesota

You can install the helm chart 1 of 2 ways. 
1. Install the chart from by cloning the edgex-examples repository
2. Install the chart from the tar zipped asset from a tag

**If you want to clone the entire edgex-examples repository:**
```bash
$ git clone https://github.com/edgexfoundry/edgex-examples.git
$ cd edgex-examples
$ cd deployment/helm
$ kubectl create namespace edgex
$ helm install edgex-minnesota -n edgex .
```

**If you are only interested in installing the helm chart as a standalone:**
1. Navigate to the tagged asset by selecting the tag you desire to use
 
![image](https://user-images.githubusercontent.com/8902109/174185451-51273981-af57-42d7-ab8d-ae913a03e1b6.png)

2. Click on **Downloads** for the tag
 
![image](https://user-images.githubusercontent.com/8902109/174185618-eff4eb77-3185-46aa-b678-169f4c8730c0.png)

3. Either right click and copy the url or click to download the `edgex-examples-helm.tar.gz` file to save locally

![image](https://user-images.githubusercontent.com/8902109/174185727-7506e740-d51e-43c7-bf80-7effeb2402cb.png)

If you're using a Linux variant, you can use the copied URL to download the file from your terminal
```console
curl -o edgex-examples-helm.tar.gz <url for file from tag>
```

You will then need to unpack the tar zipped file:
```console
tar -xvf edgex-examples-helm.tar.gz
```
You will see output similar to:
![image](https://user-images.githubusercontent.com/8902109/174187588-910e9ee7-c8e2-4083-a7c3-d2614385c42c.png)

You can now change to the unzipped directory and install the helm chart.

```console
cd helm
kubectl create namespace edgex
helm install edgex-minnesota -n edgex .
```
## Uninstallation

```bash
helm uninstall edgex-minnesota -n edgex
```

## Test EdgeX

EdgeX on kubernetes using NodePort type to expose services by default. You can use ping command to test whether the EdgeX services start successfully.

The ping command format:
```bash
http://<ExternalIP>:<ExposedPort>/api/v3/ping

```
For example, the edgex-core-data ping command format:

```bash
curl http://localhost:59880/api/v3/ping
```


## Access EdgeX UI

With a modern browser, navigate to http://\<ExternalIP\>:30400.

Use details see [EdgeX UI doc](https://github.com/edgexfoundry/edgex-ui-go)


## Helm Chart User's Guide

This section will cover some of the features of the Helm chart,
so that it may be properly configured.


### Creating Docker Image Pull Secrets

The helm chart pulls a number of standard Docker images from Docker Hub.
This may cause Docker pull limits to be exceeded.
If you have a Docker hub account,
you may create a Docker image pull secrets that allows for more generous pull limits.

To create a Docker image pull secret, create the Kubernetes namespace for EdgeX,
and run the following command:

```shell
$ kubectl create namespace edgex
$ kubectl create secret docker-registry dockerhub --namespace <namespace> --docker-server=https://index.docker.io/v1/ --docker-username=<username> --docker-password=<password> --docker-email=<email-address>
```

On the `helm` command-line, specify the secrets to be used for pull secrets, for example: `--set imagePullSecrets="{ dockerhub }"`.


### Volumes

The default value of `edgex.storage.useHostPath` is `true`.
This setting causes EdgeX data volumes to be created under `/mnt` on the host file system,
which is a reasonable choice for a single-node Kubernetes deployment only.
For all other scenarios, set `edgex.storage.useHostPath` to `false`,
and configure the following settings:

```yaml
edgex:
    storage:
        useHostPath: false
        nonSharedVolumesClassName: "TBD"
        nonSharedVolumesAccessMode: "ReadWriteOnce"
        sharedVolumesClassName: "TBD"
        sharedVolumesAccessMode: "ReadWriteMany"
```

The EdgeX helm char has been tested with both `ReadWriteOnce` and `ReadWriteMany` for both types of volumes,
and has been compatibility tested with Rook-Ceph, OpenEBS, and the Rancher LocalPath provisioner.


### Enabling Security Features

The default value of `edgex.security.enabled` is `false`.  This may change in the future.

To enable the security, add the command flag in the helm install: `--set edgex.security.enabled=true`

Setting `edgex.security.enabled` to `true` during installation (recommended)
will enable microservice-level authentication
for EdgeX peer-to-peer communcation
as well authentication to Redis, Consul, and the MQTT broker, if used.

In lieu of a standalone API gateway used by the snap- and docker-based EdgeX deployments,
the security-enabled helm chart is coded against a standard Kubernetes NGINX-based ingress controller,
and ingress routes are configured to require client-side TLS authentication.

The device USB camera by default is not enabled to run in the deployment, to start it in a pod, 
add the following command flag in the helm install: `--set edgex.replicas.device.usbcamera=1`.

#### Configuring Ingress TLS

Run the following two scripts to generate key material and install it into the cluster.

The scripts create a server-side CA so that curl can trust the server,
a client-side CA so that Nginx can trust the client,
a server-side TLS certificate for the Nginx server to present to the client,
and a client-side TLS certificate for the curl client to present to the server.

```sh
./tlsgen.sh
./tlsinstall.sh
```

To install with security features enabled,
set the `edgex.security.enabled` flag to `true` during installation.

If EdgeX was previously installed on security disabled, uninstall the non-security one first.
The helm chart is not coded to allow for dynamic switching in and out of secure mode.

```sh
helm install edgex-minnesota --set edgex.security.enabled=true -n edgex .
```

#### Creating an Authentication JWT in Kubernetes

The following job is included at the root of the helm chart as `create-proxy-user-job.yaml`.
In this example, `edgexuser` should be replaced with a username of your own choosing,
and the `image` should be updated to an appropriate release image,
such as `edgexfoundry/security-proxy-setup:3.0.0`.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: security-proxy-setup
spec:
  template:
    spec:
      affinity:
        podAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchLabels:
                  org.edgexfoundry.service: edgex-security-secretstore-setup
              topologyKey: "kubernetes.io/hostname"
      imagePullSecrets:
      - name: dockerhub
      automountServiceAccountToken: false
      containers:
      - name: security-proxy-setup
        image: nexus3.edgexfoundry.org:10004/security-proxy-setup:3.0.0
        imagePullPolicy: Always
        command: ["/edgex-init/ready_to_run_wait_install.sh"]
        args: ["/edgex/secrets-config", "proxy", "adduser", "--user", "edgexuser", "--useRootToken"]
        envFrom:
        - configMapRef:
            name: edgex-common-variables
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
        volumeMounts:
          - mountPath: /edgex-init
            name: edgex-init
          - mountPath: /vault/config
            name: vault-config
          - mountPath: /tmp/edgex/secrets
            name: edgex-secrets
      restartPolicy: Never
      securityContext:
        runAsNonRoot: false
        runAsUser: 0
        runAsGroup: 0
      volumes:
        - name: edgex-init
          persistentVolumeClaim:
            claimName: edgex-init
        - name: vault-config
          persistentVolumeClaim:
            claimName: vault-config
        - name: edgex-secrets
          persistentVolumeClaim:
            claimName: edgex-secrets
```


To get a credential, first run the job and inspect the output:

```shell
 $ kubectl apply -f create-proxy-user.job.yaml 
 $ kubectl get pods
 $ kubectl logs security-proxy-setup-m78g
 {"username":"edgexuser","password":"dZ4SX...redacted...fOT"}
```

After the passwords is saved away, be sure to delete the completed job:

```shell
$ kubectl delete job security-proxy-setup
```

Next, create a script to obtain a secret store token and run it:

```shell
username=edgexuser
password="dZ4SX...redacted...fOT"
vault_token=$(curl -s --resolve edgex:443:192.168.122.200 --cacert server-ca.pem --cert client.pem --key client.key -X POST -H "Content-Type: application/json" "https://edgex/vault/v1/auth/userpass/login/${username}" -d "{\"password\":\"${password}\"}" | jq -r '.auth.client_token')
id_token=$(curl -s --resolve edgex:443:192.168.122.200 --cacert server-ca.pem --cert client.pem --key client.key -H "Authorization: Bearer ${vault_token}" "https://edgex/vault/v1/identity/oidc/token/${username}" | jq -r '.data.token')
echo "${id_token}"
```

The output will be a JWT of the form:
```
eyJ.redacted.redacted
```

### Sending a Test Request via Ingress

Finally, test with `curl`.
Note the use of the special options to enable SNI,
client-side certificates,
and server-side certificate validation.
Replace `<MetalLB IP>` below with the external IP
that is servicing the Kubernetes ingress controller.
Use the `$id_token` above to authenticate at the microservice layer.

```sh
curl -iv --resolve edgex:443:<MetalLB IP> --cacert server-ca.pem --cert client.pem --key client.key -H"Authorization: Bearer ${id_token}" "https://edgex/core-data/api/v3/version"
```

If everything was done correctly, the output will look like:

```text
... a bunch of diagnostics ...
* Connection #0 to host edgex left intact
{"apiVersion":"v3","version":"3.0.0-dev.137","serviceName":"core-data"}
```


### Configuring Port Bindings

The helm chart exposes EdgeX services as `ClusterIP` services by default.
This can be changed by the `expose.type.<...>` settings on a per-service basis.
Since the security flag enables the ingress rules for EdgeX microservices by default,
it is unlikely that directly exposing an EdgeX service on its own port will be necessary.

The feature `edgex.features.enableHostPort` is set to `false` by default.
If enabled, EdgeX services behave as the do on the snap- and docker-based implementations:
the map service ports to localhost or the external network interface,
as specified in the hostPortXXXBind settings:

```yaml
hostPortInternalBind: 127.0.0.1
hostPortExternalBind: 0.0.0.0
```


## Some articles about deploying edgex to kubernetes

- VMware China R&D Center
https://mp.weixin.qq.com/s/ECdEkc9QdkVScn4Lvl_JJA

## Feedback

If you find a bug or want to request a new feature, please open a [GitHub Issue](https://github.com/DaveZLB/edgex-helm/issues)


