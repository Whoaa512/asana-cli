package cli

import "testing"

func TestMeCommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "me" {
			found = true
			break
		}
	}
	if !found {
		t.Error("me command should be registered")
	}
}

func TestMeCommandShort(t *testing.T) {
	if meCmd.Short == "" {
		t.Error("me command should have short description")
	}
}
