package cmd

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"
)

var resizeCmd = &cobra.Command{
	Use:   "resize <auction-name>",
	Short: "Resize images from 01_Original_Photos into 02_Resized_Photos targeting ~250-300KB",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		auctionName := args[0]
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
	},
}

func resizeAndSave(inputPath, outputPath string) error {
	src, err := imaging.Open(inputPath)
	if err != nil {
		return err
	}

	for scale := 1.0; scale >= 0.3; scale -= 0.1 {
		w := int(float64(src.Bounds().Dx()) * scale)
		h := int(float64(src.Bounds().Dy()) * scale)
		dst := imaging.Resize(src, w, h, imaging.Lanczos)

		tempFile := outputPath + ".tmp"
		out, err := os.Create(tempFile)
		if err != nil {
			return err
		}
		err = jpeg.Encode(out, dst, &jpeg.Options{Quality: 75})
		out.Close()
		if err != nil {
			return err
		}

		info, err := os.Stat(tempFile)
		if err != nil {
			return err
		}
		if info.Size() <= 300_000 {
			return os.Rename(tempFile, outputPath)
		}

		os.Remove(tempFile)
	}

	return fmt.Errorf("unable to reduce %s below target size", inputPath)
}

func init() {
	rootCmd.AddCommand(resizeCmd)
}
