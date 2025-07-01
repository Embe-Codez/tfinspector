package scanner

import (
	"path/filepath"
	"testing"
)

func TestScanDirectory(t *testing.T) {
	root := "testdata"

	projects, err := ScanDirectory(root)
	if err != nil {
		t.Fatalf("unexpected error from ScanDirectory: %v", err)
	}

	expectedFiles := map[string]struct{}{
		"azure.tf":      {},
		"multiple.tf":   {},
		"no_version.tf": {},
		"empty.tf":      {},
	}

	seen := make(map[string]bool)

	for _, p := range projects {
		_, file := filepath.Split(p.Path)
		seen[file] = true
	}

	for file := range expectedFiles {
		if !seen[file] {
			t.Errorf("expected file %q to be parsed, but it was not", file)
		}
	}
}
