package main

import (
	_ "embed"
	"os"

	"github.com/dev-hak/mini-terraform/internal/commands"
	"github.com/dev-hak/mini-terraform/internal/providers"
	awsProvider "github.com/dev-hak/mini-terraform/internal/providers/aws"
	dockerProvider "github.com/dev-hak/mini-terraform/internal/providers/docker"
	vpsProvider "github.com/dev-hak/mini-terraform/internal/providers/vps"
)

func main() {
	if len(os.Args) < 2 {
		commands.Usage()
		os.Exit(1)
	}
	// register providers
	providers.RegisterProvider("docker", dockerProvider.NewDockerProvider())
	providers.RegisterProvider("vps", vpsProvider.NewVPSProvider())
	providers.RegisterProvider("aws", awsProvider.NewAWSProvider())

	cmd := os.Args[1]
	switch cmd {
	case "init":
		commands.InitCmd()
	case "plan":
		commands.PlanCmd()
	case "apply":
		commands.ApplyCmd()
	case "destroy":
		commands.DestroyCmd()
	case "show":
		commands.ShowCmd()
	case "version":
		commands.Version()
	default:
		commands.Usage()
	}
}
