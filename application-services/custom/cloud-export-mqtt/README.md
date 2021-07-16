# cloud export mqtt

This example service is meant to demonstrate using a custom transformation function (example only) along with the MQTTSecretSender to deliver edgex data to a cloud provider over MQTT.  This is done via profiles that really only differ in the MQTT connection settings.

To run you can pass the profile on the command line, eg:

`./app-service -p aws`

or

`./app-service -p azure`
