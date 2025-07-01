package scanner

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

// ParseTerraformFile parses a Terraform file and extracts its runtime version and provider metadata.
func ParseTerraformFile(path string) (*TerraformProject, error) {
	parser := hclparse.NewParser()

	file, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, diags
	}

	project := &TerraformProject{
		Path:      path,
		Providers: []ProviderInfo{},
	}

	schema := &hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "terraform"},
			{Type: "provider"},
		},
	}

	content, diags := file.Body.Content(schema)
	if diags.HasErrors() {
		return nil, diags
	}

	// Prevent duplicated Terraform blocks from being returned.
	seen := make(map[string]bool)

	for _, block := range content.Blocks {
		switch block.Type {
		case "terraform":
			parseTerraformBlock(block, project, seen)
		case "provider":
			parseProviderBlock(block, project, seen)
		}
	}

	return project, nil
}

// parseTerraformBlock extracts the required_version and required_providers from a terraform block.
func parseTerraformBlock(block *hcl.Block, project *TerraformProject, seen map[string]bool) {
	content, _, _ := block.Body.PartialContent(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "required_version"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "required_providers"},
		},
	})

	if attr, ok := content.Attributes["required_version"]; ok {
		if val, diag := attr.Expr.Value(nil); diag == nil && val.Type().IsPrimitiveType() {
			project.RuntimeVersion = val.AsString()
		}
	}

	for _, b := range content.Blocks {
		if b.Type == "required_providers" {
			parseRequiredProviders(b, project, seen)
		}
	}
}

// parseRequiredProviders parses provider names and versions from a required_providers block.
func parseRequiredProviders(block *hcl.Block, project *TerraformProject, seen map[string]bool) {
	attrs, diags := block.Body.JustAttributes()
	if diags.HasErrors() {
		return
	}

	for name, attr := range attrs {
		val, valDiags := attr.Expr.Value(nil)
		if valDiags.HasErrors() || !val.Type().IsObjectType() {
			continue
		}

		obj := val.AsValueMap()
		version := ""
		if v, ok := obj["version"]; ok && v.Type().IsPrimitiveType() {
			version = v.AsString()
		}

		addProvider(name, version, project, seen)
	}
}

// parseProviderBlock parses individual provider blocks not declared in required_providers.
func parseProviderBlock(block *hcl.Block, project *TerraformProject, seen map[string]bool) {
	if len(block.Labels) == 0 {
		return
	}
	addProvider(block.Labels[0], "", project, seen)
}

// addProvider adds a provider to the project if it hasn't already been recorded.
func addProvider(name, version string, project *TerraformProject, seen map[string]bool) {
	if seen[name] {
		return
	}
	project.Providers = append(project.Providers, ProviderInfo{
		Name:    name,
		Version: version,
	})
	seen[name] = true
}
