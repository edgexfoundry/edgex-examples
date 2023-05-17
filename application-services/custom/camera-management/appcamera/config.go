//
// Copyright (C) 2022-2023 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

// CustomConfig holds the values for the app configuration
type CustomConfig struct {
	OnvifDeviceServiceName string
	USBDeviceServiceName   string
	EvamBaseUrl            string
	MqttAddress            string
	MqttTopic              string
	DefaultPipelineName    string
	DefaultPipelineVersion string
}

// ServiceConfig a struct that wraps CustomConfig which holds the values for driver configuration
type ServiceConfig struct {
	AppCustom CustomConfig
}

// UpdateFromRaw updates the service's full configuration from raw data received from
// the Service Provider.
func (c *ServiceConfig) UpdateFromRaw(rawConfig interface{}) bool {
	configuration, ok := rawConfig.(*ServiceConfig)
	if !ok {
		return false
	}

	*c = *configuration

	return true
}
