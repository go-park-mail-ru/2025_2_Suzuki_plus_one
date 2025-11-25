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

#### Get database data

Note that you need to download media assets for MinIO from remote storage to `testdata/minio` folder

```bash
# This will pull media assets with rclone of configured remote storage vkedu (.config/rclone/rclone.conf)
make minio-pull
```

[Databases](./database/README.md) will be filled with test data on the next step.

#### Run playbook to deploy backend

```bash
# Make sure you have access to the server via SSH key specified in deployments/ssh-key.pem
# Now you can run the playbook

# update-backend.yaml will git pull deploy branch, build and restart the services
ansible-playbook update-backend.yaml --vault-password-file=vault_password.sh -e deploy_mode=bootstrap

# The default behavior is to only update services without recreating DBs etc
ansible-playbook update-backend.yaml --vault-password-file=vault_password.sh

# Deprecated
# Read logs from: /var/log/apiserver/error.log
# ansible-playbook backend-logs.yaml --vault-password-file=vault_password.sh
```

### Deploy API documentation

> Make sure you have already set up ansible environment as described [above](#setup-ansible-environment)

```bash
ansible-playbook update-api.yaml --vault-password-file=vault_password.sh
```
