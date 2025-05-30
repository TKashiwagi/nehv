package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"configure/internal/completer"
	"configure/internal/config"
	"configure/internal/validator"
	"configure/internal/version"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

// CommandManager handles command execution and configuration management
type CommandManager struct {
	configManager *config.ConfigManager
	rl            *readline.Instance
}

// NewCommandManager creates a new CommandManager instance with the specified configuration files
func NewCommandManager(bootPath, runningPath string) (*CommandManager, error) {
	cm := config.NewConfigManager(bootPath, runningPath)
	if err := cm.Load(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "(config)# ",
		HistoryFile:     ".nehv_configure_history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		AutoComplete:    &completer.CLICompleter{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize readline: %w", err)
	}

	return &CommandManager{
		configManager: cm,
		rl:            rl,
	}, nil
}

// Close releases resources used by the CommandManager
func (cm *CommandManager) Close() {
	if cm.rl != nil {
		cm.rl.Close()
	}
}

// Execute starts the interactive configuration mode
func (cm *CommandManager) Execute() error {
	defer cm.Close()

	for {
		line, err := cm.rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			}
			continue
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := splitFields(line)
		if len(fields) == 0 {
			continue
		}

		if err := cm.HandleCommand(fields); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	return nil
}

// HandleCommand processes the command based on the input fields
func (cm *CommandManager) HandleCommand(fields []string) error {
	switch {
	case len(fields) == 1 && fields[0] == "exit":
		os.Exit(0)
	case len(fields) == 1 && (fields[0] == "help" || fields[0] == "?"):
		cm.printHelp()
	case len(fields) == 1 && fields[0] == "save":
		return cm.HandleSave()
	case len(fields) == 1 && fields[0] == "commit":
		return cm.HandleCommit()
	case len(fields) == 3 && fields[0] == "set" && fields[1] == "dns":
		return cm.HandleSetDNS(fields[2])
	case len(fields) == 3 && fields[0] == "add" && fields[1] == "dns":
		return cm.HandleAddDNS(fields[2])
	case len(fields) == 2 && fields[0] == "show" && fields[1] == "dns":
		cm.handleShowDNS()
	case len(fields) == 2 && fields[0] == "show" && fields[1] == "config":
		cm.handleShowConfig()
	case len(fields) == 2 && fields[0] == "show" && fields[1] == "interfaces":
		cm.handleShowInterfaces()
	case len(fields) == 2 && fields[0] == "show" && fields[1] == "version":
		cm.handleShowVersion()
	case len(fields) >= 4 && fields[0] == "set" && fields[1] == "interfaces":
		return cm.HandleSetInterface(fields[2:])
	case len(fields) == 5 && fields[0] == "set" && fields[1] == "ip" && fields[2] == "route" && fields[3] == "default" && fields[4] == "via":
		return cm.HandleSetDefaultRoute(fields[5:])
	default:
		return fmt.Errorf("unknown command: %s", strings.Join(fields, " "))
	}
	return nil
}

// Configuration Management Methods

// HandleSave saves the current configuration to both boot and running config files
func (cm *CommandManager) HandleSave() error {
	if err := cm.configManager.Save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	fmt.Println("Configuration saved successfully")
	return nil
}

// HandleCommit applies the current configuration to the system
func (cm *CommandManager) HandleCommit() error {
	cfg := cm.configManager.GetConfig()

	// Write resolv.conf
	if err := cm.writeResolvConf(cfg.DNS); err != nil {
		return fmt.Errorf("failed to write resolv.conf: %w", err)
	}

	// Restart services
	if err := cm.restartServices(); err != nil {
		return fmt.Errorf("failed to restart services: %w", err)
	}

	// Set default route
	if cfg.DefaultRoute != "" {
		if err := cm.setDefaultRoute(cfg.DefaultRoute); err != nil {
			return fmt.Errorf("failed to set default route: %w", err)
		}
	}

	fmt.Println("Configuration applied successfully")
	return nil
}

// DNS Configuration Methods

// HandleSetDNS sets the DNS servers
func (cm *CommandManager) HandleSetDNS(dnsAddr string) error {
	if err := validator.ValidateDNSAddress(dnsAddr); err != nil {
		return fmt.Errorf("invalid DNS address: %w", err)
	}
	cm.configManager.SetDNS([]string{dnsAddr})
	fmt.Printf("Set DNS: %s\n", dnsAddr)
	return nil
}

// HandleAddDNS adds a DNS server
func (cm *CommandManager) HandleAddDNS(dnsAddr string) error {
	if err := validator.ValidateDNSAddress(dnsAddr); err != nil {
		return fmt.Errorf("invalid DNS address: %w", err)
	}
	cm.configManager.AddDNS(dnsAddr)
	fmt.Printf("Added DNS: %s\n", dnsAddr)
	return nil
}

// Interface Configuration Methods

// HandleSetInterface sets interface parameters
func (cm *CommandManager) HandleSetInterface(fields []string) error {
	if len(fields) < 2 {
		return fmt.Errorf("missing interface parameters")
	}

	ifaceName := fields[0]
	param := fields[1]
	value := ""
	if len(fields) > 2 {
		value = fields[2]
	}

	iface := cm.configManager.GetConfig().Interfaces[ifaceName]
	switch param {
	case "address":
		if err := validator.ValidateIPAddress(value); err != nil {
			return fmt.Errorf("invalid IP address: %w", err)
		}
		iface.Address = value
	case "mac":
		if err := validator.ValidateMACAddress(value); err != nil {
			return fmt.Errorf("invalid MAC address: %w", err)
		}
		iface.MAC = value
	default:
		return fmt.Errorf("unknown interface parameter: %s", param)
	}

	cm.configManager.SetInterface(ifaceName, iface)
	fmt.Printf("Set interface %s %s to %s\n", ifaceName, param, value)
	return nil
}

