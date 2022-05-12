//
// Copyright (c) 2022 One Track Consulting
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/util"
	"github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap"
	"github.com/edgexfoundry/go-mod-messaging/v2/pkg/types"
	"github.com/nats-io/nats.go"
)

const (
	serviceKey = "app-custom-trigger-nats-rpc"
)

type rpcTrigger struct {
	tc interfaces.TriggerConfig
}

func (t *rpcTrigger) Initialize(_ *sync.WaitGroup, ctx context.Context, _ <-chan interfaces.BackgroundMessage) (bootstrap.Deferred, error) {
	ct, err := nats.Connect("demo.nats.io")

	if err != nil {
		return nil, fmt.Errorf("failed to start connect to NATS server %w", err)
	}

	sub, err := ct.Subscribe("rpc.*", func(msg *nats.Msg) {
		env := types.MessageEnvelope{
			Payload: msg.Data,
		}
		t.tc.MessageReceived(t.tc.ContextBuilder(env), env, func(ctx interfaces.AppFunctionContext, pipeline *interfaces.FunctionPipeline) error {
			return msg.Respond([]byte(fmt.Sprintf("got %s (from %s at %v)", string(msg.Data), msg.Subject, time.Now().UTC())))
		})
	})

	return func() {
		_ = sub.Drain()
		ct.Close()
	}, nil
}

func main() {
	// turn off secure mode for examples. Not recommended for production
	_ = os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// First thing to do is to create an instance of the EdgeX SDK Service, which also runs the bootstrap initialization.
	service, ok := pkg.NewAppServiceWithTargetType(serviceKey, &[]byte{})
	if !ok {
		os.Exit(-1)
	}

	service.RegisterCustomTriggerFactory("custom-rpc", func(config interfaces.TriggerConfig) (interfaces.Trigger, error) {
		return &rpcTrigger{
			tc: config,
		}, nil
	})

	var err error

	//use this to process using default pipeline only
	err = service.SetDefaultFunctionsPipeline(printToConsole)
	if err != nil {
		service.LoggingClient().Errorf("SetDefaultFunctionsPipeline returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	err = service.MakeItRun()
	if err != nil {
		service.LoggingClient().Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	service.LoggingClient().Info("Exiting service")
	// Do any required cleanup here
	os.Exit(0)
}

func printToConsole(appContext interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	input, err := util.CoerceType(data)

	if err != nil {
		appContext.LoggingClient().Error(err.Error())
		return false, err
	}

	wait := time.Millisecond * time.Duration(len(input))

	time.Sleep(wait)

	appContext.LoggingClient().Info("PrintToConsole")

	os.Stdout.WriteString(fmt.Sprintf("'%s' received %s ago\n>", string(input), wait.String()))

	return false, nil
}
