// main.go
package main

import (
	"fmt"
	"os"

	"bidzauction/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
