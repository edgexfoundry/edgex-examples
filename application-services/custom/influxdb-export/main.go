// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Systems LTD
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	influxTransforms "app-service-influx/pkg/transforms"

	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	sdkTransforms "github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
)

func main() {

	// using insecure secrets from configuration.toml.  This can be removed if setting an env via export in the OS.
	os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// 1) First thing to do is to create an instance of the EdgeX SDK, giving it a service key
	edgexSdk := &appsdk.AppFunctionsSDK{
		ServiceKey: "InfluxDBExport", // Key used by Registry (Aka Consul)
	}

	// 2) Next, we need to initialize the SDK
	if err := edgexSdk.Initialize(); err != nil {
		message := fmt.Sprintf("SDK initialization failed: %v\n", err)
		edgexSdk.LoggingClient.Error(message)
		os.Exit(-1)
	}

	// 3) Initialize the MQTT addressable/configuration
	mqttConfig := sdkTransforms.MQTTSecretConfig{}
	err := influxTransforms.LoadMQTTConfig(edgexSdk, &mqttConfig)

	if err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK MQTT Addressable initialize failed: %v\n", err))
		os.Exit(-1)
	}
	mqttSender := sdkTransforms.NewMQTTSecretSender(mqttConfig, false)

	// 4) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := edgexSdk.SetFunctionsPipeline(
		influxTransforms.NewConversion().TransformToInflux,
		mqttSender.MQTTSend,
	); err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK SetPipeline failed: %v\n", err))
		os.Exit(-1)
	}

	// 5) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events to trigger the pipeline.
	err = edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here
	os.Exit(0)
}
