// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// Package driver this package provides a simple example implementation of
// ProtocolDriver interface.
//
package driver

import (
	"fmt"
	"time"
	dsModels "github.com/edgexfoundry/device-sdk-go/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

type Driver struct {
	lc      logger.LoggingClient
	asyncCh chan<- *dsModels.AsyncValues
	gpiodevice *GPIODev
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues, deviceCh chan<- []dsModels.DiscoveredDevice) error {
	s.lc = lc
	s.asyncCh = asyncCh

	s.gpiodevice = NewGPIODev(lc)
	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *Driver) HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {

	s.lc.Debug(fmt.Sprintf("protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes))

	now := time.Now().UnixNano()

	for _, req := range reqs {
		switch req.DeviceResourceName {
		case "direction":
			{
				rxbuf, err := s.gpiodevice.GetDirection()
				if err != nil {
					return nil, err
				}
				cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, rxbuf)
				res = append(res, cv)
			}
		case "value":
			{
				rxbuf, err := s.gpiodevice.GetGPIO()
				if err != nil {
					return nil, err
				}
				cv := dsModels.NewStringValue(reqs[0].DeviceResourceName, now, rxbuf)
				res = append(res, cv)
			}
		}
	}

	return res, nil
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *Driver) HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {
	s.lc.Info(fmt.Sprintf("Driver.HandleWriteCommands: protocols: %v, resource: %v, parameters: %v", protocols, reqs[0].DeviceResourceName, params))
	for _, param := range params {
		var err error
		switch param.DeviceResourceName {
		case "export":
			{
				var num int16
				num, err = param.Int16Value()
				if err != nil {
					return err
				}
				err = s.gpiodevice.ExportGPIO(int(num))
				if err != nil {
					return err
				}
			}
		case "unexport":
			{
				var num int16
				num, err = param.Int16Value()
				if err != nil {
					return err
				}
				err = s.gpiodevice.UnexportGPIO(int(num))
				if err != nil {
					return err
				}
			}
		case "direction":
			{
				var direction string
				direction, err = param.StringValue()
				if err != nil {
					return err
				}
				err = s.gpiodevice.SetDirection(direction)
				if err != nil {
					return err
				}
			}
		case "value":
			{
				var value int8
				value, err = param.Int8Value()
				if err != nil {
					return err
				}
				err = s.gpiodevice.SetGPIO(int(value))
				if err != nil {
					return err
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *Driver) Stop(force bool) error {
	// Then Logging Client might not be initialized
	if s.lc != nil {
		s.lc.Debug(fmt.Sprintf("Driver.Stop called: force=%v", force))
	}
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (s *Driver) AddDevice(deviceName string, protocols map[string]contract.ProtocolProperties, adminState contract.AdminState) error {
	s.lc.Debug(fmt.Sprintf("a new Device is added: %s", deviceName))
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (s *Driver) UpdateDevice(deviceName string, protocols map[string]contract.ProtocolProperties, adminState contract.AdminState) error {
	s.lc.Debug(fmt.Sprintf("Device %s is updated", deviceName))
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (s *Driver) RemoveDevice(deviceName string, protocols map[string]contract.ProtocolProperties) error {
	s.lc.Debug(fmt.Sprintf("Device %s is removed", deviceName))
	return nil
}
