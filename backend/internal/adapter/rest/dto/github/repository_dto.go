package github

type GitHubRepositoryDTO struct {
	GitHubAccountID      string `json:"github_account_id"`
	RepositoryIdentifier string `json:"repository_identifier"`
	Owner                string `json:"owner"`
	Name                 string `json:"name"`
	FullName             string `json:"full_name"`
	DefaultBranch        string `json:"default_branch"`
	Private              bool   `json:"private"`
	HTMLURL              string `json:"html_url"`
}
