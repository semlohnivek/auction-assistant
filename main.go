// main.go
package main

import (
	"fmt"
	"os"

	"bidzauction/cmd"
	"bidzauction/config"
	"path/filepath"
)

func main() {

	cfgPath := filepath.Join("config.toml")
	if err := config.Load(cfgPath); err != nil {
		fmt.Printf("failed to load configuration. Could not find %v", err)
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		cmd.RunGuidedMode()
		return
	} else if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
