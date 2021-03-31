
package src

import (
	"github.com/pkg/errors"
	"os"
)

// 一个自定义的errors.Wrapper
type MyError struct {
	msg    string
	deeper error
}
//errors.
//var _ errors.Wrapper = &MyError{}

func (me *MyError) Error() string {
	return "MyError: " + me.msg
}

func (me *MyError) Unwrap() error {
	return me.deeper
}

func Foo() *MyError {
	if err := Bar(); err != nil  {
		return &MyError {
			msg:    "Foo failed",
			deeper: &MyError{
				msg:    "Bar failed",
				deeper: os.ErrNotExist,
			},
		}
	}
	return nil
}

func Bar() error {
	return errors.New("unknown error")
}

