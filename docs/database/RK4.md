# ДЗ 4 Оптимизация работы СУБД


## Установка Vegeta

```bash
go install github.com/tsenart/vegeta/...@latest
```

## Миграция базы данных

Перед запуском нагрузочного тестирования необходимо убедиться, что база данных содержит необходимые таблицы данные. 
Для этого нужно выполнить [миграции](../../migrations/), а также заполнение базы [тестовыми данными](../../testdata/postgres/):

Требуется установка `migrate`, см. [Миграции базы данных](./postgres/migrations.md)

```bash
make all-prepare
```

По ходу миграции создается сущность пользователь, а также сущность рекомендаций для данного пользователя.

## Скрипта для запуска нагрузочного тестирования

Был создан скрипт [scripts/vegeta/vegeta.sh](../../scripts/vegeta/vegeta.sh), который позволяет запускать нагрузочное тестирование с помощью Vegeta. Скрипт поддерживает следующие переменные окружения для настройки:

```bash
vegeta.sh --get-reco [--type movie|series] [--rate N] [--duration 30s] [--timeout 10s]
vegeta.sh --post-signup --users N [--rate N] [--duration 30s] [--timeout 10s]

# также через переменные окружения можно задать:
ENDPOINT=http://localhost:8080
REPORT_DIR=reports
TARGETS_DIR=targets
```

### Запуск скрипта

```bash
cd scripts/vegeta/
# Запрос на получения рекомендаций фильмов с нагрузкой 50 запросов в секунду в течение 30 секунд
ENDPOINT=https://popfilms.ru/api/v1 ./vegeta.sh --get-reco --type movie --rate 50 --duration 30s


# Создание 100_000 пользователей (с нагрузкой 100 запросов в секунду в течение 1000 секунд)
cd scripts/vegeta/
ENDPOINT=https://popfilms.ru/api/v1 ./vegeta.sh --post-signup --users 100000 --rate 100 --duration 1000s
```

## Результаты

Результаты нагрузочного тестирования сохраняются в директории, указанной в переменной окружения `REPORT_DIR` (по умолчанию [scipts/vegeta/reports](../../scripts/vegeta/reports)). В этой директории находятся файлы с результатами тестов, которые можно анализировать для оценки производительности системы.


