package redis

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テスト用の Redis サーバーをセットアップする関数
func setupTestRedis(t *testing.T) (*miniredis.Miniredis, *RedisCli) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	cli, err := NewClient(context.Background(), Options{
		Addr: mr.Addr(),
		DB:   0,
	})
	require.NoError(t, err)

	return mr, cli
}

func TestRedisCli(t *testing.T) {
	t.Parallel()
	mr, rdb := setupTestRedis(t)
	defer mr.Close()

	ctx := context.Background()

	// Test Set and Get
	err := rdb.Set(ctx, "testKey", "testValue", 0)
	require.NoError(t, err)

	val, err := rdb.Get(ctx, "testKey")
	require.NoError(t, err)
	assert.Equal(t, "testValue", string(val))

	// Test GetOrNil with existing key
	data, err := rdb.GetOrNil(ctx, "testKey")
	require.NoError(t, err)
	assert.Equal(t, []byte("testValue"), data)

	// Test GetOrNil with non-existing key
	data, err = rdb.GetOrNil(ctx, "nonExistingKey")
	require.NoError(t, err)
	assert.Nil(t, data)

	// Test Del
	err = rdb.Del(ctx, "testKey")
	require.NoError(t, err)

	val, err = rdb.Get(ctx, "testKey")
	require.Error(t, err)
	assert.Equal(t, redis.Nil, err)
	assert.Nil(t, val)
}
