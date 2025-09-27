#!/bin/bash

# Check for unencrypted .env.vault file before committing

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
REPO_ROOT=$(dirname "$(dirname "$SCRIPT_DIR")")

vault_file="$REPO_ROOT/deployments/.env.vault"


PREFIX="PRECOMMIT: LEAKS: "

hook_log() {
	echo "${PREFIX}$1"
}

hook_log "Checking for unencrypted $vault_file."

if [ -f "$vault_file" ] && ! grep -q '^$ANSIBLE_VAULT' "$vault_file"; then
    hook_log "You are probably trying to commit unencrypted $vault_file!"
    hook_log "Please encrypt it as soon as possible using:"
    hook_log "ansible-vault encrypt $vault_file --vault-password-file=vault_password.sh"
    exit 1
else
    hook_log "OK, no unencrypted vault file found."
fi
