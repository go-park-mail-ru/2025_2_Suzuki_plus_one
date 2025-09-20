#!/bin/bash
# This script installs pre-commit hooks into the .git/hooks directory of the repository.

# https://stackoverflow.com/questions/59895/how-do-i-get-the-directory-where-a-bash-script-is-located-from-within-the-script
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
REPO_ROOT=$(dirname "$SCRIPT_DIR")
HOOKS_SRC_DIR="$REPO_ROOT/githooks"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

# List of pre-commit hook scripts to install
PRECOMMIT_HOOKS=("pre-commit")
INSTALLED_HOOKS=()

echo "... Installing pre-commit hooks from $HOOKS_SRC_DIR to $HOOKS_DIR"
for hook in "${PRECOMMIT_HOOKS[@]}"; do
    if [ -f "$HOOKS_SRC_DIR/$hook" ]; then
        chmod +x "$HOOKS_SRC_DIR/$hook"
        # Prompt if the hook already exists
        if [ -f "$HOOKS_DIR/$hook" ]; then
            echo "Hook $hook already exists in $HOOKS_DIR."
            read -p "Overwrite? [y/N]: " answer
            if [[ "$answer" =~ ^[Yy]$ ]]; then
                ln -sf "$HOOKS_SRC_DIR/$hook" "$HOOKS_DIR/"
            else
                echo "Skipped $hook"
                continue
            fi
        else
            ln -s "$HOOKS_SRC_DIR/$hook" "$HOOKS_DIR/"
        fi
        INSTALLED_HOOKS+=("$hook")
    else
        echo "Warning: $hook not found in $HOOKS_SRC_DIR"
    fi
done

# Confirm installation of all hooks
if [ ${#INSTALLED_HOOKS[@]} -gt 0 ]; then
    echo "Installed pre-commit hooks: ${INSTALLED_HOOKS[*]}"
else
    echo "No pre-commit hooks were installed."
fi