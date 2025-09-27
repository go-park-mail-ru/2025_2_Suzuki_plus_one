#!/bin/bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
REPO_ROOT=$(dirname "$SCRIPT_DIR")

source "$REPO_ROOT/.env"
echo "$VAULT_PASSWORD"