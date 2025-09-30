#!/bin/bash
# This script installs pre-commit hooks into the .git/hooks directory of the repository.

set -e

# https://stackoverflow.com/questions/59895/how-do-i-get-the-directory-where-a-bash-script-is-located-from-within-the-script
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
REPO_ROOT=$(dirname "$SCRIPT_DIR")
HOOKS_SRC_DIR="$REPO_ROOT/githooks"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

# List of pre-commit hook scripts and directories to install
PRECOMMIT_HOOKS=("pre-commit")
PRECOMMIT_DIRS=() # This dirs will be symnlinked into .git/hooks
                  # ! Don't use because linked scripts can't define their relative paths correctly
INSTALLED_HOOKS=()
INSTALLED_DIRS=()

echo "Installing hooks from $HOOKS_SRC_DIR to $HOOKS_DIR"

# Function to check for overwriting and create symbolic links
link_hook() {
    local src="$1"
    local dest="$2"
    local hook_name
    hook_name="$(basename "$dest")"

    if [ -f "$dest" ]; then
        echo "Hook $hook_name already exists in $HOOKS_DIR."
        read -p "Overwrite? [y/N]: " answer
        if [[ "$answer" =~ ^[Yy]$ ]]; then
            ln -sf "$src" "$dest"
            INSTALLED_HOOKS+=("$hook_name")
            echo "Overwritten $hook_name."
        else
            echo "Skipped $hook_name."
            return
        fi
    else
        ln -s "$src" "$dest"
        INSTALLED_HOOKS+=("$hook_name")
        echo "Linked $hook_name."
    fi
}

# ? I think this is not needed
# ? Why can't we just readlink the directories?
# ? Then linking dir is obsolete

# Function to link directories
link_dir() {
    local src="$1"
    local dest="$2"
    local dir_name
    dir_name="$(basename "$dest")"

    if [ -d "$src" ]; then
        if [ -L "$dest" ] || [ -d "$dest" ]; then
            echo "Directory $dir_name already exists in $HOOKS_DIR."
            read -p "Overwrite? [y/N]: " answer
            if [[ "$answer" =~ ^[Yy]$ ]]; then
                rm "$dest"
                ln -s "$src" "$dest"
                echo "Overwritten directory $dir_name."
                INSTALLED_DIRS+=("$dir_name")
            else
                echo "Skipped directory $dir_name."
            fi
        else
            ln -s "$src" "$dest"
            echo "Linked directory $dir_name."
            INSTALLED_DIRS+=("$dir_name")
        fi
    else
        echo "Warning: Directory $src not found."
    fi
}

# ! Don't use
# # Link PRECOMMIT_DIRS
# for dir in "${PRECOMMIT_DIRS[@]}"; do
#     link_dir "$HOOKS_SRC_DIR/$dir" "$HOOKS_DIR/$dir"
#     chmod -R +x "$HOOKS_SRC_DIR/$dir"
# done

# Link PRECOMMIT_HOOKS
for hook in "${PRECOMMIT_HOOKS[@]}"; do
    if [ -f "$HOOKS_SRC_DIR/$hook" ]; then
        chmod +x "$HOOKS_SRC_DIR/$hook"
        link_hook "$HOOKS_SRC_DIR/$hook" "$HOOKS_DIR/$hook"
    else
        echo "Warning: $hook not found in $HOOKS_SRC_DIR"
    fi
done

echo "- - - - - -"

# Confirm installation of all hooks and directories
if [ ${#INSTALLED_HOOKS[@]} -gt 0 ]; then
    echo "Installed pre-commit hooks: ${INSTALLED_HOOKS[*]}"
else
    echo "No pre-commit hooks were installed."
fi

# ! Don't use
# if [ ${#INSTALLED_DIRS[@]} -gt 0 ]; then
#     echo "Installed pre-commit directories: ${INSTALLED_DIRS[*]}"
# else
#     echo "No pre-commit directories were installed."
# fi