package presentation

import (
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/gui/style"
	"github.com/jesseduffield/lazygit/pkg/i18n"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

// GetDeploymentsContent renders the environment deployment statuses shown in the
// status panel's main view.
func GetDeploymentsContent(deployments []*models.GithubDeployment, tr *i18n.TranslationSet) string {
	rows := lo.Map(deployments, func(d *models.GithubDeployment, _ int) []string {
		return deploymentRow(d)
	})

	lines, _ := utils.RenderDisplayStrings(rows, nil)
	return style.AttrBold.Sprint(tr.DeploymentsTitle) + "\n\n" + strings.Join(lines, "\n")
}

func deploymentRow(d *models.GithubDeployment) []string {
	st := deploymentStatusStyle(d.State)

	refAndSha := d.Ref
	if d.Sha != "" {
		if refAndSha != "" {
			refAndSha += " @ " + d.Sha
		} else {
			refAndSha = d.Sha
		}
	}

	updated := ""
	if !d.UpdatedAt.IsZero() {
		updated = utils.UnixToTimeAgo(d.UpdatedAt.Unix())
	}

	return []string{
		"  " + d.Environment,
		st.Sprint(deploymentStatusSymbol(d.State) + " " + humanizeState(d.State)),
		refAndSha,
		updated,
	}
}

// humanizeState turns a GraphQL state enum (e.g. "IN_PROGRESS") into something
// friendlier to read (e.g. "in progress").
func humanizeState(state string) string {
	if state == "" {
		return "unknown"
	}
	return strings.ToLower(strings.ReplaceAll(state, "_", " "))
}

func deploymentStatusSymbol(state string) string {
	switch normalizeState(state) {
	case "success", "active":
		return "✓"
	case "failure", "error", "abandoned":
		return "✗"
	case "in progress", "pending", "queued", "waiting":
		return "⟳"
	default:
		return "•"
	}
}

func deploymentStatusStyle(state string) style.TextStyle {
	switch normalizeState(state) {
	case "success", "active":
		return style.FgGreen
	case "failure", "error", "abandoned":
		return style.FgRed
	case "in progress", "pending", "queued", "waiting":
		return style.FgYellow
	default:
		return style.FgDefault
	}
}

func normalizeState(state string) string {
	return strings.ToLower(strings.ReplaceAll(state, "_", " "))
}
