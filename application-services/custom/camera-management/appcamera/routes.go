//
// Copyright (C) 2022 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"fmt"
	dtosCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/gorilla/mux"
	"net/http"
	"path"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/pkg/errors"
)

const (
	left    = -1
	right   = 1
	up      = 1
	down    = -1
	zoomIn  = 1
	zoomOut = -1

	webUIDistDir = "./web-ui/dist"

	getCamerasPath = common.ApiBase + "/cameras"
	cameraApiBase  = getCamerasPath + "/{name}"

	getProfilesPath      = cameraApiBase + "/profiles"
	cameraProfileApiBase = getProfilesPath + "/{profile}"

	getPipelinesPath        = common.ApiBase + "/pipelines"
	allPipelineStatusesPath = getPipelinesPath + "/status/all"

	startPipelinePath  = cameraApiBase + "/pipeline/start"
	stopPipelinePath   = cameraApiBase + "/pipeline/stop"
	pipelineStatusPath = cameraApiBase + "/pipeline/status"

	ptzPath        = cameraProfileApiBase + "/ptz/{action}"
	getPresetsPath = cameraProfileApiBase + "/presets"
	gotoPresetPath = cameraProfileApiBase + "/presets/{preset}"
)

func (app *CameraManagementApp) addRoutes() error {
	if err := app.addRoute(
		startPipelinePath, http.MethodPost, app.startPipelineRoute); err != nil {
		return err
	}
	if err := app.addRoute(
		stopPipelinePath, http.MethodPost, app.stopPipelineRoute); err != nil {
		return err
	}
	if err := app.addRoute(
		pipelineStatusPath, http.MethodGet, app.pipelineStatusRoute); err != nil {
		return err
	}
	if err := app.addRoute(
		allPipelineStatusesPath, http.MethodGet, app.allPipelineStatusesRoute); err != nil {
		return err
	}
	if err := app.addRoute(
		getCamerasPath, http.MethodGet, app.getCamerasRoute); err != nil {
		return err
	}

	if err := app.addRoute(
		getProfilesPath, http.MethodGet, app.getProfilesRoute); err != nil {
		return err
	}

	if err := app.addRoute(
		getPipelinesPath, http.MethodGet, app.getPipelinesRoute); err != nil {
		return err
	}

	if err := app.addRoute(
		ptzPath, http.MethodPost, app.ptzRoute); err != nil {
		return err
	}

	if err := app.addRoute(
		getPresetsPath, http.MethodGet, app.getPresetsRoute); err != nil {
		return err
	}

	if err := app.addRoute(
		gotoPresetPath, http.MethodPost, app.gotoPresetRoute); err != nil {
		return err
	}

	app.fileServer = http.FileServer(http.Dir(webUIDistDir))
	// this is a bit of a hack to get refreshing working, as the path is /home
	if err := app.addRoute("/home", http.MethodGet, app.index); err != nil {
		return err
	}
	// all other routes will be forwarded to serving the web-ui
	if err := app.addRoute("/{path:.*}", http.MethodGet, app.serveWebUI); err != nil {
		return err
	}

	return nil
}

func (app *CameraManagementApp) addRoute(path, method string, f http.HandlerFunc) error {
	if err := app.service.AddRoute(path, f, method); err != nil {
		return errors.Wrapf(err, "failed to add route, path=%s, method=%s", path, method)
	}
	return nil
}

// Routes
func (app *CameraManagementApp) index(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, path.Join(webUIDistDir, "index.html"))
}

func (app *CameraManagementApp) serveWebUI(w http.ResponseWriter, req *http.Request) {
	app.fileServer.ServeHTTP(w, req)
}

func (app *CameraManagementApp) getPresetsRoute(w http.ResponseWriter, req *http.Request) {
	rv := mux.Vars(req)
	deviceName := rv["name"]
	profileToken := rv["profile"]

	presets, err := app.getPresets(deviceName, profileToken)
	if err != nil {
		respondError(app.lc, w, http.StatusInternalServerError,
			fmt.Sprintf("Failed to get presets: %v", err))
		return
	}

	respondJson(app.lc, w, presets)
}

func (app *CameraManagementApp) getProfilesRoute(w http.ResponseWriter, req *http.Request) {
	rv := mux.Vars(req)
	deviceName := rv["name"]

	pr, err := app.getProfiles(deviceName)
	if err != nil {
		respondError(app.lc, w, http.StatusInternalServerError,
			fmt.Sprintf("Failed to get profiles: %v", err))
		return
	}

	respondJson(app.lc, w, pr)
}

func (app *CameraManagementApp) getCamerasRoute(w http.ResponseWriter, _ *http.Request) {
	devices, err := app.getDevices()
	if err != nil {
		respondError(app.lc, w, http.StatusInternalServerError,
			fmt.Sprintf("Failed to get cameras list: %v", err))
		return
	}
	respondJson(app.lc, w, devices)
}

