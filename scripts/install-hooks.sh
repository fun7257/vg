#!/bin/bash

# Script to install Git hooks

set -e

HOOKS_DIR=".git/hooks"
GITHOOKS_DIR=".githooks"

if [ ! -d "$GITHOOKS_DIR" ]; then
    echo "Error: .githooks directory not found"
    exit 1
fi

if [ ! -d "$HOOKS_DIR" ]; then
    echo "Error: .git/hooks directory not found. Are you in a git repository?"
    exit 1
fi

echo "Installing Git hooks..."

# Install pre-commit hook
if [ -f "$GITHOOKS_DIR/pre-commit" ]; then
    cp "$GITHOOKS_DIR/pre-commit" "$HOOKS_DIR/pre-commit"
    chmod +x "$HOOKS_DIR/pre-commit"
    echo "✓ Installed pre-commit hook"
else
    echo "Warning: pre-commit hook not found in .githooks/"
fi

echo ""
echo "Git hooks installed successfully!"
echo ""
echo "The pre-commit hook will run golangci-lint before each commit."
echo "To skip the hook, use: git commit --no-verify"

