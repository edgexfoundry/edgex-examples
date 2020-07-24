# App-functions-aws

App-functions-aws is a sample provided to let the user send mqtt messages to Amazon AWS using SSL

# Configuration

Client certificates (crt and key) can be copied in the folder certs

The address of the AWS host and the path to the certificates can be found in the configuration.toml:

[ApplicationSettings]

AwsIoTMQTTHost      = "***.iot.us-west-2.amazonaws.com"

AwsIoTMQTTPort      = "8883"


MQTTCert            = "./certs/user.client.crt"

MQTTKey             = "./certs/user.client.key"

SkipCertVerify	    = "false"

# Important note

If SkipCertVerify is true, TLS accepts any certificate
presented by the server and any host name in that certificate.
In this mode, TLS is susceptible to man-in-the-middle attacks.
This should be used only for testing. [golang reference](https://golang.org/pkg/crypto/tls/)
