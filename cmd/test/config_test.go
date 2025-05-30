package test

import (
	"os"
	"testing"

	"configure/internal/config"
)

// TestLoadConfig tests the loadConfig function.
func TestLoadConfig(t *testing.T) {
	env := SetupTestEnv(t)

	// Test case: Config file exists
	cfg, err := config.LoadConfig(env.BootConfig)
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
	}
	if cfg == nil {
		t.Error("Expected config to be loaded, got nil")
	}

	// Test case: Config file does not exist
	os.Remove(env.BootConfig)
	cfg, err = config.LoadConfig(env.BootConfig)
	if err == nil {
		t.Error("Expected error when config file does not exist, got nil")
	}
	if cfg != nil {
		t.Error("Expected config to be nil, got non-nil")
	}
}

// TestSaveConfig tests the saveConfig function.
func TestSaveConfig(t *testing.T) {
	env := SetupTestEnv(t)

	// Test case: Save valid config
	cfg := &config.Config{
		Hostname:   "test-router",
		Interfaces: make(map[string]config.InterfaceConfig),
		DNS:        []string{"8.8.8.8"},
	}
	if err := config.SaveConfig(cfg, env.BootConfig); err != nil {
		t.Errorf("Failed to save config: %v", err)
	}

	// Verify the saved config
	loadedCfg, err := config.LoadConfig(env.BootConfig)
	if err != nil {
		t.Errorf("Failed to load saved config: %v", err)
	}
	if loadedCfg.Hostname != cfg.Hostname {
		t.Errorf("Expected hostname %s, got %s", cfg.Hostname, loadedCfg.Hostname)
	}
	if len(loadedCfg.DNS) != len(cfg.DNS) {
		t.Errorf("Expected %d DNS servers, got %d", len(cfg.DNS), len(loadedCfg.DNS))
	}
}
