package valkey

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/valkey-io/valkey-go"
)

//go:generate go run github.com/matryer/moq -out redis_mock.go . ClientIF
type ClientIF interface {
	Set(ctx context.Context, key string, value string, expiration int64) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type Options struct {
	Addr []string
}

type Client struct {
	cli valkey.Client
}

func NewClient(ctx context.Context, o Options) (*Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: o.Addr,
	})
	if err != nil {
		return nil, err
	}
	//	defer client.Close()

	return &Client{
		cli: client,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, value string, expiration int64) error {
	err := c.cli.Do(ctx, c.cli.B().Set().Key(key).Value(value).Nx().Build()).Error()
	if err != nil {
		return err
	}
	return c.cli.Do(ctx, c.cli.B().Expire().Key(key).Seconds(expiration).Build()).Error()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	str, err := c.cli.Do(ctx, c.cli.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return "", nil
		}
		return "", err
	}
	return str, nil
}

func (c *Client) Del(ctx context.Context, key string) error {
	return c.cli.Do(ctx, c.cli.B().Del().Key(key).Build()).Error()
}
