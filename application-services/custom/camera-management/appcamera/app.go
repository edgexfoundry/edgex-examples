//
// Copyright (C) 2022-2023 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/pkg/errors"
	"net/http"
	"sync"
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
		return errors.Wrap(err, "failed to load custom configuration")
	}

	if err := app.addRoutes(); err != nil {
		return err
	}

	if err := app.queryAllPipelineStatuses(); err != nil {
		// do not exit, just log
		app.lc.Errorf("Unable to query EVAM pipeline statuses. Is EVAM running? %s", err.Error())
	}

	if err := app.service.MakeItRun(); err != nil {
		return errors.Wrap(err, "failed to run pipeline")
	}

	return nil
}
