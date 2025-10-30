# Mocking

We use uber's `gomock` and `mockgen` to generate mocks for our interfaces.

## Installation

[Gomock guide](https://github.com/uber-go/mock?tab=readme-ov-file#installation)

```bash
# Ensure you have Go installed and your GOPATH/bin is in your PATH
go install go.uber.org/mock/mockgen@latest
```

## Mocks

To generate mocks for an interface, run the following command:

```bash
go generate ./...
```

Notice that we use `//go:generate` directives in the source files to specify the mock generation commands.
This allows us to run `go generate` to automatically generate all mocks defined in the project.

## See also

- [Gomock doc](https://pkg.go.dev/go.uber.org/mock/gomock)
