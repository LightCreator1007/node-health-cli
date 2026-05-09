package k8s

type NodeStatus string

const (
	StatusHealthy  NodeStatus = "healthy"
	StatusDegraded NodeStatus = "Degraded"
	StatusFailing  NodeStatus = "Unhealthy"
)

type NodeInfo struct {
	Name           string
	Status         NodeStatus
	Roles          []string
	KubeletVersion string
	Issues         []string
	TaintCount     int
}
