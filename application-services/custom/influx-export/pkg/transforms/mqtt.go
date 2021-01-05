// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2020 IOTech Systems LTD
//
// SPDX-License-Identifier: Apache-2.0

package transforms

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
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

func LoadMQTTAddressable(sdk *appsdk.AppFunctionsSDK) (*models.Addressable, error) {

	addr := models.Addressable{}

	if sdk == nil {
		return nil, errors.New("Invalid AppFunctionsSDK")
	}

	log = sdk.LoggingClient

	appSettings := sdk.ApplicationSettings()
	if appSettings != nil {
		addr.Address = getAppSetting(appSettings, "Host")
		i, err := strconv.Atoi(getAppSetting(appSettings, "Port"))
		if err != nil {
			return nil, errors.New("MQTT Port not a number in configuration")
		}
		addr.Port = i
		addr.Protocol = getAppSetting(appSettings, "Protocol")
		addr.Publisher = getAppSetting(appSettings, "Publisher")
		addr.User = getAppSetting(appSettings, "User")
		addr.Password = getAppSetting(appSettings, "Password")
		addr.Topic = getAppSetting(appSettings, "Topic")
	} else {
		return nil, errors.New("No application-specific settings found")
	}

	return &addr, nil
}
