package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var flattenCmd = &cobra.Command{
	Use:   "flatten <auction-name>",
	Short: "Flatten images from 03_Barcoded_Lot_Photos into 04_To_Upload_Photos and 05_AI_Analysis_Photos with renamed output",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		auctionName := args[0]

		if err := flatten(auctionName); err != nil {
			fmt.Printf("Error: %v\n", err)
		}

	},
}

func flatten(auctionName string) error {
	inputRoot := filepath.Join(auctionName, "03_Barcoded_Lot_Photos")
	outputs := []string{
		filepath.Join(auctionName, "04_To_Upload_Photos"),
		filepath.Join(auctionName, "05_AI_Analysis_Photos"),
	}

	err := os.MkdirAll(outputs[0], 0755)
	if err != nil {
		return fmt.Errorf("Error creating output folder:", err)
	}
	err = os.MkdirAll(outputs[1], 0755)
	if err != nil {
		return fmt.Errorf("Error creating output folder:", err)
	}

	entries, err := os.ReadDir(inputRoot)
	if err != nil {
		return fmt.Errorf("Error reading input folder:", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		lot := entry.Name()
		lotPath := filepath.Join(inputRoot, lot)

		files, err := os.ReadDir(lotPath)
		if err != nil {
			fmt.Printf("Skipping lot %s due to error: %v\n", lot, err)
			continue
		}

		var images []os.DirEntry
		var barcode os.DirEntry

		for _, file := range files {
			name := file.Name()
			if name == "z_barcode.png" {
				barcode = file
			} else if !file.IsDir() && isImageFile(name) {
				images = append(images, file)
			}
		}

		fmt.Println(len(images))

		// Skip folders that are empty or only have a barcode
		if barcode == nil && len(images) == 0 {
			fmt.Printf("Skipping lot %s — no usable images found.\n", lot)
			continue
		}
		if barcode != nil && len(images) == 0 {
			fmt.Printf("Skipping lot %s — contains only barcode image.\n", lot)
			continue
		}

		sort.Slice(images, func(i, j int) bool {
			return images[i].Name() < images[j].Name()
		})

		imgIndex := 1

		if barcode != nil {
			flattenCopy(filepath.Join(lotPath, barcode.Name()), outputs, fmt.Sprintf("%s_img0_barcode.jpg", lot))
		}

		for _, file := range images {
			newName := fmt.Sprintf("%s_img%d.jpg", lot, imgIndex)
			srcPath := filepath.Join(lotPath, file.Name())
			flattenCopy(srcPath, outputs, newName)
			imgIndex++
		}
	}

	fmt.Println("Images flattened into Upload and AI Analysis folders.")

	return nil
}

func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}

func flattenCopy(src string, destFolders []string, filename string) {
	for _, dest := range destFolders {
		destPath := filepath.Join(dest, filename)

		in, err := os.Open(src)
		if err != nil {
			fmt.Printf("Failed to open %s: %v\n", src, err)
			continue
		}
		defer in.Close()

		out, err := os.Create(destPath)
		if err != nil {
			fmt.Printf("Failed to create %s: %v\n", destPath, err)
			in.Close()
			continue
		}

		_, err = io.Copy(out, in)
		if err != nil {
			fmt.Printf("Failed to copy %s to %s: %v\n", src, destPath, err)
		} else {
			fmt.Printf("Copied → %s\n", destPath)
		}

		in.Close()
		out.Close()
	}
}

func init() {
	rootCmd.AddCommand(flattenCmd)
}
