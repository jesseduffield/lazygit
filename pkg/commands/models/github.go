package models

type GithubPullRequest struct {
	HeadRefName         string                `json:"headRefName"`
	Number              int                   `json:"number"`
	State               string                `json:"state"` // "MERGED", "OPEN", "CLOSED"
	Url                 string                `json:"url"`
	HeadRepositoryOwner GithubRepositoryOwner `json:"headRepositoryOwner"`
}

type GithubRepositoryOwner struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}
