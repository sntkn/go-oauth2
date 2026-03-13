package repository

import (
	"context"
	"database/sql"

	"github.com/sntkn/go-oauth2/oauth2/pkg/errors"
)

type rowGetter interface {
	GetContext(ctx context.Context, dest any, query string, args ...any) error
}

func fetchAndMap[M any, D any](ctx context.Context, getter rowGetter, query string, mapper func(M) (D, error), args ...any) (D, bool, error) {
	var model M
	var zero D

	if err := getter.GetContext(ctx, &model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return zero, false, nil
		}
		return zero, false, errors.WithStack(err)
	}

	entity, err := mapper(model)
	if err != nil {
		return zero, false, err
	}
	return entity, true, nil
}
