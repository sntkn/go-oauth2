# go-oauth2

## oauth2

- gin-gonic/gin
- jmoiron/sqlx
- cockroachdb/errors

## api

- labstack/echo/v4
- gorm
- go-errors/errors

## front

- Next.js
- tailwindcss

## Sample data

```sql
> docker compose exec database psql -U app auth
> insert into oauth2_clients (id, name, redirect_uris, created_at, updated_at) values ('550e8400-e29b-41d4-a716-446655440000', 'test client', 'http://localhost:8000/callback', now(), now());
> insert into oauth2_clients (id, name, redirect_uris, created_at, updated_at) values ('684C406F-D7CA-42B0-B7AC-E2120B48B057', 'test client', 'http://localhost:3000/callback', now(), now());
> insert into users (id, name, email, password, created_at, updated_at) values ('4E77D89C-F28E-4232-BAC0-4ABB31B94590', 'test user', 'test@example.com', '$2a$10$LOzS79niq4E.hu8aib4GeuXVSII9OsYB.ReF/.BjqItfhaSnzWba6', now(), now());
```

## request

<http://localhost:8080/authorize?response_type=code&client_id=550e8400-e29b-41d4-a716-446655440000&scope=read&redirect_uri=http%3A%2F%2Flocalhost%3A8000%2Fcallback&state=ok>

<http://localhost:8080/authorize?response_type=code&client_id=684C406F-D7CA-42B0-B7AC-E2120B48B057&scope=read&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Fcallback&state=ok>

### input

email: test@example.com  
password: mypassword1234!

### tasks

```bash
mise run setup
mise run build
mise run test
mise run fmt
```

## dotenvx

```bash

# Create a ED25519 SSH KEY
cd oauth2
go run lib/generate_key.go

touch .env.development

# Encryption
dotenvx encrypt -f .env.development
```

## Others

- mise
- valkey
- Ed25991 (JWT signature algorithm)

## Reference

- <https://qiita.com/TakahikoKawasaki/items/e508a14ed960347cff11>
- <https://qiita.com/Daiius/items/9b3f26137380de74d049>
- <https://qiita.com/tatsurou313/items/ad86da1bb9e8e570b6fa#docker-compose-における-buildkit-の利用>
- <https://tech-lab.sios.jp/archives/39388#i-9>
