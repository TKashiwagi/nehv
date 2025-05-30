package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Hostname     string                     `yaml:"hostname"`
	Interfaces   map[string]InterfaceConfig `yaml:"interfaces"`
	DNS          []string                   `yaml:"dns"`
	DefaultRoute string                     `yaml:"default_route"`
}

// InterfaceConfig represents network interface configuration
type InterfaceConfig struct {
	Address string `yaml:"address"`
	MAC     string `yaml:"mac,omitempty"`
}

// ConfigManager handles configuration operations
type ConfigManager struct {
	bootConfigPath    string
	runningConfigPath string
	Config            *Config
}

// NewConfigManager creates a new ConfigManager instance
func NewConfigManager(bootPath, runningPath string) *ConfigManager {
	return &ConfigManager{
		bootConfigPath:    bootPath,
		runningConfigPath: runningPath,
		Config: &Config{
			Hostname:   "vyos-router",
			Interfaces: make(map[string]InterfaceConfig),
			DNS:        make([]string, 0),
		},
	}
}

// Load loads the configuration from the boot config file
func (cm *ConfigManager) Load() error {
	data, err := os.ReadFile(cm.bootConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config if it doesn't exist
			return cm.Save()
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cm.Config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// Save saves the current configuration to both boot and running config files
func (cm *ConfigManager) Save() error {
	data, err := yaml.Marshal(cm.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Save to boot config
	if err := os.WriteFile(cm.bootConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write boot config: %w", err)
	}

	// Save to running config
	if err := os.WriteFile(cm.runningConfigPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write running config: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *Config {
	return cm.Config
}

// SetDNS sets the DNS servers
func (cm *ConfigManager) SetDNS(servers []string) {
	cm.Config.DNS = servers
}

// AddDNS adds a DNS server
func (cm *ConfigManager) AddDNS(server string) {
	for _, existing := range cm.Config.DNS {
		if existing == server {
			return
		}
	}
	cm.Config.DNS = append(cm.Config.DNS, server)
}

// SetInterface sets interface configuration
func (cm *ConfigManager) SetInterface(name string, iface InterfaceConfig) {
	cm.Config.Interfaces[name] = iface
}

// SetDefaultRoute sets the default route
func (cm *ConfigManager) SetDefaultRoute(route string) {
	cm.Config.DefaultRoute = route
}

// Backup creates a backup of the current configuration
func (cm *ConfigManager) Backup() error {
	backupDir := "backup"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	backupPath := filepath.Join(backupDir, fmt.Sprintf("config_%s.yaml", time.Now().Format("20060102_150405")))
	data, err := yaml.Marshal(cm.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config for backup: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// Restore restores configuration from a backup file
func (cm *ConfigManager) Restore(backupPath string) error {
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %w", err)
	}

	var backupConfig Config
	if err := yaml.Unmarshal(data, &backupConfig); err != nil {
		return fmt.Errorf("failed to parse backup file: %w", err)
	}

	cm.Config = &backupConfig
	return cm.Save()
}
