package dokploy

type DokployComposeDTO struct {
	DokployServerID       string `json:"dokploy_server_id"`
	ComposeIdentifier     string `json:"compose_identifier"`
	InstanceIdentifier    string `json:"instance_identifier"`
	DisplayName           string `json:"display_name"`
	EnvironmentIdentifier string `json:"environment_identifier,omitempty"`
	Status                string `json:"status"`
}
