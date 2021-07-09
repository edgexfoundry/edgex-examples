//
// Copyright (c) 2019 Intel Corporation
// Copyright (c) 2021 One Track Consulting
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"os"

	cloudevents "github.com/cloudevents/sdk-go"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"

	localtransforms "cloud-event/pkg/transforms"
)

const (
	serviceKey = "advancedCloudEventTransform"
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")
	// First thing to do is to create an instance of the EdgeX Application Service and initialize it.
	appService, ok := pkg.NewAppService(serviceKey)

	if !ok {
		appService.LoggingClient().Errorf("App service initialization failed for `%s`", serviceKey)
		os.Exit(-1)
	}

	// Setup cloudEvent client
	ctx := context.Background()
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		appService.LoggingClient().Errorf("failed to create cloud events client, %s", err.Error())
	}

	// Start cloudEvent receiver server.
	// This would probably be a different process running on a different host but
	// this is just to illustrate sending and receiving of cloudevents events
	go func() {
		appService.LoggingClient().Info("will listen on :8080")
		appService.LoggingClient().Error(fmt.Sprintf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent)))
	}()

	// Setup pipeline.  The event our first function will get is a edgex event, this
	// will be transformed to a cloudevent and sent on the next function.  The next function, sendCloudEvent
	// will send the event.  Then the next function will transform the event back to an EdgeX event.  The last
	// function will simply print the original event to the console
	appService.SetFunctionsPipeline(
		localtransforms.NewConversion().TransformToCloudEvent,
		transforms.NewResponseData().SetResponseData,
		sendCloudEvents,
		localtransforms.NewConversion().TransformFromCloudEvent,
		printToConsole,
	)

	// Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	err = appService.MakeItRun()
	if err != nil {
		appService.LoggingClient().Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here
	os.Exit(0)
}

// function called when the cloudevent receiver gets an event
func gotEvent(ctx context.Context, event cloudevents.Event) {
	var readingData string
	event.DataAs(&readingData)
	fmt.Printf("CloudEvent received reading value: %v\n", readingData)
}

// App function to send the cloudevent to the receiver
func sendCloudEvents(exctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	exctx.LoggingClient().Info("sendCloudEvent")
	if data == nil {
		return false, errors.New("No Event Received")
	}

	events, ok := data.([]cloudevents.Event)
	if !ok {
		return false, errors.New("Cloud event not received")
	}
	ctx := cloudevents.ContextWithTarget(context.Background(), "http://localhost:8080/")
	ctx = cloudevents.ContextWithHeader(ctx, "demo", "header value")
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		exctx.LoggingClient().Errorf("failed to create client, %s", err.Error())
		return false, nil
	}
	for _, cloudevent := range events {
		if _, resp, err := c.Send(ctx, cloudevent); err != nil {
			exctx.LoggingClient().Errorf("failed to send: %s", err.Error())
			return false, nil
		} else if resp != nil {
			// don't need a response back, in this example we aren't expecting/sending one
			exctx.LoggingClient().Infof("got back a response: %s", resp)
		}
	}
	return true, events
}

func printToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	ctx.LoggingClient().Info("PrintToConsole")
	if data == nil {
		// We didn't receive a result
		return false, nil
	}
	edgexEvent := data.(models.Event)
	ctx.LoggingClient().Infof("Original EdgexEvent after it was transformed into a cloudEvent and then back to an EdgexEvent: %v", edgexEvent)
	return false, nil
}
