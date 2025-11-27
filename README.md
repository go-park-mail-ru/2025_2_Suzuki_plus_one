# Popfilms. Backend

Этот репозиторий содержит бэкенд проекта Popfilms.

Этот проект является частью учебной программы курса "Web разработка" от VK Education.

[Ссылка](https://github.com/go-park-mail-ru/lectures) на лекции

## Popfilms

Popfilms -- это стриминговый сервис, вдохновленный Netflix.

Наш сервис предоставляет пользователям доступ к просмотру библиотеки фильмов, сериалов, телевизионных шоу всевозможных жанров на различных языках. Сервис преимущественно ориентирован на распространение развлекательного контента, однако также имеет внушительное количество документальных фильмов и телепрограмм.

## Релиз

[Ссылка](https://popfilms.ru) на деплой приложения

> ДЗ по курсу баз данных можно найти в папке [docs/database/postgres](./docs/database/postgres/postgres.md)

## Фичи

[См. Swagger](https://popfilms.ru/docs/)

- Регистрация и аутентификация пользователей
- Отдача медиа контента (фильмы, сериалы и т.д.) через Minio
- ~~Управление сессиями пользователей через Redis~~
- Настройки пользователями (профиль, пароль)

### Детали реализации

- Чистейшая архитектура
- Авторизация реализована при помощи JWT Access и Refresh токенов
- База данных PostgreSQL взаимодействует с бэкендом через pgx библиотеку
- Миграции базы данных с помощью golang-migrate
- Объектное хранилище Minio
- Бэкенд задеплоен при помощи docker контейнеров, которые управляются ansible плейбуками
- Логирование с помощью Zap
- Тесты с использованием Testify и GoMock, SqlMock
- 3 микросервиса: HTTP, Auth, Search
- Взаимодействие между микросервисами через gRPC
- Мониторинг и алертинг с помощью Prometheus и Alertmanager + Grafana

### Как запустить проект

#### Установка зависимостей

```bash
# Установка docker и docker-compose
https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository

# Установка Go
https://go.dev/doc/install

# Установка migrate
echo "export PATH=$PATH:$HOME/go/bin" >> ~/.bashrc
source ~/.bashrc
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Установка protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Установка PostgreSQL клиента
sudo apt-get install -y postgresql-client
```

#### Запуск проекта в режиме прод

```bash
# Настройка окружения
cp deployments/prod.env.example .env

# Скачивание медиа ассетов
# (требуется ~/.config/rclone/rclone.conf с настройками доступа к облачному хранилищу)
cd <project root>
make minio-pull

# Запуск docker compose, миграций, наполнение тестовыми данными
make all-bootstrap
```

#### Деплой через ansible

[Настроить пароль доступа](./docs/deploy.md)

```bash
# После получения пароля и расшифровки ssh ключа

# Пересобрать образы микросервисов
make image-update

# Собрать проект с нуля (долго, так как базы пересоздаются)
make ansible-bootstrap

# Либо обновить только сервисы
make ansible-backend 
make update-deploy-backend # Для удобства: image-update + ansible-backend

# При обновлении миграций БД, docker compose, Makefile и т.д. требуется обновить ветку deploy и пересобрать проект
make update-deploy-branch # зальет dev в deploy и запушит изменения
make ansible-bootstrap
```

#### Запуск проекта на дев

```bash
# Настройка окружения
cp .env.example .env

# Запуск баз данных и наполнение их тестовыми данными
make all-prepare
# ИЛИ запуск баз без наполнения тестовыми данными
    docker compose up -d db redis minio
    # Применение миграций
    make all-migrate

# Запуск приложения
make run SERVICE=http
make run SERVICE=auth
make run SERVICE=search
```

#### Некоторые полезные команды

```bash
# Очистка старых данных
make all-wipe

## запускаем тесты
make test
## собираем и проверяем покрытие тестами
make coverage
```

## Другие ссылки

- [Репозиторий фронтэнда](https://github.com/frontend-park-mail-ru/2025_2_Suzuki_plus_one/)

## О нас

### Команда «Сузуки + 1»

- **Александр Федуков** — [github.com/sorrtory](https://github.com/sorrtory)
- **Фадеев Арсений** — [github.com/arsmfad](https://github.com/arsmfad)
- **Гилязева Динара** — [github.com/DinaraGil](https://github.com/DinaraGil)
- **Марышев Иван** — [github.com/ivanmaryshev](https://github.com/ivanmaryshev)

### Менторы

- **Володимир Коноплюк** — Go
- **Костин Глеб** — Frontend
- **Фильчаков Алексей** — СУБД
- **Даниил Хасьянов** — UX
