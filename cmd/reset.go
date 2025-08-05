// cmd/reset.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset an auction folder to initial state",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[reset] Placeholder: Cleanup folders and reset state")
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
