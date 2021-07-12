// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Systems LTD
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"errors"
	"os"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	sdkTransforms "github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"

	influxTransforms "app-service-influx/pkg/transforms"
)

const serviceKey = "app-influx-export"

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

	// 3) Load the MQTT custom configuration
	config := &ServiceConfig{}
	if err := service.LoadCustomConfig(config, "MqttSecretConfig"); err != nil {
		lc.Errorf("LoadCustomConfig failed: %s", err.Error())
		os.Exit(-1)
	}

	if err := config.Validate(); err != nil {
		lc.Errorf("Config validation failed: %s", err.Error())
		os.Exit(-1)
	}

	// 4) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := service.SetFunctionsPipeline(
		influxTransforms.NewConversion().TransformToInflux,
		sdkTransforms.NewMQTTSecretSender(config.MqttConfig, false).MQTTSend,
	); err != nil {
		lc.Errorf("SetFunctionsPipeline failed: %s", err.Error())
		os.Exit(-1)
	}

	// 5) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events to trigger the pipeline.
	if err := service.MakeItRun(); err != nil {
		lc.Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here
	os.Exit(0)
}

// Service's custom configuration which is loaded from the configuration.toml
type ServiceConfig struct {
	MqttConfig sdkTransforms.MQTTSecretConfig
}

// UpdateFromRaw updates the service's full configuration from raw data received from
// the Configuration Provider. Can just be a dummy 'return true' if never using the Configuration Provider
func (c *ServiceConfig) UpdateFromRaw(rawConfig interface{}) bool {
	configuration, ok := rawConfig.(*ServiceConfig)
	if !ok {
		return false //errors.New("unable to cast raw config to type 'ServiceConfig'")
	}

	*c = *configuration

	return true
}

func (c *ServiceConfig) Validate() error {
	if len(strings.TrimSpace(c.MqttConfig.BrokerAddress)) == 0 {
		return errors.New("configuration missing value for MqttSecretConfig.BrokerAddress")
	}

	if len(strings.TrimSpace(c.MqttConfig.ClientId)) == 0 {
		return errors.New("configuration missing value for MqttSecretConfig.ClientId")
	}

	if len(strings.TrimSpace(c.MqttConfig.Topic)) == 0 {
		return errors.New("configuration missing value for MqttSecretConfig.Topic")
	}

	if len(strings.TrimSpace(c.MqttConfig.AuthMode)) == 0 {
		return errors.New("configuration missing value for MqttSecretConfig.AuthMode")
	}

	if len(strings.TrimSpace(c.MqttConfig.SecretPath)) == 0 {
		return errors.New("configuration missing value for MqttSecretConfig.SecretPath")
	}

	return nil
}
