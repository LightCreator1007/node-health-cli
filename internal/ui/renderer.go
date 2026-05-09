package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/LightCreator1007/node-health-cli/internal/k8s"
	"github.com/charmbracelet/lipgloss"
)

var (
	colorGreen   = lipgloss.Color("#00C17C")
	colorRed     = lipgloss.Color("#FF4F4F") // Failing
	colorYellow  = lipgloss.Color("#FFD700") // Degraded
	colorCyan    = lipgloss.Color("#00D7FF") // Accent / headers
	colorGray    = lipgloss.Color("#6C7A8A") // Subdued text
	colorWhite   = lipgloss.Color("#F0F0F0") // Normal text
	colorBgDark  = lipgloss.Color("#0D1117") // Background (GitHub dark)
	colorBgPanel = lipgloss.Color("#161B22") // Panel background

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCyan).
			PaddingBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorGray).
			Italic(true)

	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorCyan).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(colorCyan)

	healthyBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(colorGreen).
			Padding(0, 1) // 0 vertical padding, 1 horizontal padding

	failingBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(colorRed).
			Padding(0, 1)

	degradedBadge = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(colorYellow).
			Padding(0, 1)

	nameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite)

	roleStyle = lipgloss.NewStyle().
			Foreground(colorCyan)

	versionStyle = lipgloss.NewStyle().
			Foreground(colorGray)

	issueStyle = lipgloss.NewStyle().
			Foreground(colorRed)

	taintStyle = lipgloss.NewStyle().
			Foreground(colorYellow)

	summaryBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorGray).
			Padding(0, 2).
			MarginTop(1)

	healthyCountStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorGreen)

	failingCountStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorRed)

	degradedCountStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorYellow)

	colWidthName    = 28
	colWidthStatus  = 14
	colWidthRoles   = 20
	colWidthVersion = 18
	colWidthIssues  = 42
)

func RenderDashboard(nodes []k8s.NodeInfo) {
	printHeader(len(nodes))
	printTableHeader()
	for _, node := range nodes {
		printNodeRow(node)
	}
	printSummary(nodes)
}

func printHeader(nodeCount int) {
	banner := `
	███╗   ██╗ ██████╗ ██████╗ ███████╗    ██╗  ██╗███████╗ █████╗ ██╗  ████████╗██╗  ██╗
	████╗  ██║██╔═══██╗██╔══██╗██╔════╝    ██║  ██║██╔════╝██╔══██╗██║  ╚══██╔══╝██║  ██║
	██╔██╗ ██║██║   ██║██║  ██║█████╗      ███████║█████╗  ███████║██║     ██║   ███████║
	██║╚██╗██║██║   ██║██║  ██║██╔══╝      ██╔══██║██╔══╝  ██╔══██║██║     ██║   ██╔══██║
	██║ ╚████║╚██████╔╝██████╔╝███████╗    ██║  ██║███████╗██║  ██║███████╗██║   ██║  ██║
	╚═╝  ╚═══╝ ╚═════╝ ╚═════╝ ╚══════╝    ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝╚═╝   ╚═╝  ╚═╝`

	fmt.Println(lipgloss.NewStyle().Foreground(colorCyan).Render(banner))
	fmt.Println()
	fmt.Println(
		titleStyle.Render(" ⎈  Kubernetes Node Health Dashboard"),
	)

	fmt.Println(
		subtitleStyle.Render(fmt.Sprintf(
			"  Cluster scan at %s  ·  %d node(s) found",
			time.Now().Format("2006-01-02 15:04:05"),
			nodeCount,
		)),
	)
	fmt.Println()

	fmt.Println(lipgloss.NewStyle().Foreground(colorGray).Render(strings.Repeat("─", 120)))
}

