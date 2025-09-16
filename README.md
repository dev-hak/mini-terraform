# mini-terra — mini Terraform-like IaC in Go (MVP)

## Installation

### With Go Direct
To install the `mini-terraform` tool directly using Go, run the following command:

```bash
go install github.com/dev-hak/mini-terraform@latest
```

### macOS
Download and install `mini-terraform` using cURL:

```bash
curl -L https://github.com/dev-hak/mini-terraform/releases/latest/download/mini-terraform-darwin-amd64 -o /usr/local/bin/mini-terraform
chmod +x /usr/local/bin/mini-terraform
```

### Windows
Install `mini-terraform` with PowerShell:

```powershell
Invoke-WebRequest -Uri https://github.com/dev-hak/mini-terraform/releases/latest/download/mini-terraform-windows-amd64.exe -OutFile mini-terraform.exe
Move-Item -Path mini-terraform.exe -Destination C:\Windows\System32
```

### Linux
Download and install `mini-terraform` using cURL:

```bash
curl -L https://github.com/dev-hak/mini-terraform/releases/latest/download/mini-terraform-linux-amd64 -o /usr/local/bin/mini-terraform
chmod +x /usr/local/bin/mini-terraform
```

Features:
- JSON config + var-files
- Commands: init, init-project, plan, apply, destroy, show, version
- Local JSON state (.mini-terra/mini-terra.state.json)
- Providers: docker (implemented), vps (ssh exec), aws (skeleton)

Build:
  go build ./cmd/mini-terra

Examples:
  ./mini-terra init
  ./mini-terra plan -config examples/config.json -var-file examples/vars.json
  ./mini-terra apply -config examples/config.json -var-file examples/vars.json
  ./mini-terra show
  ./mini-terra destroy -config examples/config.json -var-file examples/vars.json
  ./mini-terra version

Notes:
- Docker provider uses the docker CLI; ensure docker is installed and accessible.
- VPS provider expects a private key path and SSH access.
- State file contains sensitive info; do not commit to VCS.
