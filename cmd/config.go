// cmd/config.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configCreateCmd = &cobra.Command{
	Use:   "config create",
	Short: "Generate a default config.toml file in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		defaultConfig := `
[resize]
width = 1600
height = 1200
ascii_only = true
strip_metadata = true

[sftp]
host = "media.bidzauctionhouse.com"
port = 22
user = "auction"
pass = "plaintext-ok-for-now"
remote_path = "/var/www/media/uploads"

[openai]
model = "gpt-4-vision-preview"
key = "sk-..."
max_tokens = 500
temperature = 0.7

[flatten]
include_barcode_in_upload = true
include_barcode_in_analysis = true
barcode_first = true
lot_prefix = "lot"
barcode_suffix = "_img0_barcode"
image_prefix = "img"
image_ext = ".jpg"

[defaults]
auction_size = 300
`
		if err := os.WriteFile("config.toml", []byte(defaultConfig), 0644); err != nil {
			fmt.Printf("Error writing config file: %v\n", err)
			return
		}
		fmt.Println("Default config.toml created.")
	},
}

func init() {
	rootCmd.AddCommand(configCreateCmd)
}
