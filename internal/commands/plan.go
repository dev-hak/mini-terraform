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

func PlanCmd() {
	fs := flag.NewFlagSet("plan", flag.ExitOnError)
	cfgPath := fs.String("config", "config.json", "path to config json")
	varFile := fs.String("var-file", "vars.json", "path to vars json")
	fs.Parse(os.Args[2:])
	if *cfgPath == "" {
		log.Fatal("plan: -config is required")
	}
	cfg, err := config.LoadConfig(*cfgPath, *varFile)
	if err != nil {
		log.Fatalf("plan load config: %v", err)
	}
	st, err := state.LoadState(".mini-terra/mini-terra.state.json")
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
