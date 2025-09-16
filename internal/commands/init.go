package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/dev-hak/mini-terraform/internal/state"
)

func InitCmd() {
	if err := os.MkdirAll(".mini-terra", 0755); err != nil {
		log.Fatalf("init: %v", err)
	}

	statePath := ".mini-terra/mini-terra.state.json"
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		if err := state.SaveState(statePath, state.NewEmptyState()); err != nil {
			log.Fatalf("init save state: %v", err)
		}
		fmt.Println("Initialized .mini-terra and state file.")
	} else {
		fmt.Println(".mini-terra already initialized.")
	}
}
