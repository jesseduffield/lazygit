package models

// TODO: see if I need to store the head repo name in case it differs from the base repo
type GithubPullRequest struct {
	HeadRefName         string                `json:"headRefName"`
	Number              int                   `json:"number"`
	State               string                `json:"state"` // "MERGED", "OPEN", "CLOSED"
	Url                 string                `json:"url"`
	HeadRepositoryOwner GithubRepositoryOwner `json:"headRepositoryOwner"`
}

func (pr *GithubPullRequest) UserName() string {
	// e.g. 'jesseduffield'
	return pr.HeadRepositoryOwner.Login
}

func (pr *GithubPullRequest) BranchName() string {
	// e.g. 'feature/my-feature'
	return pr.HeadRefName
}

type GithubRepositoryOwner struct {
	Login string `json:"login"`
}
