package errors

import (
	"fmt"
	"log/slog"

	go_errors "github.com/go-errors/errors"
)

// go_errorsのメソッドを全て公開
var (
	Is = go_errors.Is
	As = go_errors.As
	//	Join = go_errors.Join
	//	Unwrap = go_errors.Unwrap
	New        = go_errors.New
	ParsePanic = go_errors.ParsePanic
	Wrap       = go_errors.Wrap
	WrapPrefix = go_errors.WrapPrefix
	Errorf     = go_errors.Errorf
)

func WithStack(err error) error {
	if err != nil {
		return go_errors.Wrap(err, 0)
	}
	return nil
}

func LogStackTrace(err error) slog.Attr {
	if err == nil {
		return slog.Any("stacktrace", []any{})
	}

	// go-errors/errors の Error かどうか
	goerror, ok := err.(*go_errors.Error)
	if !ok {
		return slog.String("details", fmt.Sprintf("%+v", err))
	}
	return slog.Any("stacktrace", goerror.StackFrames())
}
