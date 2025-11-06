# Redis

Мы используем Redis для кеширования сессий пользователей.

## Ключи

Сессии хранятся с использованием двух типов ключей:

```bash
access:<token> -> userID
access:user:<userID> -> Set(token1, token2, token3 ...)
```

### Просмотр ключей

```bash
❯ redis-cli
127.0.0.1:6379> KEYS *
1) "access:user:1"
2) "access:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI0NDIxMzcsInVzZXJfaWQiOjF9.7YA1In9sdH9jsMxnouMi82TC5Xbon5xOWa0_lHRw91k"
```

### Пример значений ключей

```bash
127.0.0.1:6379> TYPE access:user:1
set
127.0.0.1:6379> GET access:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI0NDIxMzcsInVzZXJfaWQiOjF9.7YA1In9sdH9jsMxnouMi82TC5Xbon5xOWa0_lHRw91k
"1"
127.0.0.1:6379> SMEMBERS access:user:1
1) "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjI0NDIxMzcsInVzZXJfaWQiOjF9.7YA1In9sdH9jsMxnouMi82TC5Xbon5xOWa0_lHRw91k"
127.0.0.1:6379>
```
