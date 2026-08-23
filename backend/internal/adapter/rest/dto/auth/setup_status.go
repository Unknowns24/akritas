package auth

type SetupStatusDTO struct {
	SetupRequired    bool `json:"setup_required"`
	RegistrationOpen bool `json:"registration_open"`
}
