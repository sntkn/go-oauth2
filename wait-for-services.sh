#!/bin/bash

# databaseコンテナのIDを取得
database_container_id=$(docker compose ps -q database)

# データベースが起動するまで待つ
until docker exec "$database_container_id" pg_isready; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"


# session コンテナのIDを取得
session_container_id=$(docker compose ps -q session)

# セッションサービスが起動するまで待つ
until docker exec "$session_container_id" redis-cli -h session ping; do
  >&2 echo "Session service is unavailable - sleeping"
  sleep 1
done

>&2 echo "Session service is up - all services are ready"
