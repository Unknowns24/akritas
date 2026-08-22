package dto

type UpdateDokployServerRequestDTO struct {
	Name          *string `json:"name"`
	BaseURL       *string `json:"base_url"`
	APICredential *string `json:"api_credential"`
}
