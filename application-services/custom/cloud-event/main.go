//
// Copyright (c) 2021 Intel Corporation
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
	"os"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/dtos"

	cloud "github.com/cloudevents/sdk-go"

	localtransforms "cloud-event/pkg/transforms"
)

const (
	serviceKey = "cloudEventTransform"
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	_ = os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// 1) First thing to do is to create an new instance of an EdgeX Application Service.
	service, ok := pkg.NewAppService(serviceKey)
	if !ok {
		os.Exit(-1)
	}

	// Leverage the built in logging service in EdgeX
	lc := service.LoggingClient()

	// Setup cloudEvent client
	ctx := context.Background()
	c, err := cloud.NewDefaultClient()
	if err != nil {
		lc.Errorf("failed to create client, %s", err.Error())
	}
	// Start cloudEvent receiver server.
	// This would probably be a different process running on a different host but
	// this is just to illustrate sending and receiving of cloud events
	go func() {
		lc.Info("will listen on :8080")
		lc.Error(fmt.Sprintf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent)))
	}()

	// Setup pipeline.  The event our first function will get is a edgex event, this
	// will be transformed to a cloudevent and sent on the next function.  The next function, sendCloudEvent
	// will send the event.  Then the next function will transform the event back to an EdgeX event.  The last
	// function will simply print the original event to the console
	if err := service.SetFunctionsPipeline(
		localtransforms.NewConversion().TransformToCloudEvent,
		transforms.NewResponseData().SetResponseData,
		sendCloudEvents,
		localtransforms.NewConversion().TransformFromCloudEvent,
		printToConsole,
	); err != nil {
		lc.Error("SetFunctionsPipeline returned error: ", err.Error())
		os.Exit(-1)
	}

	// Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	if err := service.MakeItRun(); err != nil {
		lc.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here
	os.Exit(0)
}

// function called when the cloud event receiver gets an event
func gotEvent(_ context.Context, event cloud.Event) {
	var readingData string
	event.DataAs(&readingData)
	fmt.Printf("CloudEvent received reading value: %v\n", readingData)
}

// App function to send the cloudevent to the receiver
func sendCloudEvents(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	lc := ctx.LoggingClient()
	lc.Debug("sendCloudEvents")

	if data == nil {
		return false, errors.New("sendCloudEvents: No data received")
	}

	events, ok := data.([]cloud.Event)
	if !ok {
		return false, errors.New("sendCloudEvents: didn't receive expect '[]cloud.Event' type")
	}

	sendCtx := cloud.ContextWithTarget(context.Background(), "http://localhost:8080/")
	sendCtx = cloud.ContextWithHeader(sendCtx, "demo", "header value")
	c, err := cloud.NewDefaultClient()
	if err != nil {
		lc.Error(fmt.Sprintf("failed to create client, %v", err))
		return false, nil
	}
	for _, cloudEvent := range events {
		if _, resp, err := c.Send(sendCtx, cloudEvent); err != nil {
			lc.Errorf("failed to send Cloud Event: %s", err.Error())
			return false, nil
		} else if resp != nil {
			// don't need a response back, in this example we aren't expecting/sending one
			lc.Infof("got back a response: %s", resp)
		}
	}
	return true, events
}

func printToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	lc := ctx.LoggingClient()

	lc.Debug("PrintToConsole")

	if data == nil {
		return false, errors.New("printToConsole: No data received")
	}

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("printToConsole: didn't receive expect Event type")
	}

	lc.Infof("Original Edgex Event after it was transformed into a Cloud Event and then back to an Edgex Event: %v", event)
	return false, nil
}
