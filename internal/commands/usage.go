package commands

import "fmt"

func Usage() {
	fmt.Println("mini-terra commands:")
	fmt.Println("  init               Initialize working directory (.mini-terra)")
	fmt.Println("  init-project       Create mini-terra/ project folder with template files")
	fmt.Println("  plan   -config -var-file   Show plan for config")
	fmt.Println("  apply  -config -var-file   Apply plan and update state")
	fmt.Println("  destroy -config -var-file  Destroy resources described in config")
	fmt.Println("  show               Show current state")
	fmt.Println("  version            Show version")
}
