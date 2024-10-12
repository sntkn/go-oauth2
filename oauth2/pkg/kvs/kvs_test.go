package kvs

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// テスト用の kvs サーバーをセットアップする関数
func setupTest(t *testing.T) *KVSClient {

	cli, err := NewClient(context.Background(), Options{
		Addr: []string{"127.0.0.1:6379"},
	})
	require.NoError(t, err)
	fmt.Printf("%+v", cli)

	return cli
}

func TestKVSCli(t *testing.T) {
	t.Parallel()
	cli := setupTest(t)

	ctx := context.Background()

	// Test Set and Get
	err := cli.Set(ctx, "testKey", "testValue", 10)
	require.NoError(t, err)

	val, err := cli.Get(ctx, "testKey")
	require.NoError(t, err)
	assert.Equal(t, "testValue", val)

	// Test Del
	err = cli.Del(ctx, "testKey")
	require.NoError(t, err)

	val, err = cli.Get(ctx, "testKey")
	require.NoError(t, err)
	assert.Equal(t, "", val)
}