func printTableHeader() {

	header := fmt.Sprintf("  %s  %s  %s  %s  %s",
		lipgloss.NewStyle().Bold(true).Foreground(colorCyan).Width(colWidthName).Render("NODE NAME"),
		lipgloss.NewStyle().Bold(true).Foreground(colorCyan).Width(colWidthStatus).Render("STATUS"),
		lipgloss.NewStyle().Bold(true).Foreground(colorCyan).Width(colWidthRoles).Render("ROLES"),
		lipgloss.NewStyle().Bold(true).Foreground(colorCyan).Width(colWidthVersion).Render("KUBELET VER"),
		lipgloss.NewStyle().Bold(true).Foreground(colorCyan).Render("ISSUES / NOTES"),
	)
	fmt.Println(header)
	fmt.Println(lipgloss.NewStyle().Foreground(colorGray).Render(strings.Repeat("─", 120)))
}

func printNodeRow(node k8s.NodeInfo) {

	var statusRendered string
	switch node.Status {
	case k8s.StatusHealthy:
		statusRendered = healthyBadge.Render("● HEALTHY")
	case k8s.StatusFailing:
		statusRendered = failingBadge.Render("✖ FAILING")
	case k8s.StatusDegraded:
		statusRendered = degradedBadge.Render("⚠ DEGRADED")
	default:
		statusRendered = degradedBadge.Render("? UNKNOWN")
	}

	issuesRendered := lipgloss.NewStyle().Foreground(colorGreen).Render("✓ No issues detected")
	if len(node.Issues) > 0 {
		issuesRendered = issueStyle.Render(strings.Join(node.Issues, " | "))
	}

	var taintSuffix string
	if node.TaintCount > 0 {
		taintSuffix = taintStyle.Render(fmt.Sprintf("  [%d taint(s)]", node.TaintCount))
	}

	row := fmt.Sprintf("  %s  %s  %s  %s  %s%s",
		nameStyle.Width(colWidthName).Render(node.Name),
		lipgloss.NewStyle().Width(colWidthStatus).Render(statusRendered),
		roleStyle.Width(colWidthRoles).Render(strings.Join(node.Roles, ", ")),
		versionStyle.Width(colWidthVersion).Render(node.KubeletVersion),
		issuesRendered,
		taintSuffix,
	)

	fmt.Println(row)

	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("#1E2633")).Render(strings.Repeat("╌", 120)))
}

func printSummary(nodes []k8s.NodeInfo) {

	var healthy, failing, degraded int
	for _, n := range nodes {
		switch n.Status {
		case k8s.StatusHealthy:
			healthy++
		case k8s.StatusFailing:
			failing++
		case k8s.StatusDegraded:
			degraded++
		}
	}

	var overallMsg string
	switch {
	case failing > 0:
		overallMsg = lipgloss.NewStyle().Bold(true).Foreground(colorRed).Render(
			"✖ CLUSTER ATTENTION REQUIRED — " + fmt.Sprintf("%d node(s) failing", failing),
		)
	case degraded > 0:
		overallMsg = lipgloss.NewStyle().Bold(true).Foreground(colorYellow).Render(
			"⚠ CLUSTER DEGRADED — all nodes schedulable but check taints/pressure",
		)
	default:
		overallMsg = lipgloss.NewStyle().Bold(true).Foreground(colorGreen).Render(
			"✓ ALL SYSTEMS NOMINAL — cluster is fully healthy",
		)
	}

	summaryContent := fmt.Sprintf(
		"  %s  %s  %s  %s  │  %s",
		lipgloss.NewStyle().Foreground(colorGray).Render("Summary:"),
		healthyCountStyle.Render(fmt.Sprintf("%d Healthy", healthy)),
		degradedCountStyle.Render(fmt.Sprintf("%d Degraded", degraded)),
		failingCountStyle.Render(fmt.Sprintf("%d Failing", failing)),
		overallMsg,
	)

	fmt.Println()
	fmt.Println(summaryBoxStyle.Render(summaryContent))
	fmt.Println()
}
