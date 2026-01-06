package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/models"
	"github.com/whoaa512/asana-cli/internal/output"
)

var taskCommentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage task comments",
}

var taskCommentListCmd = &cobra.Command{
	Use:   "list <task_gid>",
	Short: "List comments on a task",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskCommentList,
}

var taskCommentAddCmd = &cobra.Command{
	Use:   "add <task_gid>",
	Short: "Add a comment to a task",
	Args:  cobra.ExactArgs(1),
	RunE:  runTaskCommentAdd,
}

var (
	commentListLimit  int
	commentListOffset string
	commentAddText    string
)

func init() {
	taskCmd.AddCommand(taskCommentCmd)
	taskCommentCmd.AddCommand(taskCommentListCmd)
	taskCommentCmd.AddCommand(taskCommentAddCmd)

	taskCommentListCmd.Flags().IntVar(&commentListLimit, "limit", 50, "Max results to return")
	taskCommentListCmd.Flags().StringVar(&commentListOffset, "offset", "", "Pagination offset")

	taskCommentAddCmd.Flags().StringVar(&commentAddText, "text", "", "Comment text (required)")
	_ = taskCommentAddCmd.MarkFlagRequired("text")
}

func runTaskCommentList(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	result, err := client.ListStories(context.Background(), args[0], commentListLimit, commentListOffset)
	if err != nil {
		return err
	}

	var comments []models.Story
	for _, story := range result.Data {
		if story.Type == "comment" {
			comments = append(comments, story)
		}
	}

	filteredResult := &models.ListResponse[models.Story]{
		Data:     comments,
		NextPage: result.NextPage,
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(filteredResult)
}

func runTaskCommentAdd(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	if cfg.DryRun {
		out := output.NewJSON(os.Stdout)
		return out.Print(map[string]any{"dry_run": true, "task_gid": args[0], "text": commentAddText})
	}

	client := newClient(cfg)
	story, err := client.AddComment(context.Background(), args[0], commentAddText)
	if err != nil {
		return err
	}

	out := output.NewJSON(os.Stdout)
	return out.Print(story)
}
