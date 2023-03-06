//
// Copyright (C) 2022-2023 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IOTechSystems/onvif/media"
	"github.com/pkg/errors"
	"net/url"
	"path"
)

const (
	Aborted = "ABORTED"

	evamRtspPort = 8555
)

// Note: DLStreamer Pipeline Server / EVAM APIs can be viewed here:
// https://github.com/dlstreamer/pipeline-server/blob/master/docs/restful_microservice_interfaces.md

type PipelineInfo struct {
	// Id is the instance_id assigned by the pipeline server once it is started
	Id string `json:"id,omitempty"`
	// Name is the first part of the pipeline's full name. In the case of 'object_detection/person_vehicle_bike'
	// the name is 'object_detection'
	Name string `json:"name,omitempty"`
	// Version is the second part of the pipeline's full name. In the case of 'object_detection/person_vehicle_bike'
	// the version is 'person_vehicle_bike'
	Version string `json:"version,omitempty"`
}

type PipelineInfoStatus struct {
	Camera string       `json:"camera"`
	Info   PipelineInfo `json:"info"`
	Status interface{}  `json:"status"`
}

func (app *CameraManagementApp) addPipelineInfo(camera string, info PipelineInfo) error {
	app.pipelinesMutex.Lock()
	defer app.pipelinesMutex.Unlock()
	if _, exists := app.pipelinesMap[camera]; !exists {
		app.pipelinesMap[camera] = info
		return nil
	}
	return errors.Errorf("pipeline already running for device %v", camera)
}

func (app *CameraManagementApp) deletePipelineInfo(camera string) {
	app.pipelinesMutex.Lock()
	defer app.pipelinesMutex.Unlock()
	delete(app.pipelinesMap, camera)
}

func (app *CameraManagementApp) isPipelineRunning(camera string) bool {
	app.pipelinesMutex.RLock()
	defer app.pipelinesMutex.RUnlock()
	_, exists := app.pipelinesMap[camera]
	return exists
}

func (app *CameraManagementApp) getPipelineInfo(camera string) (PipelineInfo, bool) {
	app.pipelinesMutex.RLock()
	defer app.pipelinesMutex.RUnlock()
	// note: go will not let us return this lookup directly since it is overloaded
	val, found := app.pipelinesMap[camera]
	return val, found
}

func (app *CameraManagementApp) queryStreamUri(deviceName string, sr StartPipelineRequest) (string, error) {
	if sr.USB != nil {
		return app.getUSBStreamUri(deviceName)
	} else if sr.Onvif != nil {
		return app.getOnvifStreamUri(deviceName, sr.Onvif.ProfileToken)
	}
	return "", errors.New("missing required stream configuration")
}

func (app *CameraManagementApp) getOnvifStreamUri(deviceName string, profileToken string) (string, error) {
	req := StreamUriRequest{ProfileToken: profileToken}
	resp := media.GetStreamUriResponse{}
	err := app.issueGetCommandWithJsonForResponse(context.Background(), deviceName, streamUriCommand, req, &resp)
	if err != nil {
		return "", err
	}
	return string(resp.MediaUri.Uri), nil
}

func (app *CameraManagementApp) getUSBStreamUri(deviceName string) (string, error) {
	cmdResponse, err := app.issueGetCommand(context.Background(), deviceName, usbStreamUriCommand)
	if err != nil {
		return "", errors.Wrapf(err, "failed to issue get StreamUri command")
	}
	return cmdResponse.Event.Readings[0].Value, nil
}

func (app *CameraManagementApp) startPipeline(deviceName string, sr StartPipelineRequest) error {
	streamUri, err := app.queryStreamUri(deviceName, sr)
	if err != nil {
		return err
	}
	app.lc.Infof("Received stream uri for the device %s: %s", deviceName, streamUri)

	// if device is usb camera, start streaming first
	if sr.USB != nil {
		_, err := app.startStreaming(deviceName, *sr.USB)
		if err != nil {
			return errors.Wrapf(err, "failed to start streaming usb camera %s", deviceName)
		}
	}

	body, err := app.createPipelineRequestBody(streamUri, deviceName)
	if err != nil {
		return errors.Wrapf(err, "failed to create DLStreamer pipeline request body")
	}

	info := PipelineInfo{
		Name:    sr.PipelineName,
		Version: sr.PipelineVersion,
	}
	var res interface{}
	baseUrl, err := url.Parse(app.config.AppCustom.EvamBaseUrl)
	if err != nil {
		return err
	}
	reqPath := path.Join("pipelines", info.Name, info.Version)

	if err = issuePostRequest(context.Background(), &res, baseUrl.String(), reqPath, body); err != nil {
		err = errors.Wrap(err, "POST request to start EVAM pipeline failed")
		// if we started the streaming on usb camera, we need to stop it
		if sr.USB != nil {
			if _, err2 := app.stopStreaming(deviceName); err2 != nil {
				err = errors.Wrapf(err, "failed to stop streaming usb camera %s", deviceName)
			}
		}
		return err
	}
	info.Id = fmt.Sprintf("%v", res)

	if err = app.addPipelineInfo(deviceName, info); err != nil {
		return err
	}

	app.lc.Infof("Successfully started EVAM pipeline for the device %s", deviceName)
	app.lc.Infof("View inference results at 'rtsp://<SYSTEM_IP_ADDRESS>:%d/%s'", evamRtspPort, deviceName)

	return nil
}

