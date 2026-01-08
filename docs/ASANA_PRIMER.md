# Asana Primer

A quick introduction to Asana's concepts and data model for developers and AI agents.

## Core Hierarchy

```
Organization / Workspace
    └── Teams
         └── Projects
              └── Sections
                   └── Tasks
                        └── Subtasks
```

### Workspace

The top-level container. An organization has one workspace that contains all teams, projects, and users. Most API operations require a workspace GID.

### Teams

Groups of users. Teams own projects and control membership/permissions. Examples: "Engineering", "Design", "Marketing".

### Projects

The main organizing unit for work. A project is a collection of tasks, often representing:
- A sprint or iteration
- An initiative or feature
- A process (bug tracking, requests)
- Personal task list

Projects can be displayed as:
- **List view** - tasks as rows, sections as headers
- **Board view** - sections as columns (kanban-style)

### Sections

Subdivisions within a project. In list view, they're collapsible headers. In board view, they're columns.

Common patterns:
- Backlog → In Progress → Review → Done
- To Do → Doing → Done
- P0 → P1 → P2 → Icebox

### Tasks

The basic unit of work. A task has:
- **Name** - what needs to be done
- **Assignee** - who's responsible (one person)
- **Due date** - when it's due
- **Description** - details (notes field)
- **Completion status** - done or not
- **Projects** - which projects contain this task
- **Tags** - labels for categorization
- **Custom fields** - project-specific metadata
- **Subtasks** - breakdown into smaller pieces
- **Followers** - people watching for updates
- **Comments** - discussion and updates (called "stories")

### Subtasks

Tasks can have subtasks, which are themselves full tasks. Subtasks can have their own subtasks (unlimited nesting, but avoid going deep).

## Cross-Cutting Concepts

### Tags

Labels that can be applied to tasks across projects. Unlike sections (which are project-specific), tags work organization-wide.

Use cases:
- Bug types: `bug`, `regression`, `p0`
- Work types: `tech-debt`, `feature`, `spike`
- Status: `blocked`, `needs-review`

### Multi-Homing

A task can belong to multiple projects simultaneously. Updates sync across all projects containing the task.

Use cases:
- Personal task list + team sprint board
- Cross-team initiatives
- Rollup views

### Stories (Comments)

Activity on a task. Two types:
- **System stories** - auto-generated (assigned, completed, due date changed)
- **User stories** - comments added by people

Stories are append-only - you can edit/delete your own comments but not others'.

### Custom Fields

Project-specific metadata fields. Types:
- Text
- Number
- Enum (dropdown)
- Multi-enum
- Date
- People

Examples: Priority, Story Points, Sprint, Status.

Custom fields are defined at project level and can be shared across projects.

### Followers

People who receive notifications about a task. The assignee is automatically a follower. Add others for visibility without assignment.

## Key IDs

Asana uses GIDs (globally unique IDs) - numeric strings like `"1234567890123456"`.

Find GIDs:
- API responses include them
- Web UI URLs contain them: `app.asana.com/0/{project-gid}/{task-gid}`

## API Patterns

### Everything Returns `{ "data": ... }`

Single resource:
```json
{ "data": { "gid": "123", "name": "Task" } }
```

List:
```json
{ "data": [{ "gid": "123" }, { "gid": "456" }], "next_page": { ... } }
```

### Pagination

Lists are paginated (max 100 per page). Response includes `next_page` with offset token for subsequent requests.

### Field Selection

By default, API returns minimal fields. Use `opt_fields` to request specific fields:
```
GET /tasks/123?opt_fields=name,assignee.name,due_on
```

### Rate Limits

- ~1500 requests/minute for paid workspaces
- 429 response includes `Retry-After` header
- Complex queries cost more (deep subtask nesting, large projects)

## Common Workflows

### Personal Task Management

1. Create a personal project
2. Add sections: To Do, In Progress, Done
3. Create tasks as needed
4. Move through sections as you work
5. Complete when done

### Team Sprint Board

1. Team owns a project (board view)
2. Sections: Backlog, Sprint, In Progress, Review, Done
3. Tasks assigned to team members
4. Move cards through columns
5. Use custom fields for points/priority

### Request Intake

1. Create project with form
2. Form submissions create tasks
3. Triage in Backlog section
4. Assign and prioritize
5. Track through completion

### Cross-Team Initiative

1. Create initiative project
2. Multi-home tasks from team projects
3. Single view of all related work
4. Updates sync to team boards

## Tips for Automation

1. **Cache workspace/project GIDs** - they don't change
2. **Use context** - set project context to avoid repetition
3. **Batch reads** - `opt_fields` reduces API calls
4. **Respect rate limits** - exponential backoff on 429
5. **Tasks are the core** - most operations are task-centric
6. **Comments for updates** - post progress as task comments
7. **Sections for state** - move tasks between sections to show progress

## Glossary

| Term | Meaning |
|------|---------|
| GID | Globally unique identifier (numeric string) |
| Workspace | Top-level container (organization) |
| Project | Collection of tasks |
| Section | Subdivision within a project |
| Task | Unit of work |
| Subtask | Child task |
| Story | Comment or activity on a task |
| Tag | Cross-project label |
| Custom Field | Project-specific metadata |
| Multi-homing | Task in multiple projects |
| Assignee | Person responsible for task |
| Follower | Person watching task for updates |
| Due date | When task is due (`due_on` for date, `due_at` for datetime) |
