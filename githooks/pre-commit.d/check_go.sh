#!/bin/bash

PREFIX="PRECOMMIT: Go: "

hook_log() {
	echo "${PREFIX}$1"
}

hook_log "Starting checks..."

# Find staged Go files
files=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$')
if [ -z "$files" ]; then
	hook_log "No Go files to check. Exiting..."
    exit 0
fi


# Note: ./... means all packages in and below the current directory

# Run go fmt, automate formatting
if ! go fmt ./...; then
	hook_log "Go files not formatted:"
	hook_log "Please run 'go fmt ./...' to format your code."
	exit 1
fi

# Add formatted files to the staging area if go fmt made changes
if ! git diff --quiet; then
	git add $files
	hook_log "Go files were formatted and re-staged."
fi

# Run go vet to catch suspicious constructs
if ! go vet ./...; then
	hook_log "go vet found issues. Please fix them before committing."
	exit 1
fi

# Run go test to ensure tests pass
if ! go test ./...; then
	hook_log "go test failed. Please fix tests before committing."
	exit 1
fi

hook_log "Code reformatting and checks passed."