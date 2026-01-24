package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/whoaa512/asana-cli/internal/api"
	"github.com/whoaa512/asana-cli/internal/config"
	"github.com/whoaa512/asana-cli/internal/errors"
	"github.com/whoaa512/asana-cli/internal/models"
)

func init() {
	rootCmd.AddCommand(onboardCmd)
}

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Interactive setup wizard for asana-cli",
	Long: `Guides you through setting up asana-cli with your Asana account.

Steps:
  1. Validate ASANA_ACCESS_TOKEN
  2. Pick default workspace
  3. Optionally pick default team
  4. Optionally pick default project
  5. Save config to ~/.config/asana-cli/config.json`,
	RunE: runOnboard,
}

type onboardConfig struct {
	DefaultWorkspace string `json:"default_workspace,omitempty"`
	DefaultTeam      string `json:"default_team,omitempty"`
	DefaultProject   string `json:"default_project,omitempty"`
}

func runOnboard(_ *cobra.Command, _ []string) error {
	ctx := context.Background()

	token, err := checkToken(ctx)
	if err != nil {
		return err
	}

	cfg := &config.Config{
		AccessToken: token,
		Timeout:     config.DefaultTimeout,
	}
	client := api.NewHTTPClient(cfg)

	workspace, err := pickWorkspace(ctx, client)
	if err != nil {
		return err
	}

	var team *models.Team
	if workspace.IsOrganization {
		team, err = pickTeamOptional(ctx, client, workspace.GID)
		if err != nil {
			return err
		}
	}

	project, err := pickProjectOptional(ctx, client, workspace.GID)
	if err != nil {
		return err
	}

	result := onboardConfig{
		DefaultWorkspace: workspace.GID,
	}
	if team != nil {
		result.DefaultTeam = team.GID
	}
	if project != nil {
		result.DefaultProject = project.GID
	}

	if err := saveConfig(result); err != nil {
		return err
	}

	printSummary(workspace, team, project)
	return nil
}

func checkToken(ctx context.Context) (string, error) {
	token := os.Getenv("ASANA_ACCESS_TOKEN")
	if token == "" {
		return promptForToken(ctx)
	}

	if err := validateToken(ctx, token); err != nil {
		return "", errors.NewAuthError(fmt.Sprintf("token validation failed: %v", err))
	}

	fmt.Println("Token validated successfully.")
	return token, nil
}

func promptForToken(ctx context.Context) (string, error) {
	fmt.Println(`ASANA_ACCESS_TOKEN not set.

To get a token:
  1. Go to https://app.asana.com/0/my-apps
  2. Create a Personal Access Token
  3. Set it in your shell:

     export ASANA_ACCESS_TOKEN="your-token-here"

     Add to ~/.bashrc or ~/.zshrc for persistence.

Press Enter after setting the token...`)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	token := os.Getenv("ASANA_ACCESS_TOKEN")
	if token == "" {
		return "", errors.NewAuthError("ASANA_ACCESS_TOKEN still not set")
	}

	if err := validateToken(ctx, token); err != nil {
		return "", errors.NewAuthError(fmt.Sprintf("token validation failed: %v", err))
	}

	fmt.Println("Token validated successfully.")
	return token, nil
}

func validateToken(ctx context.Context, token string) error {
	cfg := &config.Config{
		AccessToken: token,
		Timeout:     config.DefaultTimeout,
	}
	client := api.NewHTTPClient(cfg)
	_, err := client.GetMe(ctx)
	return err
}

func pickWorkspace(ctx context.Context, client api.Client) (*models.Workspace, error) {
	resp, err := client.ListWorkspaces(ctx, 100)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, errors.NewGeneralError("no workspaces found", nil)
	}

	if len(resp.Data) == 1 {
		fmt.Printf("Using workspace: %s\n", resp.Data[0].Name)
		return &resp.Data[0], nil
	}

	return pick(resp.Data, "Select a workspace")
}

func pickTeamOptional(ctx context.Context, client api.Client, workspaceGID string) (*models.Team, error) {
	fmt.Print("\nWould you like to set a default team? [y/N] ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Text() != "y" && scanner.Text() != "Y" {
		return nil, nil
	}

	resp, err := client.ListUserTeams(ctx, api.UserTeamListOptions{
		Organization: workspaceGID,
		Limit:        100,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		fmt.Println("No teams found in this workspace.")
		return nil, nil
	}

	return pick(resp.Data, "Select a team")
}

func pickProjectOptional(ctx context.Context, client api.Client, workspaceGID string) (*models.Project, error) {
	fmt.Print("\nWould you like to set a default project? [y/N] ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Text() != "y" && scanner.Text() != "Y" {
		return nil, nil
	}

	resp, err := client.ListProjects(ctx, api.ProjectListOptions{
		Workspace: workspaceGID,
		Limit:     100,
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		fmt.Println("No projects found in this workspace.")
		return nil, nil
	}

	// Filter out archived projects
	var activeProjects []models.Project
	for _, p := range resp.Data {
		if !p.Archived {
			activeProjects = append(activeProjects, p)
		}
	}

	if len(activeProjects) == 0 {
		fmt.Println("No active projects found in this workspace.")
		return nil, nil
	}

	return pick(activeProjects, "Select a project")
}

func saveConfig(cfg onboardConfig) error {
	configDir := config.ExpandPath("~/.config/asana-cli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return errors.NewGeneralError("failed to create config directory", err)
	}

	configPath := filepath.Join(configDir, "config.json")

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return errors.NewGeneralError("failed to marshal config", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return errors.NewGeneralError("failed to write config", err)
	}

	fmt.Printf("\nConfig saved to: %s\n", configPath)
	return nil
}

func printSummary(workspace *models.Workspace, team *models.Team, project *models.Project) {
	fmt.Println("\n--- Configuration Summary ---")
	fmt.Printf("Workspace: %s (%s)\n", workspace.Name, workspace.GID)
	if team != nil {
		fmt.Printf("Team:      %s (%s)\n", team.Name, team.GID)
	}
	if project != nil {
		fmt.Printf("Project:   %s (%s)\n", project.Name, project.GID)
	}
	fmt.Println("\nYou're all set!")
}
