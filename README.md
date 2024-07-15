# go-oauth2

## Endpoint

- GET /authorize -> display authorization information
- POST /authorization -> return authorization code
- POST /token -> return token information
- GET /me -> return user information
- DELETE /token -> revoke token

## Table structure

### users

| name     | type   |
| -------- | ------ |
| id       | uuid   |
| name     | string |
| email    | string |
| password | string |

### oauth2_clients

| name          | type   |
| ------------- | ------ |
| id            | uuid   |
| name          | string |
| redirect_uris | string |

### oauth2_codes

| name         | type      |
| ------------ | --------- |
| code         | string    |
| client_id    | uuid      |
| user_id      | uuid      |
| scope        | string    |
| redirect_uri | string    |
| expires_at   | timestamp |

### oauth2_tokens

| name         | type      |
| ------------ | --------- |
| access_token | string    |
| client_id    | uuid      |
| user_id      | uuid      |
| scope        | string    |
| expires_at   | timestamp |

### oauth2_refresh_tokens

| name          | type      |
| ------------- | --------- |
| access_token  | string    |
| refresh_token | string    |
| expires_at    | timestamp |

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

## Reference

- <https://qiita.com/TakahikoKawasaki/items/e508a14ed960347cff11>
- <https://qiita.com/Daiius/items/9b3f26137380de74d049>
- <https://qiita.com/tatsurou313/items/ad86da1bb9e8e570b6fa#docker-compose-における-buildkit-の利用>
- <https://tech-lab.sios.jp/archives/39388#i-9>
