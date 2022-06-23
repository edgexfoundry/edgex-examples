//
// Copyright (C) 2022 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/http/utils"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/responses"
	"io"
	"net/http"
	"net/url"
	"path"
)

const (
	maxBodySz = 1024 * 1024 * 10
)

func parseStreamUri(res *responses.EventResponse) (string, error) {
	val := res.Event.Readings[0].ObjectValue
	js, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	sr := StreamUriResponse{}
	if err := json.Unmarshal(js, &sr); err != nil {
		return "", err
	}
	return sr.MediaURI.URI, nil
}

func (app *CameraManagementApp) issueGetCommand(ctx context.Context, deviceName string, commandName string, jsonValue interface{}) (*responses.EventResponse, error) {
	jsonStr, err := json.Marshal(jsonValue)
	if err != nil {
		return nil, err
	}

	return app.service.CommandClient().IssueGetCommandByNameWithQueryParams(ctx, deviceName, commandName,
		map[string]string{"jsonObject": base64.URLEncoding.EncodeToString(jsonStr)})
}

func issuePostRequest(ctx context.Context, res interface{}, baseUrl string, reqPath string, jsonValue []byte) (err error) {
	return utils.PostRequest(ctx, &res, baseUrl, reqPath, jsonValue, common.ContentTypeJSON)
}

func issueGetRequest(ctx context.Context, res interface{}, baseUrl string, requestPath string) (err error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}
	requestPath = path.Join(u.Path, requestPath)
	return utils.GetRequest(ctx, &res, baseUrl, requestPath, nil)
}

func issueDeleteRequest(ctx context.Context, res interface{}, baseUrl string, requestPath string) (err error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return err
	}
	requestPath = path.Join(u.Path, requestPath)
	return utils.DeleteRequest(ctx, &res, baseUrl, requestPath)
}

func respondError(lc logger.LoggingClient, w http.ResponseWriter, statusCode int, errStr string) {
	lc.Error(errStr)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	if _, writeErr := w.Write([]byte(errStr)); writeErr != nil {
		lc.Error(writeErr.Error())
	}
}

func respondJson(lc logger.LoggingClient, w http.ResponseWriter, val interface{}) {
	lc.Debugf("response: %+v\n", val)
	b, err := json.Marshal(val)
	if err != nil {
		respondError(lc, w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal response body: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(b)
	if err != nil {
		// this will most likely fail since the w.Write failed, but it will at least print the error
		respondError(lc, w, http.StatusInternalServerError, fmt.Sprintf("failed to write response body: %v", err))
		return
	}
}

// extractJSONBody implements boilerplate code for unmarshalling a http.Request body.
//
// If the basic validation and unmarshalling is successful, this returns true.
// Otherwise, it writes a status to w and returns false;
// in that case, the caller should simply return.
//
// If the body is too large or cannot be unmarshalled into v, this writes an error.
func extractJSONBody(lc logger.LoggingClient, w http.ResponseWriter, r *http.Request, v interface{}) bool {
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxBodySz))
	if err != nil {
		respondError(lc, w, http.StatusInternalServerError,
			fmt.Sprintf("failed to read request body: %v", err))
		return false
	}

	if err = json.Unmarshal(data, v); err != nil {
		respondError(lc, w, http.StatusBadRequest,
			fmt.Sprintf("failed to unmarshal request body to target type %T: %v", v, err))
		return false
	}

	return true
}
