package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
)

var taskCFCmd = &cobra.Command{
	Use:   "custom-field",
	Short: "Manage custom fields on tasks",
}

var taskCFListCmd = &cobra.Command{
	Use:   "list",
	Short: "List custom field definitions for a project",
	Long:  "List custom field settings configured on a project. Uses context project if --project is not set.",
	RunE:  runTaskCFList,
}

var taskCFSetCmd = &cobra.Command{
	Use:   "set <task> <field-gid> <value>",
	Short: "Set a custom field value on a task",
	Long: `Set a custom field value on a task.

The value format depends on the field type:
  text, date:     string value
  number:         numeric value (e.g. 42 or 3.14)
  enum:           GID of the enum option to select
  multi_enum:     comma-separated GIDs of enum options to select
  people:         comma-separated user GIDs to assign`,
	Example: `  # Set a text field
  asana task custom-field set 1234567890 9876543210 "In review"

  # Set an enum field (pass the enum option GID)
  asana task custom-field set 1234567890 9876543210 1122334455

  # Set a people field with two users
  asana task custom-field set 1234567890 9876543210 111,222

  # Clear a field (pass empty string)
  asana task custom-field set 1234567890 9876543210 ""`,
	Args: cobra.ExactArgs(3),
	RunE: runTaskCFSet,
}

var (
	taskCFListProject string
	taskCFSetPick     bool
)

func init() {
	taskCmd.AddCommand(taskCFCmd)
	taskCFCmd.AddCommand(taskCFListCmd)
	taskCFCmd.AddCommand(taskCFSetCmd)

	taskCFListCmd.Flags().StringVar(&taskCFListProject, "project", "", "Project GID (uses context project if not set)")
	taskCFSetCmd.Flags().BoolVar(&taskCFSetPick, "pick", false, "Show interactive picker if multiple task name matches")
}

func runTaskCFList(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	projectGID := taskCFListProject
	if projectGID == "" && cfg.Project != "" {
		projectGID = cfg.Project
	}
	if projectGID == "" {
		return errors.NewGeneralError("no project specified; use --project or set a context project", nil)
	}

	client := newClient(cfg)
	result, err := client.GetProjectCustomFieldSettings(context.Background(), projectGID)
	if err != nil {
		return err
	}

	out := newOutput()
	return out.Print(result)
}

func runTaskCFSet(_ *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	if err := requireAuth(cfg); err != nil {
		return err
	}

	client := newClient(cfg)
	ctx := context.Background()

	taskGID, err := resolveTaskGID(ctx, cfg, client, args[0], taskCFSetPick)
	if err != nil {
		return err
	}

	fieldGID := args[1]
	rawValue := args[2]

	field, err := client.GetCustomField(ctx, fieldGID)
	if err != nil {
		return fmt.Errorf("failed to get custom field definition: %w", err)
	}

	fieldValue, err := encodeTypedCustomFieldValue(field, rawValue)
	if err != nil {
		return err
	}

	req := models.TaskUpdateRequest{
		CustomFields: map[string]any{fieldGID: fieldValue},
	}

	if cfg.DryRun {
		out := newOutput()
		return out.Print(map[string]any{
			"dry_run":    true,
			"gid":        taskGID,
			"field_gid":  fieldGID,
			"field_name": field.Name,
			"field_type": field.Type,
			"value":      fieldValue,
		})
	}

	task, err := client.UpdateTask(ctx, taskGID, req)
	if err != nil {
		return err
	}

	out := newOutput()
	return out.Print(task)
}

// splitCustomFieldPairs parses "field-gid=value" CLI strings into a raw gid→value map.
func splitCustomFieldPairs(pairs []string) (map[string]string, error) {
	raw := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		idx := strings.IndexByte(pair, '=')
		if idx < 1 {
			return nil, fmt.Errorf("invalid --custom-field value %q: expected field-gid=value", pair)
		}
		raw[pair[:idx]] = pair[idx+1:]
	}
	return raw, nil
}

// resolveCustomFields fetches the definition for each custom field GID and encodes
// the value into the correct Go type for the Asana API.
func resolveCustomFields(ctx context.Context, client api.Client, pairs []string) (map[string]any, error) {
	if len(pairs) == 0 {
		return nil, nil
	}
	raw, err := splitCustomFieldPairs(pairs)
	if err != nil {
		return nil, err
	}
	result := make(map[string]any, len(raw))
	for gid, value := range raw {
		field, err := client.GetCustomField(ctx, gid)
		if err != nil {
			return nil, fmt.Errorf("failed to look up custom field %s: %w", gid, err)
		}
		encoded, err := encodeTypedCustomFieldValue(field, value)
		if err != nil {
			return nil, err
		}
		result[gid] = encoded
	}
	return result, nil
}

// encodeTypedCustomFieldValue converts a raw CLI string to the correct Go type for
// the Asana API based on the field definition.
//
// Asana type → wire format:
//
//	number     → float64
//	multi_enum → []string (enum option GIDs, comma-separated input)
//	people     → []string (user GIDs, comma-separated input)
//	text/enum/date → string
//	empty      → nil (clears the field)
func encodeTypedCustomFieldValue(field *models.CustomField, raw string) (any, error) {
	if raw == "" {
		return nil, nil
	}
	switch field.Type {
	case "number":
		n, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, fmt.Errorf("field %q is a number field; value must be numeric, got %q", field.Name, raw)
		}
		return n, nil
	case "multi_enum", "people":
		return strings.Split(raw, ","), nil
	default:
		return raw, nil
	}
}
