package scanner

import (
	"path/filepath"
	"testing"
)

func TestParseTerraformFile(t *testing.T) {
	cases := []struct {
		name              string
		filename          string
		expectError       bool
		expectedVersion   string
		expectedProviders []string
	}{
		{
			name:              "azure provider with version",
			filename:          "azure.tf",
			expectedVersion:   ">= 1.4.0",
			expectedProviders: []string{"azurerm"},
		},
		{
			name:              "multiple providers",
			filename:          "multiple.tf",
			expectedVersion:   "~> 1.2.0",
			expectedProviders: []string{"azurerm", "aws"},
		},
		{
			name:              "provider with no required_version",
			filename:          "no_version.tf",
			expectedVersion:   "",
			expectedProviders: []string{"azurerm"},
		},
		{
			name:        "empty file",
			filename:    "empty.tf",
			expectError: false,
		},
		{
			name:        "malformed file",
			filename:    "malformed.tf",
			expectError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join("testdata", tc.filename)

			project, err := ParseTerraformFile(path)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if project.RuntimeVersion != tc.expectedVersion {
				t.Errorf("expected runtime version %q, got %q", tc.expectedVersion, project.RuntimeVersion)
			}

			if len(tc.expectedProviders) > 0 {
				found := make(map[string]bool)
				for _, p := range project.Providers {
					found[p.Name] = true
				}
				for _, expected := range tc.expectedProviders {
					if !found[expected] {
						t.Errorf("expected provider %q not found", expected)
					}
				}
			}
		})
	}
}
