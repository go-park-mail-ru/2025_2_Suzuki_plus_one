# Githooks

To install githooks you simply symlink them into .git/hooks folder

```bash
chmod +x githooks/install.sh
./githooks/install.sh
```

## pre-commit

Calls all scripts in `githooks/pre-commit.d/`

### pre-commit.d scripts

- `check_leaks.sh`: Prevent committing files that may contain sensitive information like passwords or API keys.
- `go_checks.sh`: Runs `go fmt ./...` and `go vet ./...` to set code formatting and prevent possible errors. Also check for `go test ./...` that all tests are passing.
- `check_trails.sh`: Prevent committing files with trailing whitespace. (sample git hook)
