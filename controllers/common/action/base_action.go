package action

import (
	"github.com/go-logr/logr"
	"github.com/securesign/operator/client"
	"k8s.io/client-go/tools/record"
)

type BaseAction struct {
	Client   client.Client
	Recorder record.EventRecorder
	Logger   logr.Logger
}

func (action *BaseAction) InjectClient(client client.Client) {
	action.Client = client
}

func (action *BaseAction) InjectRecorder(recorder record.EventRecorder) {
	action.Recorder = recorder
}

func (action *BaseAction) InjectLogger(logger logr.Logger) {
	action.Logger = logger
}