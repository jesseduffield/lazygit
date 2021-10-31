package models

type GithubPullRequest struct {
	HeadRefName         string                `json:"headRefName"`
	Number              int                   `json:"number"`
	State               string                `json:"state"` // "MERGED", "OPEN", "CLOSED"
	Url                 string                `json:"url"`
	HeadRepositoryOwner GithubRepositoryOwner `json:"headRepositoryOwner"`
}

func (pr *GithubPullRequest) UserName() string {
	return pr.HeadRepositoryOwner.Login
}

func (pr *GithubPullRequest) BranchName() string {
	return pr.HeadRefName
}

type GithubRepositoryOwner struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}
