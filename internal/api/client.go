package api

import (
	"context"

	"github.com/whoaa512/asana-cli/internal/models"
)

const BaseURL = "https://app.asana.com/api/1.0"

type Client interface {
	GetMe(ctx context.Context) (*models.User, error)
	ListWorkspaces(ctx context.Context, limit int) (*models.ListResponse[models.Workspace], error)
	GetWorkspace(ctx context.Context, gid string) (*models.Workspace, error)

	ListTasks(ctx context.Context, opts TaskListOptions) (*models.ListResponse[models.Task], error)
	GetTask(ctx context.Context, gid string) (*models.Task, error)
	CreateTask(ctx context.Context, req models.TaskCreateRequest) (*models.Task, error)
	UpdateTask(ctx context.Context, gid string, req models.TaskUpdateRequest) (*models.Task, error)
	DeleteTask(ctx context.Context, gid string) error

	ListStories(ctx context.Context, taskGID string, limit int, offset string) (*models.ListResponse[models.Story], error)
	AddComment(ctx context.Context, taskGID string, text string) (*models.Story, error)

	ListSubtasks(ctx context.Context, taskGID string, limit int, offset string) (*models.ListResponse[models.Task], error)
	AddSubtask(ctx context.Context, parentGID string, name string) (*models.Task, error)

	ListDependencies(ctx context.Context, taskGID string) ([]models.Task, error)
	ListDependents(ctx context.Context, taskGID string) ([]models.Task, error)
	AddDependency(ctx context.Context, taskGID string, dependsOnGID string) error
	RemoveDependency(ctx context.Context, taskGID string, dependsOnGID string) error

	ListProjects(ctx context.Context, opts ProjectListOptions) (*models.ListResponse[models.Project], error)
	GetProject(ctx context.Context, gid string) (*models.Project, error)
	CreateProject(ctx context.Context, req models.ProjectCreateRequest) (*models.Project, error)

	ListSections(ctx context.Context, opts SectionListOptions) (*models.ListResponse[models.Section], error)
	GetSection(ctx context.Context, gid string) (*models.Section, error)
	CreateSection(ctx context.Context, projectGID string, req models.SectionCreateRequest) (*models.Section, error)
	AddTaskToSection(ctx context.Context, sectionGID string, taskGID string) error

	ListTags(ctx context.Context, opts TagListOptions) (*models.ListResponse[models.Tag], error)
	GetTag(ctx context.Context, gid string) (*models.Tag, error)

	ListTeams(ctx context.Context, opts TeamListOptions) (*models.ListResponse[models.Team], error)
	ListUserTeams(ctx context.Context, opts UserTeamListOptions) (*models.ListResponse[models.Team], error)
}
