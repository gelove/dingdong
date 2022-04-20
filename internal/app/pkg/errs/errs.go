package errs

import (
	"errors"
	"fmt"
	"io"

	"dingdong/internal/app/pkg/errs/code"
)

// Error 错误
type Error interface {
	error
	fmt.Formatter
	Unwrap() error
	Code() code.ErrorCode
	Message() string
	CodeEqual(code.ErrorCode) bool
}

// ErrorImpl 错误实现
type ErrorImpl struct {
	error
	code code.ErrorCode
}

// Unwrap 获取内部错误
func (e ErrorImpl) Unwrap() error {
	return e.error
}

func (e ErrorImpl) Code() code.ErrorCode {
	return e.code
}

func (e ErrorImpl) Error() string {
	if e.error == nil {
		return e.code.String()
	}
	return e.code.String() + " " + e.error.Error()
}

func (e ErrorImpl) Message() string {
	// 不要注释这里的代码, 运行前请先执行 make generate
	return e.code.String()
}

func (e ErrorImpl) CodeEqual(code code.ErrorCode) bool {
	return e.code == code
}

// Format 格式化打印
func (e ErrorImpl) Format(f fmt.State, verb rune) {
	if v, ok := e.error.(fmt.Formatter); ok {
		v.Format(f, verb)
	} else {
		_, _ = io.WriteString(f, e.Error())
	}
}

var (
	Is     = errors.Is
	As     = errors.As
	Unwrap = errors.Unwrap
)

// New 新建一个 Error
func New(code code.ErrorCode) Error {
	return &ErrorImpl{
		code: code,
	}
}

// Wrap 包装 error
func Wrap(code code.ErrorCode, err error) Error {
	if err == nil {
		return nil
	}
	return &ErrorImpl{
		error: err,
		code:  code,
	}
}

// WithMessage 包装 error
func WithMessage(code code.ErrorCode, msg string) Error {
	if msg == "" {
		return New(code)
	}
	return &ErrorImpl{
		error: errors.New(msg),
		code:  code,
	}
}
