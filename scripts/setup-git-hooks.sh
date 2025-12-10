#!/bin/bash
# Setup script for installing git hooks
# Run this after cloning the repository to enable local commit validation

set -e

echo "Setting up git hooks for Terraform Typesense Provider..."
echo ""

# Check if we're in a git repository
if [ ! -d .git ]; then
    echo "Error: Not in a git repository. Please run this script from the repository root."
    exit 1
fi

# Check if Node.js is installed
if ! command -v node >/dev/null 2>&1; then
    echo "Warning: Node.js is not installed."
    echo "Commit message validation requires Node.js and npm."
    echo "Please install Node.js from https://nodejs.org/"
    echo ""
    echo "Skipping git hook installation."
    exit 1
fi

# Install npm dependencies if not already installed
if [ ! -d node_modules ]; then
    echo "Installing npm dependencies..."
    npm install
    echo ""
fi

# Install commit-msg hook
if [ -f .github/hooks/commit-msg.sample ]; then
    echo "Installing commit-msg hook..."
    cp .github/hooks/commit-msg.sample .git/hooks/commit-msg
    chmod +x .git/hooks/commit-msg
    echo "✓ commit-msg hook installed"
else
    echo "Warning: commit-msg.sample not found in .github/hooks/"
fi

echo ""
echo "=========================================="
echo "Git hooks installed successfully!"
echo "=========================================="
echo ""
echo "Your commits will now be validated against the Conventional Commits specification."
echo ""
echo "Example commit messages:"
echo "  feat: add support for nested fields"
echo "  fix: resolve collection update issue"
echo "  docs: update testing documentation"
echo ""
echo "To bypass the hook (not recommended):"
echo "  git commit --no-verify -m 'message'"
echo " "
