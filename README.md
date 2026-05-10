# Node-Health-Cli

A Go CLI tool that connects to your local Kubernetes cluster and renders a
color-coded node health dashboard in the terminal.

<img width="1638" height="607" alt="Screenshot From 2026-05-09 13-04-36" src="https://github.com/user-attachments/assets/80447534-ffd2-44f5-8374-2c0fc0f3cc0b" />


## Project Structure

```text
node-health-cli/
├── cmd/
│   └── main.go              # Entry point, wires everything together, starts daemon
├── internal/
│   ├── k8s/
│   │   ├── types.go         # NodeInfo struct, NodeStatus enum
│   │   ├── client.go        # Kubernetes auth + client construction
│   │   └── nodes.go         # API calls + analysis engine
│   ├── metrics/
│   │   └── prometheus.go    # Prometheus registry, custom metrics, and HTTP server
│   └── ui/
│       └── renderer.go      # Terminal dashboard rendering
├── prometheus.yml           # Configuration for local Prometheus scraping
├── go.mod
├── go.sum
└── README.md
```
---

## Prerequisites

- Go 1.22+ ([install](https://go.dev/dl/))
- `kubectl` ([install](https://kubernetes.io/docs/tasks/tools/))
- Docker (required for running Prometheus/Grafana)
- A running local cluster: [Minikube](https://minikube.sigs.k8s.io/) or [Kind](https://kind.sigs.k8s.io/)

---

## Setup

### 1. Initialize your Go module

```bash
mkdir node-health-cli && cd node-health-cli

go mod init github.com/yourusername/node-health-cli
```

### 2. Install dependencies

```bash
# client-go: the official Kubernetes Go client
# api + apimachinery: the type definitions (Node, Pod, etc.)
# lipgloss: terminal styling library
go get k8s.io/client-go@v0.29.3
go get k8s.io/api@v0.29.3
go get k8s.io/apimachinery@v0.29.3
go get [github.com/charmbracelet/lipgloss@v0.10.0](https://github.com/charmbracelet/lipgloss@v0.10.0)

# prometheus: for exposing internal metrics to the observability pipeline
go get [github.com/prometheus/client_golang/prometheus](https://github.com/prometheus/client_golang/prometheus)
go get [github.com/prometheus/client_golang/prometheus/promhttp](https://github.com/prometheus/client_golang/prometheus/promhttp)

# Tidy: removes unused deps and writes go.sum (like Cargo.lock / package-lock.json)
go mod tidy
```

### 3. Build and run

```bash
go run cmd/main.go

go build -o node-health-cli cmd/main.go
./node-health-cli
```

---

## Observability & Metrics (Prometheus + Grafana)

To align with modern Cloud Native monitoring standards, `node-health-cli` runs as a continuous daemon and exposes internal evaluation metrics. These metrics are scraped by Prometheus and visualized in Grafana, transforming the CLI into a real-time observability pipeline.

<img width="1638" alt="Grafana Dashboard showing Node Health" src="YOUR_GRAFANA_SCREENSHOT_LINK_HERE" />

### Exported Metrics
The application exposes a `/metrics` endpoint on port `5000` containing standard Go runtime metrics, alongside the following custom Kubernetes tracking metrics:
* `k8s_api_fetches_total` (Counter): Total number of times the Kubernetes API was successfully queried.
* `node_ready_status` (Gauge): Current readiness status of the node (1 = Healthy, 0 = Failing/Degraded), labeled by `node_name` and `kubelet_version`.

### Local Architecture Flow
1. **Go CLI (Port 5000):** Continuously polls the K8s API and updates internal Prometheus registries.
2. **Prometheus (Port 9090):** Scrapes the Go CLI `/metrics` endpoint every 10 seconds.
3. **Grafana (Port 3000):** Queries Prometheus to visualize cluster health trends over time.

---

### How to Run the Observability Pipeline Locally

If you want to spin up the full dashboard locally alongside the CLI, you can use the included `prometheus.yml` configuration and Docker.

**1. Start the Go CLI Daemon:**
```bash
go run cmd/main.go
```
**2. Start Prometheus (in a new terminal):**

```bash
docker run -d \
  --name prometheus \
  -p 9090:9090 \
  --add-host host.docker.internal:host-gateway \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus
```

**3. Start Grafana:**

```bash
docker run -d \
  --name grafana \
  -p 3000:3000 \
  --add-host host.docker.internal:host-gateway \
  grafana/grafana
```

Once running, navigate to http://localhost:3000 (admin/admin), add Prometheus (http://host.docker.internal:9090) as a data source, and query the custom metrics to build your dashboard.

---

## Testing with Minikube (Recommended)

Minikube spins up a single-node Kubernetes cluster in a VM or container on your machine.

```bash
# Install minikube on Fedora
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

# Start a cluster (uses Docker as the driver,make sure Docker is running)
minikube start

# Verify it's running
kubectl get nodes
# Expected output:
# NAME       STATUS   ROLES           AGE   VERSION
# minikube   Ready    control-plane   1m    v1.29.x

# Now run the CLI,it will auto-detect your ~/.kube/config
go run cmd/main.go
```

### Simulating a failing node (for testing the red output)

You can't easily break minikube's single node, but you can add taints to test
the yellow "Degraded" state:

```bash
# Add a taint to the minikube node
kubectl taint nodes minikube test-key=test-value:NoSchedule

# Run the CLI, he node should now show as DEGRADED
go run cmd/main.go

# Remove the taint afterward
kubectl taint nodes minikube test-key=test-value:NoSchedule-
```

---

## Testing with Kind (Kubernetes-in-Docker)

Kind lets you run a multi-node cluster using Docker containers. Great for
testing multi-node scenarios.

```bash
# Install Kind
go install sigs.k8s.io/kind@latest

# Create a 3-node cluster (1 control-plane + 2 workers)
# Save this as kind-config.yaml:
cat > kind-config.yaml << 'EOF'
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
  - role: worker
  - role: worker
EOF

kind create cluster --config kind-config.yaml

# Verify
kubectl get nodes

# Run the CLI,you should see 3 nodes in the dashboard
go run cmd/main.go
```

---
