# asana-cli

Go CLI for Asana task management.

## Development

```bash
mise run build      # Build binary
mise run test       # Run tests
mise run lint       # Run linter
mise run check      # Lint + test
mise run install    # Install to $GOPATH/bin
```

## Releasing

```bash
# Interactive (opens $EDITOR for notes)
mise run release 1.3.0

# With LLM-drafted notes (skips editor)
mise run release 1.3.0 --notes "## Features
- Added new command
- Fixed bug"

# With notes from file
mise run release 1.3.0 --notes-file ./release-notes.md
```

The release script:
1. Validates clean working directory
2. Generates draft notes from commits (or uses provided notes)
3. Builds binaries for darwin/linux (arm64/amd64)
4. Creates annotated git tag
5. Pushes tag and creates GitHub release with binaries

When drafting release notes, categorize changes: New Commands, Improvements, Bug Fixes, etc.
