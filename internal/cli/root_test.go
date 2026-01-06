package cli

import "testing"

func TestRootCommandExists(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}
	if rootCmd.Use != "asana" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "asana")
	}
}

func TestGlobalFlagsRegistered(t *testing.T) {
	flags := []string{"workspace", "debug", "dry-run", "timeout", "config"}
	for _, name := range flags {
		if rootCmd.PersistentFlags().Lookup(name) == nil {
			t.Errorf("flag %q not registered", name)
		}
	}
}

func TestWorkspaceShortFlag(t *testing.T) {
	f := rootCmd.PersistentFlags().ShorthandLookup("w")
	if f == nil {
		t.Error("workspace flag should have -w shorthand")
	}
}
