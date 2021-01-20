# Kubernetes (K8s) Deployment

## Deployment and Orchestration Agnostic
EdgeX Foundry is a platform agnostic edge application.  It can run on any hardware (Intel or ARM).  It can run on any OS (Linux, Unix, macOS, Windows).  It has been constructed to be cloud agnostic (with demonstrated north side connectivity/data export to AWS, Azure, Google, and other cloud systems).  It is also an application that is **agnostic with regard to how it is deployed and orchestrated**.  The EdgeX community provides a set of Docker containers and [Docker Compose files](https://github.com/edgexfoundry/developer-scripts/tree/master/releases) to help demonstrate simple container based deployment and orchestration.  The community also provides [Ubuntu Snaps](https://snapcraft.io/) to demonstrate an alternate deployment and orchestration technique.  Members of our community use custom tools and shell scripts or OS dependent deployment package mechanisms to deploy EdgeX directly on to host platforms.

In this example, the community shows you how to deploy EdgeX via [Kubernetes](https://kubernetes.io/).  A simple Kubernetes (K8s) deployment file is provided to show how EdgeX Foundry can be deployed and orchestrated on to a single Kubernetes node.  

## K8s Installation and Setup

The installation and setup of K8s is beyond the scope of this example.  Consult the [K8s documentation](https://kubernetes.io/docs/setup/) for more details and education on how to install and configure K8s.

As the typical K8s production environment involves a control plane server, 1-to-many nodes (a.k.a. worker nodes) and many configuration options around these components, K8s environments can be a bit complicated to setup and resource intensive to run.  For the purpose of this example, [MicroK8s](https://microk8s.io/) is used.  MicroK8s is a "the simplest production-grade upstream K8s."  MicroK8s is built by the Canonical Kubernetes team.  It is a lightweight K8s environment that is easy to install and run - especially for single node examples as to be shown here. MicroK8s is available for Linux, Windows and macOS.

### MicroK8s Install and Setup

To install MicroK8s, see the [installation instructions](https://microk8s.io/) for the available platforms and the [MicroK8s documentation](https://microk8s.io/docs) for help on installation and setup.  In most cases, it is a one click and run or one line install command.

For example, on an Ubuntu OS machine, install the latest stable MicroK8s release with the `snap install` command:

```bash
sudo snap install microk8s --classic
```

### Turn on MicroK8s DNS

As the EdgeX services in K8s will need to talk to each other via network, turn on the MicroK8s DNS service.  This is done with a simple microk8s command (note user must have appropriate privileges to run this - see the the [docs](https://microk8s.io/docs) for assistance):

```bash
microk8s enable dns
```

*Note: you may also wish to install other K8s services such as the dashboard, registry, istio, etc.  See the MicroK8s documentation for other service tools and conveniences.*

The MicroK8s environment should provide a response indicating the DNS service is applied and working for your environment.

```bash
Enabling DNS
Applying manifest
serviceaccount/coredns created
configmap/coredns created
deployment.apps/coredns created
service/kube-dns created
clusterrole.rbac.authorization.k8s.io/coredns created
clusterrolebinding.rbac.authorization.k8s.io/coredns created
Restarting kubelet
DNS is enabled
```

## Apply the K8s Deployment File (Installing EdgeX)

Deploy EdgeX to your MicroK8s node with kubectl apply command with your deployment file as shown below.  In this example folder, you will find two different deployment files - one for the EdgeX Geneva release (ver 1.2.1) and the other for the Hanoi release (ver 1.3.0).  The Geneva release is shown in use below. 

```bash
microk8s.kubectl apply -f ./k8s-geneva-redis-no-secty.yml
```

*Note: Like all K8s environment, MicroK8s comes with the command line interface, kubectl, to issue K8s commands.  If you are deploying EdgeX in an alternate K8s environment, simply remove "microK8s." from the command above.*

K8s should respond, as shown below, with information about the services and deployment created.

```
configmap/common-variables created
service/edgex-core-consul created
service/edgex-redis created
service/edgex-support-notifications created
service/edgex-core-metadata created
service/edgex-core-data created
service/edgex-core-command created
service/edgex-support-scheduler created
service/edgex-app-service-configurable-rules created
service/edgex-kuiper created
service/edgex-device-virtual created
service/edgex-device-rest created
service/edgex-ui-go created
deployment.apps/edgex-core-consul created
deployment.apps/edgex-redis created
deployment.apps/edgex-support-notifications created
deployment.apps/edgex-core-metadata created
deployment.apps/edgex-core-data created
deployment.apps/edgex-core-command created
deployment.apps/edgex-support-scheduler created
deployment.apps/edgex-app-service-configurable-rules created
deployment.apps/edgex-kuiper created
deployment.apps/edgex-device-virtual created
deployment.apps/edgex-device-rest created
deployment.apps/edgex-ui-go created
```

If the deployment was successful, you can inspect the pods with the following command:

```bash
microk8s.kubectl get pods
```

This should result in an output that looks something like the following:

```
NAME                                                  READY   STATUS    RESTARTS   AGE
edgex-core-consul-5f7fbcccdd-zm9vn                    1/1     Running   0          33m
edgex-redis-564d5ff7c5-ppln8                          1/1     Running   0          33m
edgex-core-data-57d8c8f8b8-dxcx2                      1/1     Running   0          32m
edgex-app-service-configurable-rules-c4f9f76c-grgff   1/1     Running   0          32m
edgex-support-notifications-6cd46658b-qmldk           1/1     Running   0          32m
edgex-kuiper-58877bdbc4-5qxf9                         1/1     Running   0          32m
edgex-support-scheduler-65d87d8cb-rms5k               1/1     Running   0          32m
edgex-device-virtual-5bb8c5dc68-blznn                 1/1     Running   0          32m
edgex-core-command-cc6ff847c-tthx7                    1/1     Running   0          32m
edgex-device-rest-88b66679d-bf9rz                     1/1     Running   0          32m
edgex-ui-go-7fc7bb8749-8v72n                          1/1     Running   0          32m
edgex-core-metadata-847d48d68-pwgmv                   1/1     Running   0          32m
```

Likewise, to see the services, issue this command:

```bash
microk8s.kubectl get services
```

Which will result in a listing of the EdgeX services (and their IP address/ports).

```
NAME                                   TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                           AGE
kubernetes                             ClusterIP   10.152.183.1     <none>        443/TCP                           16h
edgex-core-consul                      NodePort    10.152.183.228   <none>        8500:30850/TCP,8400:30840/TCP     35m
edgex-redis                            ClusterIP   10.152.183.230   <none>        6379/TCP                          35m
edgex-support-notifications            NodePort    10.152.183.253   <none>        48060:30060/TCP                   35m
edgex-core-metadata                    NodePort    10.152.183.182   <none>        48081:30081/TCP                   35m
edgex-core-data                        NodePort    10.152.183.134   <none>        5563:31687/TCP,48080:30080/TCP    35m
edgex-core-command                     NodePort    10.152.183.209   <none>        48082:30082/TCP                   35m
edgex-support-scheduler                NodePort    10.152.183.214   <none>        48085:30085/TCP                   35m
edgex-app-service-configurable-rules   NodePort    10.152.183.184   <none>        48100:30100/TCP                   35m
edgex-kuiper                           NodePort    10.152.183.216   <none>        48075:30075/TCP,20498:30098/TCP   35m
edgex-device-virtual                   NodePort    10.152.183.85    <none>        49990:30090/TCP                   35m
edgex-device-rest                      NodePort    10.152.183.103   <none>        49986:30086/TCP                   34m
edgex-ui-go                            NodePort    10.152.183.211   <none>        4000:30040/TCP                    34m
```

## Stopping and Removing EdgeX

Should you need to remove the EdgeX deployment, perform a delete using kubectl with the same file.

```bash
microk8s.kubectl delete -f ./k8s-geneva-redis-no-secty.yml
```

### Clean up volumes

The deployment creates a number of volumes on your host.  These are, by default, located in your /consul and /data folders.  A script (clean-volumes.sh) has been provided with this example to clean out these volumes if you need to make a clean start or if you are switching between versions.

There is also a clean-cache.sh script to remove any files and cache in your local ~/.kube folder.  **Careful in that this script removes everything in your local .kube folder (whether it is part of EdgeX or not).** 

## Explore EdgeX in Kubernetes

Once EdgeX is running in Kubernetes, you can explore the APIs and services as you would explore any other EdgeX instance.  

### K8s special service port mapping

However, the port mapping to expose EdgeX services externally (outside of the Kubernetes environment) is different.  Kubernetes service node port range is between 30000-32767, so appropriate ports in this range have been selected to expose the EdgeX services externally.

| EdgeX Service | EdgeX Port | External Port for the service with K8s |
|---------------|------------|----------------------------------------|
| Consul | 8500 | 30850 |
| Consul | 8400 | 30840 |
| Notifications | 48060 | 30060 |
| Metadata | 48081 | 30081 |
| Data | 48080 | 30080 |
| Command | 48082 | 30082 |
| Scheduler | 48085 | 30085 |
| App Service configurable for the Rules Engine | 48100 | 30100 |
| Kuiper rules engine | 48075 | 30075 |
| Device Virtual | 49990 | 30090 |
| Device Rest | 49986 | 30086 |
| EdgeX UI | 4000 | 30040 |

So, for example, to hit the `ping` API on the core metadata service, you would need to issue a GET request against `http://{K8s host}:30081/api/v1/ping`.  Or, to get the count of events in core data, you would need to issue a GET request against `http://{K8s host}:30080/api/v1/event/count`.

### EdgeX Service Logs

If you need the logs for a particular service, use kubectl to get the pod name to the service you want to inspect.

``` bash
microk8s.kubectl get pods
NAME                                                  READY   STATUS    RESTARTS   AGE
edgex-core-consul-5f7fbcccdd-zm9vn                    1/1     Running   0          64m
edgex-redis-564d5ff7c5-ppln8                          1/1     Running   0          64m
edgex-core-data-57d8c8f8b8-dxcx2                      1/1     Running   0          64m
edgex-app-service-configurable-rules-c4f9f76c-grgff   1/1     Running   0          64m
edgex-support-notifications-6cd46658b-qmldk           1/1     Running   0          64m
edgex-kuiper-58877bdbc4-5qxf9                         1/1     Running   0          64m
edgex-support-scheduler-65d87d8cb-rms5k               1/1     Running   0          64m
edgex-device-virtual-5bb8c5dc68-blznn                 1/1     Running   0          64m
edgex-core-command-cc6ff847c-tthx7                    1/1     Running   0          64m
edgex-device-rest-88b66679d-bf9rz                     1/1     Running   0          64m
edgex-ui-go-7fc7bb8749-8v72n                          1/1     Running   0          64m
edgex-core-metadata-847d48d68-pwgmv                   1/1     Running   0          64m
```

So, for example, in order to get the logs from core data given the pod names listed above, you would issue a kubectl logs command like this (note the parameter to the logs request matches the identifier for core data in the pods listing above):

```bash
microk8s.kubectl logs edgex-core-data-57d8c8f8b8-dxcx2
level=INFO ts=2020-12-30T20:52:29.755220207Z app=edgex-core-data source=config.go:219 msg="Loaded configuration from /res/configuration.toml"
level=INFO ts=2020-12-30T20:52:29.770071719Z app=edgex-core-data source=environment.go:331 msg="Environment override of 'Databases.Primary.Host' by environment variable: Databases_Primary_Host=edgex-redis"
level=INFO ts=2020-12-30T20:52:29.774390978Z app=edgex-core-data source=environment.go:331 msg="Environment override of 'Clients.Metadata.Host' by environment variable: Clients_Metadata_Host=edgex-core-metadata"
...
```

## Device Services

In most circumstances, device services have to communicate with sensors, devices and other "things" via the protocol of the thing.  This is not always easily negotiated from inside a K8s environment.  Therefore, most device services will likely run outside of K8s and communicate into the other EdgeX services running inside of K8s.

Two exceptions to this are the Virtual Device Service and the REST Device Service.  These two device services have been made a part of the example K8s deployment example.  Neither of these device services communicate via anything other than TCP/IP and therefore can exist and be used conveniently in K8s.  Further, the Virtual Device Service provides a means to generate sample sensor data and explore it with the other EdgeX services - which is helpful for example sake.

If these device services are not needed, removed them from the example K8s deployment (YAML) files.