func (app *CameraManagementApp) stopPipeline(deviceName string, id string) error {
	var res interface{}

	if err := issueDeleteRequest(context.Background(), &res, app.config.AppCustom.EvamBaseUrl, path.Join("pipelines", id)); err != nil {
		return errors.Wrap(err, "DELETE request to stop EVAM pipeline failed")
	}
	app.lc.Infof("Successfully stopped EVAM pipeline for the device %s", deviceName)

	if info, found := app.getPipelineInfo(deviceName); found && info.Id == id {
		// only delete pipeline if it matches the id
		app.deletePipelineInfo(deviceName)
	}

	return nil
}

func (app *CameraManagementApp) createPipelineRequestBody(streamUri string, deviceName string) ([]byte, error) {
	uri, err := url.Parse(streamUri)
	if err != nil {
		return nil, err
	}

	if creds, err := app.tryGetCredentials(); err != nil {
		app.lc.Warnf("Error retrieving %s secret from the SecretStore: %s", CameraCredentials, err.Error())
	} else {
		uri.User = url.UserPassword(creds.Username, creds.Password)
	}

	pipelineData := PipelineRequest{
		Source: Source{
			URI:  uri.String(),
			Type: "uri",
		},
		Destination: Destination{
			Metadata: Metadata{
				Type:  "mqtt",
				Host:  app.config.AppCustom.MqttAddress,
				Topic: app.config.AppCustom.MqttTopic,
			},
			Frame: Frame{
				Type: "rtsp",
				Path: deviceName,
			},
		},
	}

	pipeline, err := json.Marshal(pipelineData)
	if err != nil {
		return pipeline, err
	}

	return pipeline, nil
}

func (app *CameraManagementApp) getPipelineStatus(deviceName string) (interface{}, error) {
	if info, found := app.getPipelineInfo(deviceName); found {
		var res interface{}
		if err := issueGetRequest(context.Background(), &res, app.config.AppCustom.EvamBaseUrl, path.Join("pipelines", "status", info.Id)); err != nil {
			return nil, errors.Wrap(err, "GET request to query EVAM pipeline status failed")
		}
		return res, nil
	}

	return nil, nil
}

// queryAllPipelineStatuses queries EVAM for all pipeline statuses, attempts to link them to devices, and then
// insert them into the pipeline map.
func (app *CameraManagementApp) queryAllPipelineStatuses() error {
	var statuses []PipelineStatus
	if err := issueGetRequest(context.Background(), &statuses, app.config.AppCustom.EvamBaseUrl, path.Join("pipelines", "status")); err != nil {
		return errors.Wrap(err, "GET request to query EVAM pipeline statuses failed")
	}

	for _, status := range statuses {
		if status.State == Aborted {
			continue // ignore stopped pipelines
		}

		var resp PipelineInformationResponse
		if err := issueGetRequest(context.Background(), &resp, app.config.AppCustom.EvamBaseUrl, path.Join("pipelines", status.Id)); err != nil {
			app.lc.Errorf("GET request to query EVAM pipeline %s info failed: %s", status.Id, err.Error())
			continue
		}

		// assume the destination streaming path is the camera name, because that is how it is when we created the pipeline instance
		deviceName := resp.Request.Destination.Frame.Path
		// ensure the device actually exists
		if _, err := app.getDeviceByName(deviceName); err != nil {
			app.lc.Warnf("Unable to determine device name from EVAM pipeline %s: %s", status.Id, err.Error())
			continue
		}

		info := PipelineInfo{
			Id:      resp.Id,
			Name:    resp.Request.Pipeline.Name,
			Version: resp.Request.Pipeline.Version,
		}
		app.deletePipelineInfo(deviceName) // delete the info in case it already exists
		// add pipeline info to map to ensure we track it
		if err := app.addPipelineInfo(deviceName, info); err != nil {
			app.lc.Errorf("Error adding pipeline info to map: %s", err.Error())
		}
	}
	return nil
}

func (app *CameraManagementApp) getAllPipelineStatuses() (map[string]PipelineInfoStatus, error) {
	response := make(map[string]PipelineInfoStatus)
	// pre-create the response object using a read lock to minimize the time we hold the lock
	app.pipelinesMutex.RLock()
	for camera, info := range app.pipelinesMap {
		response[camera] = PipelineInfoStatus{
			Camera: camera,
			Info:   info,
		}
	}
	app.pipelinesMutex.RUnlock()

	// loop through the partially filled response map to fill in the missing data. we do not need to hold the lock here.
	for camera, data := range response {
		if err := issueGetRequest(context.Background(), &data.Status, app.config.AppCustom.EvamBaseUrl, path.Join("pipelines", "status", data.Info.Id)); err != nil {
			return nil, errors.Wrap(err, "GET request to query EVAM pipeline failed")
		}
		// overwrite the changed result in the map
		response[camera] = data
	}

	return response, nil
}

func (app *CameraManagementApp) getPipelines() (interface{}, error) {
	var res interface{}
	if err := issueGetRequest(context.Background(), &res, app.config.AppCustom.EvamBaseUrl, "pipelines"); err != nil {
		return nil, errors.Wrap(err, "GET request to query all EVAM pipelines failed")
	}
	return res, nil
}
