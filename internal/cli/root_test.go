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
