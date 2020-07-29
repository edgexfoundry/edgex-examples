//
// Copyright (c) 2019 Intel Corporation
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

	"github.com/edgexfoundry/go-mod-core-contracts/models"

	cloudevents "github.com/cloudevents/sdk-go"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"

	localtransforms "cloud-event/pkg/transforms"
)

const (
	serviceKey = "cloudEventTransform"
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")
	// First thing to do is to create an instance of the EdgeX SDK and initialize it.
	edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: serviceKey}

	original, _ := disableSecureStore()
	defer resetSecureStoreEnv(edgexSdk, original)

	if err := edgexSdk.Initialize(); err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK initialization failed: %v\n", err))
		os.Exit(-1)
	}

	// Setup cloudEvent client
	ctx := context.Background()
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("failed to create client, %v", err))
	}
	// Start cloudEvent receiver server.
	// This would probably be a different process running on a different host but
	// this is just to illustrate sending and receiving of cloudevents events
	go func() {
		edgexSdk.LoggingClient.Info("will listen on :8080")
		edgexSdk.LoggingClient.Error(fmt.Sprintf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent)))
	}()

	// Setup pipeline.  The event our first function will get is a edgex event, this
	// will be transformed to a cloudevent and sent on the next function.  The next function, sendCloudEvent
	// will send the event.  Then the next function will transform the event back to an EdgeX event.  The last
	// function will simply print the original event to the console
	edgexSdk.SetFunctionsPipeline(
		localtransforms.NewConversion().TransformToCloudEvent,
		transforms.NewOutputData().SetOutputData,
		sendCloudEvents,
		localtransforms.NewConversion().TransformFromCloudEvent,
		printToConsole,
	)

	// Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	err = edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
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
func sendCloudEvents(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	edgexcontext.LoggingClient.Info("sendCloudEvent")
	if len(params) < 1 {
		return false, errors.New("No Event Received")
	}

	events, ok := params[0].([]cloudevents.Event)
	if !ok {
		return false, errors.New("Cloud event not received")
	}
	ctx := cloudevents.ContextWithTarget(context.Background(), "http://localhost:8080/")
	ctx = cloudevents.ContextWithHeader(ctx, "demo", "header value")
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		edgexcontext.LoggingClient.Error(fmt.Sprintf("failed to create client, %v", err))
		return false, nil
	}
	for _, cloudevent := range events {
		if _, resp, err := c.Send(ctx, cloudevent); err != nil {
			edgexcontext.LoggingClient.Error(fmt.Sprintf("failed to send: %v", err))
			return false, nil
		} else if resp != nil {
			// don't need a response back, in this example we aren't expecting/sending one
			edgexcontext.LoggingClient.Info(fmt.Sprintf("got back a response: %s", resp))
		}
	}
	return true, events
}

func printToConsole(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	edgexcontext.LoggingClient.Info("PrintToConsole")
	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}
	edgexEvent := params[0].(models.Event)
	edgexcontext.LoggingClient.Info(fmt.Sprintf("Original EdgexEvent after it was transformed into a cloudEvent and then back to an EdgexEvent: %v", edgexEvent))
	return false, nil
}

// helper function to disable security
func disableSecureStore() (origEnv string, err error) {
	origEnv = os.Getenv("EDGEX_SECURITY_SECRET_STORE")
	err = // turn off secure mode for examples. Not recommended for production
		os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")
	return origEnv, err
}

// helper function to reset security related env vars
func resetSecureStoreEnv(edgexSdk *appsdk.AppFunctionsSDK, origEnv string) {
	if err := os.Setenv("EDGEX_SECURITY_SECRET_STORE", origEnv); err != nil {
		edgexSdk.LoggingClient.Error("Failed to set env variable: EDGEX_SECURITY_SECRET_STORE back to original value")
	}
}
