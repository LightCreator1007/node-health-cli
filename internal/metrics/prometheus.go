package metrics

import (
	"log"
	"net/http"

	"github.com/LightCreator1007/node-health-cli/internal/k8s"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apiFetchCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "k8s_api_fetches_total",
			Help: "Total number of times the Kubernetes API was queried",
		},
	)

	nodeReadyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "node_ready_status",
			Help: "Current ready status of the node",
		},
		[]string{"node_name", "kubelet_version"},
	)
)

func init() {
	prometheus.MustRegister(apiFetchCounter)
	prometheus.MustRegister(nodeReadyGauge)
}

func StartServer() {
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Metrics Server running on http://localhost:5000/metrics")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalf("Metrics server failed to start: %v", err)
	}
}

func RecordNodeStatus(nodes []k8s.NodeInfo) {
	apiFetchCounter.Inc()

	for _, node := range nodes {
		statusValue := 0.0
		if node.Status == k8s.StatusHealthy {
			statusValue = 1.0
		}
		nodeReadyGauge.WithLabelValues(node.Name, node.KubeletVersion).Set(statusValue)
	}
}
