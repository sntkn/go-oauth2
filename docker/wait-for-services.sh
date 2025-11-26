#!/bin/bash

# 最大リトライ回数と待機間隔（秒）
max_retries=30
sleep_seconds=10

# databaseコンテナのIDを取得
database_container_id=$(docker compose ps -q database)

# データベースが起動するまで待つ
for i in $(seq 1 "$max_retries"); do
  if docker exec "$database_container_id" pg_isready >/dev/null 2>&1; then
    >&2 echo "Postgres is up - executing command"
    break
  fi

  >&2 echo "Postgres is unavailable - sleeping ($i/$max_retries)"
  sleep "$sleep_seconds"

  if [ "$i" -eq "$max_retries" ]; then
    >&2 echo "Postgres is still unavailable after $max_retries attempts"
    exit 1
  fi
done

>&2 echo "Postgres is up - executing command"


# kev コンテナのIDを取得
session_container_id=$(docker compose ps -q kvs)

# セッションサービスが起動するまで待つ（最大 max_retries 回）
for i in $(seq 1 "$max_retries"); do
  if docker exec "$session_container_id" redis-cli -h kvs ping >/dev/null 2>&1; then
    >&2 echo "Session service is up - all services are ready"
    break
  fi

  >&2 echo "Session service is unavailable - sleeping ($i/$max_retries)"
  sleep "$sleep_seconds"

  if [ "$i" -eq "$max_retries" ]; then
    >&2 echo "Session service is still unavailable after $max_retries attempts"
    exit 1
  fi
done

>&2 echo "Session service is up - all services are ready"
