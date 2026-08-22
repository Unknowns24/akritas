package dto

type CreateGitHubPATAccountRequestDTO struct {
	DisplayName         string `json:"display_name"`
	AccountType         string `json:"account_type"`
	AccountIdentifier   string `json:"account_identifier"`
	PersonalAccessToken string `json:"personal_access_token"`
}
