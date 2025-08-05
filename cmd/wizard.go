// cmd/wizard.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var wizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Run the guided wizard interface",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[wizard] Placeholder: Interactive guided mode coming soon...")
	},
}

func init() {
	rootCmd.AddCommand(wizardCmd)
}
