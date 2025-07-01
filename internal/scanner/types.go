package scanner

type ProviderInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type TerraformProject struct {
	Path           string         `json:"path"`
	RuntimeVersion string         `json:"runtime_version,omitempty"`
	Providers      []ProviderInfo `json:"providers"`
}
