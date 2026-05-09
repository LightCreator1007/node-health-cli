package main

import (
	"fmt"
	"os"

	"github.com/LightCreator1007/node-health-cli/internal/k8s"
	"github.com/LightCreator1007/node-health-cli/internal/ui"
)

func main() {
	client, err := k8s.NewClient()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to Kubernetes: %v\n", err)
		os.Exit(1)
	}

	nodes, err := k8s.FetchNodes(client)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to fetch nodes: %v\n", err)
		os.Exit(1)
	}

	ui.RenderDashboard(nodes)
}
