# Simple CBOR Filter Application Service

This **simple-cbor-filter** Application Service demonstrates end to end `CBOR` integration. It depends on the **device-simple** example device service from **Device SDK Go** to generate `CBOR` encode events.

This **simple-cbor-filter** Application Service uses two application functions:

- Built in **FilterByResourceName** function to filter just for the **Image** values.
- Custom **Process Images** function which re-encodes the `binary value` as an Image and prints stats about the image to the console.

The end result from this application service is that it shows that the Application Functions SDK is un-marshaling `CBOR` encode events sent from the **device-simple** device service. These event can be processed by functions similar to `JSON` encoded events. The only difference is the `CBOR` encode events have the `BinaryValue` field set, while the `JSON` encoded events have the `Value` field set.

#### Follow these steps to run the end to end CBOR demonstration

1. Start **EdgeX** stack

   - [ ] down load compose file from [here](https://github.com/edgexfoundry/edgex-compose/blob/main/docker-compose-no-secty.yml)

   - [ ] start edgex 
     
     ```
     docker-compose -p edgex -f docker-compose-no-secty.yml up -d
     ```

3. Build & run **simple-cbor-filter** example

   - [ ] Clone this **[Examples](https://github.com/edgexfoundry/edgex-examples)** repo
   - [ ] cd to **application-services/custom/simple-cbor-filter** folder
   - [ ] run "**make build**"
   - [ ] run "**./app-service**"

3. Run **device-simple** device service

   - [ ] Clone **<https://github.com/edgexfoundry/device-sdk-go>** repo

   - [ ] run "**make build**"

   - [ ] cd to **example/cmd/device-simple** folder

   - [ ] run "./**device-simple**"

     This sample device service will send a `png` (light bulb on) or `jpeg` (light bulb off) image every 30 seconds. The image it sends depends on the value of its `switch` resource, which is `off` (false) by default.
     
     > *Note that since the **device-simple** is running from command-line connecting to services running in Docker, the call-backs from **Core Metadata** for when the new device is added can't not be routed back to **device-simple**. The simple work around is to restart **device-simple** so it then sees the new device it added the first time it was run.*

5. Now data will be flowing due to auto-events configured in **device-simple**.

   - In the terminal that you ran **simple-cbor-filter** you will see the messages like this:

     ```text
     Received Image from Device: Simple-Device01, ResourceName: Image, Image Type: jpeg, Image Size: (1000,1307), Color in middle: {0 128 128}
     ```

     Note that the image received is a jpeg since the `switch` resource in **device-simple** is set to `off ` (false)

   - The `switch` resource can be queried and changed using commands sent via PostMan by doing the following:

     1. Start PostMan

     2. Load the postman collection from the **simple-cbor-filter** example

        `Device Simple Switch commands.postman_collection.json`

     3. This collection contains 3 commands

        - `Get Switch status`
        - `Turn Switch on`
        - `Turn Switch off`

     4. Run  `Turn Switch on`

   -  Now see how the **simple-cbor-filter** output has changed

         ```
         Received Image from Device: Simple-Device01, ReadingName: Image, Image Type: png, Image Size: (1000,1307), Color in middle: {255 246 0 255}
         ```

