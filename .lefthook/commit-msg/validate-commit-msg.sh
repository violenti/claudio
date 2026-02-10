#!/usr/bin/env bash
set -euo pipefail

COMMIT_MSG_FILE="$1"
COMMIT_MSG=$(head -1 "$COMMIT_MSG_FILE")

PATTERN="^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\(.+\))?(!)?: .{1,}"

if ! echo "$COMMIT_MSG" | grep -qE "$PATTERN"; then
    echo ""
    echo "ERROR: Commit message does not follow Conventional Commits format."
    echo ""
    echo "  Expected: <type>[optional scope]: <description>"
    echo "  Example:  feat(ai): add streaming support"
    echo "  Example:  fix: resolve nil pointer in provider"
    echo ""
    echo "  Allowed types: build, chore, ci, docs, feat, fix, perf, refactor, revert, style, test"
    echo ""
    exit 1
fi
