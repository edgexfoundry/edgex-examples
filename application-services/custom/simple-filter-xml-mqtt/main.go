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
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
)

const (
	serviceKey = "app-simple-filter-xml-mqtt"
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

	// 2) shows how to access the application's specific simple configuration settings.
	deviceNames, err := service.GetAppSettingStrings("DeviceNames")
	if err != nil {
		lc.Error(err.Error())
		os.Exit(-1)
	}

	lc.Info(fmt.Sprintf("Filtering for devices %v", deviceNames))

	// 3) Configure pipeline functions that require more complex configuration. Here we
	//    use the advanced structure custom configuration
	config := ServiceConfig{}
	if err := service.LoadCustomConfig(&config, "MqttSecretConfig"); err != nil {
		lc.Errorf("LoadCustomConfig returned error: %s", err.Error())
		os.Exit(-1)
	}

	if err := config.Validate(); err != nil {
		lc.Errorf("Custom Config failed validation: %s", err.Error())
		os.Exit(-1)
	}

	// 4) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := service.SetFunctionsPipeline(
		transforms.NewFilterFor(deviceNames).FilterByDeviceName,
		transforms.NewConversion().TransformToXML,
		printXMLToConsole,
		transforms.NewMQTTSecretSender(config.MqttExportConfig, false).MQTTSend,
	); err != nil {
		lc.Errorf("SetFunctionsPipeline returned error: %s", err.Error())
		os.Exit(-1)
	}

	// 5) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = service.MakeItRun()
	if err != nil {
		lc.Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

func printXMLToConsole(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	// Leverage the built in logging service in EdgeX
	lc := ctx.LoggingClient()

	if data == nil {
		return false, errors.New("printXMLToConsole: No data received")
	}

	xml, ok := data.(string)
	if !ok {
		return false, errors.New("printXMLToConsole: Data received is not the expected 'string' type")

	}

	lc.Debug(xml)
	ctx.SetResponseData([]byte(xml))
	return true, xml
}

// Service's custom configuration which is loaded from the configuration.toml
type ServiceConfig struct {
	MqttExportConfig transforms.MQTTSecretConfig
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
	if len(strings.TrimSpace(c.MqttExportConfig.BrokerAddress)) == 0 {
		return errors.New("Configuration missing value for MqttSecretConfig.BrokerAddress")
	}

	if len(strings.TrimSpace(c.MqttExportConfig.ClientId)) == 0 {
		return errors.New("Configuration missing value for MqttSecretConfig.ClientId")
	}

	if len(strings.TrimSpace(c.MqttExportConfig.Topic)) == 0 {
		return errors.New("Configuration missing value for MqttSecretConfig.Topic")
	}

	if len(strings.TrimSpace(c.MqttExportConfig.AuthMode)) == 0 {
		return errors.New("Configuration missing value for MqttSecretConfig.AuthMode")
	}

	return nil
}
