# Database

Мы используем postgres

## Start up the database

[Установить](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)
docker и docker compose и стартануть базу

```bash
make db-start
```

## Migration

```bash
make migrate-up
make migrate-down
```
