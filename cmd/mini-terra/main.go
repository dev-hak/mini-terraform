package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dev-hak/mini-terraform/internal/config"
	"github.com/dev-hak/mini-terraform/internal/engine"
	"github.com/dev-hak/mini-terraform/internal/providers"
	awsProvider "github.com/dev-hak/mini-terraform/internal/providers/aws"
	dockerProvider "github.com/dev-hak/mini-terraform/internal/providers/docker"
	vpsProvider "github.com/dev-hak/mini-terraform/internal/providers/vps"
	"github.com/dev-hak/mini-terraform/internal/state"
)

//go:embed template/mini-terra-config.json
var tmplConfig string

//go:embed template/mini-terra-vars.json
var tmplVars string

const Version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	// register providers
	providers.RegisterProvider("docker", dockerProvider.NewDockerProvider())
	providers.RegisterProvider("vps", vpsProvider.NewVPSProvider())
	providers.RegisterProvider("aws", awsProvider.NewAWSProvider())

	cmd := os.Args[1]
	switch cmd {
	case "init":
		initCmd()
	case "init-project":
		initProjectCmd()
	case "plan":
		planCmd()
	case "apply":
		applyCmd()
	case "destroy":
		destroyCmd()
	case "show":
		showCmd()
	case "version":
		fmt.Println("mini-terra", Version)
	default:
		usage()
	}
}

func usage() {
	fmt.Println("mini-terra commands:")
	fmt.Println("  init               Initialize working directory (.mini-terra)")
	fmt.Println("  init-project       Create mini-terra/ project folder with template files")
	fmt.Println("  plan   -config -var-file   Show plan for config")
	fmt.Println("  apply  -config -var-file   Apply plan and update state")
	fmt.Println("  destroy -config -var-file  Destroy resources described in config")
	fmt.Println("  show               Show current state")
	fmt.Println("  version            Show version")
}

func initCmd() {
	if err := os.MkdirAll(".mini-terra", 0755); err != nil {
		log.Fatalf("init: %v", err)
	}
	// initialize empty state if not exists
	statePath := ".mini-terra/terraform.tfstate.json"
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		if err := state.SaveState(statePath, state.NewEmptyState()); err != nil {
			log.Fatalf("init save state: %v", err)
		}
		fmt.Println("Initialized .mini-terra and state file.")
	} else {
		fmt.Println(".mini-terra already initialized.")
	}
}

