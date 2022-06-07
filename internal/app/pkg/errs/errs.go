package errs

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"dingdong/internal/app/pkg/errs/code"
)

var (
	Is          = errors.Is
	As          = errors.As
	Unwrap      = errors.Unwrap
	New         = errors.New
	Errorf      = errors.Errorf
	Wrap        = errors.Wrap
	Wrapf       = errors.Wrapf
	WithStack   = errors.WithStack
	WithMessage = errors.WithMessage
)

var (
	OK                    = NewCode(code.OK)                    // OK
	InternalError         = NewCode(code.InternalError)         // 内部错误
	SelectSessionFailed   = NewCode(code.SelectSessionFailed)   // 选择session文件错误
	GetUserDetailFailed   = NewCode(code.GetUserDetailFailed)   // 获取用户信息失败
	GetAddressFailed      = NewCode(code.GetAddressFailed)      // 获取收货地址失败
	GetFlowDetailFailed   = NewCode(code.GetFlowDetailFailed)   // 获取首页流水详情失败
	NoValidAddress        = NewCode(code.NoValidAddress)        // 当前没有可用的收货地址
	CheckAllFailed        = NewCode(code.CheckAllFailed)        // 购物车全选失败
	GetCartFailed         = NewCode(code.GetCartFailed)         // 获取购物车失败
	NoValidProduct        = NewCode(code.NoValidProduct)        // 当前购物车中没有可购商品
	GetReserveTimeFailed  = NewCode(code.GetReserveTimeFailed)  // 获取运力失败
	NoReserveTime         = NewCode(code.NoReserveTime)         // 当前没有可用的运力
	NoReserveTimeAndRetry = NewCode(code.NoReserveTimeAndRetry) // 当前没有可用的运力, 请稍后再试
	ReserveTimeIsDisabled = NewCode(code.ReserveTimeIsDisabled) // 您选择的送达时间已经失效, 请重新选择
	CheckOrderFailed      = NewCode(code.CheckOrderFailed)      // 订单校验失败
	SubmitOrderFailed     = NewCode(code.SubmitOrderFailed)     // 提交订单失败
	NotifyFailed          = NewCode(code.NotifyFailed)          // 通知失败
)

type IUnwrap interface {
	Unwrap() error
}

type ICode interface {
	Code() code.ErrorCode
}

type ICause interface {
	Cause() error
}

type Error interface {
	error
	ICode
}

type Code struct {
	code code.ErrorCode
}

func (e *Code) Error() string {
	return e.code.String()
}

func (e *Code) Code() code.ErrorCode {
	return e.code
}

func NewCode(code code.ErrorCode) Error {
	return &Code{code}
}

type withCode struct {
	error
	code code.ErrorCode
}

func (e *withCode) Error() string {
	// 不要注释这里的代码, 运行前请先执行 make generate
	return e.code.String()
}

func (e *withCode) Message() string {
	// 不要注释这里的代码, 运行前请先执行 make generate
	return e.code.String() + ": " + e.error.Error()
}

func (e *withCode) Code() code.ErrorCode {
	return e.code
}

func (e *withCode) Is(err error) bool {
	return Cause(err).Code() == e.Code()
}

func (e *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// _, _ = io.WriteString(s, e.code.String()+"\n")
			_, _ = fmt.Fprintf(s, "%s\n", e.code)
			_, _ = fmt.Fprintf(s, "%+v\n", e.error)
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(s, e.Error())
	}
}

// WithCode 新建一个 Error
func WithCode(err error, codes ...code.ErrorCode) Error {
	if err == nil {
		return nil
	}
	if len(codes) < 1 {
		return &withCode{err, code.InternalError}
	}
	return &withCode{
		error: err,
		code:  codes[0],
	}
}

func Cause(err error) Error {
	if err == nil {
		return OK
	}
	if ec, ok := err.(Error); ok {
		return ec
	}
	if ec, ok := errors.Cause(err).(Error); ok {
		return ec
	}
	return WithCode(err, code.InternalError)
}

func PickError(err error) Error {
	for err != nil {
		if v, ok := err.(Error); ok {
			return v
		}
		u, ok := err.(IUnwrap)
		if !ok {
			return nil
		}
		err = u.Unwrap()
	}
	return nil
}
