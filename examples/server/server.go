package main

import (
	"encoding/json"
	v1 "github.com/containerd/nri/types/v1"
	"github.com/prometheus/common/log"
	"net"
	"net/http"
	"nri-hook/pkg/event"
	"os"
)

const SOCKET = "/host/tmp/nri-server.sock"

func main() {
	os.Remove(SOCKET)
	sock, err := net.Listen("unix", SOCKET)
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Handler: http.HandlerFunc(demoFunc),
	}
	if err := server.Serve(sock); err != nil {
		panic(err)
	}
}

func demoFunc(writer http.ResponseWriter, request *http.Request) {
	e := event.Event{}
	if err := json.NewDecoder(request.Body).Decode(&e); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	switch v1.State(e.State) {
	case v1.Create:
		// TODO: snapshot pod volume ...
		log.Infof("pod %s/%s container %s creating", e.PodNamespace, e.PodName, e.ContainerName)
	case v1.Delete:
		// TODO: snapshot pod volume again
		log.Infof("pod %s/%s container %s creating", e.PodNamespace, e.PodName, e.ContainerName)
	}

	writer.WriteHeader(http.StatusOK)
}