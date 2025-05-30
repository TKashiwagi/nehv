package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy running.config.yaml to boot.config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		source := "running.config.yaml"
		dest := "boot.config.yaml"
		input, err := os.ReadFile(source)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", source, err)
			os.Exit(1)
		}
		if err := os.WriteFile(dest, input, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", dest, err)
			os.Exit(1)
		}
		fmt.Printf("Copied %s to %s\n", source, dest)
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
