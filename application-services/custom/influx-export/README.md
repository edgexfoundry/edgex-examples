# EdgeX Foundry application service InfluxDB

EdgeX Foundry application service to export readings to InfluxDB via MQTT.

The pipeline simply consists of two functions:
- conversion of EdgeX Event/Reading objects to Influx [line protocol](https://docs.influxdata.com/influxdb/v2.0/reference/syntax/line-protocol/)
- export of results to MQTT Sender, which is a MQTT topic that Telegraf / Influx is positioned to listen to and ingest data from

In this example, the MQTT Broker is assumed to exist on a different host and is accessed via port 1884.  The target topic is edgex/EdgeXEvents.  Modify the configuration.toml file to point the MQTT Sender to an alternate broker or topic.

No filtering or other transformation functions occur in this example.
No additional MQTT options (such as QOS or retry has been provided for with this example)

InfluxDB installation and tooling is not provided with this example.

![image](./app-service-influx.png)
