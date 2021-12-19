package main

import (
	"context"
	"fmt"
	"nri-hook/pkg/plugins"
	"os"

	"github.com/containerd/nri/skel"
)

func main() {
	if err := skel.Run(context.Background(), &plugins.LogHook{}); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}
