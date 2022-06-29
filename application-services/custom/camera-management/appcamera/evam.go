//
// Copyright (C) 2022 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"path"
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
	// Profile is the ProfileToken for the specific stream
	Profile string `json:"profile,omitempty"`
}

type PipelineInfoStatus struct {
	Camera string       `json:"camera"`
	Info   PipelineInfo `json:"info"`
	Status interface{}  `json:"status"`
}

func (info PipelineInfo) getPipelineUrl(evamBaseUrl string) (string, error) {
	uri, err := url.Parse(evamBaseUrl)
	if err != nil {
		return "", err
	}
	uri.Path = path.Join(uri.Path, "pipelines", info.Name, info.Version, info.Id)
	return uri.String(), nil
}

func (app *CameraManagementApp) addPipelineInfo(camera string, info PipelineInfo) error {
	app.pipelinesMutex.Lock()
	defer app.pipelinesMutex.Unlock()
	if _, exists := app.pipelinesMap[camera]; !exists {
		app.pipelinesMap[camera] = info
		return nil
	}
	return errors.Errorf("Pipeline already running for device %v", camera)
}

func (app *CameraManagementApp) deletePipelineInfo(camera string) {
	app.pipelinesMutex.Lock()
	defer app.pipelinesMutex.Unlock()
	if _, exists := app.pipelinesMap[camera]; exists {
		delete(app.pipelinesMap, camera)
	}
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

func (app *CameraManagementApp) startPipeline(deviceName, profileToken, pipelineName, pipelineVersion string) error {
	streamUri, err := app.queryStreamUri(deviceName, profileToken)
	if err != nil {
		return err
	}
	app.lc.Infof("Received stream uri for the device %s: %s", deviceName, streamUri)

	body, err := app.createPipelineRequestBody(streamUri, deviceName)
	if err != nil {
		return errors.Wrapf(err, "failed to create DLStreamer pipeline request body")
	}

	info := PipelineInfo{
		Name:    pipelineName,
		Version: pipelineVersion,
		Profile: profileToken,
	}
	var res interface{}
	baseUrl, err := url.Parse(app.config.AppCustom.EvamBaseUrl)
	if err != nil {
		return err
	}
	reqPath := path.Join("pipelines", info.Name, info.Version)

	if err = issuePostRequest(context.Background(), &res, baseUrl.String(), reqPath, body); err != nil {
		return errors.Wrap(err, "POST request to start EVAM pipeline failed")
	}
	info.Id = fmt.Sprintf("%v", res)

	if err = app.addPipelineInfo(deviceName, info); err != nil {
		return err
	}

	app.lc.Infof("Successfully started EVAM pipeline for the device %s", deviceName)
	app.lc.Infof("View inference results at 'rtsp://<SYSTEM_IP_ADDRESS>:8554/%s'", deviceName)

	return nil
}

func (app *CameraManagementApp) stopPipeline(deviceName string) error {
	var res interface{}

	if info, found := app.getPipelineInfo(deviceName); found {
		pipelineUrl, err := info.getPipelineUrl(app.config.AppCustom.EvamBaseUrl)
		if err != nil {
			return err
		}
		if err := issueDeleteRequest(context.Background(), &res, pipelineUrl, ""); err != nil {
			return errors.Wrap(err, "DELETE request to stop EVAM pipeline failed")
		}
		app.deletePipelineInfo(deviceName)
		app.lc.Infof("Successfully stopped EVAM pipeline for the device %s", deviceName)
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
		pipelineUrl, err := info.getPipelineUrl(app.config.AppCustom.EvamBaseUrl)
		if err != nil {
			return nil, err
		}
		if err := issueGetRequest(context.Background(), &res, pipelineUrl, "status"); err != nil {
			return nil, errors.Wrap(err, "GET request to query EVAM pipeline failed")
		}
		return res, nil
	}

	return nil, nil
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
		pipelineUrl, err := data.Info.getPipelineUrl(app.config.AppCustom.EvamBaseUrl)
		if err != nil {
			return nil, err
		}
		if err = issueGetRequest(context.Background(), &data.Status, pipelineUrl, "status"); err != nil {
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
