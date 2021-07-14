# Example Advanced App Functions Service

This **advanced-filter-convert-publish** Application Service example depends on the Device Virtual Go device service to be generating random number events. It uses the following functions in its pipeline:

- Built in **Filter by Value Descriptor** function to filter just for the random **float32** & **float64** values.
- Custom **Value Converter** function which converts the encoded float values to human readable string values.
- Custom **Print to Console** function which simply prints the human readable strings to the console.
- Custom **Publish** function which prepares the modified event with converted values and outputs it to be published back to the message bus using the configured publish host/topic.

The end result from this application service is random float values in human readable format are published back to the message bus for another App Service to consume.

#### End to end Edgex integration proof point for App Functions

Using the following setup, this example advanced **App Functions Service** can be used to demonstrate an EdgeX end to end proof point with **App Functions**.

1. Start **EdgeX** stack

   - [ ] download Ireland non-secure compose file from [here](https://github.com/edgexfoundry/edgex-compose/blob/ireland/docker-compose-no-secty.yml)

   - [ ] start edgex which includes Device Virtual

       ```
       docker-compose -p edgex -f docker-no-secty.yml up -d
       ```

2. Build & run **Advanced App Functions** example

    - [ ] Clone this **[Examples](https://github.com/edgexfoundry/edgex-examples)** repo
    - [ ] cd to **application-services/custom/advanced-filter-convert-publish** folder
    - [ ] run "**make build**"
    - [ ] run "./**app-service**"

3. Configure and Run **Simple Filter XML**  example

   - [ ] cd **application-services/custom/simple-filter-xml** folder

   - [ ] edit **res/configuration.toml** so the **Port** and **SubscribeTopics** sections are as follows:

     ```toml
     [Service]
     HealthCheckInterval = "10s"
     Host = "localhost"
     Port = 59781 <= Change to avoid conflict
     
     [Trigger]
     Type="edgex-messagebus"
       [Trigger.EdgexMessageBus]
       Type = "redis"
         [Trigger.EdgexMessageBus.SubscribeHost]
         Host = "localhost"
         Port = 6379
         Protocol = "redis"
         SubscribeTopics="converted" <= Change so receives Events from this example
     ```
     
   - [ ] run "**make build**"

   - [ ] run "./**app-service**"

4. Now data will be flowing due to the auto-events configured in Device Virtual.

   - In the terminal that you ran **advanced-filter-convert-publish** you will see the random float values printed.

        ```text
        Float64 readable value from Random-Float-Device is 2.1742
        Float32 readable value from Random-Float-Device is 1.3577
        ```

   - In the terminal that you ran **simple-filter-xml** you will see the xml representation of the events printed. Note the human readable float values in the event XML.
        ```xml
        <Event><ApiVersion>v2</ApiVersion><Id>8dc2ee9e-7824-4e57-a4a9-6ceb21229126</Id><DeviceName>Random-Float-Device</DeviceName><ProfileName>MyProfile</ProfileName><SourceName>MySource</SourceName><Origin>1626300284231075300</Origin><Readings><Id>1c1f399b-7cdd-47e8-9bbc-22efe0798ad0</Id><Origin>1626300284231075300</Origin><DeviceName>Random-Float-Device</DeviceName><ResourceName>Float64</ResourceName><ProfileName>Random-Float-Device</ProfileName><ValueType>Float64</ValueType><BinaryValue></BinaryValue><MediaType></MediaType><Value>2.1742</Value></Readings></Event>
        ```
