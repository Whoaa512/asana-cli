# Asana CLI MVP Specification

## Executive Summary

The `asana-cli` is a Go-based command-line interface designed primarily for AI agents that need to track and document work across sessions, projects, and repositories. It provides structured JSON output by default, supports local per-repo context via `.asana.json`, and includes session management features for logging work progress directly to Asana tasks.

---

## Design Principles

| Principle | Rationale |
|-----------|-----------|
| **AI-first, JSON default** | Primary users are AI agents; JSON enables reliable parsing |
| **Local context awareness** | `.asana.json` in repo root binds commands to project/task context |
| **Minimal flags for common ops** | Context inheritance reduces flag repetition |
| **Predictable exit codes** | Machine-parseable success/failure states |
| **Session-based logging** | Captures work across agent invocations into Asana stories |
| **Single binary, zero deps** | Go static binary; no runtime requirements |

---

## Command Reference

```
asana
├── task
│   ├── create    --name --workspace --project --section --assignee --due --notes
│   ├── get       <gid>
│   ├── list      --workspace --project --section --assignee --completed
│   ├── update    <gid> --name --assignee --due --notes --completed
│   ├── delete    <gid>
│   ├── complete  <gid>
│   └── comment   <gid> --text
│
├── subtask
│   ├── create    <parent-gid> --name --assignee --due --notes
│   └── list      <parent-gid>
│
├── project
│   ├── list      --workspace --team --archived
│   ├── get       <gid>
│   ├── create    --workspace --team --name --notes --layout
│   └── update    <gid> --name --notes --archived
│
├── section
│   ├── list      --project
│   ├── get       <gid>
│   ├── create    --project --name
│   └── add-task  <section-gid> <task-gid>
│
├── workspace
│   ├── list
│   ├── get       <gid>
│   └── set       <gid>
│
├── tag
│   ├── list      --workspace
│   ├── get       <gid>
│   ├── create    --workspace --name --color
│   └── tasks     <gid>
│
├── search        --workspace --text --assignee --project --completed --due-before --due-after --limit
│
├── session
│   ├── start     [--task <gid>] [--project <gid>]
│   ├── end       [--summary <text>]
│   ├── status
│   └── log       <text>
│
├── ctx
│   ├── task      [<gid>]
│   ├── project   [<gid>]
│   └── clear
│
├── config
│   ├── show
│   ├── set       <key> <value>
│   └── init
│
├── me
│
└── version
```

### Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `json` | Output format: `json`, `jsonl`, `table`, `quiet` |
| `--fields` | `-f` | (resource default) | Comma-separated fields to include |
| `--workspace` | `-w` | (from config) | Override workspace |
| `--config` | | `~/.config/asana-cli/config.yaml` | Config file path |

### Quick Aliases

| Alias | Expands To |
|-------|------------|
| `asana log <text>` | `asana session log <text>` |
| `asana done` | `asana task complete <ctx-task>` |
| `asana note <text>` | `asana task comment <ctx-task> --text <text>` |

---

## Configuration

### File Locations (Resolution Order)

1. **CLI flags** (highest priority)
2. **Environment variables** (`ASANA_*`)
3. **Local `.asana.json`** (repo-specific, in cwd or parent)
4. **Global `~/.config/asana-cli/config.yaml`** (user defaults)

### Global Config (`~/.config/asana-cli/config.yaml`)

```yaml
default_workspace: "1234567890"
default_output: json
presets:
  my-tasks:
    assignee: me
    completed: false
  overdue:
    due_before: today
    completed: false
```

### Local Context (`.asana.json`)

```json
{
  "workspace": "1234567890",
  "project": "9876543210",
  "task": "1111111111"
}
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `ASANA_ACCESS_TOKEN` | **Required.** Personal access token |
| `ASANA_WORKSPACE` | Default workspace GID |
| `ASANA_OUTPUT` | Default output format |

---

## Output Formats

### `json` (default)
```json
{
  "gid": "1234567890",
  "name": "Implement feature X",
  "completed": false,
  "due_on": "2024-01-15",
  "assignee": {"gid": "111", "name": "Agent"}
}
```

### `jsonl` (for streaming/piping)
```
{"gid":"1234567890","name":"Task 1","completed":false}
{"gid":"1234567891","name":"Task 2","completed":true}
```

### `table` (human-readable)
```
GID           NAME                 DUE        STATUS
1234567890    Implement feature X  2024-01-15 pending
1234567891    Write tests          2024-01-16 complete
```

### `quiet` (GIDs only, for scripting)
```
1234567890
1234567891
```

---

## Authentication

Set `ASANA_ACCESS_TOKEN` environment variable with a Personal Access Token.

Generate at: https://app.asana.com/0/developer-console

```bash
export ASANA_ACCESS_TOKEN="1/1234567890:abcdef..."
```

No OAuth flow in MVP. Single-user PAT only.

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments / usage error |
| 3 | Authentication failure |
| 4 | Resource not found |
| 5 | Rate limited (includes Retry-After in stderr) |

---

## Session Management

Sessions track work across agent invocations and post summaries to Asana.

### Session Data (`~/.config/asana-cli/session.json`)

```json
{
  "task_gid": "1234567890",
  "project_gid": "9876543210",
  "started_at": "2024-01-15T10:00:00Z",
  "repo": "/path/to/repo",
  "branch": "feature/foo",
  "logs": [
    {"ts": "2024-01-15T10:05:00Z", "text": "Started implementation"},
    {"ts": "2024-01-15T11:30:00Z", "text": "Added tests"}
  ]
}
```

### Workflow

```bash
# Start session (auto-captures git info)
asana session start --task 1234567890

