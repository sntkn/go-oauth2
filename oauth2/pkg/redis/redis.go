package redis

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-redis/redis/v8"
)

type Options struct {
	Addr     string // Redisのアドレスとポート番号に合わせて変更してください
	Password string // Redisにパスワードが設定されている場合は設定してください
	DB       int    // データベース番号
}

type RedisCli struct {
	cli *redis.Client
}

func NewClient(ctx context.Context, o Options) (*RedisCli, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     o.Addr,     // Redisのアドレスとポート番号に合わせて変更してください
		Password: o.Password, // Redisにパスワードが設定されている場合は設定してください
		DB:       o.DB,       // データベース番号
	})

	// ピングしてRedis接続を確認
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &RedisCli{
		cli: cli,
	}, nil
}

func (r *RedisCli) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.cli.Set(ctx, key, value, expiration).Err()
}

func (r *RedisCli) Get(ctx context.Context, key string) ([]byte, error) {
	return r.cli.Get(ctx, key).Bytes()
}

func (r *RedisCli) Del(ctx context.Context, key string) error {
	return r.cli.Del(ctx, key).Err()
}

func (r *RedisCli) GetOrNil(ctx context.Context, key string) ([]byte, error) {
	b, err := r.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return b, nil
}
