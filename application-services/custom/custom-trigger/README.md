# Custom-Trigger

Custom-Trigger provides an example for using a custom trigger with the app functions SDK.

## Overview

In this example we introduce a trigger that listens for input on os.Stdin and sends it through the function pipeline.  The function pipeline is a single function that introduce a slight delay based on the length of the input string and then print it to the console.

To run:

`make docker && docker run -it edgexfoundry/docker-custom-trigger-example:dev`