package k8s

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchNodes(client *kubernetes.Clientset) ([]NodeInfo, error) {
	nodeList, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]NodeInfo, 0, len(nodeList.Items))

	for _, node := range nodeList.Items {
		info := analyzeNode(node)
		result = append(result, info)
	}

	return result, nil
}

func analyzeNode(node corev1.Node) NodeInfo {
	info := NodeInfo{
		Name:           node.Name,
		KubeletVersion: node.Status.NodeInfo.KubeletVersion,
	}

	const roleLabelPrefix = "node-role.kubernetes.io/"
	for labelKey := range node.Labels {
		if strings.HasPrefix(labelKey, roleLabelPrefix) {
			role := strings.TrimPrefix(labelKey, roleLabelPrefix)
			info.Roles = append(info.Roles, role)
		}
	}

	if len(info.Roles) == 0 {
		info.Roles = []string{"<none>"}
	}

	info.TaintCount = len(node.Spec.Taints)

	info.Status = StatusHealthy

	for _, condition := range node.Status.Conditions {
		switch condition.Type {
		case corev1.NodeReady:
			if condition.Status != corev1.ConditionTrue {
				info.Status = StatusFailing
				msg := "Not Ready"
				if condition.Message != "" {
					msg = fmt.Sprintf("Not ready: %s", condition.Message)
				}
				info.Issues = append(info.Issues, msg)
			}
		case corev1.NodeMemoryPressure:
			if condition.Status == corev1.ConditionTrue {
				info.Status = StatusFailing
				info.Issues = append(info.Issues, "MemoryPressure: node is running low on RAM")
			}
		case corev1.NodeDiskPressure:
			if condition.Status == corev1.ConditionTrue {
				info.Status = StatusFailing
				info.Issues = append(info.Issues, "DiskPressure:Node is running low on disk space")
			}
		case corev1.NodePIDPressure:
			if condition.Status == corev1.ConditionTrue {
				if info.Status == StatusHealthy {
					info.Status = StatusDegraded
				}
				info.Issues = append(info.Issues, "PIDPressure: Node is running low on process IDs")
			}
		}

	}

	if info.TaintCount > 0 && info.Status == StatusHealthy {
		info.Status = StatusDegraded
	}

	return info
}
