package plugins

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"nri-hook/pkg/config"
	"nri-hook/pkg/event"

	v1 "github.com/containerd/nri/types/v1"
)

type LogHook struct {
	config  config.NRIConfig
	request v1.Request
}

func (l *LogHook) Type() string {
	return config.PLUGIN_NAME
}

func (l *LogHook) Invoke(ctx context.Context, request *v1.Request) (res *v1.Result, err error) {
	if l.config, err = config.FromRequest(request); err != nil {
		return &v1.Result{
			Plugin:   config.PLUGIN_NAME,
			Version:  config.PLUGIN_VERSION,
			Error:    "bad request",
			Metadata: nil,
		}, err
	}
	l.request = *request
	if !l.interested() {
		return &v1.Result{
			Plugin:   config.PLUGIN_NAME,
			Version:  config.PLUGIN_VERSION,
			Error:    "",
			Metadata: nil,
		}, nil
	}

	return l.process()
}

func (l *LogHook) interested() bool {
	// is container?
	return l.request.ID != l.request.SandboxID &&
		// state interested?
		l.config.InterestedState(string(l.request.State)) &&
		// label interested?
		l.config.InterestedLabels(l.request.Labels)
}

func (l *LogHook) process() (res *v1.Result, err error) {
	res = &v1.Result{}
	e := event.EventFromNRIRequest(l.request)
	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", l.config.Address)
			},
		},
	}
	resp, err := client.Post(l.config.URI, "application/json", bytes.NewBuffer(e.ToJson()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%d error from server: %s", resp.StatusCode, string(buf))
	}
	return
}
