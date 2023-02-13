//
// Copyright (C) 2022 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IOTechSystems/onvif/ptz"
	"github.com/IOTechSystems/onvif/xsd/onvif"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	dtosCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/pkg/errors"
)

const (
	relativeMoveCommand      = "RelativeMove"
	gotoPresetCommand        = "GotoPreset"
	streamUriCommand         = "StreamUri"
	profilesCommand          = "Profiles"
	getPresetsCommand        = "GetPresets"
	getConfigurationsCommand = "GetConfigurations"
	startStreamingCommand    = "StartStreaming"
	stopStreamingCommand     = "StopStreaming"
	usbStreamUriCommand      = "StreamURI"

	zoomScale = 1
)

func (app *CameraManagementApp) getProfiles(deviceName string) (ProfilesResponse, error) {
	profiles, err := app.issueGetCommandWithJson(context.Background(), deviceName, profilesCommand, struct{}{})
	if err != nil {
		return ProfilesResponse{}, errors.Wrapf(err, "failed to issue get Profiles command")
	}

	val := profiles.Event.Readings[0].ObjectValue
	js, err := json.Marshal(val)
	if err != nil {
		return ProfilesResponse{}, errors.Wrapf(err, "failed to marshal profiles json object")
	}
	pr := ProfilesResponse{}
	err = json.Unmarshal(js, &pr)
	return pr, err
}

func (app *CameraManagementApp) doPTZ(deviceName, profileToken string, x, y, zoom float64) (dtosCommon.BaseResponse, error) {
	trans := ptz.Vector{
		PanTilt: &onvif.Vector2D{
			X: x,
			Y: y,
		},
	}
	if zoom != 0 {
		trans.Zoom = &onvif.Vector1D{
			X: zoom * zoomScale,
		}
	}

	cmd := &ptz.RelativeMove{
		ProfileToken: onvif.ReferenceToken(profileToken),
		Translation:  trans,
	}
	return app.sendPutCommand(deviceName, relativeMoveCommand, cmd)
}

func (app *CameraManagementApp) getPresets(deviceName string, profileToken string) (GetPresetsResponse, error) {
	cmd := &ptz.GetPresets{
		ProfileToken: onvif.ReferenceToken(profileToken),
	}

	presets, err := app.issueGetCommandWithJson(context.Background(), deviceName, getPresetsCommand, cmd)
	if err != nil {
		return GetPresetsResponse{}, errors.Wrapf(err, "failed to issue get presets command")
	}

	val := presets.Event.Readings[0].ObjectValue
	js, err := json.Marshal(val)
	if err != nil {
		return GetPresetsResponse{}, errors.Wrapf(err, "failed to marshal presets json object")
	}
	pr := GetPresetsResponse{}
	err = json.Unmarshal(js, &pr)
	return pr, err
}

func (app *CameraManagementApp) getPTZConfiguration(deviceName string) (GetPTZConfigurationsResponse, error) {
	cmd := &ptz.GetConfigurations{}

	config, err := app.issueGetCommandWithJson(context.Background(), deviceName, getConfigurationsCommand, cmd)
	if err != nil {
		return GetPTZConfigurationsResponse{}, errors.Wrapf(err, "failed to issue get configurations command")
	}

	val := config.Event.Readings[0].ObjectValue
	js, err := json.Marshal(val)
	if err != nil {
		return GetPTZConfigurationsResponse{}, errors.Wrapf(err, "failed to marshal configurations json object")
	}
	pr := GetPTZConfigurationsResponse{}
	err = json.Unmarshal(js, &pr)
	return pr, err
}

func (app *CameraManagementApp) gotoPreset(deviceName string, profile string, preset string) (dtosCommon.BaseResponse, error) {
	cmd := &ptz.GotoPreset{
		ProfileToken: (*onvif.ReferenceToken)(&profile),
		PresetToken:  (*onvif.ReferenceToken)(&preset),
	}

	return app.sendPutCommand(deviceName, gotoPresetCommand, cmd)
}

func (app *CameraManagementApp) startStreaming(deviceName string, req USBStartStreamingRequest) (dtosCommon.BaseResponse, error) {
	return app.sendPutCommand(deviceName, startStreamingCommand, req)
}

func (app *CameraManagementApp) stopStreaming(deviceName string) (dtosCommon.BaseResponse, error) {
	return app.sendPutCommand(deviceName, stopStreamingCommand, true)
}

func (app *CameraManagementApp) sendPutCommand(deviceName string, commandName string, commandValue interface{}) (dtosCommon.BaseResponse, error) {
	app.lc.Infof("Sending PUT command: %s, %+v", commandName, commandValue)
	return app.service.CommandClient().IssueSetCommandByNameWithObject(context.Background(), deviceName, commandName,
		map[string]interface{}{
			// note: we are using the actual name of the command as the key
			commandName: commandValue,
		})
}

func (app *CameraManagementApp) getAllDevices() ([]dtos.Device, error) {
	response1, err1 := app.service.DeviceClient().DevicesByServiceName(context.Background(), app.config.AppCustom.OnvifDeviceServiceName, 0, -1)
	response2, err2 := app.service.DeviceClient().DevicesByServiceName(context.Background(), app.config.AppCustom.USBDeviceServiceName, 0, -1)

	// if both failed, throw an error
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("failed to get devices for the device services: %v, %v", err1, err2)
	}

	var devices []dtos.Device
	if err1 == nil {
		// if the first one succeeded, just overwrite the slice
		devices = response1.Devices
	}
	if err2 == nil {
		// if the second one succeeded, append all items
		for _, d := range response2.Devices {
			devices = append(devices, d)
		}
	}

	if len(devices) <= 0 {
		return nil, errors.Errorf("no devices registered yet for the device services %s or %s",
			app.config.AppCustom.OnvifDeviceServiceName, app.config.AppCustom.USBDeviceServiceName)
	}

	return devices, nil
}
