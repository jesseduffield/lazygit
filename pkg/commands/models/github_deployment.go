package models

import "time"

// GithubDeployment is the most recent deployment to a given environment, as
// reported by the GitHub Deployments API.
type GithubDeployment struct {
	Environment string
	// State is the latest deployment-status state (e.g. "SUCCESS", "FAILURE",
	// "IN_PROGRESS"), falling back to the deployment's own state when no status
	// has been reported yet.
	State       string
	Ref         string
	Sha         string // abbreviated commit hash
	Description string
	UpdatedAt   time.Time
}