func initProjectCmd() {
	outDir := "mini-terra"
	if _, err := os.Stat(outDir); err == nil {
		fmt.Printf("'%s' already exists — aborting\n", outDir)
		return
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("create project dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(outDir, "config.json"), []byte(tmplConfig), 0644); err != nil {
		log.Fatalf("write config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "vars.json"), []byte(tmplVars), 0644); err != nil {
		log.Fatalf("write vars: %v", err)
	}

	runSh := `#!/usr/bin/env bash
BIN="./mini-terra/mini-terra"
if [ -f "./mini-terra" ]; then
  ./mini-terra "$@"
  exit $?
fi
if [ -f "$BIN" ]; then
  "$BIN" "$@"
  exit $?
fi
if command -v mini-terra >/dev/null 2>&1; then
  mini-terra "$@"
  exit $?
fi
echo "mini-terra binary not found. Build it: go build ./cmd/mini-terra"
exit 1
`
	if err := os.WriteFile(filepath.Join(outDir, "run.sh"), []byte(runSh), 0755); err != nil {
		log.Fatalf("write run.sh: %v", err)
	}

	runBat := `@echo off
if exist mini-terra (
  mini-terra %*
  exit /B %ERRORLEVEL%
)
where mini-terra >nul 2>nul
if %ERRORLEVEL%==0 (
  mini-terra %*
  exit /B %ERRORLEVEL%
)
echo mini-terra binary not found. Build it: go build ./cmd/mini-terra
exit /B 1
`
	if err := os.WriteFile(filepath.Join(outDir, "run.bat"), []byte(runBat), 0644); err != nil {
		log.Fatalf("write run.bat: %v", err)
	}

	readme := `# mini-terra (project folder)
This folder contains mini-terra configuration for this project.

Usage:
  ./run.sh plan -config config.json -var-file vars.json
  ./run.sh apply -config config.json -var-file vars.json

Don't commit secrets or state files.
`
	if err := os.WriteFile(filepath.Join(outDir, "README.md"), []byte(readme), 0644); err != nil {
		log.Fatalf("write README: %v", err)
	}

	fmt.Println("Created 'mini-terra/' with example config and run scripts.")
}

func planCmd() {
	fs := flag.NewFlagSet("plan", flag.ExitOnError)
	cfgPath := fs.String("config", "", "path to config json")
	varFile := fs.String("var-file", "", "path to vars json")
	fs.Parse(os.Args[2:])
	if *cfgPath == "" {
		log.Fatal("plan: -config is required")
	}
	cfg, err := config.LoadConfig(*cfgPath, *varFile)
	if err != nil {
		log.Fatalf("plan load config: %v", err)
	}
	st, err := state.LoadState(".mini-terra/terraform.tfstate.json")
	if err != nil {
		log.Fatalf("plan load state: %v", err)
	}
	plan, err := engine.GeneratePlan(cfg, st)
	if err != nil {
		log.Fatalf("plan: %v", err)
	}
	fmt.Println("Plan:")
	for _, op := range plan.Operations {
		fmt.Printf("  - %s: %s.%s\n", op.Action, op.Resource.Type, op.Resource.Name)
	}
}

func applyCmd() {
	fs := flag.NewFlagSet("apply", flag.ExitOnError)
	cfgPath := fs.String("config", "", "path to config json")
	varFile := fs.String("var-file", "", "path to vars json")
	fs.Parse(os.Args[2:])
	if *cfgPath == "" {
		log.Fatal("apply: -config is required")
	}
	cfg, err := config.LoadConfig(*cfgPath, *varFile)
	if err != nil {
		log.Fatalf("apply load config: %v", err)
	}
	st, err := state.LoadState(".mini-terra/terraform.tfstate.json")
	if err != nil {
		log.Fatalf("apply load state: %v", err)
	}
	plan, err := engine.GeneratePlan(cfg, st)
	if err != nil {
		log.Fatalf("apply plan: %v", err)
	}
	newState, err := engine.Apply(plan, st)
	if err != nil {
		log.Fatalf("apply execute: %v", err)
	}
	if err := state.SaveState(".mini-terra/terraform.tfstate.json", newState); err != nil {
		log.Fatalf("apply save state: %v", err)
	}
	fmt.Println("Apply complete. State saved to .mini-terra/terraform.tfstate.json")
}

func destroyCmd() {
	fs := flag.NewFlagSet("destroy", flag.ExitOnError)
	cfgPath := fs.String("config", "", "path to config json")
	varFile := fs.String("var-file", "", "path to vars json")
	fs.Parse(os.Args[2:])
	if *cfgPath == "" {
		log.Fatal("destroy: -config is required")
	}
	cfg, err := config.LoadConfig(*cfgPath, *varFile)
	if err != nil {
		log.Fatalf("destroy load config: %v", err)
	}
	st, _ := state.LoadState(".mini-terra/terraform.tfstate.json")
	plan, err := engine.GeneratePlanForDestroy(cfg, st)
	if err != nil {
		log.Fatalf("destroy plan: %v", err)
	}
	newState, err := engine.Apply(plan, st)
	if err != nil {
		log.Fatalf("destroy apply: %v", err)
	}
	if err := state.SaveState(".mini-terra/terraform.tfstate.json", newState); err != nil {
		log.Fatalf("destroy save state: %v", err)
	}
	fmt.Println("Destroy complete. State updated.")
}

func showCmd() {
	st, err := state.LoadState(".mini-terra/terraform.tfstate.json")
	if err != nil {
		fmt.Println("No state found.")
		return
	}
	b, _ := state.PrettyJSON(st)
	fmt.Println(string(b))
}
