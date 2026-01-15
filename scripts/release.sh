#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-}"
if [[ -z "$VERSION" ]]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 1.0.0"
  exit 1
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

LAST_TAG=$(git tag --sort=-version:refname | head -1)
if [[ -z "$LAST_TAG" ]]; then
  echo "No previous tag found, using initial commit"
  LAST_TAG=$(git rev-list --max-parents=0 HEAD)
fi

echo "Generating release notes from $LAST_TAG to HEAD..."
NOTES_FILE=$(mktemp)

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
gh-public release create "$TAG" \
  dist/asana-darwin-arm64 \
  dist/asana-darwin-amd64 \
  dist/asana-linux-arm64 \
  dist/asana-linux-amd64 \
  --title "$TAG" \
  --notes-file "$NOTES_FILE"

rm -f "$NOTES_FILE"
echo "Done! https://github.com/Whoaa512/asana-cli/releases/tag/$TAG"
