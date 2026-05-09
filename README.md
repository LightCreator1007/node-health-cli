# Node-Health-Cli

A Go CLI tool that connects to your local Kubernetes cluster and renders a
color-coded node health dashboard in the terminal.


## Project Structure

```
node-health-cli/
├── cmd/
│   └── main.go              # Entry point,wires everything together
├── internal/
│   ├── k8s/
│   │   ├── types.go         # NodeInfo struct, NodeStatus enum
│   │   ├── client.go        # Kubernetes auth + client construction
│   │   └── nodes.go         # API calls + analysis engine
│   └── ui/
│       └── renderer.go      # Terminal dashboard rendering
├── go.mod
├── go.sum
└── README.md
```
---

## Prerequisites

- Go 1.22+ ([install](https://go.dev/dl/))
- `kubectl` ([install](https://kubernetes.io/docs/tasks/tools/))
- A running local cluster: [Minikube](https://minikube.sigs.k8s.io/) or [Kind](https://kind.sigs.k8s.io/)

---

## Setup

### 1. Initialize your Go module

```bash
# Clone or create the project directory
mkdir node-health-cli && cd node-health-cli

# Initialize the Go module. This creates go.mod (like Cargo.toml in Rust,
# or package.json in Node). The module path is how other packages import yours.
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
go get github.com/charmbracelet/lipgloss@v0.10.0

# Tidy: removes unused deps and writes go.sum (like Cargo.lock / package-lock.json)
go mod tidy
```

### 3. Build and run

```bash
# Run directly (compiles + executes in one step, like `cargo run`)
go run cmd/main.go

# OR: Compile a binary first (like `cargo build --release`)
go build -o node-health-cli cmd/main.go
./node-health-cli
```

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
