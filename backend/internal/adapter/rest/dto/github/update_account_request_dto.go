package github

type UpdateGitHubAccountRequestDTO struct {
	DisplayName         *string `json:"display_name"`
	PersonalAccessToken *string `json:"personal_access_token"`
}
