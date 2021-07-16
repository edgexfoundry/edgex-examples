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
* Download this [configuration.toml](./res/ibm-mqtt-export/configuration.toml) and edit the **[Writable.Pipeline.Functions.MQTTExport]** and **[Writable.InsecureSecrets.mqtt]** sections.
* Enter your Watson IoT Organization (6 character \<orgid\>), \<DeviceType\>, \<Device ID\> and \<Authentication token\>.
    ```yaml
        [Writable.Pipeline.Functions.MQTTExport]
          [Writable.Pipeline.Functions.MQTTExport.Parameters]
          # TODO - Change <orgid> placeholder
          BrokerAddress = "tcps://<orgid>.messaging.internetofthings.ibmcloud.com:1883"
          Topic = "iot-2/evt/status/fmt/json"
          # TODO - Change <orgid>, <devicetype> and <deviceid> placeholders
          ClientId = "d:<orgid>:<devicetype>:<deviceid>"
          QOS="0"
          AutoReconnect="false"
          KeepAlive = "" # Empty value means use default setting
          ConnectionTimeout = "" # Empty value means use default setting
          Retain="false"
          SkipVerify = "false"
          PersistOnError = "false"
          AuthMode = "usernamepassword"
          SecretPath = "/mqtt"
      ...
    
        [Writable.InsecureSecrets.mqtt]
        path = "mqtt"
          [Writable.InsecureSecrets.mqtt.Secrets]
          username = "use-token-auth"
          # TODO - Change <Authentication-Token> placeholder
          password = "<Authentication-Token>"
          cacert = ""
          clientcert = ""
          clientkey = ""
    ```

