// cmd/shift.go
package cmd

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code39"
	"github.com/spf13/cobra"
)

var shiftCmd = &cobra.Command{
	Use:   "shift <auction name> <start lot> <shift amount>",
	Short: "Shift lot numbers and image filenames by a given offset",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		auctionName := args[0]

		// TODO: Needs to check args index are in bounds
		shift(auctionName, []string{args[1], args[2]})

	},
}

func shift(auctionName string, args []string) error {
	startLot, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid start lot")
	}
	shiftAmount, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("Invalid shift amount")
	}
	if startLot < 0 || shiftAmount <= 0 {
		return fmt.Errorf("Start lot must be >= 0 and shift amount must be > 0")
	}

	lotDir := filepath.Join(auctionName, "03_Barcoded_Lot_Photos")
	entries, err := os.ReadDir(lotDir)
	if err != nil {
		return fmt.Errorf("Error reading lot photo folder: %v\n", err)
	}

	lotFolders := []string{}
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) == 3 && isDigits(entry.Name()) {
			lotFolders = append(lotFolders, entry.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(lotFolders))) // process highest first

	for _, oldLot := range lotFolders {
		oldNum, _ := strconv.Atoi(oldLot)
		if oldNum < startLot {
			continue
		}
		oldPath := filepath.Join(lotDir, oldLot)
		newLot := fmt.Sprintf("%03d", oldNum+shiftAmount)
		newPath := filepath.Join(lotDir, newLot)
		if err := os.Rename(oldPath, newPath); err != nil {
			fmt.Printf("Error renaming folder %s to %s: %v\n", oldLot, newLot, err)
			continue
		}

		inner, _ := os.ReadDir(newPath)
		for _, f := range inner {
			oldName := f.Name()
			if strings.HasPrefix(oldName, oldLot+"_") {
				suffix := strings.TrimPrefix(oldName, oldLot+"_")
				newName := newLot + "_" + suffix
				os.Rename(filepath.Join(newPath, oldName), filepath.Join(newPath, newName))
			}
		}

		// Always regenerate barcode for the shifted lot
		barcodeValue := strconv.Itoa(oldNum + shiftAmount)
		code, err := code39.Encode(barcodeValue, false, false)
		if err != nil {
			fmt.Printf("Error re-generating barcode for %s: %v\n", newLot, err)
			continue
		}
		scaled, err := barcode.Scale(code, 400, 100)
		if err != nil {
			fmt.Printf("Error scaling barcode for %s: %v\n", newLot, err)
			continue
		}
		barcodePath := filepath.Join(newPath, "z_barcode.png")
		file, err := os.Create(barcodePath)
		if err != nil {
			fmt.Printf("Error writing barcode for %s: %v\n", newLot, err)
			continue
		}
		png.Encode(file, scaled)
		file.Close()
	}

	// Fill in new blank folders with barcodes
	for i := 0; i < shiftAmount; i++ {
		lotNum := startLot + i
		lot := fmt.Sprintf("%03d", lotNum)
		path := filepath.Join(lotDir, lot)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.Mkdir(path, 0755)

			barcodeValue := strconv.Itoa(lotNum)
			code, err := code39.Encode(barcodeValue, false, false)
			if err != nil {
				fmt.Printf("Error generating barcode for %s: %v\n", lot, err)
				continue
			}
			scaled, err := barcode.Scale(code, 400, 100)
			if err != nil {
				fmt.Printf("Error scaling barcode for %s: %v\n", lot, err)
				continue
			}
			file, err := os.Create(filepath.Join(path, "z_barcode.png"))
			if err != nil {
				fmt.Printf("Error writing barcode file for %s: %v\n", lot, err)
				continue
			}
			png.Encode(file, scaled)
			file.Close()
		}
	}

	fmt.Println("Shift completed successfully.")

	return nil
}

func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func init() {
	rootCmd.AddCommand(shiftCmd)
}
