package cli

import (
	"testing"
)

func TestGetRootCmd(t *testing.T) {
	cmd := GetRootCmd()
	if cmd == nil {
		t.Error("GetRootCmd returned nil")
	}
	if cmd.Use != "david" {
		t.Errorf("Expected Use to be 'david', got '%s'", cmd.Use)
	}
}

func TestRootCommand(t *testing.T) {
	cmd := GetRootCmd()
	if cmd.Run == nil {
		t.Error("Root command Run function is nil")
	}
}

func TestServerCommand(t *testing.T) {
	serverCmd := GetServerCmd()
	if serverCmd == nil {
		t.Error("GetServerCmd returned nil")
	}
	if serverCmd.Short != "Start the WebDAV server" {
		t.Errorf("Expected Short to be 'Start the WebDAV server', got '%s'", serverCmd.Short)
	}
}

func TestServerFlags(t *testing.T) {
	serverCmd := GetServerCmd()
	if serverCmd == nil {
		t.Fatal("GetServerCmd returned nil")
	}

	flags := []string{"config", "host", "port", "debug", "production"}
	for _, flag := range flags {
		if serverCmd.Flag(flag) == nil {
			t.Errorf("Flag --%s not found", flag)
		}
	}
}
