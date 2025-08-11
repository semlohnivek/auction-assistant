// cmd/root.go
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "auction",
	Short: "Bidz Auction Utility",
	Long:  `A command-line tool to help prepare and manage photo-based auctions.`,
	Run: func(cmd *cobra.Command, args []string) {
		// if len(args) == 0 {
		// 	cmd.Println("Launching guided wizard mode...")
		// 	RunGuidedMode()
		// 	//wizardCmd.Run(cmd, args)
		// } else {
		cmd.Help()
		// }
	},
}

func Execute() error {
	return rootCmd.Execute()
}
