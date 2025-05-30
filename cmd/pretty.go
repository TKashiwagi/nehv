package cmd

import (
	"configure/internal/config"
	"fmt"
)

// 制約: 単発のコマンドを発行した場合は、プロンプトが戻ったら評価する

// prettyPrintConfig prints the config in a human-friendly YAML-like way
func prettyPrintConfig(cfg *config.Config) {
	fmt.Println("---")
	fmt.Printf("hostname: %s\n", cfg.Hostname)
	fmt.Println("interfaces:")
	for name, iface := range cfg.Interfaces {
		fmt.Printf("  %s:\n", name)
		fmt.Printf("    address: %s\n", iface.Address)
	}
	if len(cfg.DNS) > 0 {
		fmt.Println("dns:")
		for _, dns := range cfg.DNS {
			fmt.Printf("  - %s\n", dns)
		}
	}
}
