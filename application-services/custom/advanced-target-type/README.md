# Advanced Target Type  

This **advanced-target-type** Application Service demonstrates how to create an Application Service that expects a custom type to feed to the functions pipeline. For more detail refer to the Application Functions SDK documentation section on [TargetType](https://docs.edgexfoundry.org/latest/microservices/application/AdvancedTopics/#target-type)

To run this example:

1.  Clone **[edgex-examples](https://github.com/edgexfoundry/edgex-examples)** repo

2. cd `edgex-examples/application-services/custom/advanced-target-type`

3. run `make build`

4. run `./app-service`

5. Start PostMan

6. Load `Post Person to Trgger.postman_collection.json` collection in PostMan

7. Run the `Person Trigger` request

   - The following XML will be printed to the console by the Application Service and will be returned as the trigger HTTP response in PostMan.

     ```
     <Person>
        <FirstName>Sam</FirstName>
        <LastName>Smith</LastName>
        <Phone>
           <CountryCode>1</CountryCode>
           <AreaCode>480</AreaCode>
           <LocalPrefix>970</LocalPrefix>
           <LocalNumber>3476</LocalNumber>
        </Phone>
        <PhoneDisplay>+01(480) 970-3476</PhoneDisplay>
     </Person>
     ```

   - Note that the PhoneDisplay field that is not present in the XML sent from PostMan is now present and filled out.

