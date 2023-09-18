# go-oauth2

## Endpoint

- GET /authorize -> display authorization information
- POST /authorization -> return authorization code
- POST /token -> return token information

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
> insert into oauth2_clients (id, name, redirect_uris, created_at, updated_at) values ('550e8400-e29b-41d4-a716-446655440000', 'test client', 'http://localhost:8080/callback', now(), now());
```

## Reference

- <https://qiita.com/TakahikoKawasaki/items/e508a14ed960347cff11>
