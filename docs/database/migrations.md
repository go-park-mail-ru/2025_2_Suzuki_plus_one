# Database migrations

Мы используем [golang-migrate/migrate](https://github.com/golang-migrate/migrate) для управления миграциями базы данных.

## Installation

См. официальный [гайд](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Create a new migration

См. официальный [гайд](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md#postgresql-tutorial-for-beginners)

```bash
# cd 2025_2_Suzuki_plus_one/  # Stay in project root folder
make migrate-create NAME=<migration_name>
```

## Apply migrations

```bash
make migrate-up     # to apply all up migrations
make migrate-down   # to apply all down migrations
```
