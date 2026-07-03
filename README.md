# cloudexec

[![Go Version](https://img.shields.io/github/go-mod/go-version/n40y/cloudexec)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`cloudexec` is a fast, concurrent, multi-cloud security scanner and offensive audit framework written in Go. Designed for cloud security engineers and pentester specialists, it combines automated OSINT, signature-based template scanning, and a multi-threaded secrets finder with automated real-time pivot and post-authentication enumeration capability.

---

## Legal Disclaimer

> [!IMPORTANT]
> **Legal Disclaimer**: This tool is developed strictly for educational purposes, authorized penetration testing, and security auditing. The author (`n40y`) assumes no liability for any unauthorized misuse, damage, or illegal activities caused by this tool. Usage of `cloudexec` for attacking targets without prior mutual consent is illegal. Users are solely responsible for complying with all applicable local and international laws.

## Key Features

* **⚡ High-Performance Secrets Scanner**: Multi-threaded worker pool to hunt down leaked credentials (AWS, GCP, Slack, etc.) across local directories.
* **🔄 Automated Pivot Mode**: Validates discovered secrets on-the-fly against live cloud provider APIs during the filesystem scan without blocking the pipeline.
* **🕵️ Post-Auth Enumeration**: Automatically maps out accessible services, permissions, and attack paths the moment a valid credential is confirmed.
* **🛠️ YAML Template-Driven Core**: Decoupled engine logic using modular YAML signatures for flexible multi-cloud checks.
* **☁️ Multi-Cloud Support**:
    * **AWS**: STS identity verification (`GetCallerIdentity`) and S3 exposure auditing.
    * **Azure**: Domain OSINT, Tenant ID enumeration, break-glass account identification, and blob storage mapping.
    * **GCP**: Google Workspace public infrastructure discovery and API Key validation (`Identity Toolkit`).

---

## 🛠️ Customizing Templates (Playbooks)

`cloudexec` uses a modular YAML template engine. You can easily extend the scanner's capabilities by adding your own templates into the provider subdirectories (e.g., `templates/gcp/`, `templates/aws/`).

### Template Structure

Every template must follow this anatomy:

```yaml
id: custom-gcp-check
engine: gcp                 # Target engine: aws, azure, gcp, cloudflare, recon
action: gcp:ApiKeyCheck     # The internal method to execute
info:
  name: "Custom API Key Validation Workflow"
  description: "Triggers a validation request and maps downstream permissions."
  severity: "high"
```

## Installation

Ensure you have [Go](https://go.dev/doc/install) installed (version 1.20+ recommended).

```bash
# Clone the repository
git clone [https://github.com/n40y/cloudexec.git](https://github.com/n40y/cloudexec.git)

cd cloudexec
# Build the binary
go build
```

## 🚀 Usage Guide

### 1. Global Secrets Scanner (with Auto-Pivot)
Scan a directory for leaked credentials, automatically intercept them, and validate them against cloud provider APIs in real-time.
```bash
# Basic scan
./cloudexec secrets --path .

# Scan with auto-pivot validation and JSON export
./cloudexec secrets --path /path/to/project -p -o report.json
```

### 2. AWS Audit Engine

Executes AWS-specific templates (identity verification via **STS**, S3 bucket enumeration). It automatically utilizes credentials configured in your **config.yaml** or standard local AWS environment variables.

```bash
./cloudexec aws
```

### 3. Azure OSINT & Enumeration Engine

Performs passive and active reconnaissance on a Microsoft 365 / Azure tenant using a target domain. Extracts Tenant ID, authentication mechanics, identity provider details, and lists potential break-glass accounts.

```bash
./cloudexec azure -d targetcompany.com
```

### 4. GCP Audit Engine

Validates Google Cloud API keys and checks if a domain is mapped to a public Google Workspace infrastructure.

```bash
# Check both Workspace infrastructure and a specific API Key
./cloudexec gcp -d targetcompany.com --apikey AIzaSy...

# Check an isolated API key only
./cloudexec gcp --apikey AIzaSy...
```

### 5. General Recon Engine

Triggers passive multi-cloud discovery templates, including historical DNS lookups, certificate transparency logs via **crt.sh**, and cross-provider storage bucket detection.

```bash
./cloudexec recon -d targetcompany.com
```


## 📁 Project Structure

| Path | Description |
| :--- | :--- |
| **`cmd/`** | Contains Cobra CLI command definitions (`aws`, `azure`, `gcp`, `recon`, `secrets`). |
| **`pkg/engines/`** | Core execution logic for cloud provider validation and post-auth enumeration. |
| **`pkg/templates/`** | YAML template parsing and signature matching engines. |
| **`pkg/utils/`** | Thread-safe logging utilities and console formatting. |
| **`templates/`** | Ready-to-use security playbooks and detection signatures (.yaml). |
| **`main.go`** | Application entry point. |


## Configuration

For authenticated modules, you can maintain a local config.yaml file at the root.

    [!WARNING]
    Never commit your config.yaml to a public repository. Ensure it is included in your .gitignore.

```yaml
aws:
  access_key: "AKIA..."
  secret_key: "..."
  region: "eu-west-3"
```

## License

Distributed under the MIT License. See LICENSE for more information.
