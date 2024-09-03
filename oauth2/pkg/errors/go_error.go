package errors

import go_errors "github.com/go-errors/errors"

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
	return go_errors.Wrap(err, 0)
}