# Log progress during work
asana log "Implemented core parsing logic"
asana log "Fixed edge case in validation"

# End session (posts summary as story/comment)
asana session end --summary "Completed feature X with tests"
```

Session summary posted as Asana story with git repo/branch, duration, all log entries, and optional summary.

---

## MVP Scope

### IN Scope

- [x] Task CRUD (create, get, list, update, delete, complete)
- [x] Task comments (stories)
- [x] Subtask create/list
- [x] Project list/get/create
- [x] Section list/get/create/add-task
- [x] Workspace list/get/set
- [x] Tag list/get/create
- [x] Basic search (text, assignee, project, completed, due filters)
- [x] Session management (start, end, status, log)
- [x] Local context (`.asana.json`)
- [x] Global config
- [x] Output formats: json, jsonl, table, quiet
- [x] PAT authentication
- [x] Standard exit codes

### OUT of Scope (Post-MVP)

- [ ] OAuth flow
- [ ] Attachments
- [ ] Custom fields
- [ ] Portfolios
- [ ] Goals
- [ ] Webhooks
- [ ] Team management
- [ ] User management (beyond `me`)
- [ ] Shell completions
- [ ] Saved search presets
- [ ] Batch operations
- [ ] Interactive mode

---

## Implementation Phases

### Phase 1: Foundation
1. Project structure + CLI framework (cobra)
2. Config loading (viper): global config, env vars, local `.asana.json`
3. API client with auth + rate limit handling
4. Output formatting system (json, jsonl, table, quiet)
5. `me` and `workspace list/get/set` commands

### Phase 2: Core Task Operations
1. `task create/get/list/update/delete`
2. `task complete` + `task comment`
3. `subtask create/list`
4. Context resolution (`--workspace`, `--project` from config/flags)

### Phase 3: Organization
1. `project list/get/create/update`
2. `section list/get/create/add-task`
3. `tag list/get/create/tasks`
4. Local context commands (`ctx task/project/clear`)

### Phase 4: Search + Session
1. `search` with filters
2. `session start/end/status/log`
3. Git metadata capture
4. Session summary posting
5. Quick aliases (log, done, note)

### Phase 5: Polish
1. Error handling refinement
2. Rate limit retry logic
3. `--fields` support for all resources
4. Documentation + tests

---

## Architecture

```
asana-cli/
├── cmd/
│   └── asana/
│       └── main.go              # entry point
├── internal/
│   ├── cli/
│   │   ├── root.go              # root command + global flags
│   │   ├── task.go              # task subcommands
│   │   ├── project.go           # project subcommands
│   │   ├── section.go           # section subcommands
│   │   ├── workspace.go         # workspace subcommands
│   │   ├── tag.go               # tag subcommands
│   │   ├── search.go            # search command
│   │   ├── session.go           # session subcommands
│   │   ├── ctx.go               # context commands
│   │   └── config.go            # config commands
│   ├── api/
│   │   ├── client.go            # HTTP client, auth, rate limits
│   │   ├── tasks.go             # task API methods
│   │   ├── projects.go          # project API methods
│   │   ├── sections.go          # section API methods
│   │   ├── workspaces.go        # workspace API methods
│   │   ├── tags.go              # tag API methods
│   │   ├── stories.go           # stories (comments) API
│   │   ├── search.go            # search API
│   │   └── users.go             # user API (me)
│   ├── config/
│   │   ├── config.go            # config loading + resolution
│   │   ├── context.go           # local .asana.json handling
│   │   └── session.go           # session state management
│   └── output/
│       ├── formatter.go         # output format interface
│       ├── json.go              # json/jsonl output
│       ├── table.go             # table output
│       └── quiet.go             # GID-only output
├── docs/
│   ├── ASANA_API_REFERENCE.md
│   └── MVP_SPEC.md
├── go.mod
└── README.md
```

### Key Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/spf13/viper` | Config management |
| `github.com/olekukonko/tablewriter` | Table output |
