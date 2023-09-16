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

## Reference

- <https://qiita.com/TakahikoKawasaki/items/e508a14ed960347cff11>
