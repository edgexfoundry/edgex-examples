# app-service-configurable-ibm
EdgeX **app-service-configurable** Profile for IBM Watson IoT Platform

The [configuration.toml](/res/ibm-mqtt-export/configuration.toml) provided in this repository defines the [EdgeX app-service-configurable](https://github.com/edgexfoundry/app-service-configurable) **profile** required to send MQTT data to [IBM Watson IoT Platform](https://cloud.ibm.com/catalog/services/internet-of-things-platform#about)

## Prerequisites

This tutorial can be completed using an IBM Cloud Lite account.

* Create an [IBM Cloud account](https://ibm.biz/BdzgKN)
* Log into [IBM Cloud](https://cloud.ibm.com/login)
* Create a [Watson IoT Service instance](https://cloud.ibm.com/catalog/services/internet-of-things-platform)

## Install EdgeX

You can build the EdgeX Foundry services using the open source code on [GitHub](https://github.com/edgexfoundry), but more often than now you just need to get these services running so that you can connect your own services to them. To support that, the project publishes Docker images based on the latest stable release of the open source code, as well as docker-compose.yml files that will run all the necessary services together on your development machine. Learn more at https://www.edgexfoundry.org/get-started/

## Send MQTT Data to IBM Watson IoT platform

If you have an Edge device generating MQTT data, you might want to send the IoT data to the cloud for real-time alerts, time series database storage, analytics and modeling. IBM Cloud provides an IoT ingestion service, Watson AI services, Cloud Object Storage and  Watson Studio data science portal that can help developers manage their IoT data and find insights.

EdgeX provides an [App-Service-Configurable](https://github.com/edgexfoundry/app-service-configurable/blob/master/README.md) service as an easy way to get started with processing data flowing through EdgeX. This service leverages the App Functions SDK and provides a way for developers to use configuration instead of having to compile standalone services to utilize built in functions in the SDK. For a full list of supported/built-in functions view the README located in the App Functions SDK repository.

### Watson IoT Configuration Example:

* Create a [Watson IoT Service instance](https://cloud.ibm.com/catalog/services/internet-of-things-platform)
* Create a Devicetype, Device ID and a secure Authentication token.
* In this example, we use token authentication instead of TLS.
  * Settings \> Security \> Connection Security \> Default Connection Security Level \> TLS Optional
* Download this [configuration.toml](/res/ibm-mqtt-export/configuration.toml) and edit the **[Writable.Pipeline.Functions.MQTTSend.Addressable]** section.
* Enter your Watson IoT Organization (6 character \<orgid\>), \<DeviceType\>, \<Device ID\> and \<Authentication token\>.
```yaml
[Writable.Pipeline.Functions.MQTTSend.Addressable]
  Address=   "<orgid>.messaging.internetofthings.ibmcloud.com"
  Port=      1883
  Protocol=  "tcp"
  Publisher= "d:<orgid>:<devicetype>:<deviceid>"
  User=      "use-token-auth"
  Password=  "<Authentication-Token>"
  Topic=     "iot-2/evt/status/fmt/json"
```  

* After you have made the appropriate modifications for your Watson IoT account settings, save this file. Make sure the file is saved somewhere that can be easily volume mounted into your container.
* Below is a snippet of how to configure the app-service-configurable service in an [edgex docker-compose file](https://github.com/edgexfoundry/developer-scripts/tree/master/releases):
```yaml
app-service-configurable:
    image: nexus3.edgexfoundry.org:10004/docker-app-service-configurable:latest
    command: --profile=docker
    ports:
      - "48095:48095"
    container_name: edgex-app-service-configurable
    hostname: edgex-app-service-configurable
    networks:
      edgex-network:
        aliases:
          - edgex-app-service-configurable
    depends_on:
      - data
      - command
  ```

You'll need to make a couple changes to the docker-compose.yaml file:

1. Mount the directory that has the `configuration.toml` file that you've modified above into the `/res` directory:
```yaml
app-service-configurable:
  ...
  volumes:
    - ./res/ibm-mqtt-export/:/res/ibm-mqtt-export/
  ...
```
2. Edit the `docker-compose.yaml` in the `app-service-configurable:` section and change the `command:` line to incorporate the IBM MQTT profile. This needs to match the name of the folder that you've mounted into the `/res` directory above.
```yaml
app-service-configurable:
  ...
  command: --profile=ibm-mqtt-export
  ...
```

## Example docker-compose.yaml snippet
```yaml
app-service-configurable:
    image: nexus3.edgexfoundry.org:10004/docker-app-service-configurable:latest
    command: --profile=ibm-mqtt-export
    volumes:
      - ./res/ibm-mqtt-export/:/res/ibm-mqtt-export/
    ports:
      - "48095:48095"
    container_name: edgex-app-service-configurable
    hostname: edgex-app-service-configurable
    networks:
      edgex-network:
        aliases:
          - edgex-app-service-configurable
    depends_on:
      - data
      - command
  ```

* Finally, remember to restart the app-service-configurable container to pick up the changes
```
docker-compose restart app-service-configurable
```
