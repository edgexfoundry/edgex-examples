//
// Copyright (C) 2022-2023 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"net/http"
	"sync"
	"time"

	"github.com/edgexfoundry/app-functions-sdk-go/v3/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
)

const (
	maxRetries = 10
)

type CameraManagementApp struct {
	service        interfaces.ApplicationService
	lc             logger.LoggingClient
	config         *ServiceConfig
	pipelinesMap   map[string]PipelineInfo
	pipelinesMutex sync.RWMutex
	ptzRangeMap    map[string]PTZRange
	ptzRangeMutex  sync.RWMutex
	fileServer     http.Handler
}

func NewCameraManagementApp(service interfaces.ApplicationService) *CameraManagementApp {
	return &CameraManagementApp{
		service:      service,
		lc:           service.LoggingClient(),
		config:       &ServiceConfig{},
		pipelinesMap: make(map[string]PipelineInfo),
		ptzRangeMap:  make(map[string]PTZRange),
	}
}

func (app *CameraManagementApp) Run() error {
	if err := app.service.LoadCustomConfig(app.config, "AppCustom"); err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to load custom configuration", err)
	}

	var err error
	for i := 0; i < maxRetries; i++ {
		app.lc.Infof("Querying EVAM pipeline statuses.")
		if err = app.queryAllPipelineStatuses(); err != nil {
			app.lc.Errorf("Unable to query EVAM pipeline statuses. Is EVAM running? %s", err.Error())
			time.Sleep(time.Second)
		} else {
			break // no error, so lets continue
		}
	}
	if err != nil {
		app.lc.Errorf("Unable to query EVAM pipeline statuses after %d tries. .")
		return err // exit. we do not want to run if evam is not accessible
	}

	if err := app.addRoutes(); err != nil {
		return err
	}

	// Subscribe to events.
	if err := app.service.SetDefaultFunctionsPipeline(app.processEdgeXDeviceSystemEvent); err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to set default pipeline to processEdgeXEvent", err)
	}

	devices, err := app.getAllDevices()
	if err != nil {
		app.lc.Errorf("no devices found: %s", err.Error())
	} else {
		for _, device := range devices {
			if err = app.startDefaultPipeline(device); err != nil {
				app.lc.Errorf("Error starting default pipeline for %s, %v", device.Name, err)
			}
		}
	}

	if err = app.service.Run(); err != nil {
		return errors.NewCommonEdgeX(errors.KindServerError, "failed to run pipeline", err)
	}

	return nil
}
