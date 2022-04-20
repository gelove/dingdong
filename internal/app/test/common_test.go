package test

import (
	"testing"

	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/date"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/pkg/js"
	"dingdong/pkg/json"
)

func TestJsCall(t *testing.T) {
	headers := session.GetHeaders()
	headers["accept-language"] = "en-us"
	params := session.GetParams(headers)
	params["tab_type"] = "1"
	params["page"] = "1"
	res, err := js.Call(jsFile, "sign", json.MustEncodeToString(params))
	if err != nil {
		t.Error("js parser error =>", err)
		return
	}
	// nars 对应就可以
	// sesi 可以不用管, 依赖 nars 与 一个随机字符串, 每次计算应都不同, 但是在JS虚拟机中伪随机数似乎不变, 每次都会得到同一个值
	t.Log("value =>", res.String())
}

func TestJsonGet(t *testing.T) {
	conf := config.Get()
	bs := json.MustEncode(conf)
	t.Log(json.Get(bs, "headers", "cookie").ToString())
	t.Log(json.Get(bs, "headers").Get("cookie").ToString())
}

func TestSnapUpTime(t *testing.T) {
	t.Log(date.FirstSnapUpTime())
	t.Log(date.SecondSnapUpTime())
}
