package dokploy

type DokployComposeServiceDTO struct {
	DokployServerID   string `json:"dokploy_server_id"`
	ComposeIdentifier string `json:"compose_identifier"`
	ServiceName       string `json:"service_name"`
}
