# EdgeX Foundry on Kubernetes

A [Helm](https://helm.sh/) chart to easily deploy the EdgeX IoT project on Kubernetes.
Based on EdgeX [Jakarta](https://github.com/edgexfoundry/edgex-compose/tree/kamakura) version.

## Prerequisites

- Kubernetes cluster 1.10+
- [Helm](https://helm.sh/) 3.7.0+

## Installation

Install the EdgeX helm chart with a release name edgex-kamakura

You can install the helm chart 1 of 2 ways. 
1. Install the chart from by cloning the edgex-examples repository
2. Install the chart from the tar zipped asset from a tag

**If you want to clone the entire edgex-examples repository:**
```bash
$ git clone https://github.com/edgexfoundry/edgex-examples.git
$ cd edgex-examples
$ cd deployment/helm
$ kubectl create namespace edgex
$ helm install edgex-kamakura -n edgex .
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
curl -o edgex-examples-heml.tar.gz <url for file from tag>
```

You will then need to unpack the tar zipped file:
```console
tar -xvf edgex-examples-heml.tar.gz
```
You will see output similar to:
![image](https://user-images.githubusercontent.com/8902109/174187588-910e9ee7-c8e2-4083-a7c3-d2614385c42c.png)

You can now change to the unzipped directory and install the helm chart.

```console
cd helm
kubectl create namespace edgex
helm install edgex-kamakura -n edgex .
```
## Uninstallation

```bash
helm uninstall edgex-kamakura -n edgex
```

## Test EdgeX

EdgeX on kubernetes using NodePort type to expose services by default. You can use ping command to test whether the EdgeX services start successfully.

The ping command format:
```bash
http://<ExternalIP>:<ExposedPort>/api/v2/ping

```
For example, the edgex-core-data ping command format:

```bash
curl http://localhost:59880/api/v2/ping
```


## Access EdgeX UI

With a modern browser, navigate to http://\<ExternalIP\>:30400.

Use details see [EdgeX UI doc](https://github.com/edgexfoundry/edgex-ui-go)

## Tips

- This project is based on [docker-compose-no-secty.yml](https://github.com/edgexfoundry/edgex-compose/blob/kamakura/docker-compose-no-secty.yml),
you can implement your customized version based on this.
- Since the EdgeX pods communicates with each other through the kubernetes service name, make sure the kubernetes DNS is enabled.
- Since other EdgeX services need to rely on consul to obtain configuration or register themselves to consul, other services cannot run normally until consul starts successfully.
- Unlike the docker-compose files for this release (which use a separate Docker volume container), the manifest files mount host based volumes as follows:

1、edgex-core-consul's /consul/config directory is mapped to the host's /mnt/edgex-consul-config directory.

2、edgex-core-consul's /consul/data directory is mapped to the host's /mnt/edgex-consul-data directory.

3、edgex-db's /data/db directory is mapped to the host's /mnt/edgex-db directory.

4、edgex-kuiper's /kuiper/data directory is mapped to the host's /mnt/edgex-kuiper-data directory.

- NodePort is enabled by default. According to default NodePort range(30000～32767), EdgeX NodePort mappings are as follows. 

| EdgeX Service Name          | Exposed Port 
| :-------------------------- | ------------- 
| edgex-core-data             | 59880         
| edgex-core-metadata         | 59881         
| edgex-core-command          | 59882         
| edgex-support-notifications | 59860         
| edgex-support-scheduler     | 59861         
| edgex-app-rules-engine      | 59701         
| edgex-kuiper                | 59720         
| edgex-device-rest           | 59986         
| edgex-device-virtual        | 59900         
| edgex-ui                    | 4000          
| edgex-sys-mgmt-agent        | 58890         
| edgex-redis                 | Not Exposed 
| edgex-core-consul           | 8500          

## Enabling security features

The helm chart uses an Kubernetes ingress controller in lieu of a Kong API gateway.
The ingress routes are configured to require client-side TLS authentication,
which replaces the Kong JWT authentication method.

### Prerequisites

Before starting, make sure you have curl and openssl installed locally.
These tools are need to generate TLS assets and test the configuration.

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
(which should be passed using the TLS Server Name Identification (SNI) feature).
This can also be customized as well using the `edgex.security.tlsHost` setting.

### Installation with security features

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

If necessary, uninstall the non-security one first.
The helm chart is not coded to allow for dynamic switching in and out of secure mode.

```sh
helm install edgex-kamakura --set edgex.security.enabled=true -n edgex .
```

Finally, test with `curl`.
Note the use of the special options to enable SNI,
client-side certificates,
and server-side certificate validation.
Replace `<MetalLB IP>` below with the external IP
that is servicing the Kubernetes ingress controller.

```sh
curl -iv --resolve edgex:443:<MetalLB IP> --cacert server-ca.pem --cert client.pem --key client.key "https://edgex/core-data/api/v2/ping"
```

```text
... a bunch of diagnostics ...
* Connection #0 to host edgex left intact
{"apiVersion":"v2","timestamp":"Wed Feb  2 18:32:57 UTC 2022"}
```


## Some articles about deploying edgex to kubernetes

- VMware China R&D Center
https://mp.weixin.qq.com/s/ECdEkc9QdkVScn4Lvl_JJA

## Feedback

If you find a bug or want to request a new feature, please open a [GitHub Issue](https://github.com/DaveZLB/edgex-helm/issues)


