# Базы данных

Мы используем postgres как основную базу, redis как кеш и minio как хранилище медиафайлов.

Миграции базы данных находятся в папке [migrations/](../../migrations/).

Тестовые данные находятся в папке [testdata/](../../testdata/). См. [minio](./minio/minio.md#наполнение).

Дальнейшее описания каждой из баз данных находится в соответствующих файлах:

- [Postgres](./postgres/postgres.md)
- [Redis](./redis/redis.md)
- [MinIO](./minio/minio.md)

## Быстрый старт всех БД

```bash
make all-prepare
```

# Домашки находятся рядом:

- [RK3](./RK3.md)
- [RK4](./RK4.md)
