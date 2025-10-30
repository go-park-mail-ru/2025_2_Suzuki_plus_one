# Style of this project

We use [hooks](./githooks.md) to enforce every commit has passed `go fmt`, `go vet` and `go test`.

## Project structure

Follow [standard](https://github.com/golang-standards/project-layout) project layout.

- `api` - API definitions (OpenAPI, protobuf, etc)
- `cmd` - application entry points
- `docs` - documentation
- `githooks` - git hooks (pre-commit)
- `internal` - application code
- `config` - configuration
- `deployments` - deployment files (docker, k8s. TODO: separate to `deployments/docker`, `deployments/k8s`, `deployments/ansible`)

### Clean architecture

Follow [clean architecture](https://github.com/evrone/go-clean-template) template.

- `cmd/<app-name>` - application entry point
- `internal/<app-name>/controller` - HTTP, gRPC handlers. Implement controller interfaces.
- `internal/<app-name>/usecase` - business logic layer, tiny and clear.
  Defines application flow. Defines controller and adapter interfaces. Depends on entity.
- `internal/<app-name>/entity` - domain specific entities.
  Defines domain constants, models and their validation. Depends on nothing.
- `internal/<app-name>/adapter` - data access layer. Implements adapter interfaces.
- `internal/<app-name>/dto` - data transfer objects. Depends on entity.

Follow Golang rule: Declare interfaces where you use them.

### Microservices

TODO: Add `popfilms` folder, where all microservices will be located.
Note that, each microservice should follow clean architecture too.

TODO: Edit `makefile` and `docker-compose.yml` to support microservices.

Use `pkg` folder for shared code between microservices.

## Code style

Follow [uber](https://github.com/uber-go/guide/blob/master/style.md) code style guide.

## Commit message style

Follow [conventional](https://gist.github.com/qoomon/5dfcdf8eec66a051ecd85625518cfd13) commits style.

### Commit message format

```
<type>(<scope>): <subject>
```

The next part describes types and scopes we should use.

#### Commit message types

- `feat` - refer to Notion's task tracker if applicable
- `fix` - mention certain [issue](https://github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/issues) if applicable
- `docs` - documentation only changes
- `style` - formatting, missing semi colons, etc; no code change
- `refactor` - refactoring production code, eg. renaming a variable
- `test` - adding missing tests, refactoring tests; no production code change
- `chore` - miscellaneous tasks, that better off not mentioned
- `ci` - changes to our CI configuration files and scripts

#### Commit message scopes

- `api` - changes to API definitions
- `githooks` - changes to git hooks

Other scopes should be the name of the affected package or file.

#### Examples

```
feat(auth): add auth entry point
chore(server): remove unused files
refactor(db): move db connection logic to internal/db
docs(api): update auth service description
feat(api): add auth service end point
```

## Branching model

Follow [conventional](https://conventional-branch.github.io/) branching model.

### Branches

- `main` - production branch, always deployable. Nobody should commit directly to this branch.
- `dev` - development branch. Nobody should commit directly to this branch.
- `feat/<TICKET_ID>-<feature-name>` - feature branches, created off `dev` branch, where `<TICKET_ID>` is the Notion's ID and `<feature-name>` is a short name of the feature being developed.

### Pull Requests

- Pull Requests should be made to `dev` branch
- Pull Requests should be reviewed by at least one team member
- Pull Request title should be in the format `<TICKET_ID>: <description>`, where `<TICKET_ID>` is the Notion's ID and `<description>` is a short summary of the work done.

### Releases

- Before each release, `dev` branch is merged into `main` branch
- Before each midterm exam, releases should be tagged with exam number, e.g. `RK 1`, `RK 2`, etc.
