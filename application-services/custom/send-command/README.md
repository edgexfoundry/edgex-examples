# Send Command Example #

## Overview ##

This Application Service example demonstrates how to use the `Command` client to send `Set` and `Get` commands to a device. In this case we are using the a device from the Device Virtual service.

## Prerequisites ##

* Obtain the code from the https://github.com/edgexfoundry/edgex-examples/application-services/custom/send-command 

* Ensure that EdgeX is running including Device Virtual Service. Run the follow command to achieve this

  ```bash
  curl https://raw.githubusercontent.com/edgexfoundry/edgex-compose/ireland/docker-compose-no-secty.yml -o docker-compose.yml; docker-compose -p edgex up -d
  ```

- Install PostMan (https://www.postman.com/)

## Steps

1. Build the Send Command service 

```
    make build
```

2. Run the Send Command service 

   ```
   ./app-service
   ```

3. Start PostMan and import the `Send Command Triggers.postman_collection.json` file

   This collection has two requests. One to trigger a `Set` command and one to trigger a `Get` command. Both send a custom `ActionRequest` object to the application service's HTTP trigger. This `ActionRequest` contains the information needed to send the commands. Take a look at the JSON being sent and play around with the values to send different commands to devices defined by Device Virtual . 

   See here for full list of devices: https://github.com/edgexfoundry/device-virtual-go/blob/v2.0.0/cmd/res/devices/devices.toml and here for all the profiles that define the resources for those devices: https://github.com/edgexfoundry/device-virtual-go/tree/v2.0.0/cmd/res/profiles

4. Run the `Trigger Set Action` request from PostMan

   Response should be:

   ```json
   {
       "apiVersion": "v2",
       "statusCode": 200
   }
   ```

5. Run the `Trigger Get Action` request from PostMan

   Response should be:

   ```json
   {
       "apiVersion": "v2",
       "statusCode": 200,
       "event": {
           "apiVersion": "v2",
           "id": "2369f3f1-df20-4522-bd2a-1fbc20e7049e",
           "deviceName": "Random-Integer-Device",
           "profileName": "Random-Integer-Device",
           "sourceName": "Int8",
           "origin": 1625868240254928800,
           "readings": [
               {
                   "id": "b37c67c6-1678-4d18-b07c-a6f148aeaefd",
                   "origin": 1625868240254928800,
                   "deviceName": "Random-Integer-Device",
                   "resourceName": "Int8",
                   "profileName": "Random-Integer-Device",
                   "valueType": "Int8",
                   "binaryValue": null,
                   "mediaType": "",
                   "value": "101"
               }
           ]
       }
   }
   ```