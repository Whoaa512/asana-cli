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
git tag -a "$TAG" -m "Release $TAG"

echo "Pushing tag..."
git push origin "$TAG"

echo "Creating GitHub release..."
gh-public release create "$TAG" \
  dist/asana-darwin-arm64 \
  dist/asana-darwin-amd64 \
  dist/asana-linux-arm64 \
  dist/asana-linux-amd64 \
  --title "$TAG" \
  --generate-notes

echo "Done! https://github.com/Whoaa512/asana-cli/releases/tag/$TAG"
