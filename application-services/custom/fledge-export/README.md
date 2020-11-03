# Fledge Export

This example shows how to get data (readings) from an EdgeX Foundry instance to [LF Edge](https://www.lfedge.org/) sister project [Fledge](https://www.lfedge.org/projects/fledge/).

Fledge offers North and South plugins.  In this example, EdgeX Foundry is sending sensor/device data via a custom [application service](https://docs.edgexfoundry.org/1.2/microservices/application/ApplicationServices/) to the [Fledge South HTTP plugin](https://fledge-iot.readthedocs.io/en/v1.8.1/plugins/fledge-south-http/index.html#).

## Fledge Version

This example was created with Fledge version 1.8.  

## Fledge Setup and Configuration

See the [Fledge documentation](https://fledge-iot.readthedocs.io/en/v1.8.1/quick_start.html) on how to setup and configure Fledge.

See the [Fledge Plugin documentation](https://fledge-iot.readthedocs.io/en/v1.8.1/plugins/fledge-south-http/index.html) for information on how to install and configure the South HTTP plugin.

## Build and Run

A Makefile has been provide to easily create and execute this service.  In order to build the micro service executable run the `make build` from the root directory of this example.

Once the micro service has successfully been compiled, run the executable created in the root directory with `./app-service`.

## Configuration

In order to supply data from your EdgeX instance to your Fledge instance, you must provide the REST endpoint for the Fledge South HTTP plugin.  Open the `configuration.toml` file in the `res` folder and change the address associated `FledgeSouthHTTPEndpoint` in the `ApplicationSettings` configuration section.

``` toml
    [ApplicationSettings]
    FledgeSouthHTTPEndpoint = "http://192.168.0.10:6683/sensor-reading"
```

Unless you made changes to the default configuration when installing the Fledge South HTTP plugin, you should only need to change the IP address of your Fledge instance. 

## Data Payload

EdgeX uses the Fledge payload format to export EdgeX Readings (per each EdgeX Event) to Fledge.  The South HTTP Plugin requires the data to be sent in the following way:

``` JSON
    [
        {
            "timestamp" : "2020-07-08 16:16:07.263657+00:00",
            "asset" : "motor1",
            "readings" : {
                "voltage" : 239.4,
                "current" : 1003,
                "rpm" : 120147
            }
        }
    ]
```

The `asset` is set with the EdgeX device name.  The `readings` hold the EdgeX Reading name and value pairs.  The `timestamp` is set from the Event's created timestamp.