package event

import (
	"encoding/json"
	"io"
	"nri-hook/pkg/config"

	v1 "github.com/containerd/nri/types/v1"
)

// NRI-Request -> Event -> json -> [log-agent server] -> json -> Result -> NRI-Result

type Event struct {
	State         string
	PodNamespace  string
	PodName       string
	PodID         string
	ContainerName string
	ContainerID   string
	ContainerRoot string
}

func EventFromNRIRequest(request v1.Request) Event {
	return Event{
		State:         string(request.State),
		PodNamespace:  request.Labels[config.LABEL_POD_NAMESPACE],
		PodName:       request.Labels[config.LABEL_POD_NAME],
		PodID:         request.Labels[config.LABEL_POD_UID],
		ContainerName: request.Spec.Annotations[config.ANNO_CONTAINER_NAME],
		ContainerID:   request.ID,
		ContainerRoot: "",
	}
}

func (e Event) ToJson() []byte {
	bytes, _ := json.Marshal(e)
	return bytes
}

// Deprecated
type Result struct {
	Message string
}

func (r *Result) ToNRIResult() *v1.Result {
	return &v1.Result{
		Plugin:   config.PLUGIN_NAME,
		Version:  config.PLUGIN_VERSION,
		Error:    "",
		Metadata: nil,
	}
}

func (r *Result) FromJson(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(r)
}
