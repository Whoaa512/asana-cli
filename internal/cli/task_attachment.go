package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/output"
)

var taskAttachmentCmd = &cobra.Command{
	Use:   "attachment",
	Short: "Manage task attachments",
}

var taskAttachmentListCmd = &cobra.Command{
	Use:   "list <task_gid>",
	Short: "List attachments on a task",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskAttachmentList,
}

var taskAttachmentGetCmd = &cobra.Command{
	Use:   "get <attachment_gid>",
	Short: "Get attachment details",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskAttachmentGet,
}

var taskAttachmentUploadCmd = &cobra.Command{
	Use:   "upload <task_gid>",
	Short: "Upload an attachment to a task",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskAttachmentUpload,
}

var taskAttachmentDeleteCmd = &cobra.Command{
	Use:   "delete <attachment_gid>",
	Short: "Delete an attachment",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskAttachmentDelete,
}

var (
	attachmentListLimit  int
	attachmentListOffset string
	attachmentUploadFile string
	attachmentUploadName string
)

func init() {
	taskCmd.AddCommand(taskAttachmentCmd)
	taskAttachmentCmd.AddCommand(taskAttachmentListCmd)
	taskAttachmentCmd.AddCommand(taskAttachmentGetCmd)
	taskAttachmentCmd.AddCommand(taskAttachmentUploadCmd)
	taskAttachmentCmd.AddCommand(taskAttachmentDeleteCmd)

	taskAttachmentListCmd.Flags().IntVar(&attachmentListLimit, "limit", 50, "Max results to return")
	taskAttachmentListCmd.Flags().StringVar(&attachmentListOffset, "offset", "", "Pagination offset")

	taskAttachmentUploadCmd.Flags().StringVar(&attachmentUploadFile, "file", "", "Path to file (required)")
	taskAttachmentUploadCmd.Flags().StringVar(&attachmentUploadName, "name", "", "Attachment name override")
	_ = taskAttachmentUploadCmd.MarkFlagRequired("file")
}

func runTaskAttachmentList(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	result, err := client.ListAttachments(context.Background(), args[0], attachmentListLimit, attachmentListOffset)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(result)
}

func runTaskAttachmentGet(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	attachment, err := client.GetAttachment(context.Background(), args[0])
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(attachment)
}

func runTaskAttachmentUpload(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{
			"dry_run":  true,
			"task_gid": args[0],
			"file":     attachmentUploadFile,
			"name":     attachmentUploadName,
		})
	}

	client := newClient(cfg)
	attachment, err := client.UploadAttachment(context.Background(), args[0], attachmentUploadFile, attachmentUploadName)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(attachment)
}

func runTaskAttachmentDelete(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "gid": args[0], "action": "delete"})
	}

	client := newClient(cfg)
	if err := client.DeleteAttachment(context.Background(), args[0]); err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(map[string]any{"deleted": true, "gid": args[0]})
}
