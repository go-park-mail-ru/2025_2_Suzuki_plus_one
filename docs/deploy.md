# Deploy

## Setup Ansible environment

Access ansible vault to get SSH keys and vault password

```bash
cd deployments
# Decrypt ansible.env.vault
ansible-vault view ansible.env.vault > ansible.env # Type vault password when prompted
# (make sure ansible.env is in .gitignore)
# Make password script executable before using (it sources .env file etc)
chmod +x vault_password.sh
```

Now on you can use `--vault-password-file=vault_password.sh` flag with ansible commands
to escape typing vault password every time.

### Decrypt SSH key

```bash
ansible-vault view ssh-key.pem.vault --vault-password-file=vault_password.sh > ssh-key.pem
# Make ssh-key.pem file only readable by user (required by ssh)
chmod 400 ssh-key.pem
```

## Playbooks

### Deploy backend

> In first place, create top-level .env file before running playbook according to .env.example

#### Setup production environment file

```bash
cd deployments
# Decrypt prod.env
ansible-vault view prod.env.vault --vault-password-file=vault_password.sh > prod.env
# (make sure prod.env is in .gitignore)
# Note: systemd expects the strict key=value format
```

#### Run playbook to deploy backend

```bash
# Make sure you have access to the server via SSH key specified in deployments/ssh-key.pem
# Now you can run the playbook
# backend.yaml will git pull deploy branch, build and restart the systemd service
# Which is placed in /lib/systemd/system/APIserver.service
ansible-playbook backend.yaml --vault-password-file=vault_password.sh
# Read logs from: /var/log/apiserver/error.log
ansible-playbook backend-logs.yaml --vault-password-file=vault_password.sh
```

### Deploy API documentation

> Make sure you have already set up ansible environment as described [above](#setup-ansible-environment)

```bash
ansible-playbook update-api.yaml --vault-password-file=vault_password.sh
```