func (app *CameraManagementApp) getPipelinesRoute(w http.ResponseWriter, _ *http.Request) {
	pl, err := app.getPipelines()
	if err != nil {
		respondError(app.lc, w, http.StatusInternalServerError,
			fmt.Sprintf("Failed to get pipelines: %v", err))
		return
	}

	respondJson(app.lc, w, pl)
}

func (app *CameraManagementApp) startPipelineRoute(w http.ResponseWriter, req *http.Request) {
	rv := mux.Vars(req)
	deviceName := rv["name"]

	sr := StartPipelineRequest{}
	if !extractJSONBody(app.lc, w, req, &sr) {
		return
	}

	if app.isPipelineRunning(deviceName) {
		respondError(app.lc, w, http.StatusBadRequest, fmt.Sprintf("pipeline already running for camera: %s", deviceName))
		return
	}

	if err := app.startPipeline(deviceName, sr.ProfileToken, sr.PipelineName, sr.PipelineVersion); err != nil {
		respondError(app.lc, w, http.StatusInternalServerError, fmt.Sprintf("Failed to start pipeline: %v", err))
		return
	}
}

func (app *CameraManagementApp) pipelineStatusRoute(w http.ResponseWriter, req *http.Request) {
	rv := mux.Vars(req)
	deviceName := rv["name"]
	res, err := app.getPipelineStatus(deviceName)
	if err != nil {
		respondError(app.lc, w, http.StatusInternalServerError,
			fmt.Sprintf("failed to get pipeline status: %v", err))
		return
	}
	if res == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	respondJson(app.lc, w, res)
}

func (app *CameraManagementApp) allPipelineStatusesRoute(w http.ResponseWriter, _ *http.Request) {
	res, err := app.getAllPipelineStatuses()
	if err != nil {
		respondError(app.lc, w, http.StatusInternalServerError,
			fmt.Sprintf("failed to get all pipeline statuses: %v", err))
		return
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	respondJson(app.lc, w, res)
}

func (app *CameraManagementApp) stopPipelineRoute(w http.ResponseWriter, req *http.Request) {
	rv := mux.Vars(req)
	deviceName := rv["name"]
	if err := app.stopPipeline(deviceName); err != nil {
		respondError(app.lc, w, http.StatusInternalServerError,
			fmt.Sprintf("failed to stop pipeline: %v", err))
		return
	}
}

func (app *CameraManagementApp) gotoPresetRoute(w http.ResponseWriter, req *http.Request) {
	rv := mux.Vars(req)
	deviceName := rv["name"]
	profileToken := rv["profile"]
	preset := rv["preset"]

	var res dtosCommon.BaseResponse
	var err error

	res, err = app.gotoPreset(deviceName, profileToken, preset)
	if err != nil {
		respondError(app.lc, w, http.StatusInternalServerError,
			fmt.Sprintf("Failed to do gotoPreset: %v", err))
		return
	}
	respondJson(app.lc, w, res)
}

func (app *CameraManagementApp) ptzRoute(w http.ResponseWriter, req *http.Request) {
	rv := mux.Vars(req)
	deviceName := rv["name"]
	profileToken := rv["profile"]
	action := rv["action"]

	var res dtosCommon.BaseResponse
	var err error

	switch action {
	case "left":
		res, err = app.doPTZ(deviceName, profileToken, left, 0, 0)
	case "right":
		res, err = app.doPTZ(deviceName, profileToken, right, 0, 0)

	case "up":
		res, err = app.doPTZ(deviceName, profileToken, 0, up, 0)
	case "up-left":
		res, err = app.doPTZ(deviceName, profileToken, left, up, 0)
	case "up-right":
		res, err = app.doPTZ(deviceName, profileToken, right, up, 0)

	case "down":
		res, err = app.doPTZ(deviceName, profileToken, 0, down, 0)
	case "down-left":
		res, err = app.doPTZ(deviceName, profileToken, left, down, 0)
	case "down-right":
		res, err = app.doPTZ(deviceName, profileToken, right, down, 0)

	case "zoom-in":
		res, err = app.doPTZ(deviceName, profileToken, 0, 0, zoomIn)
	case "zoom-out":
		res, err = app.doPTZ(deviceName, profileToken, 0, 0, zoomOut)

	default:
		err = fmt.Errorf("unknown ptz action: %s", action)
	}

	if err != nil {
		respondError(app.lc, w, http.StatusInternalServerError,
			fmt.Sprintf("Failed to do ptz: %v", err))
		return
	}
	_, err = w.Write([]byte(res.Message))
	if err != nil {
		app.lc.Error(err.Error())
	}
}
