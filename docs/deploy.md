# Deploy

## Setup .env file with Ansible Vault

```bash
cd deployments

# Decrypt .env
ansible-vault decrypt .env.vault # Type vault password when prompted
# Copy decrypted data to .env (make sure .env is in .gitignore)
cp .env.vault .env
# Make password script executable before using (it sources .env file etc)
chmod +x vault_password.sh
# IMPORTANT: Encrypt .env.vault back
ansible-vault encrypt .env.vault  --vault-password-file=vault_password.sh
```

## Run Playbook

```bash
# Create top .env file before running playbook according to .env.example
cd deployments
# Decrypt prod.env
ansible-vault decrypt prod.env.vault --vault-password-file=vault_password.sh
# Copy decrypted data to prod.env (make sure prod.env is in .gitignore)
cp prod.env.vault prod.env # Note: systemd expects the strict key=value format
# IMPORTANT: Encrypt prod.env.vault back
ansible-vault encrypt prod.env.vault --vault-password-file=vault_password.sh

# Make sure you have access to the server via SSH key specified in deployments/.env file
# This will git pull deploy branch, build and restart the systemd service
# Which is placed in /lib/systemd/system/APIserver.service
ansible-playbook backend.yaml --vault-password-file=vault_password.sh
# Read logs from: /var/log/apiserver/error.log
ansible-playbook backend-logs.yaml --vault-password-file=vault_password.sh 
```
