module custom-trigger

go 1.15

require (
	github.com/cloudevents/sdk-go v1.1.2
	github.com/edgexfoundry/app-functions-sdk-go/v2 v2.0.0-dev.27
	github.com/edgexfoundry/go-mod-bootstrap/v2 v2.0.0-dev.14
	github.com/edgexfoundry/go-mod-core-contracts/v2 v2.0.0-dev.42
	github.com/edgexfoundry/go-mod-messaging/v2 v2.0.0-dev.3
	github.com/stretchr/testify v1.7.0
)

replace github.com/edgexfoundry/app-functions-sdk-go/v2 => ../../../../app-functions-sdk-go
