package test

import (
	"os"
	"path/filepath"
	"testing"

	"configure/internal/config"
)

// TestEnv represents the test environment
type TestEnv struct {
	TempDir       string
	BootConfig    string
	RunningConfig string
	ConfigManager *config.ConfigManager
}

// SetupTestEnv sets up the test environment
func SetupTestEnv(t *testing.T) *TestEnv {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "configure_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create test configuration files
	bootConfig := filepath.Join(tempDir, "boot.config.yaml")
	runningConfig := filepath.Join(tempDir, "running.config.yaml")

	// Create initial configuration
	cfg := &config.Config{
		Hostname:   "test-router",
		Interfaces: make(map[string]config.InterfaceConfig),
		DNS:        make([]string, 0),
	}

	// Save initial configuration
	cm := config.NewConfigManager(bootConfig, runningConfig)
	cm.Config = cfg
	if err := cm.Save(); err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// Clean up after test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return &TestEnv{
		TempDir:       tempDir,
		BootConfig:    bootConfig,
		RunningConfig: runningConfig,
		ConfigManager: cm,
	}
}

// LoadConfig loads the configuration
func (env *TestEnv) LoadConfig(t *testing.T) {
	if err := env.ConfigManager.Load(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
}

// SaveConfig saves the configuration
func (env *TestEnv) SaveConfig(t *testing.T) {
	if err := env.ConfigManager.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
}
