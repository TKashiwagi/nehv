package main

import (
	"fmt"
	"os"

	"configure/cmd"
)

func main() {
	// Execute configure command
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
