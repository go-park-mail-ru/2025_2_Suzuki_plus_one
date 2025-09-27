#!/bin/bash
# This script installs pre-commit hooks into the .git/hooks directory of the repository.

# https://stackoverflow.com/questions/59895/how-do-i-get-the-directory-where-a-bash-script-is-located-from-within-the-script
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
REPO_ROOT=$(dirname "$SCRIPT_DIR")
HOOKS_SRC_DIR="$REPO_ROOT/githooks"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

# List of pre-commit hook scripts and directories to install
PRECOMMIT_HOOKS=("pre-commit")
PRECOMMIT_DIRS=("pre-commit.d")
INSTALLED_HOOKS=()
INSTALLED_DIRS=()

# Function to check for overwriting and create symbolic links
link_hook() {
    local src="$1"
    local dest="$2"

    if [ -f "$dest" ]; then
        echo "Hook $(basename "$dest") already exists in $HOOKS_DIR."
        read -p "Overwrite? [y/N]: " answer
        if [[ "$answer" =~ ^[Yy]$ ]]; then
            ln -sf "$src" "$dest"
            echo "Overwritten $(basename "$dest")."
        else
            echo "Skipped $(basename "$dest")."
            return
        fi
    else
        ln -s "$src" "$dest"
        echo "Linked $(basename "$dest")."
    fi
}

# Function to link directories
link_dir() {
    local src="$1"
    local dest="$2"

    if [ -d "$src" ]; then
        ln -s "$src" "$dest"
        echo "Linked directory $(basename "$dest")."
        INSTALLED_DIRS+=("$(basename "$dest")")
    else
        echo "Warning: Directory $src not found."
    fi
}

# Ensure pre-commit directories are linked
for dir in "${PRECOMMIT_DIRS[@]}"; do
    link_dir "$HOOKS_SRC_DIR/$dir" "$HOOKS_DIR/$dir"
done

# Install pre-commit hooks
echo "... Installing pre-commit hooks from $HOOKS_SRC_DIR to $HOOKS_DIR"
for hook in "${PRECOMMIT_HOOKS[@]}"; do
    if [ -f "$HOOKS_SRC_DIR/$hook" ]; then
        chmod +x "$HOOKS_SRC_DIR/$hook"
        link_hook "$HOOKS_SRC_DIR/$hook" "$HOOKS_DIR/$hook"
        INSTALLED_HOOKS+=("$hook")
    else
        echo "Warning: $hook not found in $HOOKS_SRC_DIR"
    fi
done

# Confirm installation of all hooks and directories
if [ ${#INSTALLED_HOOKS[@]} -gt 0 ]; then
    echo "Installed pre-commit hooks: ${INSTALLED_HOOKS[*]}"
else
    echo "No pre-commit hooks were installed."
fi

if [ ${#INSTALLED_DIRS[@]} -gt 0 ]; then
    echo "Installed pre-commit directories: ${INSTALLED_DIRS[*]}"
else
    echo "No pre-commit directories were installed."
fi