package auth

type SetupStatus struct {
	SetupRequired    bool `json:"setup_required"`
	RegistrationOpen bool `json:"registration_open"`
}

type SetupStatusResponse struct {
	Data SetupStatus `json:"data"`
}
