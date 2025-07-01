<p align="center">
  <img src="assets/mascot.png" alt="tfinspector mascot" width="200"/>
</p>

<h1 align="center">tfinspector</h1>

<p align="center">
  A simple CLI tool to inspect Terraform files for runtime and provider information.
</p>

<p align="center">
  <em>Lightweight. Minimal. Actively developed.</em>
</p>

<p align="center">
  <a href="https://github.com/Embe-Codez/tfinspector/actions">
    <img alt="CI" src="https://github.com/Embe-Codez/tfinspector/actions/workflows/ci.yml/badge.svg">
  </a>
  <a href="https://golang.org/">
    <img alt="Go version" src="https://img.shields.io/badge/go-1.23.8-blue.svg">
  </a>
  <a href="LICENSE">
    <img alt="License: Apache 2.0" src="https://img.shields.io/badge/license-Apache%202.0-green.svg">
  </a>
  <img alt="Status: Active Dev" src="https://img.shields.io/badge/status-active--development-orange">
</p>

---

## 📚 Table of Contents

- [Why tfinspector?](#-why-tfinspector)
- [Features](#-features)
- [Installation](#-installation)
- [Example Usage](#-example-usage)
- [Planned Features](#-planned-features)
- [Development](#-development)
- [License](#-license)

---

**tfinspector** is a CLI tool that scans Terraform configuration files (`.tf`) and extracts metadata like:

- Required Terraform version
- Provider names and versions

It’s designed to be lightweight, CI-friendly, and easy to extend.

> ⚠️ This project is under active development. Expect rapid iteration and new features.

---

## ❓ Why tfinspector?

Checking Terraform files by hand takes time and can lead to mistakes — especially in large projects.

**tfinspector** makes this easy. It scans `.tf` files and gives you clear, structured info about your setup.

Useful for:

- Reviewing your infrastructure  
- Running checks in CI/CD  
- Managing many Terraform projects  



## ✨ Features

- Recursively scan `.tf` files in a local directory
- Detect `required_version` in `terraform` blocks
- Extract providers from `required_providers` and `provider` blocks
- Output results as:
  - Plain text
  - JSON
  - YAML
  - CSV
---

## 📦 Installation

```bash
Clone and build locally:
git clone https://github.com/your-username/tfinspector.git
cd tfinspector
make install

# Now you can run it from anywhere:
tfinspector scan ./my-tf-dir

Or use Docker:
make docker
docker run --rm tfinspector:latest
```

## 🔍 Example Usage
```
# Basic scan
tfinspector scan ./infra

# Export as JSON
tfinspector scan ./infra --output json

# Save results to file
tfinspector scan ./infra --output csv --out report.csv
```

## 🔭 Planned Features

### 🔍 Smarter Scanning
- Check provider versions against the latest from the Terraform Registry
- Show update suggestions for outdated providers

### 📤 More Output Options
- Export to more formats like:
  - PDF
  - Excel (XLSX)

### 🛡️ Policy Rules
- Support custom rules in a `.tfinspector.rules.yaml` file
- Let users define required or banned providers, versions, etc.

### 📁 Multi-Repo Support
- Scan multiple Terraform repos using:
  - CSV or YAML files
  - GitHub organizations

- Ignore files or folders using `.tfinspectorignore`
- Add flags for:
  - Excluding paths
  - Limiting recursion
  - Filtering results

### ⚙️ CI/CD Ready
- Output structured data (JSON, YAML, CSV)
- Use exit codes to enforce rules in pipelines
- Easy to integrate with GitHub Actions, GitLab, and others

## 🛠 Development

- Written in Go
- Clean, testable structure
- Static binary build (`make build`)
- Docker support (`make docker`)
- Development Makefile tasks:
  - `make fmt` – format code
  - `make vet` – run static analysis
  - `make test` – run unit tests

> ⚠️ Currently tested on **Linux only**.  
> For more details, see [CODE_OVERVIEW.md](CODE_OVERVIEW.md) and [DEVELOPMENT.md](DEVELOPMENT.md).

## 📄 License

Licensed under the [Apache 2.0 License](LICENSE).