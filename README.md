mini-terra — mini Terraform-like IaC in Go (MVP)

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
