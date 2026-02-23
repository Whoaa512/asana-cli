package cli

import "testing"

func TestTaskAttachmentCommandRegistered(t *testing.T) {
	foundTask := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "task" {
			foundTask = true
			break
		}
	}
	if !foundTask {
		t.Fatal("task command should be registered")
	}

	foundAttachment := false
	for _, cmd := range taskCmd.Commands() {
		if cmd.Use == "attachment" {
			foundAttachment = true
			break
		}
	}
	if !foundAttachment {
		t.Fatal("task attachment command should be registered")
	}
}

func TestTaskAttachmentSubcommandsRegistered(t *testing.T) {
	subcommands := map[string]bool{
		"list <task_gid>":         false,
		"get <attachment_gid>":    false,
		"upload <task_gid>":       false,
		"delete <attachment_gid>": false,
	}

	for _, cmd := range taskAttachmentCmd.Commands() {
		if _, ok := subcommands[cmd.Use]; ok {
			subcommands[cmd.Use] = true
		}
	}

	for use, found := range subcommands {
		if !found {
			t.Errorf("subcommand %q should be registered", use)
		}
	}
}

func TestTaskAttachmentUploadFileFlagRequired(t *testing.T) {
	flag := taskAttachmentUploadCmd.Flags().Lookup("file")
	if flag == nil {
		t.Fatal("upload command should have --file flag")
	}

	if taskAttachmentUploadCmd.Flag("file").DefValue != "" {
		t.Errorf("upload --file default value = %q, want empty", taskAttachmentUploadCmd.Flag("file").DefValue)
	}
}
