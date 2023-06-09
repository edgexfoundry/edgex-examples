# app-service-configurable-ibm
EdgeX **app-service-configurable** Profile for IBM Watson IoT Platform

The [configuration.yaml](/res/ibm-mqtt-export/configuration.yaml) provided in this repository defines the [EdgeX app-service-configurable](https://github.com/edgexfoundry/app-service-configurable/tree/minnesota) **profile** required to send MQTT data to [IBM Watson IoT Platform](https://cloud.ibm.com/catalog/services/internet-of-things-platform#about)

## Prerequisites

This tutorial can be completed using an IBM Cloud Lite account.

* Create an [IBM Cloud account](https://ibm.biz/BdzgKN)
* Log into [IBM Cloud](https://cloud.ibm.com/login)
* Create a [Watson IoT Service instance](https://cloud.ibm.com/catalog/services/internet-of-things-platform)

## Install EdgeX

You can build the EdgeX Foundry services using the open source code on [GitHub](https://github.com/edgexfoundry), but more often than now you just need to get these services running so that you can connect your own services to them. To support that, the project publishes Docker images based on the latest stable release of the open source code, as well as docker-compose.yml files that will run all the necessary services together on your development machine. Learn more at https://www.edgexfoundry.org/get-started/

## Send MQTT Data to IBM Watson IoT platform

If you have an Edge device generating MQTT data, you might want to send the IoT data to the cloud for real-time alerts, time series database storage, analytics and modeling. IBM Cloud provides an IoT ingestion service, Watson AI services, Cloud Object Storage and  Watson Studio data science portal that can help developers manage their IoT data and find insights.

EdgeX provides an [App-Service-Configurable](https://github.com/edgexfoundry/app-service-configurable/blob/minnesota/README.md) service as an easy way to get started with processing data flowing through EdgeX. This service leverages the App Functions SDK and provides a way for developers to use configuration instead of having to compile standalone services to utilize built in functions in the SDK. For a full list of supported/built-in functions view the README located in the App Functions SDK repository.

### Watson IoT Configuration Example:

* Create a [Watson IoT Service instance](https://cloud.ibm.com/catalog/services/internet-of-things-platform)
* Create a Devicetype, Device ID and a secure Authentication token.
* In this example, we use token authentication instead of TLS.
  * Settings \> Security \> Connection Security \> Default Connection Security Level \> TLS Optional
* Download this [configuration.yaml](./res/ibm-mqtt-export/configuration.yaml) and edit the **[Writable.Pipeline.Functions.MQTTExport]** and **[Writable.InsecureSecrets.mqtt]** sections.
  * Enter your Watson IoT Organization (6 character \<orgid\>), \<DeviceType\>, \<Device ID\> and \<Authentication token\>.
      ```yaml
        MQTTExport:
          Parameters:
            # TODO - Change <orgid> placeholder
            BrokerAddress: "tcps://<orgid>.messaging.internetofthings.ibmcloud.com:1883"
            # TODO - Change <orgid>, <devicetype> and <deviceid> placeholders
            Topic: iot-2/evt/status/fmt/json
            ClientId: "d:<orgid>:<devicetype>:<deviceid>"
            QOS: "0"
            AutoReconnect: "false"
            KeepAlive: ""
            ConnectionTimeout: ""
            Retain: "false"
            SkipVerify: "false"
            PersistOnError: "false"
            AuthMode: usernamepassword
            SecretPath: mqtt
       
        InsecureSecrets:
          mqtt:
            path: mqtt
            Secrets:
              username: use-token-auth
              # TODO - Change <Authentication-Token> placeholder
              password: <Authentication-Token>
              cacert: ""
              clientcert: ""
              clientkey: ""
      ```

* After you have made the appropriate modifications for your Watson IoT account settings, save this file. Make sure the file is saved somewhere that can be easily volume mounted into your container. i.e. next to the compose file in the next step.
* Download the [no security compose file](https://github.com/edgexfoundry/edgex-compose/blob/minnesota/docker-compose-no-secty.yml) and rename it to `docker-compose.yml`.
* Add the following snippet to the compose file:
    ```yaml
      app-ibm-mqtt-export:
        container_name: edgex-app-ibm-mqtt-export
        depends_on:
          - consul
          - data
        environment:
          EDGEX_PROFILE: ibm-mqtt-export
          EDGEX_SECURITY_SECRET_STORE: "false"
          SERVICE_HOST: edgex-app-ibm-mqtt-export
        hostname: edgex-app-ibm-mqtt-export
        image: edgexfoundry/app-service-configurable:3.0.0
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
    docker compose -p edgex up -d
    ```

- View the running containers.

  ```bash
  docker compose -p edgex ps
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
  ```

- Check the logs for `app-ibm-mqtt-export`

  ```bash
  docker compose -p edgex logs app-ibm-mqtt-export
  ```

  ```bash
  level=INFO ts=2023-05-03T21:05:24.789823262Z app=app-<profile> source=service.go:490 msg="Starting app-ibm-mqtt-export 3.0.0 "
  level=INFO ts=2023-05-03T21:05:24.789949262Z app=app-ibm-mqtt-export source=config.go:629 msg="Using Configuration provider (consul) from: http://localhost:8500 with base path of edgex/v3/core-common-config-
  level=INFO ts=2023-05-03T21:05:24.930048273Z app=app-ibm-mqtt-export source=config.go:240 msg="listening for private config changes"
  level=INFO ts=2023-05-03T21:05:24.930073473Z app=app-ibm-mqtt-export source=config.go:242 msg="listening for all services common config changes"
  level=INFO ts=2023-05-03T21:05:24.930104873Z app=app-ibm-mqtt-export source=config.go:245 msg="listening for application service common config changes"
  level=INFO ts=2023-05-03T21:05:24.930142273Z app=app-ibm-mqtt-export source=messaging.go:66 msg="Setting options for secure MessageBus with AuthMode='usernamepassword' and SecretName='redisdb"
  level=INFO ts=2023-05-03T21:05:24.930227273Z app=app-ibm-mqtt-export source=messaging.go:104 msg="Connected to redis Message Bus @ redis://localhost:6379 with AuthMode='usernamepassword'"
  level=INFO ts=2023-05-03T21:05:24.930260473Z app=app-ibm-mqtt-export source=clients.go:164 msg="Using REST for 'core-metadata' clients @ http://localhost:59881"
  level=INFO ts=2023-05-03T21:05:24.931320578Z app=app-ibm-mqtt-export source=manager.go:127 msg="Metrics Manager started with a report interval of 30s"
  level=INFO ts=2023-05-03T21:05:24.931377278Z app=app-ibm-mqtt-export source=bootstrap.go:251 msg="SecuritySecretsRequested metric registered and will be reported (if enabled)"
  level=INFO ts=2023-05-03T21:05:24.931390678Z app=app-ibm-mqtt-export source=bootstrap.go:251 msg="SecuritySecretsStored metric registered and will be reported (if enabled)"
  level=INFO ts=2023-05-03T21:05:24.931544879Z app=app-ibm-mqtt-export source=configupdates.go:48 msg="Waiting for App Service configuration updates..."
  level=INFO ts=2023-05-03T21:05:24.931535079Z app=app-ibm-mqtt-export source=server.go:88 msg="Registering standard routes..."
  level=INFO ts=2023-05-03T21:05:24.932444283Z app=app-ibm-mqtt-export source=service.go:554 msg="Service started in: 142.663921ms"
  level=INFO ts=2023-05-03T21:05:24.932493583Z app=app-ibm-mqtt-export source=main.go:38 msg="Loading Configurable Pipeline..."
  level=INFO ts=2023-05-03T21:05:24.932617584Z app=app-ibm-mqtt-export source=runtime.go:197 msg="PipelineMessagesProcessed-default-pipeline metric has been registered and will be reported (if enabled)"       
  level=INFO ts=2023-05-03T21:05:24.932660284Z app=app-ibm-mqtt-export source=runtime.go:197 msg="PipelineMessageProcessingTime-default-pipeline metric has been registered and will be reported (if enabled)"   
  level=INFO ts=2023-05-03T21:05:24.932677484Z app=app-ibm-mqtt-export source=runtime.go:197 msg="PipelineProcessingErrors-default-pipeline metric has been registered and will be reported (if enabled)"        
  level=INFO ts=2023-05-03T21:05:24.932690384Z app=app-ibm-mqtt-export source=runtime.go:122 msg="Transforms set for `default-pipeline` pipeline"
  level=INFO ts=2023-05-03T21:05:24.932711384Z app=app-ibm-mqtt-export source=triggermessageprocessor.go:90 msg="MessagesReceived metric has been registered and will be reported"
  level=INFO ts=2023-05-03T21:05:24.932738584Z app=app-ibm-mqtt-export source=triggermessageprocessor.go:96 msg="InvalidMessagesReceived metric has been registered and will be reported (if enabled)"
  level=INFO ts=2023-05-03T21:05:24.932765384Z app=app-ibm-mqtt-export source=triggerfactory.go:51 msg="EdgeX MessageBus trigger selected"
  level=INFO ts=2023-05-03T21:05:24.932795485Z app=app-ibm-mqtt-export source=messaging.go:66 msg="Initializing EdgeX Message Bus Trigger for 'redis'"
  level=INFO ts=2023-05-03T21:05:24.932808785Z app=app-ibm-mqtt-export source=messaging.go:83 msg="Subscribing to topic: edgex/events/#"
  level=INFO ts=2023-05-03T21:05:24.932817885Z app=app-ibm-mqtt-export source=messaging.go:98 msg="Publish topic not set for Trigger. Response data, if set, will not be published"
  level=INFO ts=2023-05-03T21:05:24.932962685Z app=app-ibm-mqtt-export source=messaging.go:106 msg="Waiting for messages from the MessageBus on the 'edgex/events/#' topic"
  level=INFO ts=2023-05-03T21:05:24.932803285Z app=app-ibm-mqtt-export source=server.go:181 msg="Starting HTTP Web Server on address localhost:59780"
  level=INFO ts=2023-05-03T21:05:24.933519288Z app=app-ibm-mqtt-export source=service.go:208 msg="StoreAndForward disabled. Not running retry loop."
  level=INFO ts=2023-05-03T21:05:24.933559088Z app=app-ibm-mqtt-export source=service.go:211 msg="app-ibm-mqtt-export has started"
  ```

  
