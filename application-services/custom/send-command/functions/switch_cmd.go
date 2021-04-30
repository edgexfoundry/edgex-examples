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

package functions

import (
	"context"
	"errors"
	"fmt"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
)

const (
	jsonSwitchOn  = "{\"SwitchButton\": \"true\"}"
	jsonSwitchOff = "{\"SwitchButton\": \"false\"}"
)

type Switch struct {
	Status string `json:"status"`
}

type SwitchCommand struct {
	deviceName  string
	commandName string
}

func NewSwitchCommand(deviceName, commandName string) *SwitchCommand {
	return &SwitchCommand{
		deviceName:  deviceName,
		commandName: commandName,
	}
}

func (s *SwitchCommand) SendSwitchCommand(funcCtx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	lc := funcCtx.LoggingClient()

	lc.Debug("Sending Switch Command")

	if data == nil {
		return false, errors.New("SendSwitchCommand: No data received")
	}

	if funcCtx.CommandClient == nil {
		return false, errors.New("SendSwitchCommand: Command client is available")
	}

	sw, ok := data.(Switch)
	if !ok {
		return false, errors.New("SendSwitchCommand: Data received is not the expected 'Switch' type")
	}

	var commandBody string

	ctx := context.WithValue(context.Background(), "", "")

	switch status := sw.Status; status {
	case "on":
		lc.Info("Switch On")
		commandBody = jsonSwitchOn
	case "off":
		lc.Info("Switch Off")
		commandBody = jsonSwitchOff
	default:
		lc.Error("Invalid switch status: " + status)
		return false, nil
	}

	lc.Infof("Sending command '%s' for device '%s'", s.deviceName, s.commandName)
	r, err := funcCtx.CommandClient().PutDeviceCommandByNames(ctx, s.deviceName, s.commandName, commandBody)

	if err == nil {
		lc.Debug("Response : " + r)
	} else {
		return false, fmt.Errorf("Error sending command request: %s", err.Error())
	}

	return true, commandBody
}
