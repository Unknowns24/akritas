package github

type GitHubManifestRegistrationRequestDTO struct {
	DisplayName  string `json:"display_name"`
	OwnerType    string `json:"owner_type"`
	Organization string `json:"organization"`
}
