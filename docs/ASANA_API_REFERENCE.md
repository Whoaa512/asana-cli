# Asana REST API Reference

A comprehensive reference for the Asana REST API, useful for understanding the full API surface when building CLI tools or integrations.

**Base URL:** `https://app.asana.com/api/1.0`

**API Version:** 1.0

---

## Table of Contents

- [Authentication](#authentication)
- [Common Patterns](#common-patterns)
- [Rate Limits](#rate-limits)
- [Core Resources](#core-resources)
  - [Tasks](#tasks)
  - [Projects](#projects)
  - [Workspaces](#workspaces)
  - [Users](#users)
  - [Teams](#teams)
  - [Sections](#sections)
  - [Tags](#tags)
- [Collaboration Resources](#collaboration-resources)
  - [Stories (Comments)](#stories-comments)
  - [Attachments](#attachments)
  - [Status Updates](#status-updates)
- [Organization Resources](#organization-resources)
  - [Portfolios](#portfolios)
  - [Goals](#goals)
  - [Goal Relationships](#goal-relationships)
  - [Custom Fields](#custom-fields)
- [Template Resources](#template-resources)
  - [Project Templates](#project-templates)
  - [Task Templates](#task-templates)
- [Membership Resources](#membership-resources)
  - [Memberships (Generic)](#memberships-generic)
  - [Workspace Memberships](#workspace-memberships)
  - [Team Memberships](#team-memberships)
  - [Portfolio Memberships](#portfolio-memberships)
  - [Project Memberships](#project-memberships-deprecated)
- [Time & Planning Resources](#time--planning-resources)
  - [Time Tracking Entries](#time-tracking-entries)
  - [Time Periods](#time-periods)
  - [Allocations](#allocations)
- [Integration Resources](#integration-resources)
  - [Webhooks](#webhooks)
  - [Events](#events)
  - [Rules](#rules)
- [Utility Resources](#utility-resources)
  - [Batch API](#batch-api)
  - [Jobs](#jobs)
  - [Typeahead](#typeahead)
  - [Organization Exports](#organization-exports)
  - [Audit Log API](#audit-log-api)
- [Additional Resources](#additional-resources)
  - [User Task Lists](#user-task-lists)
  - [Project Briefs](#project-briefs)
  - [Project Statuses](#project-statuses-deprecated)

---

## Authentication

### Personal Access Token (PAT)
Simplest authentication method. Generate in the Asana developer console. Ideal for scripts and single-user applications.

```
Authorization: Bearer <your_personal_access_token>
```

### Service Account
Enterprise-only tokens with org-wide access. Created by super admins via admin console.

### OAuth 2.0
For multi-user applications. Standard OAuth flow with user consent.

### OpenID Connect (OIDC)
Built on OAuth 2.0 for single sign-on scenarios.

---

## Common Patterns

### Request Format
- JSON or form-encoded request bodies
- All responses wrapped in `{ "data": ... }`

### Input/Output Options

| Parameter | Description |
|-----------|-------------|
| `opt_fields` | Comma-separated list of fields to return |
| `opt_pretty` | Format JSON response with indentation |

**GET request example:**
```
GET /tasks/12345?opt_fields=name,assignee,due_on&opt_pretty=true
```

**POST/PUT request example:**
```json
{
  "data": { ... },
  "options": {
    "fields": ["name", "assignee"],
    "pretty": true
  }
}
```

### Nested Field Selection
Use dot notation: `opt_fields=assignee.name,followers.email`

Grouping: `opt_fields=this.(followers|assignee).email`

### Pagination

| Parameter | Description |
|-----------|-------------|
| `limit` | Results per page (1-100) |
| `offset` | Token from previous response for next page |

Response includes `next_page` object with `offset`, `path`, and `uri` for subsequent requests. Returns `null` when no more pages exist.

**Recommendation:** Always paginate requests for potentially large datasets. Non-paginated queries may be capped at ~1,000 objects.

### Error Response Format
```json
{
  "errors": [
    {
      "message": "workspace: Missing input",
      "phrase": "6 sad squid snuggle softly"
    }
  ]
}
```

### HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created (includes `Location` header) |
| 400 | Bad Request (invalid parameters) |
| 401 | Unauthorized (invalid token) |
| 402 | Payment Required (premium feature) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 429 | Too Many Requests (rate limited) |
| 451 | Unavailable For Legal Reasons |
| 500 | Internal Server Error |

---

## Rate Limits

### Request Limits (per minute)
| Tier | Limit |
|------|-------|
| Free domains | 150 requests/min |
| Paid domains | 1,500 requests/min |
| Search API | 60 requests/min |
| Duplication/Export jobs | 5 concurrent per user |

### Concurrent Request Limits
| Type | Limit |
|------|-------|
| GET requests | 50 concurrent |
| POST/PUT/PATCH/DELETE | 15 concurrent |

### Retry-After Header
When rate limited (429), the response includes `Retry-After` header with seconds to wait.

### Cost-Based Limits
Complex queries (deep graph traversal) consume computational "cost." Excessive cost depletes quota and causes rejections.

**Avoid:**
- Deeply nested subtasks
- Projects with >1,000 tasks
- Numerous unreadable tags
- Undeleted webhooks

---

## Core Resources

### Tasks

Tasks are the basic unit of work in Asana.

#### Core Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks` | Get multiple tasks (requires workspace, project, assignee, tag, or section filter) |
| POST | `/tasks` | Create a task |
| GET | `/tasks/{task_gid}` | Get a task |
| PUT | `/tasks/{task_gid}` | Update a task |
| DELETE | `/tasks/{task_gid}` | Delete a task |
| POST | `/tasks/{task_gid}/duplicate` | Duplicate a task |

#### Task Retrieval by Context

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects/{project_gid}/tasks` | Get tasks from a project |
| GET | `/sections/{section_gid}/tasks` | Get tasks from a section |
| GET | `/tags/{tag_gid}/tasks` | Get tasks from a tag |
| GET | `/user_task_lists/{user_task_list_gid}/tasks` | Get tasks from user's task list |
| GET | `/tasks/custom_id/{custom_id}` | Get task by custom ID |

#### Task Search

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/workspaces/{workspace_gid}/tasks/search` | Search tasks in workspace |

**Search supports filters for:** assignee, projects, sections, tags, completion status, due dates, custom fields, text, and more.

#### Subtasks & Parent

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks/{task_gid}/subtasks` | Get subtasks |
| POST | `/tasks/{task_gid}/subtasks` | Create a subtask |
| POST | `/tasks/{task_gid}/setParent` | Set parent of task |

#### Dependencies

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks/{task_gid}/dependencies` | Get task dependencies |
| POST | `/tasks/{task_gid}/addDependencies` | Add dependencies |
| POST | `/tasks/{task_gid}/removeDependencies` | Remove dependencies |
| GET | `/tasks/{task_gid}/dependents` | Get task dependents |
| POST | `/tasks/{task_gid}/addDependents` | Add dependents |
| POST | `/tasks/{task_gid}/removeDependents` | Remove dependents |

#### Task Organization

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/tasks/{task_gid}/addProject` | Add task to project |
| POST | `/tasks/{task_gid}/removeProject` | Remove task from project |
| POST | `/tasks/{task_gid}/addTag` | Add tag to task |
| POST | `/tasks/{task_gid}/removeTag` | Remove tag from task |
| POST | `/tasks/{task_gid}/addFollowers` | Add followers |
| POST | `/tasks/{task_gid}/removeFollower` | Remove follower |

#### Key Task Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Task name |
| `completed` | boolean | Completion status |
| `assignee` | object/gid | Assigned user |
| `due_on` | date | Due date (YYYY-MM-DD) |
| `due_at` | datetime | Due datetime (ISO 8601) |
| `start_on` | date | Start date |
| `start_at` | datetime | Start datetime |
| `notes` | string | Task description (plain text) |
| `html_notes` | string | Task description (HTML) |
| `projects` | array | Associated projects |
| `workspace` | object | Parent workspace |
| `parent` | object | Parent task (for subtasks) |
| `custom_fields` | array | Custom field values |

---

### Projects

Projects organize tasks and can be displayed as lists or boards.

#### Core Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects` | Get multiple projects |
| POST | `/projects` | Create a project |
| GET | `/projects/{project_gid}` | Get a project |
| PUT | `/projects/{project_gid}` | Update a project |
| DELETE | `/projects/{project_gid}` | Delete a project |
| POST | `/projects/{project_gid}/duplicate` | Duplicate a project |

#### Team & Workspace Scoped

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/teams/{team_gid}/projects` | Get team's projects |
| POST | `/teams/{team_gid}/projects` | Create project in team |
| GET | `/workspaces/{workspace_gid}/projects` | Get workspace projects |
| POST | `/workspaces/{workspace_gid}/projects` | Create project in workspace |
| GET | `/tasks/{task_gid}/projects` | Get projects containing task |

#### Project Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects/{project_gid}/task_counts` | Get task counts |
| POST | `/projects/{project_gid}/members` | Add members |
| POST | `/projects/{project_gid}/members/{user_gid}/remove` | Remove member |
| POST | `/projects/{project_gid}/followers` | Add followers |
| POST | `/projects/{project_gid}/followers/{user_gid}/remove` | Remove follower |
| POST | `/projects/{project_gid}/save_as_template` | Save as template |

#### Custom Fields on Projects

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/projects/{project_gid}/custom_fields` | Add custom field |
| POST | `/projects/{project_gid}/custom_fields/{field_gid}/remove` | Remove custom field |
| GET | `/projects/{project_gid}/custom_field_settings` | Get custom field settings |

---

### Workspaces

Workspaces are the highest-level organizational unit.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/workspaces` | Get all workspaces |
| GET | `/workspaces/{workspace_gid}` | Get a workspace |
| PUT | `/workspaces/{workspace_gid}` | Update a workspace |
| POST | `/workspaces/{workspace_gid}/addUser` | Add user to workspace |
| POST | `/workspaces/{workspace_gid}/removeUser` | Remove user from workspace |

#### Key Workspace Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Workspace name |
| `is_organization` | boolean | Whether this is an organization |
| `email_domains` | array | Associated email domains |

---

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/users` | Get multiple users |
| GET | `/users/{user_gid}` | Get a user |
| PUT | `/users/{user_gid}` | Update a user |
| GET | `/users/{user_gid}/favorites` | Get user's favorites |
| GET | `/users/me` | Get current user |
| GET | `/workspaces/{workspace_gid}/users` | Get users in workspace |
| GET | `/teams/{team_gid}/users` | Get users in team |

**Note:** The special identifier `me` can be used in place of any user GID to reference the authenticated user.

---

### Teams

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/teams` | Create a team |
| GET | `/teams/{team_gid}` | Get a team |
| PUT | `/teams/{team_gid}` | Update a team |
| GET | `/workspaces/{workspace_gid}/teams` | Get teams in workspace |
| GET | `/users/{user_gid}/teams` | Get user's teams |
| POST | `/teams/{team_gid}/addUser` | Add user to team |
| POST | `/teams/{team_gid}/removeUser` | Remove user from team |

#### Team Visibility Options
- `public` - Visible to everyone in organization
- `request_to_join` - Users can request to join
- `secret` - Only visible to members

---

### Sections

Sections organize tasks within projects (headers in list view, columns in board view).

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/sections/{section_gid}` | Get a section |
| PUT | `/sections/{section_gid}` | Update a section |
| DELETE | `/sections/{section_gid}` | Delete a section |
| GET | `/projects/{project_gid}/sections` | Get sections in project |
| POST | `/projects/{project_gid}/sections` | Create section in project |
| POST | `/projects/{project_gid}/sections/insert` | Move/reorder sections |
| POST | `/sections/{section_gid}/addTask` | Add task to section |

---

### Tags

Tags provide cross-project categorization (no ordering on associated tasks).

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tags` | Get multiple tags |
| POST | `/tags` | Create a tag |
| GET | `/tags/{tag_gid}` | Get a tag |
| PUT | `/tags/{tag_gid}` | Update a tag |
| DELETE | `/tags/{tag_gid}` | Delete a tag |
| GET | `/tasks/{task_gid}/tags` | Get tags on task |
| GET | `/workspaces/{workspace_gid}/tags` | Get tags in workspace |
| POST | `/workspaces/{workspace_gid}/tags` | Create tag in workspace |

#### Tag Colors
18 color options available: `dark-pink`, `dark-green`, `dark-blue`, `dark-red`, `dark-teal`, `dark-brown`, `dark-orange`, `dark-purple`, `dark-warm-gray`, `light-pink`, `light-green`, `light-blue`, `light-red`, `light-teal`, `light-brown`, `light-orange`, `light-purple`, `light-warm-gray`

---

## Collaboration Resources

### Stories (Comments)

Stories represent activity on tasks. System-generated for actions; user-created for comments.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/stories/{story_gid}` | Get a story |
| PUT | `/stories/{story_gid}` | Update a story |
| DELETE | `/stories/{story_gid}` | Delete a story |
| GET | `/tasks/{task_gid}/stories` | Get stories on task |
| POST | `/tasks/{task_gid}/stories` | Create story (comment) on task |

#### Key Story Fields

| Field | Type | Description |
|-------|------|-------------|
| `text` | string | Comment text (create only) |
| `html_text` | string | HTML formatted text |
| `is_pinned` | boolean | Whether story is pinned |
| `sticker_name` | string | Emoji sticker name |

---

### Attachments

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/attachments/{attachment_gid}` | Get an attachment |
| DELETE | `/attachments/{attachment_gid}` | Delete an attachment |
| GET | `/tasks/{task_gid}/attachments` | Get attachments on task |
| POST | `/tasks/{task_gid}/attachments` | Upload attachment to task |

**Supported hosts:** asana, dropbox, gdrive, onedrive, box, vimeo, external

---

### Status Updates

Progress updates sent to followers.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/status_updates/{status_gid}` | Get a status update |
| DELETE | `/status_updates/{status_gid}` | Delete a status update |
| GET | `/{object_gid}/status_updates` | Get status updates on object |
| POST | `/{object_gid}/status_updates` | Create status update |

**Status Types:** `on_track`, `at_risk`, `off_track`, `on_hold`, `achieved`, `complete`, `dropped`, `missed`, `partial`

**Note:** Status updates cannot be modified after creation.

---

## Organization Resources

### Portfolios

Portfolios group and manage multiple projects.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/portfolios` | Get multiple portfolios |
| POST | `/portfolios` | Create a portfolio |
| GET | `/portfolios/{portfolio_gid}` | Get a portfolio |
| PUT | `/portfolios/{portfolio_gid}` | Update a portfolio |
| DELETE | `/portfolios/{portfolio_gid}` | Delete a portfolio |
| GET | `/portfolios/{portfolio_gid}/items` | Get portfolio items |
| POST | `/portfolios/{portfolio_gid}/items` | Add item to portfolio |
| POST | `/portfolios/{portfolio_gid}/removeItem` | Remove item from portfolio |
| POST | `/portfolios/{portfolio_gid}/addMembers` | Add members |
| POST | `/portfolios/{portfolio_gid}/removeMembers` | Remove members |
| POST | `/portfolios/{portfolio_gid}/addCustomFieldSetting` | Add custom field |
| POST | `/portfolios/{portfolio_gid}/removeCustomFieldSetting` | Remove custom field |

---

### Goals

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/goals` | Get multiple goals |
| POST | `/goals` | Create a goal |
| GET | `/goals/{goal_gid}` | Get a goal |
| PUT | `/goals/{goal_gid}` | Update a goal |
| DELETE | `/goals/{goal_gid}` | Delete a goal |
| POST | `/goals/{goal_gid}/metrics` | Create goal metric |
| POST | `/goals/{goal_gid}/metrics/{metric_gid}` | Update goal metric |
| POST | `/goals/{goal_gid}/followers` | Add followers |
| POST | `/goals/{goal_gid}/followers/{follower_gid}` | Remove follower |
| GET | `/goals/{goal_gid}/parentGoals` | Get parent goals |
| GET | `/goals/{goal_gid}/custom_field_settings` | Get custom fields |

---

### Goal Relationships

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/goal_relationships` | Get goal relationships |
| GET | `/goal_relationships/{goal_relationship_gid}` | Get a relationship |
| PUT | `/goal_relationships/{goal_relationship_gid}` | Update a relationship |
| POST | `/goals/{goal_gid}/addSupportingRelationship` | Add supporting relationship |
| POST | `/goals/{goal_gid}/removeSupportingRelationship` | Remove supporting relationship |

**Relationship Subtypes:**
- `subgoal` - Goal supporting another goal
- `supporting_work` - Project, task, or portfolio supporting a goal

---

### Custom Fields

Custom fields are a premium feature.

#### Core Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/custom_fields` | Create a custom field |
| GET | `/custom_fields/{custom_field_gid}` | Get a custom field |
| PUT | `/custom_fields/{custom_field_gid}` | Update a custom field |
| DELETE | `/custom_fields/{custom_field_gid}` | Delete a custom field |
| GET | `/workspaces/{workspace_gid}/custom_fields` | Get workspace custom fields |

#### Enum Options

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/custom_fields/{custom_field_gid}/enum_options` | Create enum option |
| POST | `/custom_fields/{custom_field_gid}/enum_options/insert` | Reorder enum options |
| PUT | `/enum_options/{enum_option_gid}` | Update enum option |

#### Custom Field Settings (Read-Only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/projects/{project_gid}/custom_field_settings` | Project custom fields |
| GET | `/portfolios/{portfolio_gid}/custom_field_settings` | Portfolio custom fields |
| GET | `/teams/{team_gid}/custom_field_settings` | Team custom fields |
| GET | `/goals/{goal_gid}/custom_field_settings` | Goal custom fields |

---

## Template Resources

### Project Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/project_templates` | Get multiple templates |
| GET | `/project_templates/{project_template_gid}` | Get a template |
| DELETE | `/project_templates/{project_template_gid}` | Delete a template |
| GET | `/teams/{team_gid}/project_templates` | Get team's templates |
| POST | `/project_templates/{project_template_gid}/instantiate` | Create project from template |

---

### Task Templates

Available to Premium, Business, and Enterprise customers.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/task_templates` | Get multiple templates |
| GET | `/task_templates/{task_template_gid}` | Get a template |
| DELETE | `/task_templates/{task_template_gid}` | Delete a template |
| POST | `/task_templates/{task_template_gid}/instantiate` | Create task from template |

---

## Membership Resources

### Memberships (Generic)

New unified endpoint for managing memberships across goals, projects, portfolios, and custom fields.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/memberships` | Get multiple memberships |
| POST | `/memberships` | Create a membership |
| GET | `/memberships/{membership_id}` | Get a membership |
| PUT | `/memberships/{membership_id}` | Update a membership |
| DELETE | `/memberships/{membership_id}` | Delete a membership |

---

### Workspace Memberships

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/workspace_memberships/{workspace_membership_gid}` | Get a membership |
| GET | `/users/{user_gid}/workspace_memberships` | Get user's memberships |
| GET | `/workspaces/{workspace_gid}/workspace_memberships` | Get workspace memberships |

#### Key Fields
- `is_active` - Currently associated with workspace
- `is_admin` - Admin status
- `is_guest` - Guest access
- `is_view_only` - View-only license
- `vacation_dates` - Start/end dates

---

### Team Memberships

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/team_memberships` | Get multiple memberships |
| GET | `/team_memberships/{team_membership_gid}` | Get a membership |
| GET | `/teams/{team_gid}/memberships` | Get team memberships |
| GET | `/users/{user_gid}/team_memberships` | Get user's team memberships |

#### Key Fields
- `is_guest` - Guest access
- `is_limited_access` - Limited access
- `is_admin` - Admin status

---

### Portfolio Memberships

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/portfolio_memberships` | Get multiple memberships |
| GET | `/portfolio_memberships/{id}` | Get a membership |
| GET | `/portfolios/{portfolio_id}/memberships` | Get portfolio memberships |

#### Access Levels
`admin`, `editor`, `viewer` (no commenter access for portfolios)

---

### Project Memberships (Deprecated)

**Note:** Use the new Memberships API instead.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/project_memberships/{project_membership_gid}` | Get a membership |
| GET | `/projects/{project_gid}/project_memberships` | Get project memberships |

#### Access Levels
`admin`, `editor`, `commenter`, `viewer`

---

## Time & Planning Resources

### Time Tracking Entries

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/time_tracking_entries` | Get multiple entries |
| POST | `/time_tracking_entries` | Create an entry |
| GET | `/time_tracking_entries/{time_tracking_entry_gid}` | Get an entry |
| PUT | `/time_tracking_entries/{time_tracking_entry_gid}` | Update an entry |
| DELETE | `/time_tracking_entries/{time_tracking_entry_gid}` | Delete an entry |
| GET | `/tasks/{task_gid}/time_tracking_entries` | Get entries for task |

#### Key Fields
- `duration_minutes` - Time tracked
- `entered_on` - Date logged
- `approval_status` - `APPROVED`, `DRAFT`, `REJECTED`, `SUBMITTED`
- `billable_status` - `billable`, `nonBillable`, `notApplicable`

---

### Time Periods

Read-only. Used for goals.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/time_periods` | Get multiple time periods |
| GET | `/time_periods/{time_period_gid}` | Get a time period |

#### Period Codes
`FY`, `H1`, `H2`, `Q1`, `Q2`, `Q3`, `Q4`

---

### Allocations

Enterprise feature for resource planning.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/allocations` | Get multiple allocations |
| POST | `/allocations` | Create an allocation |
| GET | `/allocations/{allocation_gid}` | Get an allocation |
| PUT | `/allocations/{allocation_gid}` | Update an allocation |
| DELETE | `/allocations/{allocation_gid}` | Delete an allocation |

#### Effort Types
- `hours` - Absolute hours
- `percent` - Percentage of capacity

---

## Integration Resources

### Webhooks

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/webhooks` | Get multiple webhooks |
| POST | `/webhooks` | Create a webhook |
| GET | `/webhooks/{webhook_gid}` | Get a webhook |
| PUT | `/webhooks/{webhook_gid}` | Update a webhook |
| DELETE | `/webhooks/{webhook_gid}` | Delete a webhook |

#### Webhook Configuration
- `resource` - Resource to monitor
- `target` - URL to receive POST requests
- `filters[]` - Whitelist of event types
  - `resource_type` - Type being monitored
  - `action` - `changed`, `added`, `removed`, `deleted`, `undeleted`
  - `fields[]` - Specific fields (for `changed` only)

#### Health Metrics
- `created_at`
- `last_success_at`
- `last_failure_at`
- `delivery_retry_count`

---

### Events

Streaming interface for resource changes.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/events` | Get events for a resource |

#### Parameters
- `resource` (required) - Resource ID
- `sync` - Token for pagination

#### Behavior
- Initial request without `sync` returns 412 with sync token
- Events available for 24 hours
- "At most once" delivery
- Changes bubble up (project events include task changes)

---

### Rules

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/rules/{rule_trigger_gid}/trigger` | Trigger a rule |

Enables cross-application workflows via incoming web requests.

---

## Utility Resources

### Batch API

Submit parallel requests.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/batch` | Submit batch request |

#### Response Structure
```json
{
  "status_code": 200,
  "headers": { "location": "/tasks/1234" },
  "body": { "data": { ... } }
}
```

---

### Jobs

Track asynchronous operations (duplication, export, etc.).

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/jobs/{job_gid}` | Get job status |

Only the creator can access job status.

---

### Typeahead

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/workspaces/{workspace_gid}/typeahead` | Search by name prefix |

Returns `AsanaNamedResource` objects with `gid`, `resource_type`, `name`.

---

### Organization Exports

Enterprise+ only.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/organization_exports` | Create export request |
| GET | `/organization_exports/{id}` | Get export status |

Poll until state is `finished`, then download. URLs valid for ~1 hour.

---

### Audit Log API

Enterprise+ only. Read-only.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/workspaces/{workspace_gid}/audit_log_events` | Get audit events |

Returns events with:
- `event_type`, `event_category`
- `actor` - Who triggered
- `resource` - What was affected
- `context` - IP address, auth method

---

## Additional Resources

### User Task Lists

Personal "My Tasks" lists.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/user_task_lists/{user_task_list_gid}` | Get a task list |
| GET | `/users/{user_gid}/user_task_list` | Get user's task list |

---

### Project Briefs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/project_briefs/{project_brief_gid}` | Get a brief |
| PUT | `/project_briefs/{project_brief_gid}` | Update a brief |
| DELETE | `/project_briefs/{project_brief_gid}` | Delete a brief |
| POST | `/projects/{project_gid}/project_briefs` | Create a brief |

---

### Project Statuses (Deprecated)

**Note:** Use Status Updates instead.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/project_statuses/{id}` | Get a status |
| DELETE | `/project_statuses/{id}` | Delete a status |
| GET | `/projects/{id}/statuses` | Get project statuses |
| POST | `/projects/{id}/statuses` | Create project status |

**Colors:** `blue`, `complete`, `green`, `red`, `yellow`

---

## Quick Reference: All Resources

| Resource | CRUD | Notes |
|----------|------|-------|
| Tasks | CRUD | Core work unit |
| Projects | CRUD | Organize tasks |
| Workspaces | R-U- | Top-level container |
| Users | R-U- | Account info |
| Teams | CRU- | Group users |
| Sections | CRUD | Organize within projects |
| Tags | CRUD | Cross-project labels |
| Stories | CRUD | Comments & activity |
| Attachments | CR-D | File attachments |
| Status Updates | CR-D | Cannot update after create |
| Portfolios | CRUD | Group projects |
| Goals | CRUD | Objectives tracking |
| Goal Relationships | -RU- | Plus add/remove |
| Custom Fields | CRUD | Premium feature |
| Project Templates | -R-D | Plus instantiate |
| Task Templates | -R-D | Plus instantiate |
| Memberships | CRUD | Generic membership |
| Workspace Memberships | -R-- | Read only |
| Team Memberships | -R-- | Read only |
| Portfolio Memberships | -R-- | Read only |
| Time Tracking | CRUD | Time entries |
| Time Periods | -R-- | Goal time ranges |
| Allocations | CRUD | Enterprise |
| Webhooks | CRUD | Event notifications |
| Events | -R-- | Streaming changes |
| Rules | ---T | Trigger only |
| Batch | ---B | Parallel requests |
| Jobs | -R-- | Async status |
| Typeahead | -R-- | Name search |
| Org Exports | CR-- | Enterprise+ |
| Audit Log | -R-- | Enterprise+ |
| User Task Lists | -R-- | My Tasks |
| Project Briefs | CRUD | Project descriptions |

---

## Further Resources

- [Official API Documentation](https://developers.asana.com/reference/rest-api-reference)
- [Client Libraries](https://developers.asana.com/docs/official-client-libraries)
- [API Explorer](https://developers.asana.com/explorer)
- [Changelog](https://developers.asana.com/docs/changelog)
