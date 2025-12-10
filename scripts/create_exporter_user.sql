\echo Creating exporter role :exporter_user in database :dbname

-- на dev можно смело пересоздавать
DROP ROLE IF EXISTS :"exporter_user";

CREATE ROLE :"exporter_user" WITH
    LOGIN
    PASSWORD :'exporter_password'
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    NOINHERIT
    NOREPLICATION;

-- можно ограничиться правами только на статистику
GRANT CONNECT ON DATABASE :dbname TO :"exporter_user";

-- В Postgres 10+ уже есть преднастроенные роли:
-- pg_monitor даёт достаточно прав для мониторинга
GRANT pg_monitor TO :"exporter_user";

\echo Exporter role :exporter_user created
