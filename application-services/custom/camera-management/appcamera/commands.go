//
// Copyright (C) 2022-2023 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"context"
	"fmt"

	"github.com/IOTechSystems/onvif/device"
	"github.com/IOTechSystems/onvif/media"
	"github.com/IOTechSystems/onvif/ptz"
	"github.com/IOTechSystems/onvif/xsd/onvif"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/dtos"
	dtosCommon "github.com/edgexfoundry/go-mod-core-contracts/v3/dtos/common"
	"github.com/pkg/errors"
)

const (
	relativeMoveCommand       = "PTZRelativeMove"
	gotoPresetCommand         = "PTZGotoPreset"
	streamUriCommand          = "StreamUri"
	profilesCommand           = "MediaProfiles"
	getPresetsCommand         = "PTZPresets"
	getCapabilitiesCommand    = "Capabilities"
	getConfigurationsCommand  = "PTZConfigurations"
	startStreamingCommand     = "StartStreaming"
	stopStreamingCommand      = "StopStreaming"
	usbStreamUriCommand       = "StreamURI"
	usbStreamingStatusCommand = "StreamingStatus"
	usbImageFormatsCommand    = "ImageFormats"
)

func (app *CameraManagementApp) getImageFormats(deviceName string) (interface{}, error) {
	resp, err := app.issueGetCommand(context.Background(), deviceName, usbImageFormatsCommand)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to issue get ImageFormats command")
	}
	return resp.Event.Readings[0].ObjectValue, nil
}

func (app *CameraManagementApp) getProfiles(deviceName string) (media.GetProfilesResponse, error) {
	resp := media.GetProfilesResponse{}
	err := app.issueGetCommandForResponse(context.Background(), deviceName, profilesCommand, &resp)
	return resp, err
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
			X: zoom,
		}
	}

	cmd := &ptz.RelativeMove{
		ProfileToken: onvif.ReferenceToken(profileToken),
		Translation:  trans,
	}
	return app.sendPutCommand(deviceName, relativeMoveCommand, cmd)
}

func (app *CameraManagementApp) getCameraFeatures(deviceName string) (CameraFeatures, error) {
	var err error
	features := CameraFeatures{}
	dev, err := app.getDeviceByName(deviceName)
	if err != nil {
		return CameraFeatures{}, err
	}

	switch dev.ServiceName {
	case app.config.AppCustom.OnvifDeviceServiceName:
		features.CameraType = Onvif
		caps, err := app.getCapabilities(deviceName)
		if err != nil {
			return CameraFeatures{}, errors.Wrapf(err, "unable to get device capabilities")
		}
		if caps.Capabilities.PTZ.XAddr != "" {
			features.PTZ = true

			if ptzConfigs, err := app.getPTZConfiguration(deviceName); err != nil {
				app.lc.Errorf("Error calling Get PTZ Configuration for device %s: %s", deviceName, err.Error())
			} else {
				features.Zoom = ptzConfigs.PTZConfiguration[0].ZoomLimits != nil
			}
		}
	case app.config.AppCustom.USBDeviceServiceName:
		features.CameraType = USB
	default:
		features.CameraType = Unknown
	}

	return features, nil
}

func (app *CameraManagementApp) getCapabilities(deviceName string) (device.GetCapabilitiesResponse, error) {
	cmd := &device.GetCapabilities{
		Category: []onvif.CapabilityCategory{onvif.CapabilityCategory("All")},
	}

	resp := device.GetCapabilitiesResponse{}
	err := app.issueGetCommandWithJsonForResponse(context.Background(), deviceName, getCapabilitiesCommand, cmd, &resp)
	return resp, err
}

func (app *CameraManagementApp) getPresets(deviceName string, profileToken string) (ptz.GetPresetsResponse, error) {
	cmd := &ptz.GetPresets{
		ProfileToken: onvif.ReferenceToken(profileToken),
	}

	resp := ptz.GetPresetsResponse{}
	err := app.issueGetCommandWithJsonForResponse(context.Background(), deviceName, getPresetsCommand, cmd, &resp)
	return resp, err
}

func (app *CameraManagementApp) getPTZConfiguration(deviceName string) (ptz.GetConfigurationsResponse, error) {
	cmd := &ptz.GetConfigurations{}

	resp := ptz.GetConfigurationsResponse{}
	err := app.issueGetCommandWithJsonForResponse(context.Background(), deviceName, getConfigurationsCommand, cmd, &resp)
	return resp, err
}

func (app *CameraManagementApp) gotoPreset(deviceName string, profile string, preset string) (dtosCommon.BaseResponse, error) {
	cmd := &ptz.GotoPreset{
		ProfileToken: (*onvif.ReferenceToken)(&profile),
		PresetToken:  (*onvif.ReferenceToken)(&preset),
	}

	return app.sendPutCommand(deviceName, gotoPresetCommand, cmd)
}

func (app *CameraManagementApp) isStreaming(deviceName string) (bool, error) {
	resp := StreamingStatusResponse{}
	err := app.issueGetCommandForResponse(context.Background(), deviceName, usbStreamingStatusCommand, &resp)
	if err != nil {
		return false, err
	}
	return resp.IsStreaming, nil
}

func (app *CameraManagementApp) startStreaming(deviceName string, req USBStartStreamingRequest) (dtosCommon.BaseResponse, error) {
	isStreaming, err := app.isStreaming(deviceName)
	if err == nil && isStreaming {
		// skip if already streaming
		return dtosCommon.BaseResponse{}, nil
	}
	if err != nil {
		app.lc.Errorf(err.Error())
	}
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
	onvifResponse, err1 := app.service.DeviceClient().DevicesByServiceName(context.Background(), app.config.AppCustom.OnvifDeviceServiceName, 0, -1)
	usbResponse, err2 := app.service.DeviceClient().DevicesByServiceName(context.Background(), app.config.AppCustom.USBDeviceServiceName, 0, -1)

	// if both failed, throw an error
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("failed to get devices for the device services: %v, %v", err1, err2)
	}

	var devices []dtos.Device
	if err1 == nil {
		// if the first one succeeded, just overwrite the slice
		devices = onvifResponse.Devices
	}
	if err2 == nil {
		// if the second one succeeded, append all items
		for _, d := range usbResponse.Devices {
			devices = append(devices, d)
		}
	}

	if len(devices) <= 0 {
		return nil, errors.Errorf("no devices registered yet for the device services %s or %s",
			app.config.AppCustom.OnvifDeviceServiceName, app.config.AppCustom.USBDeviceServiceName)
	}

	return devices, nil
}
