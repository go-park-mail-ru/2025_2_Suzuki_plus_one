# Githooks

To install githooks you simply copy them into .git folder

```bash
chmod +x githooks/install.sh
sh githooks/install.sh
```

## pre-commit

Runs `go fmt ./...` and `go vet ./...` to set code formatting and prevent possible errors.
Also check for `go test ./...` that all tests are passing.
