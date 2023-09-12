CREATE USER app;
CREATE DATABASE auth;
--ユーザーにDBの権限をまとめて付与
GRANT ALL PRIVILEGES ON DATABASE auth TO app;
--ユーザーを切り替え
\c auth
--テーブルを作成
--テーブルにデータを挿入