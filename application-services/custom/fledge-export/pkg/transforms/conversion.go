package transforms

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
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
func (f Conversion) TransformToFledge(ctx interfaces.AppFunctionContext, data interface{}) (continuePipeline bool, stringType interface{}) {
	lc := ctx.LoggingClient()

	lc.Debug("Transforming to Fledge format")

	if data == nil {
		return false, errors.New("TransformToFledge: No data received")
	}

	event, ok := data.(dtos.Event)
	if !ok {
		return false, errors.New("TransformToFledge: didn't receive expect Event type")
	}

	payload := make([]FledgeReading, 1)
	fReading := newFledgeReading(event.Origin, event.DeviceName)

	for _, reading := range event.Readings {
		fReading.Readings[reading.ResourceName] = reading.Value
	}
	payload[0] = fReading

	msg, err := json.Marshal(payload)
	if err != nil {
		return false, errors.New(fmt.Sprintf("Failed to transform Fledge data: %s", err))
	}

	lc.Debugf("Fledge Payload: %s", msg)

	return true, string(msg)
}
