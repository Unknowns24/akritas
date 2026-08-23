package dokploy

type DokploySourceDTO struct {
	Type               string `json:"type"`
	DokployServerID    string `json:"dokploy_server_id"`
	ResourceIdentifier string `json:"resource_identifier"`
	ServiceName        string `json:"service_name,omitempty"`
	InstanceIdentifier string `json:"instance_identifier"`
	DisplayName        string `json:"display_name"`
	Environment        string `json:"environment,omitempty"`
	Status             string `json:"status"`
	RuntimeType        string `json:"runtime_type,omitempty"`
}
