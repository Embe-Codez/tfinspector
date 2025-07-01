package cmd

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/embe-codez/tfinspector/internal/scanner"
)

type outputType int

const (
	outputText outputType = iota
	outputJSON
	outputYAML
	outputCSV
)

var (
	outputFormat string
	outputFile   string
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan a directory for Terraform files",
	Long: `The scan command recursively scans a directory for .tf files and extracts
provider and runtime metadata. Supports structured output formats.`,
	Example: `  tfinspector scan ./infra --output json --out report.json`,
	Args:    cobra.ExactArgs(1),
	RunE:    runScan,
}

func init() {
	scanCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format: text, json, yaml, csv")
	scanCmd.Flags().StringVarP(&outputFile, "out", "", "", "Write output to file instead of stdout")
	RootCmd.AddCommand(scanCmd)
}

// runScan handles the CLI command and calls RunScan,
// which does the actual work. This keeps the CLI code
// separate from the main logic, so it's easier to test and reuse.
func runScan(cmd *cobra.Command, args []string) error {
	return RunScan(args[0], outputFormat, outputFile)
}

func RunScan(path, formatStr, outFile string) error {
	projects, err := scanner.ScanDirectory(path)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}
	if len(projects) == 0 {
		return errors.New("no Terraform files found")
	}

	format, err := parseFormat(formatStr)
	if err != nil {
		return err
	}

	out, err := openOutputFile(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	return renderOutput(projects, out, format)
}

func parseFormat(formatStr string) (outputType, error) {
	switch strings.ToLower(formatStr) {
	case "text":
		return outputText, nil
	case "json":
		return outputJSON, nil
	case "yaml", "yml":
		return outputYAML, nil
	case "csv":
		return outputCSV, nil
	default:
		return outputText, fmt.Errorf("invalid output format: %s (choose text, json, yaml, csv)", formatStr)
	}
}

func openOutputFile(path string) (io.WriteCloser, error) {
	if path == "" {
		return os.Stdout, nil
	}
	return os.Create(path)
}

func renderOutput(projects []scanner.TerraformProject, out io.Writer, format outputType) error {
	writers := map[outputType]func([]scanner.TerraformProject, io.Writer) error{
		outputText: writeText,
		outputJSON: writeJSON,
		outputYAML: writeYAML,
		outputCSV:  writeCSV,
	}

	if writer, ok := writers[format]; ok {
		return writer(projects, out)
	}
	return errors.New("unsupported output format")
}

func writeText(projects []scanner.TerraformProject, out io.Writer) error {
	writer := bufio.NewWriter(out)

	for _, p := range projects {
		if _, err := fmt.Fprintf(writer, "Path: %s\n", p.Path); err != nil {
			return err
		}
		if p.RuntimeVersion != "" {
			if _, err := fmt.Fprintf(writer, "  Runtime Version: %s\n", p.RuntimeVersion); err != nil {
				return err
			}
		}
		for _, provider := range p.Providers {
			line := fmt.Sprintf("  Provider: %s", provider.Name)
			if provider.Version != "" {
				line += fmt.Sprintf(" (version %s)", provider.Version)
			}
			if _, err := fmt.Fprintln(writer, line); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(writer); err != nil {
			return err
		}
	}

	return writer.Flush()
}

func writeJSON(projects []scanner.TerraformProject, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(projects)
}

func writeYAML(projects []scanner.TerraformProject, out io.Writer) error {
	enc := yaml.NewEncoder(out)
	defer enc.Close()
	return enc.Encode(projects)
}

func writeCSV(projects []scanner.TerraformProject, out io.Writer) error {
	writer := csv.NewWriter(out)
	defer writer.Flush()

	if err := writer.Write([]string{"path", "provider", "version"}); err != nil {
		return err
	}

	for _, p := range projects {
		for _, provider := range p.Providers {
			row := []string{p.Path, provider.Name, provider.Version}
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}
	return writer.Error()
}
