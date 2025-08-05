// cmd/init.go
package cmd

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	"bidzauction/config"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code39"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new auction folder structure",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		auctionName := "Auction"
		lotCount := config.Current.Defaults.AuctionSize

		if len(args) > 0 {
			auctionName = args[0]
		}
		if len(args) > 1 {
			if count, err := strconv.Atoi(args[1]); err == nil {
				lotCount = count
			}
		}

		fmt.Printf("[init] Creating auction '%s' with %d lots...\n", auctionName, lotCount)

		base := auctionName
		dirs := []string{
			"01_Original_Photos",
			"02_Resized_Photos",
			"03_Barcoded_Lot_Photos",
			"04_To_Upload_Photos",
			"05_AI_Analysis_Photos",
			"state",
		}

		for _, dir := range dirs {
			path := filepath.Join(base, dir)
			if err := os.MkdirAll(path, 0755); err != nil {
				fmt.Printf("Error creating %s: %v\n", path, err)
			}
		}

		barcodeRoot := filepath.Join(base, "03_Barcoded_Lot_Photos")

		for i := 1; i <= lotCount; i++ {
			barcodeFilename := "z_barcode.png"
			barcodeValue := fmt.Sprintf("%d", i)
			lotName := fmt.Sprintf("%03d", i)
			lotDir := filepath.Join(barcodeRoot, lotName)
			if err := os.MkdirAll(lotDir, 0755); err != nil {
				fmt.Printf("Error creating %s: %v\n", lotDir, err)
			}

			// Generate Code 39 barcode
			code, err := code39.Encode(barcodeValue, false, false)
			if err != nil {
				fmt.Printf("Error generating barcode for %s: %v\n", lotName, err)
				continue
			}
			codeScaled, err := barcode.Scale(code, 400, 100)
			if err != nil {
				fmt.Printf("Error scaling barcode for %s: %v\n", lotName, err)
				continue
			}

			barcodePath := filepath.Join(lotDir, fmt.Sprintf(barcodeFilename))
			file, err := os.Create(barcodePath)
			if err != nil {
				fmt.Printf("Error creating barcode file %s: %v\n", barcodePath, err)
				continue
			}
			defer file.Close()
			if err := png.Encode(file, codeScaled); err != nil {
				fmt.Printf("Error encoding barcode PNG for %s: %v\n", lotName, err)
			}
		}

		fmt.Println("Auction folder structure and barcodes created.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
