# Popfilms. Backend

Этот репозиторий содержит бэкенд проекта Popfilms.

## Popfilms

<span style="color: cyan;">Popfilms</span> -- это стриминговый сервис, вдохновленный Netflix.

Наш сервис предоставляет пользователям доступ к просмотру библиотеки фильмов, сериалов, телевизионных шоу всевозможных жанров на различных языках. Сервис преимущественно ориентирован на распространение развлекательного контента, однако также имеет внушительное количество документальных фильмов и телепрограмм.

## Релиз

[Ссылка](http://217.16.18.125/) на деплой приложения

> ДЗ по курсу баз данных можно найти в папке [docs/database/postgres](./docs/database/postgres/postgres.md)

## Реализованные функции

[Swagger](http://217.16.18.125/docs/)

- Регистрация и аутентификация пользователей
- Отдача замоканных объектов фильмов

### Детали реализации

- Бэкенд задеплоен при помощи systemd сервисов, которые управляются с помощью ansible плейбуков
- Чистейшая архитектура
- Авторизация реализована при помощи JWT Access и Refresh токенов
- База данных PostgreSQL взаимодействует с бэкендом через pgx библиотеку
- Объектное хранилище Minio
- Кеширование сессий пользователей в Redis
- Логирование с помощью Zap
- Тесты с использованием Testify и GoMock, SqlMock

### Как запустить проект локально

#### Установка

```bash
# Установка docker и docker-compose
[Ссылка](https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository)

# Установка Go
[Ссылка](https://go.dev/doc/install)

# Установка migrate
echo "export PATH=$PATH:$HOME/go/bin" >> ~/.bashrc
source ~/.bashrc
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Установка PostgreSQL клиента
sudo apt-get install -y postgresql-client
```

#### Запуск проекта

```bash
# Инициализация проекта
## копируем .env.example в .env
cp .env.example .env

# Запуск docker контейнеров и миграций
make all-bootstrap
```

#### Некоторые полезные команды

```bash
# Очистка старых данных
docker compose down -v

# Запуск баз данных
make all-prepare

## собираем и запускаем приложение без контейнера
make run



## запускаем тесты
make test
## собираем и проверяем покрытие тестами
make coverage

# Запуск докер компоса
make all-deploy

# Заполнение базы тестовыми данными и запуск приложения
make all-bootstrap
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
