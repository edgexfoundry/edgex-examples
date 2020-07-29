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

package functions

import (
	"context"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
)

const (
	jsonSwitchOn  = "{\"SwitchButton\": \"true\"}"
	jsonSwitchOff = "{\"SwitchButton\": \"false\"}"

	appConfigDeviceID  = "DeviceID"
	appConfigCommandID = "CommandID"
)

type Switch struct {
	Status string `json:"status"`
}

func SendSwitchCommand(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	edgexcontext.LoggingClient.Debug("Sending Switch Command")

	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	if edgexcontext.CommandClient == nil {
		edgexcontext.LoggingClient.Error("Command client is available")
		return false, nil
	}

	sw, ok := params[0].(Switch)

	if !ok {
		edgexcontext.LoggingClient.Error("Invalid switch")
		return false, nil
	}

	deviceId := edgexcontext.Configuration.ApplicationSettings[appConfigDeviceID]
	commandId := edgexcontext.Configuration.ApplicationSettings[appConfigCommandID]
	var cmd string

	ctx := context.WithValue(context.Background(), "", "")

	switch status := sw.Status; status {
	case "on":
		edgexcontext.LoggingClient.Info("Switch On")
		cmd = jsonSwitchOn
	case "off":
		edgexcontext.LoggingClient.Info("Switch Off")
		cmd = jsonSwitchOff
	default:
		edgexcontext.LoggingClient.Error("Invalid switch status: " + status)
		return false, nil
	}

	edgexcontext.LoggingClient.Info("Device ID: " + deviceId)
	edgexcontext.LoggingClient.Info("Command ID: " + commandId)
	r, err := edgexcontext.CommandClient.Put(ctx, deviceId, commandId, cmd)

	if err == nil {
		edgexcontext.LoggingClient.Debug("Response : " + r)
	} else {
		edgexcontext.LoggingClient.Error("Error sending request: " + err.Error())

	}

	return true, cmd
}
