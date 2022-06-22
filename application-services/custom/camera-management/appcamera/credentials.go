//
// Copyright (C) 2022 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v2/config"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
)

const (
	CameraCredentials = "CameraCredentials"
	UsernameKey       = "username"
	PasswordKey       = "password"
)

// tryGetCredentials will attempt one time to get the camera credentials from the
// secret provider and return them, otherwise return an error.
func (app *CameraManagementApp) tryGetCredentials() (config.Credentials, errors.EdgeX) {
	secretData, err := app.service.GetSecret(CameraCredentials, UsernameKey, PasswordKey)
	if err != nil {
		return config.Credentials{}, errors.NewCommonEdgeXWrapper(err)
	}
	return config.Credentials{
		Username: secretData[UsernameKey],
		Password: secretData[PasswordKey],
	}, nil
}
