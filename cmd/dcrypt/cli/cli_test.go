package cli

import (
	"testing"
)

func TestGetRootCmd(t *testing.T) {
	cmd := GetRootCmd()
	if cmd == nil {
		t.Error("GetRootCmd returned nil")
	}
	if cmd.Use != "dcrypt" {
		t.Errorf("Expected Use to be 'dcrypt', got '%s'", cmd.Use)
	}
}

func TestRootCommand(t *testing.T) {
	cmd := GetRootCmd()
	// dcrypt root command has no Run function, it delegates to subcommands
	if cmd.Run != nil {
		t.Log("Root command has Run function (optional)")
	}
}

func TestPasswdCommand(t *testing.T) {
	passwdCmd := GetPasswdCmd()
	if passwdCmd == nil {
		t.Error("GetPasswdCmd returned nil")
	}
	if passwdCmd.Short != "Generate a password hash" {
		t.Errorf("Expected Short to be 'Generate a password hash', got '%s'", passwdCmd.Short)
	}
}

func TestPasswdFlags(t *testing.T) {
	passwdCmd := GetPasswdCmd()
	if passwdCmd == nil {
		t.Fatal("GetPasswdCmd returned nil")
	}

	flags := []string{"password", "cost"}
	for _, flag := range flags {
		if passwdCmd.Flag(flag) == nil {
			t.Errorf("Flag --%s not found", flag)
		}
	}
}

func TestReadPassword(t *testing.T) {
	readPw := GetReadPassword()
	if readPw == nil {
		t.Error("GetReadPassword returned nil")
	}
}
