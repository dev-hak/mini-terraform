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

func DestroyCmd() {
	fs := flag.NewFlagSet("destroy", flag.ExitOnError)
	cfgPath := fs.String("config", "config.json", "path to config json")
	varFile := fs.String("var-file", "vars.json", "path to vars json")
	fs.Parse(os.Args[2:])
	if *cfgPath == "" {
		log.Fatal("destroy: -config is required")
	}
	cfg, err := config.LoadConfig(*cfgPath, *varFile)
	if err != nil {
		log.Fatalf("destroy load config: %v", err)
	}
	st, _ := state.LoadState(".mini-terra/mini-terra.state.json")
	plan, err := engine.GeneratePlanForDestroy(cfg, st)
	if err != nil {
		log.Fatalf("destroy plan: %v", err)
	}
	newState, err := engine.Apply(plan, st)
	if err != nil {
		log.Fatalf("destroy apply: %v", err)
	}
	if err := state.SaveState(".mini-terra/mini-terra.state.json", newState); err != nil {
		log.Fatalf("destroy save state: %v", err)
	}
	fmt.Println("Destroy complete. State updated.")
}
