# Code Overview

This document explains the main parts of **tfinspector** and how they work together.

---

## 📁 Project Structure

- **main.go**  
  Entry point of the application. Starts the CLI.

- **cmd/**  
  Contains the CLI commands using [Cobra](https://github.com/spf13/cobra).  
  Handles command parsing, flags, and user input.

  - `scan.go`  
    Main command that:
    - Scans a directory for `.tf` files
    - Calls the scanner logic
    - Formats and outputs the results (text, JSON, YAML, CSV)

- **internal/**  
  Core logic, organized by responsibility:

  - **scanner/**  
    - Walks directories and finds `.tf` files
    - Uses the HCL parser to extract:
      - Required Terraform version
      - Providers and versions
    - Defines main types:
      - `TerraformProject`
      - `ProviderInfo`

> 🚧 This structure may change as the project grows and new features are added.

---

## Key Data Structures

These are defined in the `scanner` package:

```go
type ProviderInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type TerraformProject struct {
	Path           string         `json:"path"`
	RuntimeVersion string         `json:"runtime_version,omitempty"`
	Providers      []ProviderInfo `json:"providers"`
}
```

### TerraformProject 
Represents a scanned .tf file or Terraform project:
 - Path: location of the file.
 - RuntimeVersion: required Terraform version (if declared).
 - Providers: list of providers and their version constraints.

### ProviderInfo
Holds metadata about a single provider:
 - Name: provider name (e.g., aws, azurerm).
 - Version: optional version constraint (e.g., >= 3.0.0).

## How It Works

1. **CLI layer** (`cmd/`) accepts commands like `scan` and a directory path.  
2. The CLI calls the **scanner** to process the directory.  
3. **Scanner** walks the directory tree, finds `.tf` files, and sends each to the **parser**.  
4. **Parser** reads each file’s content using the HashiCorp HCL parser and extracts:
   - Terraform runtime version
   - Provider names and versions  
5. The scanner collects all this info into **TerraformProject** data structures.  
6. The CLI receives these results and prints or exports them (future feature).

## Extending the Code

- To add new CLI commands, create a new file in cmd/ and add the command using RootCmd.AddCommand(...) in the init() function.
- To support exporting results, implement new exporters in `internal/exporter/`.  
- To query the Terraform Registry for latest provider versions, add code under `internal/resolver/`.  
- To improve parsing (e.g., add module support), enhance the parser logic in `internal/scanner/parser.go`.

---

This high-level overview should help new contributors and maintainers understand how tfinspector is organized and where to focus when adding features or fixing bugs.