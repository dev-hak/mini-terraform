package commands

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dev-hak/mini-terraform/internal/config"
	"github.com/dev-hak/mini-terraform/internal/engine"
	"github.com/dev-hak/mini-terraform/internal/state"
)

func ApplyCmd() {
	fs := flag.NewFlagSet("apply", flag.ExitOnError)
	cfgPath := fs.String("config", "config.json", "path to config json")
	varFile := fs.String("var-file", "vars.json", "path to vars json")
	fs.Parse(os.Args[2:])
	if *cfgPath == "" {
		log.Fatal("apply: -config is required")
	}
	cfg, err := config.LoadConfig(*cfgPath, *varFile)
	if err != nil {
		log.Fatalf("apply load config: %v", err)
	}
	st, err := state.LoadState(".mini-terra/mini-terra.state.json")
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
	if err := state.SaveState(".mini-terra/mini-terra.state.json", newState); err != nil {
		log.Fatalf("apply save state: %v", err)
	}
	fmt.Println("Apply complete. State saved to .mini-terra/mini-terra.state.json")
}
