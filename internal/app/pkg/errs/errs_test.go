package errs

import (
	"testing"

	"dingdong/internal/app/pkg/errs/code"
)

func TestWithCode(t *testing.T) {
	// 逻辑层返回原错误和自定义错误

	// 控制器层需要错误码响应时包装原错误(包含堆栈), 系统error和第三方库error的错误码统一为 code.InternalError
	err := WithCode(New("test"), code.InternalError)
	if !Is(err, InternalError) {
		t.Error("expected not to be code.InternalError")
		return
	}
	t.Log(err)
	t.Logf(err.Error())
	t.Logf("%+v", err)

	// 主要用来自定义错误(包含堆栈)
	err1 := WithStack(InternalError)
	if !Is(err1, InternalError) {
		t.Error("expected not to be code.InternalError")
		return
	}
	t.Log(err1)
	t.Logf(err1.Error())
	t.Logf("%+v", err1)

	err2 := Wrap(InternalError, "test")
	if !Is(err2, InternalError) {
		t.Error("expected not to be code.InternalError")
		return
	}
	t.Log(err2)
	t.Logf(err2.Error())
	t.Logf("%+v", err2)
}
