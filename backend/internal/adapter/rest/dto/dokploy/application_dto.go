package dokploy

type DokployApplicationDTO struct {
	DokployServerID       string `json:"dokploy_server_id"`
	ApplicationIdentifier string `json:"application_identifier"`
	InstanceIdentifier    string `json:"instance_identifier"`
	DisplayName           string `json:"display_name"`
	Environment           string `json:"environment,omitempty"`
	Status                string `json:"status,omitempty"`
}
