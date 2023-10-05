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
	DB       int64  // データベース番号
}

type RedisCli struct {
	cli *redis.Client
}

func NewClient(ctx context.Context, o Options) (*RedisCli, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redisのアドレスとポート番号に合わせて変更してください
		Password: "",               // Redisにパスワードが設定されている場合は設定してください
		DB:       0,                // データベース番号
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

func (r *RedisCli) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	return r.cli.Set(ctx, key, value, expiration)
}

func (r *RedisCli) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.cli.Get(ctx, key)
}

func (r *RedisCli) Del(ctx context.Context, key string) *redis.IntCmd {
	return r.cli.Del(ctx, key)
}

func (r *RedisCli) GetOrNil(ctx context.Context, key string) ([]byte, error) {
	ret := r.cli.Get(ctx, key)
	d, err := ret.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return d, nil
}