* After you have made the appropriate modifications for your Watson IoT account settings, save this file. Make sure the file is saved somewhere that can be easily volume mounted into your container. i.e. next to the compose file in the next step.
* Download the [no security compose file](https://github.com/edgexfoundry/edgex-compose/blob/ireland/docker-compose-no-secty.yml) and rename it to `docker-compose.yml`.
* Add the following snippet to the compose file:
    ```yaml
      app-ibm-mqtt-export:
        container_name: edgex-app-ibm-mqtt-export
        depends_on:
          - consul
          - data
        environment:
          EDGEX_PROFILE: ibm-mqtt-export
          SERVICE_HOST: edgex-app-ibm-mqtt-export
          CLIENTS_CORE_COMMAND_HOST: edgex-core-command
          CLIENTS_CORE_DATA_HOST: edgex-core-data
          CLIENTS_CORE_METADATA_HOST: edgex-core-metadata
          CLIENTS_SUPPORT_NOTIFICATIONS_HOST: edgex-support-notifications
          CLIENTS_SUPPORT_SCHEDULER_HOST: edgex-support-scheduler
          DATABASES_PRIMARY_HOST: edgex-redis
          EDGEX_SECURITY_SECRET_STORE: "false"
          MESSAGEQUEUE_HOST: edgex-redis
          REGISTRY_HOST: edgex-core-consul
          TRIGGER_EDGEXMESSAGEBUS_PUBLISHHOST_HOST: edgex-redis
          TRIGGER_EDGEXMESSAGEBUS_SUBSCRIBEHOST_HOST: edgex-redis
        hostname: edgex-app-ibm-mqtt-export
        image: edgexfoundry/app-service-configurable:2.0.0
        networks:
          edgex-network: { }
        ports:
          - 127.0.0.1:59780:59780/tcp
        read_only: true
        security_opt:
          - no-new-privileges:true
        user: 2002:2001
        volumes:
          - ./ibm-mqtt-export/:/res/ibm-mqtt-export/
    ```

* Note the volume mount of the `ibm-mqtt-export` profile at the bottom of the above snippet.
* Start the EdgeX services including your new `app-ibm-mqtt-export` Application Service
    ```bash
    docker-compose -p edgex up -d
    ```

- View the running containers.

  ```bash
  docker-compose -p edgex ps
  ```

  ```bash
             Name                          Command               State                                               Ports
  ---------------------------------------------------------------------------------------------------------------------------------------------------------------------
  edgex-app-ibm-mqtt-export     /app-service-configurable  ...   Up      48095/tcp, 127.0.0.1:59780->59780/tcp
  edgex-app-rules-engine        /app-service-configurable  ...   Up      48095/tcp, 127.0.0.1:59701->59701/tcp
  edgex-core-command            /core-command -cp=consul.h ...   Up      127.0.0.1:59882->59882/tcp
  edgex-core-consul             docker-entrypoint.sh agent ...   Up      8300/tcp, 8301/tcp, 8301/udp, 8302/tcp, 8302/udp, 127.0.0.1:8500->8500/tcp, 8600/tcp, 8600/udp
  edgex-core-data               /core-data -cp=consul.http ...   Up      127.0.0.1:5563->5563/tcp, 127.0.0.1:59880->59880/tcp
  edgex-core-metadata           /core-metadata -cp=consul. ...   Up      127.0.0.1:59881->59881/tcp
  edgex-device-rest             /device-rest --cp=consul:/ ...   Up      127.0.0.1:59986->59986/tcp
  edgex-device-virtual          /device-virtual --cp=consu ...   Up      127.0.0.1:59900->59900/tcp
  edgex-kuiper                  /usr/bin/docker-entrypoint ...   Up      20498/tcp, 127.0.0.1:59720->59720/tcp, 9081/tcp
  edgex-redis                   docker-entrypoint.sh redis ...   Up      127.0.0.1:6379->6379/tcp
  edgex-support-notifications   /support-notifications -cp ...   Up      127.0.0.1:59860->59860/tcp
  edgex-support-scheduler       /support-scheduler -cp=con ...   Up      127.0.0.1:59861->59861/tcp
  edgex-sys-mgmt-agent          /sys-mgmt-agent -cp=consul ...   Up      127.0.0.1:58890->58890/tcp
  ```

- Check the logs for `app-ibm-mqtt-export`

  ```bash
  docker-compose -p edgex logs app-ibm-mqtt-export
  ```

  ```bash
  Attaching to edgex-app-ibm-mqtt-export
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.5387139Z app=app-<profile> source=service.go:385 msg="Starting app-ibm-mqtt-export 2.0.0 "
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.5388474Z app=app-ibm-mqtt-export source=variables.go:352 msg="Variables override of '-p/-profile' by environment variable: EDGEX_PROFILE=ibm-mqtt-export"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.5523933Z app=app-ibm-mqtt-export source=config.go:359 msg="Loaded service configuration from /res/ibm-mqtt-export/configuration.toml"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.5541507Z app=app-ibm-mqtt-export source=variables.go:352 msg="Variables override of 'Registry.Host' by environment variable: REGISTRY_HOST=edgex-core-consul"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.5542439Z app=app-ibm-mqtt-export source=variables.go:352 msg="Variables override of 'Clients.core-metadata.Host' by environment variable: CLIENTS_CORE_METADATA_HOST=edgex-core-metadata"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.5542698Z app=app-ibm-mqtt-export source=variables.go:352 msg="Variables override of 'Service.Host' by environment variable: SERVICE_HOST=edgex-app-ibm-mqtt-export"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.5542852Z app=app-ibm-mqtt-export source=variables.go:352 msg="Variables override of 'Trigger.EdgexMessageBus.PublishHost.Host' by environment variable: TRIGGER_EDGEXMESSAGEBUS_PUBLISHHOST_HOST=edgex-redis"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.5542994Z app=app-ibm-mqtt-export source=variables.go:352 msg="Variables override of 'Trigger.EdgexMessageBus.SubscribeHost.Host' by environment variable: TRIGGER_EDGEXMESSAGEBUS_SUBSCRIBEHOST_HOST=edgex-redis"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.554616Z app=app-ibm-mqtt-export source=config.go:156 msg="Using Config Provider access token of length 0"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.5546915Z app=app-ibm-mqtt-export source=config.go:334 msg="Using Configuration provider (consul) from: http://edgex-core-consul:8500 with base path of edgex/appservices/2.0/app-ibm-mqtt-export"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7250966Z app=app-ibm-mqtt-export source=config.go:494 msg="Configuration has been pushed to into Configuration Provider (5 envVars overrides applied)"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7252173Z app=app-ibm-mqtt-export source=registry.go:57 msg="Using Registry access token of length 0"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7252479Z app=app-ibm-mqtt-export source=registry.go:73 msg="Using Registry (consul) from http://edgex-core-consul:8500"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.743961Z app=app-ibm-mqtt-export source=telemetry.go:78 msg="Starting CPU Usage Average loop"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7464146Z app=app-ibm-mqtt-export source=server.go:78 msg="Registering standard routes..."
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7465365Z app=app-ibm-mqtt-export source=configupdates.go:55 msg="Waiting for App Service configuration updates..."
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7467022Z app=app-ibm-mqtt-export source=service.go:437 msg="Service started in: 208.0569ms"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7467583Z app=app-ibm-mqtt-export source=main.go:38 msg="Loading Configurable Pipeline..."
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7469221Z app=app-ibm-mqtt-export source=triggerfactory.go:104 msg="EdgeX MessageBus trigger selected"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7469808Z app=app-ibm-mqtt-export source=messaging.go:64 msg="Initializing Message Bus Trigger for 'redis'"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7471749Z app=app-ibm-mqtt-export source=messaging.go:96 msg="Subscribing to topic(s): 'edgex/events/#' @ redis://edgex-redis:6379"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7472306Z app=app-ibm-mqtt-export source=messaging.go:105 msg="Publishing to topic: '' @ ://edgex-redis:0"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7472606Z app=app-ibm-mqtt-export source=service.go:190 msg="StoreAndForward disabled. Not running retry loop."
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.74728Z app=app-ibm-mqtt-export source=service.go:193 msg="app-mqtt-export has started"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7473174Z app=app-ibm-mqtt-export source=server.go:147 msg="Starting HTTP Web Server on address edgex-app-ibm-mqtt-export:59780"
  edgex-app-ibm-mqtt-export | level=INFO ts=2021-07-16T21:18:11.7473287Z app=app-ibm-mqtt-export source=messaging.go:117 msg="Waiting for messages from the MessageBus on the 'edgex/events/#' topic"
  ```

  
