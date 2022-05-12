# Custom-Trigger-Nats-RPC

Custom-Trigger provides an example to create an RPC binding to an edgex app service using NATS.

## Overview

In this example we introduce an RPC trigger offering request/reply semantics over the NATS protocol.  This provides a synchronous alternative to the builtin HTTP trigger.

It is hard-coded to use the nats demo server at `demo.nats.io:4222`

To run:

```console
make build
./app-service
level=INFO ts=2022-05-12T02:42:19.204390053Z app=app-custom-trigger-nats-rpc source=service.go:496 msg="Starting app-custom-trigger-nats-rpc 0.0.0 "
level=INFO ts=2022-05-12T02:42:19.204770966Z app=app-custom-trigger-nats-rpc source=config.go:391 msg="Loaded service configuration from ./res/configuration.toml"
level=INFO ts=2022-05-12T02:42:19.205405876Z app=app-custom-trigger-nats-rpc source=config.go:551 msg="Using local configuration from file (0 envVars overrides applied)"
level=INFO ts=2022-05-12T02:42:19.205427811Z app=app-custom-trigger-nats-rpc source=version.go:79 msg="Skipping version compatibility check for SDK Beta version or running in debugger: version=0.0.0"
level=INFO ts=2022-05-12T02:42:19.2054498Z app=app-custom-trigger-nats-rpc source=metrics.go:64 msg="0 specified for metrics reporting interval. Setting to max duration to effectively disable reporting."
level=INFO ts=2022-05-12T02:42:19.205459582Z app=app-custom-trigger-nats-rpc source=manager.go:122 msg="Metrics Manager started with a report interval of 2562047h47m16.854775807s"
level=INFO ts=2022-05-12T02:42:19.205468999Z app=app-custom-trigger-nats-rpc source=server.go:86 msg="Registering standard routes..."
level=INFO ts=2022-05-12T02:42:19.205512516Z app=app-custom-trigger-nats-rpc source=service.go:559 msg="Service started in: 1.137204ms"
level=INFO ts=2022-05-12T02:42:19.205548988Z app=app-custom-trigger-nats-rpc source=runtime.go:157 msg="PipelineMessagesProcessed-default-pipeline metric has been registered and will be reported (if enabled)"
level=INFO ts=2022-05-12T02:42:19.205565557Z app=app-custom-trigger-nats-rpc source=runtime.go:165 msg="PipelineMessageProcessingTime-default-pipeline metric has been registered and will be reported (if enabled)"
level=INFO ts=2022-05-12T02:42:19.205574175Z app=app-custom-trigger-nats-rpc source=runtime.go:117 msg="Transforms set for `default-pipeline` pipeline"
level=INFO ts=2022-05-12T02:42:19.205585151Z app=app-custom-trigger-nats-rpc source=triggerfactory.go:53 msg="MessagesReceived metric has been registered and will be reported"
level=WARN ts=2022-05-12T02:42:19.205594473Z app=app-custom-trigger-nats-rpc source=triggerfactory.go:107 msg="MessagesReceived metric failed to register and will not be reported: duplicate metric: MessagesReceived"
level=INFO ts=2022-05-12T02:42:19.205620638Z app=app-custom-trigger-nats-rpc source=server.go:162 msg="Starting HTTP Web Server on address localhost:59780"
level=INFO ts=2022-05-12T02:42:19.205898997Z app=app-custom-trigger-nats-rpc source=configupdates.go:48 msg="Waiting for App Service configuration updates..."
level=INFO ts=2022-05-12T02:42:19.206010871Z app=app-custom-trigger-nats-rpc source=telemetry.go:78 msg="Starting CPU Usage Average loop"
level=INFO ts=2022-05-12T02:42:19.346342198Z app=app-custom-trigger-nats-rpc source=service.go:202 msg="StoreAndForward disabled. Not running retry loop."
level=INFO ts=2022-05-12T02:42:19.346387167Z app=app-custom-trigger-nats-rpc source=service.go:205 msg="Custom Trigger Example"
```

You can then run nats-box (or any nats client) and call the service using request/reply semantics eg:

```shell
alex@computer:~g$ docker run -it natsio/nats-box
             _             _               
 _ __   __ _| |_ ___      | |__   _____  __
| '_ \ / _` | __/ __|_____| '_ \ / _ \ \/ /
| | | | (_| | |_\__ \_____| |_) | (_) >  < 
|_| |_|\__,_|\__|___/     |_.__/ \___/_/\_\
                                           
nats-box v0.11.0
e21a2f9352a1:~# nats context add nats --server demo.nats.io:4222 --description "NATS Demo" --select
NATS Configuration Context "nats"

      Description: NATS Demo
      Server URLs: demo.nats.io:4222
             Path: /nsc/.config/nats/context/nats.json
       Connection: OK

e21a2f9352a1:~# nats request rpc.testtopic "testing"
02:42:31 Sending request on "rpc.testtopic"
02:42:31 Received with rtt 96.010289ms
got testing (from rpc.testtopic at 2022-05-12 02:42:31.364262479 +0000 UTC)

e21a2f9352a1:~# nats request rpc.testtopic6 "testing6"
02:42:35 Sending request on "rpc.testtopic6"
02:42:35 Received with rtt 103.899353ms
got testing6 (from rpc.testtopic6 at 2022-05-12 02:42:35.748212233 +0000 UTC)

```
