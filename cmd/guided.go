package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"bidzauction/config"
)

func RunGuidedMode() {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()

		fmt.Println("BidzAuction — Guided Mode")
		fmt.Println("--------------------------")
		fmt.Println("1) init   - Create a new auction folder structure")
		fmt.Println("2) resize - Resize photos (01 -> 02)")
		fmt.Println("3) flatten- Prepare upload & AI folders")
		fmt.Println("4) analyze- SFTP + AI titles/descriptions + TSV")
		fmt.Println("5) shift  - Shift lots to open space")
		fmt.Println("Q) quit")
		fmt.Print("\nSelect an option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		switch choice {
		case "1", "init":
			runWizardInit(reader)
		case "2", "resize":
			runWizardResize(reader)
		case "3", "flatten":
			runWizardFlatten(reader)
		case "4", "analyze":
			runWizardAnalyze(reader)
		case "5", "shift":
			runWizardShift(reader)
		case "q", "quit", "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown choice. Press Enter to continue...")
			reader.ReadString('\n')
		}
	}
}

func runWizardInit(reader *bufio.Reader) {
	// Auction name (default: timestamp)
	ts := time.Now().Format("20060102_1504")
	defName := "Auction_" + ts
	fmt.Printf("\nAuction name [%s]: ", defName)
	auctionName, _ := reader.ReadString('\n')
	auctionName = strings.TrimSpace(auctionName)
	if auctionName == "" {
		auctionName = defName
	}

	// Lot count (default from config)
	defCount := config.Current.Defaults.AuctionSize
	fmt.Printf("Lot count [%d]: ", defCount)
	countStr, _ := reader.ReadString('\n')
	countStr = strings.TrimSpace(countStr)
	lotCount := defCount
	if countStr != "" {
		if n, err := strconv.Atoi(countStr); err == nil && n > 0 {
			lotCount = n
		} else {
			fmt.Printf("Invalid number. Using default %d.\n", defCount)
		}
	}

	fmt.Printf("\nCreate auction '%s' with %d lots? [Y/n]: ", auctionName, lotCount)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm == "n" || confirm == "no" {
		fmt.Println("Cancelled. Press Enter to return to menu...")
		reader.ReadString('\n')
		return
	}

	// Call the same helper your init command uses
	if err := createAuctionScaffold(auctionName, lotCount); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	} else {
		fmt.Println("Auction created successfully.")
	}
	fmt.Print("Press Enter to return to menu...")
	reader.ReadString('\n')

}

func runWizardResize(reader *bufio.Reader) {
	fmt.Println("\nResize photos")
	fmt.Println("-------------")
	fmt.Print("Auction name: ")
	auctionName, _ := reader.ReadString('\n')
	auctionName = strings.TrimSpace(auctionName)

	if auctionName == "" {
		fmt.Println("Auction name is required. Press Enter...")
		reader.ReadString('\n')
		return
	}

	// Reuse existing command by setting args
	if err := resize(auctionName); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print("Press Enter to return to menu...")
	reader.ReadString('\n')
}

func runWizardFlatten(reader *bufio.Reader) {
	fmt.Println("\nFlatten")
	fmt.Println("-------")
	fmt.Print("Auction name: ")
	auctionName, _ := reader.ReadString('\n')
	auctionName = strings.TrimSpace(auctionName)

	if auctionName == "" {
		fmt.Println("Auction name is required. Press Enter...")
		reader.ReadString('\n')
		return
	}

	if err := flatten(auctionName); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print("Press Enter to return to menu...")
	reader.ReadString('\n')
}

func runWizardAnalyze(reader *bufio.Reader) {
	fmt.Println("\nAnalyze")
	fmt.Println("-------")
	fmt.Print("Auction name: ")
	auctionName, _ := reader.ReadString('\n')
	auctionName = strings.TrimSpace(auctionName)

	if auctionName == "" {
		fmt.Println("Auction name is required. Press Enter...")
		reader.ReadString('\n')
		return
	}

	if err := analyze(auctionName); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print("Press Enter to return to menu...")
	reader.ReadString('\n')
}

func runWizardShift(reader *bufio.Reader) {
	fmt.Println("\nShift")
	fmt.Println("-----")
	fmt.Print("Auction name: ")
	auctionName, _ := reader.ReadString('\n')
	auctionName = strings.TrimSpace(auctionName)
	if auctionName == "" {
		fmt.Println("Auction name is required. Press Enter...")
		reader.ReadString('\n')
		return
	}

	fmt.Print("Start lot [1]: ")
	startStr, _ := reader.ReadString('\n')
	startStr = strings.TrimSpace(startStr)
	if len(startStr) == 0 {
		startStr = "1"
	}

	fmt.Print("Shift amount [1]: ")
	shiftStr, _ := reader.ReadString('\n')
	shiftStr = strings.TrimSpace(shiftStr)
	if len(shiftStr) == 0 {
		shiftStr = "1"
	}

	if err := shift(auctionName, []string{startStr, shiftStr}); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Print("Press Enter to return to menu...")
	reader.ReadString('\n')
}

func clearScreen() {
	// Cheap cross-platform "clear": just print a bunch of newlines.
	fmt.Print("\n\n\n\n\n\n\n\n\n\n")
}
