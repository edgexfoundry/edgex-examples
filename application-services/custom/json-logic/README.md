# JSON Logic Examples

This example demonstrates a few different ways to leverage JSON logic to assist in a bit more flexible and advanced way to perform filtering operations without the need to write your own functions. If you have many rules or want an even more feature rich rules engine check out [eKuiper](https://github.com/lf-edge/ekuiper) here.

> Note: Only  operators that return true or false are supported. See http://jsonlogic.com/operations.html# for the complete list paying attention to return values. Any operator that returns manipulated data is currently not supported. 

> Tip: Leverage http://jsonlogic.com/play.html to get your rule right before implementing in code. JSON can be a bit tricky to get right in code with all the escaped double quotes.  

## Filter By Device Name - (Same as provided SDK Function)
This first example starts out simple to demonstrate how JSON logic can be used to accomplish the same thing provided by the FilterByDeviceName transform.

``` go
jsonlogicrule := "{ \"in\" : [{ \"var\" : \"deviceName\" }, [\"Random-Integer-Device\",\"Random-Float-Device\"] ] }"

edgexSdk.SetFunctionsPipeline(		
    transforms.NewJSONLogic(jsonlogicrule).Evaluate,
    transforms.NewConversion().TransformToXML,
	printXMLToConsole,
)
```

> Note: In the code we pull device names from the configuration instead of hard-coding them as shown here. 

But what if we wanted to perform an not operation? 

## Filter OUT by Device Name
The following rule demonstrates how to filter OUT the list of devices instead of filtering FOR a specific device. 

``` go
jsonlogicrule := "{ \"!\" : {\"in\" : [{ \"var\" : \"deviceName\" }, [\"Random-Integer-Device\"] ] }}"

edgexSdk.SetFunctionsPipeline(		
    transforms.NewJSONLogic(jsonlogicrule).Evaluate,
    transforms.NewConversion().TransformToXML,
	printXMLToConsole,
)
```
With this rule, all devices will flow through EXCEPT for "Random-Integer-Device".

## Filter for readings that are greater than 0
Sometime we'll need to access data inside the readings of an EdgeX Event. Using Random-Float-Device as an example, we first need to convert the float value to be numerical instead of base64 encoded. In the example, we have provided the `ConvertToReadableFloatValues` function. The following rule can be used to perform this filter:

``` go
jsonlogicrule := "{ \"all\" : [ { \"var\" : \"readings\" } , {  \">\" : [ {\"var\":\"value\"}, 0 ] } ] }"

edgexSdk.SetFunctionsPipeline(		
    ConvertToReadableFloatValues,
    transforms.NewJSONLogic(jsonlogicrule).Evaluate,
    transforms.NewConversion().TransformToXML,
	printXMLToConsole,
)
```
 By using the "all" operator, we are asking that ALL values in the readings array be greater than 0. If one reading is not greater than 0, then the pipeline execution will stop and not continue. You can leverage the "some" operator to allow pipeline execution to continue if at least one reading is greater than 0.

