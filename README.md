# mini-terra — mini Terraform-like IaC in Go (MVP)

## Installation

### With Go Direct
To install the `mini-terra` tool directly using Go, run the following command:

```bash
go install github.com/dev-hak/mini-terraform/cmd/mini-terra@latest
```

### macOS
Download and install `mini-terra` for AMD64 using cURL:

```bash
curl -sL https://github.com/dev-hak/mini-terraform/releases/download/v0.1.0/mini-terra_0.1.0_macOS_.amd64.tar.gz \
  | tar xz
sudo mv mini-terra /usr/local/bin/
```

Download and install `mini-terra` for ARM64 using cURL:

```bash
curl -sL https://github.com/dev-hak/mini-terraform/releases/download/v0.1.0/mini-terra_0.1.0_macOS_.arm64.tar.gz \
  | tar xz
sudo mv mini-terra /usr/local/bin/
```



### Windows
Install `mini-terra` for AMD64 with PowerShell:

```powershell
Invoke-WebRequest https://github.com/dev-hak/mini-terraform/releases/download/v0.1.0/mini-terra_0.1.0_Windows_.amd64.tar.gz -OutFile mini-terra.zip
Expand-Archive mini-terra.zip -DestinationPath .
Move-Item mini-terra.exe C:\Windows\System32\
```

Install `mini-terra` for ARM64 with PowerShell:

```powershell
Invoke-WebRequest https://github.com/dev-hak/mini-terraform/releases/download/v0.1.0/mini-terra_0.1.0_Windows_.arm64.tar.gz -OutFile mini-terra.zip
Expand-Archive mini-terra.zip -DestinationPath .
Move-Item mini-terra.exe C:\Windows\System32\
```

### Linux
Download and install `mini-terra` for AMD64 using cURL:

```bash
curl -sL https://github.com/dev-hak/mini-terraform/releases/download/v0.1.0/mini-terra_0.1.0_Linux_.amd64.tar.gz \
  | tar xz
sudo mv mini-terra /usr/local/bin/
```
Download and install `mini-terra` for ARM64 using cURL:

```bash
curl -sL https://github.com/dev-hak/mini-terraform/releases/download/v0.1.0/mini-terra_0.1.0_Linux_.arm64.tar.gz \
  | tar xz
sudo mv mini-terra /usr/local/bin/
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
