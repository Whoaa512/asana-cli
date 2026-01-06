# Asana CLI MVP Specification

## Executive Summary

The `asana-cli` is a Go-based command-line interface designed primarily for AI agents that need to track and document work across sessions, projects, and repositories. It provides JSON output only (MVP), supports local per-repo context via `.asana.json`, and includes session management features for logging work progress directly to Asana tasks.

---

## Design Principles

| Principle | Rationale |
|-----------|-----------|
| **AI-first, JSON only** | Primary users are AI agents; JSON enables reliable parsing. Table output deferred post-MVP. |
| **Local context awareness** | `.asana.json` in repo root binds commands to project/task context |
| **Minimal flags for common ops** | Context inheritance reduces flag repetition |
| **Predictable exit codes** | Machine-parseable success/failure states |
| **Session-based logging** | Captures work across agent invocations into Asana stories |
| **Single binary, zero deps** | Go static binary; no runtime requirements |
| **Fail fast, fail loud** | Missing config/auth fails immediately with actionable error |

---

## Command Reference

```
asana
├── task
│   ├── create    --name --project --section --assignee --due --notes [--parent]
│   ├── get       <gid>
│   ├── list      --project --section --assignee --tag --completed --limit
│   ├── update    <gid> --name --assignee --due --notes --completed
│   ├── delete    <gid> [--force]
│   ├── complete  <gid>
│   ├── assign    <gid> <assignee>
│   └── comment   <gid> --text
│
├── project
│   ├── list      --workspace --team --archived --limit
│   └── get       <gid>
│
├── section
│   ├── list      --project --limit
│   ├── get       <gid>
│   ├── create    --project --name
│   └── add-task  <section-gid> <task-gid>
│
├── workspace
│   ├── list
│   ├── get       <gid>
│   └── use       <gid> [--global]
│
├── tag
│   ├── list      --workspace --limit
│   └── get       <gid>
│
├── session
│   ├── start     [--task <gid>] [--project <gid>] [--force]
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
│
└── version
```

### Global Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--workspace` | `-w` | (from config) | Override workspace GID |
| `--debug` | | `false` | Print HTTP requests/responses to stderr |
| `--dry-run` | | `false` | Preview mutations without executing |
| `--timeout` | | `30s` | HTTP request timeout |
| `--config` | | `~/.config/asana-cli/config.json` | Config file path |

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
3. **Local `.asana.json`** (walk up from cwd, stop at git root or home dir)
4. **Global `~/.config/asana-cli/config.json`** (user defaults)

### Walk-Up Rules for `.asana.json`

Starting from current working directory, walk up parent directories looking for `.asana.json`. Stop at:
- First `.asana.json` found (use it)
- Directory containing `.git/` (git root boundary)
- User's home directory
- Filesystem root

### Global Config (`~/.config/asana-cli/config.json`)

```json
{
  "default_workspace": "1234567890",
  "timeout": "30s"
}
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
| `ASANA_DEBUG` | Enable debug output (`1` or `true`) |

---

## Output Format

MVP supports **JSON only**. Table/quiet formats deferred to post-MVP.

### JSON Output

```json
{
  "gid": "1234567890",
  "name": "Implement feature X",
  "completed": false,
  "due_on": "2024-01-15",
  "assignee": {"gid": "111", "name": "Agent"}
}
```

### List Output

```json
{
  "data": [...],
  "next_page": {"offset": "...", "uri": "..."}
}
```

### Error Output

```json
{
  "error": {
    "message": "Task not found",
    "code": "NOT_FOUND",
    "exit_code": 4
  }
}
```

---

## Pagination

All list commands support pagination:

| Flag | Default | Description |
|------|---------|-------------|
| `--limit` | `50` | Max results to return (1-100) |
| `--all` | `false` | Fetch all pages (use with caution) |

Behavior:
- Default: Return up to `--limit` results, include `next_page` in output if more exist
- With `--all`: Auto-paginate, stream all results (respects rate limits)
- API page size capped at 100 per request

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
| 6 | Network/timeout error |

---

## Session Management

Sessions track work across agent invocations and post summaries to Asana.

### Multi-Repo Design

Sessions are stored **per-repo**, not globally. This allows concurrent work in different repos.

Session file location: `.asana-cli/session.json` in the git root of the current repo.

If not in a git repo: `~/.config/asana-cli/session.json` (fallback).

### Session Data

```json
{
  "task_gid": "1234567890",
  "project_gid": "9876543210",
  "started_at": "2024-01-15T10:00:00Z",
  "repo": "github.com/org/repo",
  "start_branch": "feature/foo",
  "end_branch": null,
  "logs": [
    {"ts": "2024-01-15T10:05:00Z", "type": "progress", "text": "Started implementation"},
    {"ts": "2024-01-15T11:30:00Z", "type": "decision", "text": "Using JWT over sessions"}
  ]
}
```

### Session Lifecycle

```bash
# Start session (auto-captures git info)
asana session start --task 1234567890

