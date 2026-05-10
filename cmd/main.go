package main

import (
	"fmt"
	"os"
	"time"

	"github.com/LightCreator1007/node-health-cli/internal/k8s"
	"github.com/LightCreator1007/node-health-cli/internal/metrics"
	"github.com/LightCreator1007/node-health-cli/internal/ui"
)

func main() {
	client, err := k8s.NewClient()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to Kubernetes: %v\n", err)
		os.Exit(1)
	}

	go metrics.StartServer()
	pollInterval := 10 * time.Second

	for {
		nodes, err := k8s.FetchNodes(client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to fetch nodes: %v\n", err)
			time.Sleep(pollInterval)
			continue
		}
		metrics.RecordNodeStatus(nodes)
		fmt.Print("\033[H\033[2J")
		ui.RenderDashboard(nodes)
		time.Sleep(pollInterval)
	}
}
