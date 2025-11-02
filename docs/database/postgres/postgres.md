# Database

Мы используем postgres

> При обновлении БД, необходимо также обновлять и [схему](./schema.md)

## Start up the database

[Установить](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)
docker и docker compose и стартануть базу

```bash
# Удалить старую базу (если есть)
# make db-wipe
make db-start
```

## Migration

```bash
# Создать схему базы
make migrate-up
```

## Seed with test data

```bash
# Заполнить тестовыми данными
make db-fill
```