# Log progress during work
asana log "Implemented core parsing logic"
asana log --type decision "Using retry with exponential backoff"
asana log --type blocker "Waiting on API access"

# End session (posts summary as story/comment)
asana session end --summary "Completed feature X with tests"
```

### Crash Recovery

| Scenario | Behavior |
|----------|----------|
| `session start` when session exists | Error with hint to use `--force` |
| `session start --force` | Discards existing session, starts new |
| Session older than 24h | Warn on next command, suggest `session end --discard` |
| `session end` with empty logs | Skip posting, just clear session |
| `session end --discard` | Clear session without posting to Asana |
| Task deleted mid-session | `session end` warns but doesn't fail, clears session |

### Posted Summary Format

```markdown
## Work Session

**Duration:** 2h 15m
**Branch:** feature/foo → feature/foo
**Repo:** github.com/org/repo

### Progress
- Implemented core parsing logic
- Fixed edge case in validation

### Decisions
- Using JWT over sessions for statelessness

### Blockers
- Waiting on API access (resolved)

---
*Posted via asana-cli*
```

---

## API Client Design

### Rate Limiting

- Respect `Retry-After` header on 429 responses
- Exponential backoff: 1s, 2s, 4s, 8s, max 30s
- Max 3 retries per request
- Log rate limit hits to stderr

### Timeouts

- Default: 30s per request
- Configurable via `--timeout` flag or config
- Context cancellation propagated

### Error Handling

- Parse Asana error responses, surface `message` field
- Map HTTP status to exit codes
- Network errors → exit code 6

---

## MVP Scope

### IN Scope

- [x] Task CRUD (create, get, list, update, delete, complete, assign)
- [x] Task comments (stories)
- [x] Subtasks via `task create --parent`
- [x] Project list/get (read-only)
- [x] Section list/get/create/add-task
- [x] Workspace list/get/use
- [x] Tag list/get (read-only)
- [x] Task list with `--tag` filter
- [x] Session management (start, end, status, log)
- [x] Local context (`.asana.json`)
- [x] Global config
- [x] JSON output only
- [x] PAT authentication
- [x] Standard exit codes
- [x] Pagination (`--limit`, `--all`)
- [x] `--debug` flag
- [x] `--dry-run` flag

### OUT of Scope (Post-MVP)

- [ ] Table/quiet output formats
- [ ] Tag create/update
- [ ] Project create/update
- [ ] Search command (complex, defer)
- [ ] Presets/saved queries
- [ ] OAuth flow
- [ ] Attachments
- [ ] Custom fields
- [ ] Portfolios/Goals
- [ ] Webhooks
- [ ] Shell completions
- [ ] Batch operations

---

## Implementation Phases

### Phase 0: Foundation
1. Repository setup (CI, linting with golangci-lint, test harness)
2. Project structure creation
3. Makefile with build/test/lint targets

### Phase 1: Core Infrastructure
1. CLI framework (cobra) + global flags
2. Config loading (simple JSON + env, **skip viper**)
3. API client interface + HTTP implementation
4. Error types with exit code mapping
5. JSON output formatter
6. `me` command (auth validation)

### Phase 2: Workspaces & Context
1. `workspace list/get/use`
2. `.asana.json` context loading (walk-up logic)
3. `ctx show/task/project/clear`
4. Context resolution in commands

### Phase 3: Task Operations
1. `task create/get/list/update/delete`
2. `task complete` + `task assign`
3. `task comment`
4. Subtask support (`--parent` flag)
5. Pagination for `task list`

### Phase 4: Organization
1. `project list/get`
2. `section list/get/create/add-task`
3. `tag list/get`
4. `--tag` filter on `task list`

### Phase 5: Session Management
1. Session file management (per-repo)
2. `session start/end/status/log`
3. Git metadata capture
4. Session summary posting to Asana
5. Crash recovery logic
6. Quick aliases (log, done, note)

### Phase 6: Polish
1. `--debug` implementation
2. `--dry-run` implementation
3. Rate limit retry logic
4. Request timeout handling
5. Documentation
6. Integration tests

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
│   │   ├── session.go           # session subcommands
│   │   ├── ctx.go               # context commands
│   │   ├── config.go            # config commands
│   │   └── me.go                # me command
│   ├── api/
│   │   ├── client.go            # Client interface
│   │   ├── http.go              # HTTP implementation
│   │   ├── tasks.go             # task API methods
│   │   ├── projects.go          # project API methods
│   │   ├── sections.go          # section API methods
│   │   ├── workspaces.go        # workspace API methods
│   │   ├── tags.go              # tag API methods
│   │   ├── stories.go           # stories (comments) API
│   │   └── users.go             # user API (me)
│   ├── models/
│   │   ├── task.go              # Task struct
│   │   ├── project.go           # Project struct
│   │   ├── workspace.go         # Workspace struct
│   │   ├── user.go              # User struct
│   │   └── common.go            # shared types (AsanaResource, etc.)
│   ├── errors/
│   │   └── errors.go            # error types with exit codes
│   ├── config/
│   │   ├── config.go            # config loading + resolution
│   │   └── context.go           # local .asana.json handling
│   ├── session/
│   │   ├── session.go           # session state management
│   │   └── git.go               # git metadata capture
│   └── output/
│       └── json.go              # JSON output formatter
├── testdata/                    # test fixtures
├── docs/
│   ├── ASANA_API_REFERENCE.md
│   └── MVP_SPEC.md
├── Makefile
├── go.mod
└── README.md
```

### Key Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `gopkg.in/yaml.v3` | Config parsing (if YAML needed) |
| Standard library | HTTP client, JSON, testing |

**Note:** Viper intentionally excluded. Simple JSON config + env vars is sufficient and more predictable.

### Go Patterns

| Pattern | Implementation |
|---------|----------------|
| **Interface-first API client** | `api.Client` interface enables mock testing |
| **Context propagation** | All API methods take `context.Context` |
| **Error wrapping** | `fmt.Errorf("context: %w", err)` with sentinel errors |
| **No globals** | Pass config/client via struct, not package vars |
| **Struct tags** | Consistent `json:"snake_case"` for Asana API |

---

## Known Issues & Footguns

| Issue | Mitigation |
|-------|------------|
| Session file corruption from concurrent writes | File locking on write; warn if lock fails |
| `.asana.json` in home dir overrides everything | Walk-up stops at home dir, won't read `.asana.json` there |
| `--all` on large workspace causes OOM | Warn if >1000 results; suggest `--limit` |
| Rate limit thundering herd in scripts | Shared backoff state within process; jitter on retry |
| Token in shell history | Token only via env var; never accept as flag |
| Stale session blocks new work | 24h expiry warning; `--force` to override |
