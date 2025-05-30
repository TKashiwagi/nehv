package test

import (
	"os"
	"testing"

	"configure/cmd"
	"configure/internal/config"
)

// TestHandleSetDNS tests the handleSetDNS function.
func TestHandleSetDNS(t *testing.T) {
	env := SetupTestEnv(t)
	cm, err := cmd.NewCommandManager(env.BootConfig, env.RunningConfig)
	if err != nil {
		t.Fatalf("Failed to create command manager: %v", err)
	}

	// Test case: Valid DNS address
	dnsAddr := "8.8.8.8"
	if err := cm.HandleSetDNS(dnsAddr); err != nil {
		t.Errorf("HandleSetDNS failed: %v", err)
	}

	cfg := cm.GetConfig()
	if len(cfg.DNS) == 0 {
		t.Error("Expected DNS to be set, got empty slice")
	} else if cfg.DNS[0] != dnsAddr {
		t.Errorf("Expected DNS to be set to %s, got %s", dnsAddr, cfg.DNS[0])
	}
}

// TestHandleAddDNS tests the handleAddDNS function.
func TestHandleAddDNS(t *testing.T) {
	env := SetupTestEnv(t)
	cm, err := cmd.NewCommandManager(env.BootConfig, env.RunningConfig)
	if err != nil {
		t.Fatalf("Failed to create command manager: %v", err)
	}

	// Test case: Valid DNS address
	dnsAddr := "1.1.1.1"
	if err := cm.HandleAddDNS(dnsAddr); err != nil {
		t.Errorf("HandleAddDNS failed: %v", err)
	}

	cfg := cm.GetConfig()
	found := false
	for _, dns := range cfg.DNS {
		if dns == dnsAddr {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected DNS %s to be added, but it was not found", dnsAddr)
	}
}

// TestHandleSetInterface tests the handleSetInterface function.
func TestHandleSetInterface(t *testing.T) {
	env := SetupTestEnv(t)
	cm, err := cmd.NewCommandManager(env.BootConfig, env.RunningConfig)
	if err != nil {
		t.Fatalf("Failed to create command manager: %v", err)
	}

	// Test case: Valid interface parameters
	fields := []string{"eth0", "address", "192.168.1.1"}
	if err := cm.HandleSetInterface(fields); err != nil {
		t.Errorf("HandleSetInterface failed: %v", err)
	}

	cfg := cm.GetConfig()
	iface, exists := cfg.Interfaces["eth0"]
	if !exists {
		t.Error("Expected interface eth0 to exist")
	} else if iface.Address != "192.168.1.1" {
		t.Errorf("Expected interface address to be set to 192.168.1.1, got %s", iface.Address)
	}
}

// TestHandleSetDefaultRoute tests the handleSetDefaultRoute function.
func TestHandleSetDefaultRoute(t *testing.T) {
	env := SetupTestEnv(t)
	cm, err := cmd.NewCommandManager(env.BootConfig, env.RunningConfig)
	if err != nil {
		t.Fatalf("Failed to create command manager: %v", err)
	}

	// Test case: Valid default route
	fields := []string{"192.168.1.1"}
	if err := cm.HandleSetDefaultRoute(fields); err != nil {
		t.Errorf("HandleSetDefaultRoute failed: %v", err)
	}

	cfg := cm.GetConfig()
	if cfg.DefaultRoute != "192.168.1.1" {
		t.Errorf("Expected default route to be set to 192.168.1.1, got %s", cfg.DefaultRoute)
	}
}

// TestHandleSave tests the handleSave function.
func TestHandleSave(t *testing.T) {
	env := SetupTestEnv(t)
	cm, err := cmd.NewCommandManager(env.BootConfig, env.RunningConfig)
	if err != nil {
		t.Fatalf("Failed to create command manager: %v", err)
	}

	// Test case: Save configuration
	if err := cm.HandleSave(); err != nil {
		t.Errorf("HandleSave failed: %v", err)
	}

	// Verify that the files are identical
	bootData, err := os.ReadFile(env.BootConfig)
	if err != nil {
		t.Errorf("Failed to read boot config: %v", err)
	}
	runningData, err := os.ReadFile(env.RunningConfig)
	if err != nil {
		t.Errorf("Failed to read running config: %v", err)
	}
	if string(bootData) != string(runningData) {
		t.Error("Expected boot config and running config to be identical")
	}
}

// TestHandleCommit tests the handleCommit function.
func TestHandleCommit(t *testing.T) {
	env := SetupTestEnv(t)
	cm, err := cmd.NewCommandManager(env.BootConfig, env.RunningConfig)
	if err != nil {
		t.Fatalf("Failed to create command manager: %v", err)
	}

	// Set up test configuration
	cm.SetDNS([]string{"8.8.8.8"})
	cm.SetDefaultRoute("192.168.1.1")
	cm.SetInterface("eth0", config.InterfaceConfig{
		Address: "192.168.1.100",
		MAC:     "00:11:22:33:44:55",
	})

	// Test case: Commit configuration
	if err := cm.HandleCommit(); err != nil {
		t.Errorf("HandleCommit failed: %v", err)
	}
}
