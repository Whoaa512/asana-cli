#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<EOF
Usage: $0 <version> [options]

Options:
  -n, --notes STRING       Release notes as a string
  -f, --notes-file PATH    Path to file containing release notes
  -h, --help               Show this help

Examples:
  $0 1.0.0
  $0 1.0.0 --notes "Fixed bugs"
  $0 1.0.0 --notes-file ./notes.md
EOF
  exit 0
}

VERSION=""
NOTES_STRING=""
NOTES_FILE_PATH=""

while [[ $# -gt 0 ]]; do
  case $1 in
    -h|--help)
      usage
      ;;
    -n|--notes)
      NOTES_STRING="$2"
      shift 2
      ;;
    -f|--notes-file)
      NOTES_FILE_PATH="$2"
      shift 2
      ;;
    *)
      if [[ -z "$VERSION" ]]; then
        VERSION="$1"
        shift
      else
        echo "Error: Unknown argument: $1"
        usage
      fi
      ;;
  esac
done

if [[ -z "$VERSION" ]]; then
  echo "Error: Version required"
  echo ""
  usage
fi

if [[ "$VERSION" == v* ]]; then
  VERSION="${VERSION#v}"
fi

TAG="v$VERSION"

if git rev-parse "$TAG" >/dev/null 2>&1; then
  echo "Error: Tag $TAG already exists"
  exit 1
fi

if [[ -n "$(git status --porcelain)" ]]; then
  echo "Error: Working directory not clean"
  exit 1
fi

NOTES_FILE=$(mktemp)
CLEANUP_NOTES_FILE=true

if [[ -n "$NOTES_STRING" ]]; then
  echo "Using provided notes string..."
  echo "$NOTES_STRING" > "$NOTES_FILE"
elif [[ -n "$NOTES_FILE_PATH" ]]; then
  echo "Using notes from file: $NOTES_FILE_PATH"
  if [[ ! -f "$NOTES_FILE_PATH" ]]; then
    echo "Error: Notes file not found: $NOTES_FILE_PATH"
    exit 1
  fi
  cp "$NOTES_FILE_PATH" "$NOTES_FILE"
else
  LAST_TAG=$(git tag --sort=-version:refname | head -1)
  if [[ -z "$LAST_TAG" ]]; then
    echo "No previous tag found, using initial commit"
    LAST_TAG=$(git rev-list --max-parents=0 HEAD)
  fi

  echo "Generating release notes from $LAST_TAG to HEAD..."

  cat > "$NOTES_FILE" <<EOF
## New Commands

$(git log "$LAST_TAG..HEAD" --oneline --no-merges | sed 's/^[a-f0-9]* /- /')

## Install

\`\`\`bash
go install github.com/whoaa512/asana-cli/cmd/asana@$TAG
\`\`\`

Or download pre-built binaries from the release assets.
EOF

  echo ""
  echo "Opening release notes in editor..."
  echo "Edit the notes to categorize changes properly."
  "${EDITOR:-vim}" "$NOTES_FILE"
fi

echo "Building binaries for $TAG..."
rm -rf dist
mkdir -p dist

PLATFORMS=(
  "darwin:arm64"
  "darwin:amd64"
  "linux:arm64"
  "linux:amd64"
)

for platform in "${PLATFORMS[@]}"; do
  GOOS="${platform%%:*}"
  GOARCH="${platform##*:}"
  OUTPUT="dist/asana-${GOOS}-${GOARCH}"
  echo "  Building $OUTPUT..."
  GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags "-X main.version=$VERSION" -o "$OUTPUT" ./cmd/asana
done

echo "Creating tag $TAG..."
git tag -a "$TAG" -F "$NOTES_FILE"

echo "Pushing tag..."
git push origin "$TAG"

echo "Creating GitHub release..."
gh release create "$TAG" \
  dist/asana-darwin-arm64 \
  dist/asana-darwin-amd64 \
  dist/asana-linux-arm64 \
  dist/asana-linux-amd64 \
  --title "$TAG" \
  --notes-file "$NOTES_FILE"

rm -f "$NOTES_FILE"
echo "Done! https://github.com/Whoaa512/asana-cli/releases/tag/$TAG"
