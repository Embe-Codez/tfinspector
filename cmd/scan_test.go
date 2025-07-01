package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/embe-codez/tfinspector/internal/scanner"
)

func mockProject() []scanner.TerraformProject {
	return []scanner.TerraformProject{
		{
			Path:           "main.tf",
			RuntimeVersion: ">= 1.0.0",
			Providers: []scanner.ProviderInfo{
				{Name: "aws", Version: "5.0.0"},
				{Name: "azurerm"},
			},
		},
	}
}

func TestRenderOutput_Text(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer

	err := renderOutput(mockProject(), &buf, outputText)
	if err != nil {
		t.Fatalf("renderOutput (text) failed: %v", err)
	}

	output := buf.String()
	assertContains(t, output, "Path: main.tf")
	assertContains(t, output, "Runtime Version: >= 1.0.0")
	assertContains(t, output, "Provider: aws (version 5.0.0)")
	assertContains(t, output, "Provider: azurerm")
}

func TestWriteJSON(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer

	err := writeJSON(mockProject(), &buf)
	if err != nil {
		t.Fatalf("writeJSON failed: %v", err)
	}

	output := buf.String()
	assertContains(t, output, `"path": "main.tf"`)
	assertContains(t, output, `"runtime_version": "\u003e= 1.0.0"`)
	assertContains(t, output, `"name": "aws"`)
}

func TestWriteYAML(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer

	err := writeYAML(mockProject(), &buf)
	if err != nil {
		t.Fatalf("writeYAML failed: %v", err)
	}

	output := buf.String()
	assertContains(t, output, "path: main.tf")
	assertContains(t, output, "runtimeversion:")
	assertContains(t, output, ">= 1.0.0")
	assertContains(t, output, "name: aws")
}

func TestWriteCSV(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer

	err := writeCSV(mockProject(), &buf)
	if err != nil {
		t.Fatalf("writeCSV failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 rows), got: %d\n%s", len(lines), buf.String())
	}

	assertContains(t, lines[0], "path,provider,version")
	assertContains(t, lines[1], "main.tf,aws,5.0.0")
	assertContains(t, lines[2], "main.tf,azurerm,")
}

func TestParseFormat_Valid(t *testing.T) {
	t.Parallel()
	tests := map[string]outputType{
		"text": outputText,
		"json": outputJSON,
		"yaml": outputYAML,
		"csv":  outputCSV,
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			got, err := parseFormat(input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != expected {
				t.Errorf("expected %v, got %v", expected, got)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	t.Parallel()
	_, err := parseFormat("xml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	assertContains(t, err.Error(), "invalid output format")
}

func assertContains(t *testing.T, output, want string) {
	t.Helper()
	if !strings.Contains(output, want) {
		t.Errorf("expected output to contain %q\nGot:\n%s", want, output)
	}
}
