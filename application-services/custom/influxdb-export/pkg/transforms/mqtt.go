// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Systems LTD
//
// SPDX-License-Identifier: Apache-2.0

package transforms

import (
	"errors"
	"fmt"

	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	sdkTransforms "github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
)

var log logger.LoggingClient

func getAppSetting(settings map[string]string, name string) string {
	value, ok := settings[name]

	if ok {
		log.Debug(value)
		return value
	}
	log.Error(fmt.Sprintf("application setting %s not found", name))
	return ""
}

func LoadMQTTConfig(sdk *appsdk.AppFunctionsSDK, cfg *sdkTransforms.MQTTSecretConfig) error {

	if sdk == nil {
		return errors.New("Invalid AppFunctionsSDK")
	}

	log = sdk.LoggingClient
	appSettings := sdk.ApplicationSettings()
	if appSettings != nil {
		cfg.BrokerAddress = getAppSetting(appSettings, "BrokerAddress")
		cfg.Topic = getAppSetting(appSettings, "Topic")
		cfg.ClientId = getAppSetting(appSettings, "Publisher")
		cfg.AuthMode = "usernamepassword"
		cfg.SecretPath = "mqtt"
	} else {
		return errors.New("No application-specific settings found")
	}

	return nil
}
