# go-oauth2

## Endpoint

- GET /oauth2/authorize -> authorization endpoit
- POST /oauth2/authorization -> return authorization code
- POST /oauth2/token -> return token information
- GET /me -> return user information
- DELETE /oauth2/token -> revoke token

- GET|POST /client/signin
- GET|POST /client/signup

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
