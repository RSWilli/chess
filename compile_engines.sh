#!/usr/bin/env bash

#
# This script is used to compile all engines from all tags in this repo.
# The compiled engines can then be used from the UI
#

# --- Configuration ---
OUTPUT_DIR="$(pwd)/engines"
BINARY_PREFIX="engine_"
# ---------------------

echo "🚀 Starting build process using Git worktrees..."

for tag in $(git tag); do
  output_file="$OUTPUT_DIR/${BINARY_PREFIX}${tag}"
  
  # Idempotency: Skip if already built
  if [ -f "$output_file" ]; then
    echo "⏭️  Skipping $tag: Binary already exists"
    continue
  fi

  echo "🔨 Building tag: $tag..."
  
  # 1. Create a temporary directory for the worktree
  TEMP_DIR=$(mktemp -d)
  
  # 2. Add a worktree for the specific tag (silently)
  git worktree add --detach "$TEMP_DIR" "$tag" --quiet
  
  # 3. Navigate to the temp directory, build, and capture the exit code
  (cd "$TEMP_DIR" && go build -o "$output_file" ./cmd/engine)
  BUILD_STATUS=$?
  
  # 4. Clean up the worktree and temp directory
  git worktree remove "$TEMP_DIR" --force
  
  if [ $BUILD_STATUS -eq 0 ]; then
    echo "✅ Successfully built $tag"
  else
    echo "❌ Failed to build $tag"
  fi
done

echo "🎉 All done! Binaries are located in $OUTPUT_DIR"