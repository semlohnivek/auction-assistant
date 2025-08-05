// cmd/analyze.go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"bidzauction/config"
	"bidzauction/internal"

	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze <auction name>",
	Short: "Upload to SFTP and generate AI titles/descriptions",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		auctionName := args[0]
		aiDir := filepath.Join(auctionName, "05_AI_Analysis_Photos")

		entries, err := os.ReadDir(aiDir)
		if err != nil {
			fmt.Printf("Error reading AI analysis folder: %v\n", err)
			return
		}

		lotImages := map[string][]string{}
		for _, f := range entries {
			if !f.IsDir() {
				name := strings.ToLower(f.Name())
				if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".png") {
					lotNum := strings.SplitN(name, "_", 2)[0]
					lotImages[lotNum] = append(lotImages[lotNum], name)
				}
			}
		}

		remoteDir := auctionName
		var allImages []string
		for _, imgs := range lotImages {
			allImages = append(allImages, imgs...)
		}
		sort.Strings(allImages)
		fmt.Printf("Found %d images across %d lots to analyze in %s\n", len(allImages), len(lotImages), aiDir)

		err = internal.UploadImagesToSFTP(aiDir, remoteDir, allImages)
		if err != nil {
			fmt.Printf("SFTP upload failed: %v\n", err)
			return
		}
		fmt.Println("All images uploaded successfully. Starting analysis...")

		cfg := config.Current
		lotOrder := []string{}
		lotResults := map[string]internal.LotDetails{}

		for lotNum, images := range lotImages {
			sort.Strings(images)
			lotOrder = append(lotOrder, lotNum)

			var urls []string
			for _, name := range images {
				if strings.Contains(name, "barcode") {
					continue
				}
				url := fmt.Sprintf("%s/%s/%s", cfg.OpenAI.ImageBaseUrl, auctionName, name)
				urls = append(urls, url)
			}
			if len(urls) == 0 {
				fmt.Printf("Lot %s: No suitable images for analysis\n", lotNum)
				continue
			}

			details, err := internal.AnalyzeImageURLs(urls)
			if err != nil {
				fmt.Printf("Lot %s: ERROR - %v\n", lotNum, err)
				continue
			}

			fmt.Printf("Lot %s\n  Title: %s\n  Desc: %s\n  Make: %s\n  Model: %s\n  Condition: %s\n  Year: %s\n  Country: %s\n\n",
				lotNum, details.Title, details.Description, details.Make, details.Model, details.Condition, details.Year, details.CountryOfOrigin)

			lotResults[lotNum] = details
		}

		// Write TSV
		outPath := filepath.Join(auctionName, "lots.tsv")
		f, err := os.Create(outPath)
		if err != nil {
			fmt.Printf("Error creating TSV: %v\n", err)
			return
		}
		defer f.Close()
		writer := bufio.NewWriter(f)
		writer.WriteString("Lot Number\tItem Title\tSeller ID\tDescription\n")

		sort.Strings(lotOrder)
		for _, lotNum := range lotOrder {
			d := lotResults[lotNum]
			title := sanitize(d.Title)
			desc := sanitize(d.Description)

			// These lines add make, model, condition, year, and country
			// to the description that Open AI generated, if they're populated.
			// But it was looking cluttered, so commenting for now.

			// make := sanitize(d.Make)
			// model := sanitize(d.Model)
			// cond := sanitize(d.Condition)
			// year := sanitize(d.Year)
			// country := sanitize(d.CountryOfOrigin)

			// extra := []string{}
			// if make != "" {
			// 	extra = append(extra, "Make: "+make)
			// }
			// if model != "" {
			// 	extra = append(extra, "Model: "+model)
			// }
			// if cond != "" {
			// 	extra = append(extra, "Condition: "+cond)
			// }
			// if year != "" {
			// 	extra = append(extra, "Year: "+year)
			// }
			// if country != "" {
			// 	extra = append(extra, "Country: "+country)
			// }
			// if len(extra) > 0 {
			// 	desc += " (" + strings.Join(extra, ", ") + ")"
			// }

			line := fmt.Sprintf("%s\t%s\t\t%s", lotNum, title, desc)

			writer.WriteString(line + "\n")
		}
		writer.Flush()
		fmt.Println("lots.tsv written successfully.")
	},
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "—", "-") // em dash
	s = strings.ReplaceAll(s, "“", "\"")
	s = strings.ReplaceAll(s, "”", "\"")
	s = strings.ReplaceAll(s, "’", "'")
	s = strings.ReplaceAll(s, "‘", "'")
	s = strings.Map(func(r rune) rune {
		if r > unicode.MaxASCII || r == '\t' || r == '\n' || r == '\r' {
			return -1
		}
		return r
	}, s)
	return strings.TrimSpace(s)
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
