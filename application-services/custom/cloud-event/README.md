# Cloud-Event-Transform

Cloud-Event-Transforms is an example extension of the [pkg/transforms/conversion](https://github.com/edgexfoundry/app-functions-sdk-go/blob/master/pkg/transforms/conversion.go) found in the App Functions SDK to illustrate the transformation of Edgex events to [cloud-events](https://github.com/cloudevents/spec) and back.

## Overview

In this example we start off by using the included `transforms.NewConversion().TransformToCloudEvent` function to transform an Edgex Event into a cloud-event.  Then in the following app function we send the cloud-event, then transform it back to an Edgex Event using `transforms.NewConversion().TransformFromCloudEvent` and print out the result.  

You can trigger this pipeline by using the http trigger from a tool like postman and sending an edgex event in the JSON body as payload. 
