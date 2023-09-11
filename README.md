# go-oauth2

## Endpoint

- GET /authorize -> display authorization information
- POST /authorization -> return authorization code
- POST /token -> return token information

## Table structure

### oauth2_clients

| name | type   |
| ---- | ------ |
| id   | uuid   |
| name | string |

### oauth2_codes

| name      | type   |
| --------- | ------ |
| code      | string |
| client_id | uuid   |
| user_id   | uuid   |

### oauth2_tokens

| name         | type     |
| ------------ | -------- |
| client_id    | uuid     |
| user_id      | uuid     |
| access_token | string   |
| expires      | datetime |

### oauth2_refresh_tokens

| name          | type     |
| ------------- | -------- |
| access_token  | string   |
| refresh_token | string   |
| expires       | datetime |

## Reference

- <https://qiita.com/TakahikoKawasaki/items/e508a14ed960347cff11>
