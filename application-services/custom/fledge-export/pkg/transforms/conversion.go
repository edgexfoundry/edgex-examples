package transforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type Conversion struct {
}

func NewConversion() Conversion {
	return Conversion{}
}

type FledgeReading struct {
	Timestamp string                 `json:"timestamp"`
	Asset     string                 `json:"asset"`
	Readings  map[string]interface{} `json:"readings,omitempty"`
}

func newFledgeReading(stamp int64, asset string) FledgeReading {
	tm := time.Unix(0, stamp*int64(time.Millisecond))
	reading := FledgeReading{Timestamp: tm.String(), Asset: asset}
	reading.Readings = make(map[string]interface{})
	return reading
}

// TransformToFledge ...
func (f Conversion) TransformToFledge(edgexcontext *appcontext.Context, params ...interface{}) (continuePipeline bool, stringType interface{}) {
	if len(params) < 1 {
		return false, errors.New("No Event Received")
	}

	edgexcontext.LoggingClient.Debug("Transforming to Fledge format")

	if event, ok := params[0].(models.Event); ok {
		payload := make([]FledgeReading, 1)
		fReading := newFledgeReading(event.Created, event.Device)

		for _, reading := range event.Readings {
			fReading.Readings[reading.Name] = reading.Value
		}
		payload[0] = fReading

		msg, err := json.Marshal(payload)
		if err != nil {
			return false, errors.New(fmt.Sprintf("Failed to transform Fledge data: %s", err))
		}

		edgexcontext.LoggingClient.Debug(fmt.Sprintf("Fledge Payload: %s", msg))

		return true, string(msg)
	}

	return false, errors.New("Unexpected type received")
}
