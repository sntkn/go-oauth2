package valkey

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/valkey-io/valkey-go"
)

//go:generate go run github.com/matryer/moq -out valkey_mock.go . ClientIF
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

func NewClient(_ context.Context, o Options) (*Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: o.Addr,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	//	defer client.Close()

	return &Client{
		cli: client,
	}, nil
}

func (c *Client) Set(ctx context.Context, key string, value string, expiration int64) error {
	if err := c.cli.Do(ctx, c.cli.B().Set().Key(key).Value(value).Build()).Error(); err != nil {
		return errors.WithStack(err)
	}
	if err := c.cli.Do(ctx, c.cli.B().Expire().Key(key).Seconds(expiration).Build()).Error(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	str, err := c.cli.Do(ctx, c.cli.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if errors.Is(err, valkey.Nil) {
			return "", nil
		}
		return "", errors.WithStack(err)
	}
	return str, nil
}

func (c *Client) Del(ctx context.Context, key string) error {
	if err := c.cli.Do(ctx, c.cli.B().Del().Key(key).Build()).Error(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
