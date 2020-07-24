package transforms

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	sdkTransforms "github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

const (
	awsIoTMQTTHost           = "AwsIoTMQTTHost"
	awsIoTMQTTPort           = "AwsIoTMQTTPort"
	awsIoTThingName          = "awsIoTThingName"
	awsIoTRootCAFilename     = "CaCertPath"
	awsIoTCertFilename       = "MQTTCert"
	awsIoTPrivateKeyFilename = "MQTTKey"
	awsSkipCertVerify        = "SkipCertVerify"
	awsPersistOnError        = "PersistOnError"
	awsDeviceNames           = "AwsDeviceNames"
	user                     = "someUser"
	topic                    = "topic"
)

var log logger.LoggingClient

// AWSMQTTConfig holds AWS IoT specific information
type AWSMQTTConfig struct {
	MQTTConfig     sdkTransforms.MqttConfig
	IoTHost        string
	IoTPort        string
	IoTDevice      string
	IoTTopic       string
	DeviceNames    string
	PersistOnError bool
	KeyCertPair    *sdkTransforms.KeyCertPair
}

func getNewClient(skipVerify bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	return &http.Client{Timeout: 10 * time.Second, Transport: tr}
}

func getAppSetting(settings map[string]string, name string) string {
	value, ok := settings[name]

	if ok {
		log.Debug(value)
		return value
	}
	log.Error(fmt.Sprintf("ApplicationName application setting %s not found", name))
	return ""

}

// LoadAWSMQTTConfig Loads the mqtt configuration necessary to connect to AWS cloud
func LoadAWSMQTTConfig(sdk *appsdk.AppFunctionsSDK) (*AWSMQTTConfig, error) {
	if sdk == nil {
		return nil, errors.New("Invalid AppFunctionsSDK")
	}

	log = sdk.LoggingClient

	var ioTHost, iotPort, iotDevice, mqttCert, mqttKey, ioTTopic, deviceNames string
	var skipCertVerify, persistOnError bool
	var errSkip, errPersist error

	appSettings := sdk.ApplicationSettings()
	if appSettings != nil {
		ioTHost = getAppSetting(appSettings, awsIoTMQTTHost)
		iotPort = getAppSetting(appSettings, awsIoTMQTTPort)
		iotDevice = getAppSetting(appSettings, awsIoTThingName)
		mqttCert = getAppSetting(appSettings, awsIoTCertFilename)
		mqttKey = getAppSetting(appSettings, awsIoTPrivateKeyFilename)
		ioTTopic = getAppSetting(appSettings, topic)
		deviceNames = getAppSetting(appSettings, awsDeviceNames)
		skipCertVerify, errSkip = strconv.ParseBool(getAppSetting(appSettings, awsSkipCertVerify))
		persistOnError, errPersist = strconv.ParseBool(getAppSetting(appSettings, awsPersistOnError))
		if errSkip != nil {
			log.Error("Unable to parse " + awsSkipCertVerify + " value")
		}
		if errPersist != nil {
			log.Error("Unable to parse " + awsPersistOnError + " value")
		}
	} else {
		return nil, errors.New("No application-specific settings found")
	}

	config := AWSMQTTConfig{}

	config.IoTHost = ioTHost
	config.IoTPort = iotPort
	config.IoTDevice = iotDevice
	config.IoTTopic = ioTTopic
	config.DeviceNames = deviceNames
	config.PersistOnError = persistOnError

	pair := &sdkTransforms.KeyCertPair{
		KeyFile:  mqttKey,
		CertFile: mqttCert,
	}

	mqttConfig := sdkTransforms.MqttConfig{
		SkipCertVerify: skipCertVerify,
	}

	log.Debug(fmt.Sprintf("Read SkipCertVerify from configuration: %t", config.MQTTConfig.SkipCertVerify))
	log.Debug(fmt.Sprintf("Read PersistOnError from configuration: %t", config.PersistOnError))

	config.KeyCertPair = pair
	config.MQTTConfig = mqttConfig

	return &config, nil
}

// NewAWSMQTTSender return a mqtt sender capable of sending the event's value to the given MQTT broker
func NewAWSMQTTSender(logging logger.LoggingClient, config *AWSMQTTConfig) *sdkTransforms.MQTTSender {

	logging.Debug(config.IoTTopic)

	port, err := strconv.Atoi(config.IoTPort)
	if err != nil {
		// falling back to default AWS IoT port
		port = 8883
	}

	addressable := models.Addressable{
		Address:   config.IoTHost,
		Port:      port,
		Protocol:  "tls",
		Publisher: config.IoTDevice,
		User:      "",
		Password:  "",
		Topic:     config.IoTTopic,
	}

	mqttSender := sdkTransforms.NewMQTTSender(logging, addressable, config.KeyCertPair, config.MQTTConfig, config.PersistOnError)

	return mqttSender
}
