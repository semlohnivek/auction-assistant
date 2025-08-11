package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"

	"bidzauction/config"
)

var resizeCmd = &cobra.Command{
	Use:   "resize <auction-name>",
	Short: "Resize images from 01_Original_Photos into 02_Resized_Photos targeting ~250-300KB",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		auctionName := args[0]

		if err := resize(auctionName); err != nil {
			fmt.Printf("Error: %v\n", err)
		}

	},
}

func resize(auctionName string) error {
	originalRoot := filepath.Join(auctionName, "01_Original_Photos")
	resizedRoot := filepath.Join(auctionName, "02_Resized_Photos")

	err := filepath.Walk(originalRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		filename := info.Name()
		if strings.HasPrefix(filename, ".trash") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil
		}

		relPath, _ := filepath.Rel(originalRoot, path)
		outputPath := filepath.Join(resizedRoot, strings.TrimSuffix(relPath, ext)+".jpg")

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return err
		}

		if err := resizeAndSave(path, outputPath); err != nil {
			fmt.Printf("Failed to process %s: %v\n", path, err)
		} else {
			fmt.Printf("Resized %s → %s\n", path, outputPath)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Resize operation failed:", err)
	}

	fmt.Println("Images resized for auction ", auctionName)

	return nil
}

func resizeAndSave(inputPath, outputPath string) error {
	src, err := imaging.Open(inputPath)
	if err != nil {
		return err
	}

	for scale := 1.0; scale >= 0.3; scale -= 0.1 {
		w := src.Bounds().Dx()
		h := src.Bounds().Dy()

		resizedW := 0
		resizedH := config.Current.Resize.MaxDimension

		if w > h {
			resizedW = 1200
			resizedH = 0
		}

		dst := imaging.Resize(src, resizedW, resizedH, imaging.Lanczos)

		return imaging.Save(dst, outputPath, imaging.JPEGQuality(80))

	}

	return fmt.Errorf("unable to reduce %s below target size", inputPath)
}

func init() {
	rootCmd.AddCommand(resizeCmd)
}
