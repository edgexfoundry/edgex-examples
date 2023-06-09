//
// Copyright (C) 2022-2023 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0

package appcamera

import (
	"github.com/edgexfoundry/go-mod-bootstrap/v3/bootstrap/secret"
	"github.com/edgexfoundry/go-mod-bootstrap/v3/config"
	"github.com/edgexfoundry/go-mod-core-contracts/v3/errors"
)

const (
	rtspauth  = "rtspauth"
	onvifauth = "onvifauth"
)

// tryGetCredentials will attempt one time to get the camera credentials from the
// secret provider and return them, otherwise return an error.
func (app *CameraManagementApp) tryGetCredentials(secretName string) (config.Credentials, errors.EdgeX) {
	secretData, err := app.service.SecretProvider().GetSecret(secretName, secret.UsernameKey, secret.PasswordKey)
	if err != nil {
		return config.Credentials{}, errors.NewCommonEdgeXWrapper(err)
	}
	return config.Credentials{
		Username: secretData[secret.UsernameKey],
		Password: secretData[secret.PasswordKey],
	}, nil
}