// HandleSetDefaultRoute sets the default route
func (cm *CommandManager) HandleSetDefaultRoute(fields []string) error {
	if len(fields) == 0 {
		return fmt.Errorf("missing IP address for default route")
	}

	ipAddr := fields[0]
	if err := validator.ValidateIPAddress(ipAddr); err != nil {
		return fmt.Errorf("invalid IP address: %w", err)
	}

	cm.configManager.SetDefaultRoute(ipAddr)
	fmt.Printf("Set default route via %s\n", ipAddr)
	return nil
}

// Display Methods

// handleShowDNS displays the current DNS settings
func (cm *CommandManager) handleShowDNS() {
	cfg := cm.configManager.GetConfig()
	fmt.Printf("Current DNS: %v\n", cfg.DNS)
}

// handleShowConfig displays the current configuration
func (cm *CommandManager) handleShowConfig() {
	cfg := cm.configManager.GetConfig()
	cm.prettyPrintConfig(cfg)
}

// handleShowInterfaces displays the current interface status
func (cm *CommandManager) handleShowInterfaces() {
	cfg := cm.configManager.GetConfig()
	fmt.Println("Interfaces:")
	for name, iface := range cfg.Interfaces {
		fmt.Printf("  %s:\n", name)
		fmt.Printf("    address: %s\n", iface.Address)
		if iface.MAC != "" {
			fmt.Printf("    mac: %s\n", iface.MAC)
		}
	}
}

// handleShowVersion displays the version information
func (cm *CommandManager) handleShowVersion() {
	fmt.Printf("Version: %s\n", version.Version)
	fmt.Printf("Build Date: %s\n", version.BuildDate)
	fmt.Printf("Author: %s\n", version.Author)
}

// System Operation Methods

// writeResolvConf writes DNS settings to /etc/resolv.conf
func (cm *CommandManager) writeResolvConf(dnsServers []string) error {
	content := "nameserver " + strings.Join(dnsServers, "\nnameserver ")
	return os.WriteFile("/etc/resolv.conf", []byte(content), 0644)
}

// restartServices restarts the necessary services
func (cm *CommandManager) restartServices() error {
	if err := exec.Command("sudo", "systemctl", "restart", "resolvconf.service").Run(); err != nil {
		return err
	}
	return exec.Command("sudo", "systemctl", "restart", "systemd-resolved.service").Run()
}

// setDefaultRoute sets the default route in the system
func (cm *CommandManager) setDefaultRoute(route string) error {
	return exec.Command("sudo", "ip", "route", "add", "default", "via", route).Run()
}

// Helper Methods

// printHelp prints the help message
func (cm *CommandManager) printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  set dns <address>           Set DNS address")
	fmt.Println("  add dns <address>           Add DNS address")
	fmt.Println("  set interfaces <iface> <param> <value>  Set interface parameters")
	fmt.Println("    Parameters:")
	fmt.Println("      address <ip/mask>       Set interface IP address")
	fmt.Println("      mac <address>           Set interface MAC address")
	fmt.Println("  set ip route default via <ip>  Set default route")
	fmt.Println("  show dns                     Show current DNS settings")
	fmt.Println("  show config                  Show current configuration")
	fmt.Println("  show interfaces              Show interface status")
	fmt.Println("  show version                 Show version information")
	fmt.Println("  save                         Save current configuration")
	fmt.Println("  commit                       Apply current configuration")
	fmt.Println("  exit                         Exit configuration mode")
	fmt.Println("  help, ?                      Show this help message")
}

// prettyPrintConfig prints the configuration in a readable format
func (cm *CommandManager) prettyPrintConfig(cfg *config.Config) {
	fmt.Printf("Hostname: %s\n", cfg.Hostname)
	fmt.Println("Interfaces:")
	for name, iface := range cfg.Interfaces {
		fmt.Printf("  %s:\n", name)
		fmt.Printf("    Address: %s\n", iface.Address)
		if iface.MAC != "" {
			fmt.Printf("    MAC: %s\n", iface.MAC)
		}
	}
	fmt.Println("DNS servers:")
	for _, dns := range cfg.DNS {
		fmt.Printf("  %s\n", dns)
	}
	if cfg.DefaultRoute != "" {
		fmt.Printf("Default route: %s\n", cfg.DefaultRoute)
	}
}

// splitFields splits a line into fields by spaces
func splitFields(s string) []string {
	var res []string
	field := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			if field != "" {
				res = append(res, field)
				field = ""
			}
		} else {
			field += string(s[i])
		}
	}
	if field != "" {
		res = append(res, field)
	}
	return res
}

// Test Helper Methods

// GetConfig returns the current configuration
func (cm *CommandManager) GetConfig() *config.Config {
	return cm.configManager.GetConfig()
}

// SetDNS sets the DNS servers
func (cm *CommandManager) SetDNS(servers []string) {
	cm.configManager.SetDNS(servers)
}

// SetDefaultRoute sets the default route
func (cm *CommandManager) SetDefaultRoute(route string) {
	cm.configManager.SetDefaultRoute(route)
}

// SetInterface sets interface configuration
func (cm *CommandManager) SetInterface(name string, iface config.InterfaceConfig) {
	cm.configManager.SetInterface(name, iface)
}

// Root Command

var rootCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure network settings",
	Long:  `A command line tool for configuring network settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		cm, err := NewCommandManager("boot.config.yaml", "running.config.yaml")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if err := cm.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
