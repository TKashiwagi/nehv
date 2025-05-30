package test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"configure/cmd"
)

// TestWriteResolvConf tests the WriteResolvConf function.
func TestWriteResolvConf(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skip: not running on Linux/WSL")
	}
	env := SetupTestEnv(t)
	cmd.LoadedConfig = env.LoadedConfig
	// DNSをセット
	cmd.LoadedConfig.DNS = []string{"8.8.8.8"}

	// Create a temporary resolv.conf file
	resolvConf := filepath.Join(env.TempDir, "resolv.conf")
	os.Setenv("RESOLV_CONF", resolvConf)
	defer os.Unsetenv("RESOLV_CONF")

	// Test case: Write resolv.conf
	err := cmd.WriteResolvConf()
	if err != nil {
		t.Errorf("Failed to write resolv.conf: %v", err)
	}

	// Verify the content
	content, err := os.ReadFile(resolvConf)
	if err != nil {
		t.Errorf("Failed to read resolv.conf: %v", err)
	}
	expected := "nameserver 8.8.8.8"
	if string(content) != expected {
		t.Errorf("Expected resolv.conf content to be %s, got %s", expected, string(content))
	}
}

// TestRestartServices tests the RestartServices function.
func TestRestartServices(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skip: not running on Linux/WSL")
	}
	env := SetupTestEnv(t)
	cmd.LoadedConfig = env.LoadedConfig

	// Test case: Restart services
	err := cmd.RestartServices()
	if err != nil {
		t.Errorf("Failed to restart services: %v", err)
	}
}

// TestSetDefaultRoute tests the SetDefaultRoute function.
func TestSetDefaultRoute(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Skip: not running on Linux/WSL")
	}
	env := SetupTestEnv(t)
	cmd.LoadedConfig = env.LoadedConfig

	// Test case: Set default route
	err := cmd.SetDefaultRoute()
	if err != nil {
		t.Errorf("Failed to set default route: %v", err)
	}
}
