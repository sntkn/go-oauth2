package logs

import (
	"fmt"
	"log/slog"

	"github.com/cockroachdb/errors"
)

func Error(err error) {
	slog.Error(fmt.Sprintf("%+\v\n", err))
}

func ErrorWithWrap(err error, msg string) {
	slog.Error(fmt.Sprintf("%+\v\n", errors.Wrap(err, msg)))
}

func ErrorWithStack(err error) {
	slog.Error(fmt.Sprintf("%+\v\n", errors.WithStack(err)))
}
