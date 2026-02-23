# Attachment Support Spec

## Overview
Add attachment support to asana-cli: list, get, upload, and delete attachments on tasks.

## Asana API Reference
- **List attachments on a task:** `GET /tasks/{task_gid}/attachments` → returns `ListResponse[Attachment]`
- **Get attachment:** `GET /attachments/{attachment_gid}` → returns `Attachment`
- **Upload attachment:** `POST /tasks/{task_gid}/attachments` — **multipart/form-data** (not JSON). Fields: `file` (the file), `name` (optional override). Returns `Attachment`.
- **Delete attachment:** `DELETE /attachments/{attachment_gid}` → returns `{data: {}}`

## Files to Create/Modify

### 1. `internal/models/attachment.go` (new)
```go
type Attachment struct {
    GID         string `json:"gid"`
    Name        string `json:"name"`
    ResourceType string `json:"resource_type"`
    CreatedAt   string `json:"created_at,omitempty"`
    DownloadURL string `json:"download_url,omitempty"`
    Host        string `json:"host,omitempty"`
    ViewURL     string `json:"view_url,omitempty"`
    Size        int64  `json:"size,omitempty"`
}
```

### 2. `internal/api/client.go` — add to Client interface:
```go
ListAttachments(ctx context.Context, taskGID string, limit int, offset string) (*models.ListResponse[models.Attachment], error)
GetAttachment(ctx context.Context, gid string) (*models.Attachment, error)
UploadAttachment(ctx context.Context, taskGID string, filePath string, name string) (*models.Attachment, error)
DeleteAttachment(ctx context.Context, gid string) error
```

### 3. `internal/api/attachments.go` (new)
Implement all 4 methods. **Important:** `UploadAttachment` must use `multipart/form-data` — cannot use the existing `do()` helper which sends JSON. Build a new `doMultipart` helper or construct the request manually:
- Open the file
- Create `multipart.Writer`
- Add the file part as `file`
- If name is non-empty, add a `name` field
- Set `Content-Type` to the writer's FormDataContentType()
- Send with auth header, parse response

For the other 3, use the existing `do()` / `doWithRetry()` pattern.

### 4. `internal/cli/task_attachment.go` (new)
Add `task attachment` subcommand group with:

```
asana task attachment list <task_gid> [--limit N] [--offset S]
asana task attachment get <attachment_gid>
asana task attachment upload <task_gid> --file <path> [--name <override>]
asana task attachment delete <attachment_gid>
```

Follow the same patterns as `task_comment.go`:
- Register under `taskCmd`
- Use `loadConfig()`, `requireAuth()`, `newClient()`
- JSON output via `output.NewJSON`
- Support `--dry-run` for upload and delete

### 5. Tests
- `internal/api/attachments_test.go` — unit tests with httptest server for all 4 endpoints (follow patterns in `http_test.go`)
- `internal/cli/task_attachment_test.go` — basic CLI integration test (optional but nice)

## Patterns to Follow
- Look at `internal/api/stories.go` + `internal/cli/task_comment.go` for the closest analog
- Use the existing error types from `internal/errors/`
- Respect `--debug` flag (print request/response to stderr)
- Respect `--dry-run` for mutating operations

## Edge Cases
- Upload should validate file exists before making API call
- Large files: no special handling needed (Asana handles up to 100MB)
- If `--name` not provided, use the file's basename
