# asana-cli

A CLI for Asana, designed for AI agents and automation. JSON output only, context-aware, session-based work logging.

## Installation

```bash
# Install via go install (requires Go 1.21+)
go install github.com/whoaa512/asana-cli/cmd/asana@latest
```

Or download a pre-built binary from the [releases page](https://github.com/Whoaa512/asana-cli/releases).

## Quick Start

```bash
# 1. Get a Personal Access Token from https://app.asana.com/0/developer-console

# 2. Set the token
export ASANA_ACCESS_TOKEN="1/1234567890:abcdef..."

# 3. Verify it works
asana me
```

## Configuration

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `ASANA_ACCESS_TOKEN` | Yes | Personal access token |
| `ASANA_WORKSPACE` | No | Default workspace GID |
| `ASANA_DEBUG` | No | Enable debug output (`1` or `true`) |

### Global Config (`~/.config/asana-cli/config.json`)

```json
{
  "default_workspace": "1234567890"
}
```

### Local Context (`.asana.json` in repo/project root)

```json
{
  "workspace": "1234567890",
  "project": "9876543210",
  "task": "1111111111"
}
```

Context resolution order: CLI flags > `.asana.json` > global config > env vars

### CLAUDE.md Instruction Set

Here's the instruction block for your global CLAUDE.md (modeled after [Beads](https://github.com/steveyegge/beads)):
Add .asana.json and .asana-cli/ to .gitignore.

### Asana (Personal Task Tracker)
> For cross-repo work tracking via personal Asana. Configure sections in `.asana.json`.

```bash
# Context setup
asana ctx project <gid>             # Set project for this repo
asana ctx task <gid>                # Set active task
asana ctx show                      # View context

# Core workflow
asana prime                         # AI context dump (start of session)
asana ready --assignee me           # Find unblocked work
asana task start <gid>              # Claim task (move to in_progress section)
asana log "progress note"           # Add session log
asana done                          # Complete context task

# Dependencies
asana task dep add <task> <blocked-by>
asana task dep list <task>
asana task dep rm <task> <blocked-by>
asana blocked                       # Show blocked tasks

# Search & explore
asana search "query"                # Text search
asana task get <gid>                # Task details
```

Section config (.asana.json):
```json
{
  "project": "<project-gid>",
  "sections": {
    "planning": "<section-gid>",
    "in_progress": "<section-gid>",
    "blocked": "<section-gid>",
    "done": "<section-gid>"
  }
}
```


## Usage

### Find Your Workspace and Project

```bash
# List workspaces
asana workspace list

# List projects in workspace
asana project list --workspace <workspace-gid>

# List sections in a project
asana section list --project <project-gid>
```

### Task Operations

```bash
# Create a task
asana task create --name "Fix bug in auth" --project <project-gid>

# Create a subtask
asana task create --name "Subtask" --parent <parent-task-gid>

# List tasks
asana task list --project <project-gid>
asana task list --project <gid> --completed false --limit 20

# Get task details
asana task get <task-gid>

# Update a task
asana task update <task-gid> --name "New name" --due-on 2024-12-31

# Complete a task
asana task complete <task-gid>

# Move task to a section
asana task move <task-gid> --section <section-gid>
asana task start <task-gid>   # Move to in_progress section
asana task block <task-gid>   # Move to blocked section

# Add a comment
asana task comment add <task-gid> --text "Status update here"

# Delete a task
asana task delete <task-gid>
```

### Context Management

Set context to avoid repeating GIDs:

```bash
# View current context
asana ctx show

# Set project context
asana ctx project <project-gid>

# Set task context (for comments, completion)
asana ctx task <task-gid>

# Clear context
asana ctx clear
```

With context set, commands inherit it:

```bash
asana ctx project 123456
asana task create --name "New task"  # auto-uses project 123456
asana task list                       # auto-uses project 123456
```

### Session-Based Work Logging

Track work across agent invocations:

```bash
# Start a session (links to a task)
asana session start <task-gid>

# Log progress as you work
asana log "Implemented parsing logic"
asana log --type decision "Using JWT for auth"
asana log --type blocker "Waiting on API access"

# Check session status
asana session status

# End session (posts formatted summary as comment)
asana session end --summary "Completed feature with tests"

# Discard session without posting
asana session end --discard
```

Sessions capture git branch info and format a summary comment on the task.

### Quick Aliases

```bash
asana log <text>    # → asana session log <text>
asana done          # → asana task complete <context-task>
asana note <text>   # → asana task comment <context-task> --text <text>
```

## Command Reference

```
asana
├── prime         --project --limit --include-completed    # AI context dump
├── ready         --project --assignee --limit             # Find unblocked tasks
├── blocked       --project --assignee --limit             # Show blocked tasks
├── search        <query> --project --assignee --completed --limit --offset
├── done                                                   # Complete context task
├── log           <text> [--type]                          # Session log alias
├── note          <text>                                   # Task comment alias
│
├── task
│   ├── create    --name --project --assignee --due-on --notes [--parent]
│   ├── get       <gid>
│   ├── list      --project --section --assignee --tag --completed --limit --offset
│   ├── update    <gid> --name --assignee --due-on --notes
│   ├── delete    <gid>
│   ├── complete  <gid>
│   ├── assign    <gid> <assignee>
│   ├── move      <gid> --section (required)
│   ├── start     <gid>                                    # Move to in_progress
│   ├── block     <gid>                                    # Move to blocked
│   ├── plan      <gid>                                    # Move to planning
│   ├── duplicate <gid> [--name] [--include-subtasks] [--include-attachments]
│   ├── set-parent <gid> <parent_gid> | --clear            # Reparent or make top-level
│   ├── subtask
│   │   ├── list  <task_gid> --limit --offset
│   │   └── add   <task_gid> --name
│   ├── comment
│   │   ├── list  <task_gid> --limit --offset
│   │   └── add   <task_gid> --text
│   ├── dep
│   │   ├── add   <task_gid> <depends_on_gid>
│   │   ├── list  <task_gid>
│   │   └── rm    <task_gid> <depends_on_gid>
│   ├── follower
│   │   ├── add   <task_gid> <follower_gid>
│   │   └── rm    <task_gid> <follower_gid>
│   ├── tag
│   │   ├── add   <task_gid> <tag_gid>
│   │   └── rm    <task_gid> <tag_gid>
│   └── project
│       ├── list  <task_gid>                               # List projects task belongs to
│       ├── add   <task_gid> <project_gid>
│       └── rm    <task_gid> <project_gid>
│
├── project
│   ├── list      --archived --limit --offset
│   ├── get       <gid>
│   └── create    --name --notes --color --team
│
├── section
│   ├── list      --project --limit --offset
│   ├── get       <gid>
│   ├── create    --project --name
│   ├── update    <gid> --name
│   ├── delete    <gid>
│   ├── insert    <section-gid> --project [--before <gid>] [--after <gid>]
│   └── add-task  <section-gid> <task-gid>
│
├── workspace
│   ├── list      --limit
│   ├── get       <gid>
│   └── use       <gid> [--global]
│
├── tag
│   ├── list      --limit --offset
│   ├── get       <gid>
│   └── create    --name --workspace [--color <color>]
│
├── team
│   ├── list      --limit --offset
│   ├── get       <gid>
│   └── me        --limit --offset
│
├── session
│   ├── start     [<task-gid>] [--force]
│   ├── end       [--summary <text>] [--discard]
│   ├── status
│   └── log       <text> [--type progress|decision|blocker]
│
├── ctx
│   ├── show
│   ├── task      [<gid> | --clear]
│   ├── project   [<gid> | --clear]
│   └── clear
│
├── config
│   ├── show
│   └── init
│
├── me
│   ├── teams     --limit --offset
│   ├── projects  --limit --offset
│   └── tasks     --limit --offset --completed
│
└── version
```

### Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--workspace` | `-w` | from config | Override workspace GID |
| `--config` | | `~/.config/asana-cli/config.json` | Config file path |
| `--debug` | | `false` | Print HTTP requests/responses |
| `--dry-run` | | `false` | Preview without executing |
| `--timeout` | | `30s` | HTTP timeout |

## Output Format

All output is JSON. Examples:

```bash
# Single resource
asana task get 123 | jq '.name'

# List with pagination
asana task list --project 123 | jq '.data[].name'

# Check for more pages
asana task list --project 123 | jq '.next_page'
```

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Authentication failure |
| 4 | Resource not found |
| 5 | Rate limited |
| 6 | Network error |

## Development

```bash
mise run build     # Build binary
mise run test      # Run tests
mise run lint      # Run linter
mise run check     # Lint + test
mise run install   # Install to GOPATH/bin
mise run clean     # Clean artifacts
```

## FAQ

**Q: How do I find GIDs?**

GIDs are Asana's unique identifiers. Find them via:
- `asana workspace list` → workspace GIDs
- `asana project list` → project GIDs
- `asana task list --project <gid>` → task GIDs
- Web UI: open any resource, GID is in the URL

**Q: How do I set up for a specific repo?**

Create `.asana.json` in your repo root:
```json
{
  "project": "<project-gid>"
}
```

Now commands in that repo auto-use that project.

**Q: Can a task be in multiple projects?**

Yes, Asana calls this "multi-homing". Use `task project add <task_gid> <project_gid>` to add a task to additional projects, and `task project list <task_gid>` to see which projects a task belongs to.

**Q: How do sessions work?**

Sessions are per-repo (stored in `.asana-cli/session.json`). They:
1. Link to a task
2. Capture git branch at start
3. Collect log entries as you work
4. Post a formatted summary comment when ended

Great for AI agents that work across multiple invocations.

**Q: What if I get rate limited?**

The CLI auto-retries with backoff. For 429 errors, check the `Retry-After` in stderr. Avoid `--all` on large projects.

**Q: What happens when I delete a section?**

When you delete a section, tasks in that section are moved to the default/first section of the project. Tasks are NOT deleted. The Asana API handles this automatically.

**Q: Why JSON only?**

Primary use case is AI agents and automation. JSON is reliably parseable. Pipe to `jq` for formatting.

## Docs

- [Asana Primer](docs/ASANA_PRIMER.md) - Asana concepts and taxonomy
- [API Reference](docs/ASANA_API_REFERENCE.md) - Full Asana REST API docs
- [MVP Spec](docs/MVP_SPEC.md) - Design decisions and architecture

## License

MIT
