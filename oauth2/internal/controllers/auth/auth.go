package auth

import "github.com/sntkn/go-oauth2/oauth2/internal/redis"

type UseCase struct {
	redisCli *redis.RedisCli
}
