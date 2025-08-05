// cmd/root.go
package cmd

import (
	"fmt"
	"path/filepath"

	"bidzauction/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "auction",
	Short: "Bidz Auction Utility",
	Long:  `A command-line tool to help prepare and manage photo-based auctions.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := filepath.Join("config.toml")
		if err := config.Load(cfgPath); err != nil {
			return fmt.Errorf("failed to load configuration: %v", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Println("Launching guided wizard mode...")
			wizardCmd.Run(cmd, args)
		} else {
			cmd.Help()
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}
