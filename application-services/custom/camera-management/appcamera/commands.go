//
// Copyright (C) 2022 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"context"
	"encoding/json"
	"github.com/IOTechSystems/onvif/ptz"
	"github.com/IOTechSystems/onvif/xsd/onvif"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	dtosCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/pkg/errors"
)

const (
	relativeMoveCommand = "RelativeMove"
	gotoPresetCommand   = "GotoPreset"
	streamUriCommand    = "StreamUri"
	profilesCommand     = "Profiles"
	getPresetsCommand   = "GetPresets"

	moveScaleX = 10
	moveScaleY = 5
	zoomScale  = 1
)

func (app *CameraManagementApp) getProfiles(deviceName string) (ProfilesResponse, error) {
	profiles, err := app.issueGetCommand(context.Background(), deviceName, profilesCommand, struct{}{})
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

func (app *CameraManagementApp) queryStreamUri(deviceName, profileToken string) (string, error) {
	req := StreamUriRequest{ProfileToken: profileToken}
	cmdResponse, err := app.issueGetCommand(context.Background(), deviceName, streamUriCommand, req)
	if err != nil {
		return "", errors.Wrapf(err, "failed to issue get StreamUri command")
	}
	streamUri, errUri := parseStreamUri(cmdResponse)
	if errUri != nil {
		return "", errors.Wrapf(errUri, "failed to get stream Uri from the device %s", deviceName)
	}

	return streamUri, nil
}

func (app *CameraManagementApp) doPTZ(deviceName, profileToken string, x, y, zoom float64) (dtosCommon.BaseResponse, error) {
	trans := ptz.Vector{
		PanTilt: &onvif.Vector2D{
			X: x * moveScaleX,
			Y: y * moveScaleY,
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

	presets, err := app.issueGetCommand(context.Background(), deviceName, getPresetsCommand, cmd)
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

func (app *CameraManagementApp) gotoPreset(deviceName string, profile string, preset string) (dtosCommon.BaseResponse, error) {
	cmd := &ptz.GotoPreset{
		ProfileToken: (*onvif.ReferenceToken)(&profile),
		PresetToken:  (*onvif.ReferenceToken)(&preset),
	}

	return app.sendPutCommand(deviceName, gotoPresetCommand, cmd)
}

func (app *CameraManagementApp) sendPutCommand(deviceName string, commandName string, commandValue interface{}) (dtosCommon.BaseResponse, error) {
	app.lc.Infof("Sending PUT command: %s, %+v", commandName, commandValue)
	return app.service.CommandClient().IssueSetCommandByNameWithObject(context.Background(), deviceName, commandName,
		map[string]interface{}{
			// note: we are using the actual name of the command as the key
			commandName: commandValue,
		})
}

func (app *CameraManagementApp) getDevices() ([]dtos.Device, error) {
	response, err := app.service.DeviceClient().DevicesByServiceName(context.Background(), app.config.AppCustom.DeviceServiceName, 0, -1)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get devices for the device service %s", app.config.AppCustom.DeviceServiceName)
	}
	if len(response.Devices) <= 0 {
		return nil, errors.Errorf("no devices registered yet for the device service %s", app.config.AppCustom.DeviceServiceName)
	}

	// filter out the control-plane device
	var devices []dtos.Device
	for _, d := range response.Devices {
		if d.Name != d.ServiceName {
			devices = append(devices, d)
		}
	}
	return devices, nil
}
