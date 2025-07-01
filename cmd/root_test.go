package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmd_HelpOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs([]string{"--help"})

	err := RootCmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "tfinspector") {
		t.Errorf("expected output to mention tfinspector, got:\n%s", output)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("expected output to contain Usage section, got:\n%s", output)
	}
}
