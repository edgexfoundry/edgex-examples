package transforms

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dghubble/sling"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	sdkTransforms "github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

const (
	serviceKey         = "AzureExport"
	appConfigIoTHub    = "IoTHub"
	appConfigIoTDevice = "IoTDevice"
	appConfigMQTTCert  = "MQTTCert"
	appConfigMQTTKey   = "MQTTKey"
	appConfigTokenPath = "TokenPath"
	appConfigVaultHost = "VaultHost"
	appConfigVaultPort = "VaultPort"
	appConfigCertPath  = "CertPath"
	mqttPort           = 8883
	vaultToken         = "X-Vault-Token"
)

type AzureMQTTConfig struct {
	MQTTConfig  *sdkTransforms.MqttConfig
	IoTHub      string
	IoTDevice   string
	KeyCertPair *sdkTransforms.KeyCertPair
}

// global logger
var log logger.LoggingClient

type certCollect struct {
	Pair certPair `json:"data"`
}

type certPair struct {
	Cert string `json:"cert,omitempty"`
	Key  string `json:"key,omitempty"`
}

type auth struct {
	Token string `json:"root_token"`
}

func getAppSetting(settings map[string]string, name string) string {
	value, ok := settings[name]

	if ok {
		return value
	} else {
		log.Error(fmt.Sprintf("ApplicationName application setting %s not found", name))
		return ""
	}
}

func retrieveKeyCertPair(tokenPath string, vaultHost string, vaultPort string, certPath string) (*sdkTransforms.KeyCertPair, error) {
	a := auth{}
	content, err := ioutil.ReadFile(tokenPath)

	if err != nil {
		log.Error("Failed to read token file", err.Error())
		return nil, err
	}

	err = json.Unmarshal(content, &a)

	// we have a.Token here
	s := sling.New().Set(vaultToken, a.Token)
	vaultUrl := fmt.Sprintf("https://%s:%s/", vaultHost, vaultPort)
	req, err := s.New().Base(vaultUrl).Get(certPath).Request()

	if err != nil {
		log.Error("Failed to create request", err.Error())
		return nil, err
	}

	res, err := getNewClient(true).Do(req)

	if err != nil {
		log.Error("Client request failed", err.Error())
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unsuccessful HTTP request to Vault, status code is %d", res.StatusCode)
	}

	cc := certCollect{}
	err = json.NewDecoder(res.Body).Decode(&cc)

	if err != nil || len(cc.Pair.Key) == 0 || len(cc.Pair.Cert) == 0 {
		return nil, errors.New("Failed to load key/cert pair from Vault")
	}

	pair := &sdkTransforms.KeyCertPair{
		KeyPEMBlock:  []byte(cc.Pair.Key),
		CertPEMBlock: []byte(cc.Pair.Cert),
	}

	log.Info("Successfully loaded key/cert pair from Vault")

	return pair, nil
}

func getNewClient(skipVerify bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	return &http.Client{Timeout: 10 * time.Second, Transport: tr}
}

func LoadAzureMQTTConfig(sdk *appsdk.AppFunctionsSDK) (*AzureMQTTConfig, error) {
	if sdk == nil {
		return nil, errors.New("Invalid AppFunctionsSDK")
	}

	log = sdk.LoggingClient

	appSettings := sdk.ApplicationSettings()

	if appSettings == nil {
		return nil, errors.New("No application-specific settings found")
	}

	iotHub := getAppSetting(appSettings, appConfigIoTHub)
	iotDevice := getAppSetting(appSettings, appConfigIoTDevice)
	mqttCert := getAppSetting(appSettings, appConfigMQTTCert)
	mqttKey := getAppSetting(appSettings, appConfigMQTTKey)
	tokenPath := getAppSetting(appSettings, appConfigTokenPath)
	vaultHost := getAppSetting(appSettings, appConfigVaultHost)
	vaultPort := getAppSetting(appSettings, appConfigVaultPort)
	certPath := getAppSetting(appSettings, appConfigCertPath)

	if len(iotHub) == 0 || len(iotDevice) == 0 {
		return nil, errors.New("Required configurations " + appConfigIoTHub + " or " + appConfigIoTDevice + " are missing")
	}

	config := AzureMQTTConfig{}

	config.IoTHub = iotHub
	config.IoTDevice = iotDevice
	config.MQTTConfig = &sdkTransforms.MqttConfig{}

	// Retrieve key/cert pair from Vault
	pair, err := retrieveKeyCertPair(tokenPath, vaultHost, vaultPort, certPath)

	// Fall back to local key/cert files
	if err != nil {
		log.Error(fmt.Sprintf("Failed to load key/cert from Vault (%v), use local key/cert files instead", err))

		pair = &sdkTransforms.KeyCertPair{
			KeyFile:  mqttKey,
			CertFile: mqttCert,
		}
	}

	config.KeyCertPair = pair

	return &config, nil
}

func NewAzureMQTTSender(logging logger.LoggingClient, config *AzureMQTTConfig) *sdkTransforms.MQTTSender {
	// Generate Azure-specific host, user amd topic
	host := fmt.Sprintf("%s.azure-devices.net", config.IoTHub)
	user := fmt.Sprintf("%s/%s/?api-version=2018-06-30", host, config.IoTDevice)
	topic := fmt.Sprintf("devices/%s/messages/events/", config.IoTDevice)

	addressable := models.Addressable{
		Address:   host,
		Port:      mqttPort,
		Protocol:  "tls",
		Publisher: config.IoTDevice, // must be the same as the device name
		User:      user,
		Password:  "",
		Topic:     topic,
	}

	mqttSender := sdkTransforms.NewMQTTSender(logging, addressable, config.KeyCertPair, *config.MQTTConfig, false)

	return mqttSender
}
