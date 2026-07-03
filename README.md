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

## Installation

Ensure you have [Go](https://go.dev/doc/install) installed (version 1.20+ recommended).

```bash
# Clone the repository
git clone [https://github.com/n40y/cloudexec.git](https://github.com/n40y/cloudexec.git)

cd cloudexec
# Build the binary
go build
```

## Usage Examples


### 1. Scan for Secrets with Auto-Pivot & Post-Auth Verification

Scan a local repository, intercept Google API keys or AWS credentials, and validate them immediately against official endpoints:

```bash
./cloudexec secrets --path /path/to/sourcecode -p
```

### 2. Google Cloud Platform Audit

Run Google Workspace tenant mapping and check a specific API key:

```bash
./cloudexec gcp -d targetcompany.com --apikey AIzaSy...
```

### 3. Save Audit Reports

Export results into a structured JSON report for post-assessment processing:

```bash
./cloudexec secrets --path . --output report.json
```


## Project Structure

├── cmd/                # Cobra CLI commands definition
├── pkg/
│   ├── engines/        # Provider logic (AWS, Azure, GCP, Post-Auth)
│   ├── templates/      # YAML parser and signature matchers
│   └── utils/          # Thread-safe logging utilities
├── templates/          # Ready-to-use multi-cloud YAML scanner templates
└── main.go             # Application entry point


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
