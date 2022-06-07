package test

import (
	"testing"
	"time"

	"dingdong/internal/app/config"
	_ "dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/ios_session"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/service"
	"dingdong/internal/app/service/ios_service"
	"dingdong/internal/app/service/meituan"
	"dingdong/pkg/js"
	"dingdong/pkg/json"
)

func TestSign(t *testing.T) {
	headers := session.GetHeaders()
	headers["accept-language"] = "en-us"
	params := session.GetParams(headers)
	params["tab_type"] = "1"
	params["page"] = "1"
	res, err := js.Call("js/sign.js", "sign", json.MustEncodeToString(params))
	if err != nil {
		t.Error("js parser error =>", err)
		return
	}
	// nars 对应就可以
	// sesi 可以不用管, 依赖 nars 与 一个随机字符串, 每次计算应都不同, 但是在JS虚拟机中伪随机数似乎不变, 每次都会得到同一个值
	t.Log("value =>", res.String())
}

func TestIosSign(t *testing.T) {
	// 注意修改 time 和 ImSecret
	params := ios_session.TakeParams(1651583332, nil)
	t.Logf("params => %#v", params)
	str := json.MustEncodeToString(params)
	t.Logf("str => %s", str)
	conf := config.GetDingDong()
	t.Logf("params: %s, %#v", conf.ImSecret, str)
	res, err := js.Call("js/ios_sign.js", "ios_sign", conf.ImSecret, str)
	if err != nil {
		t.Error("js parser error =>", err)
		return
	}
	t.Log("value =>", res.String())
}

func TestMeiTuanReserveTime(t *testing.T) {
	res, err := meituan.GetMultiReserveTime()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(json.MustEncodePrettyString(res))
}

func TestGetHomeFlowDetail(t *testing.T) {
	list, err := service.GetHomeFlowDetail()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(json.MustEncodePrettyString(list))
}

func TestGetUser(t *testing.T) {
	user, err := session.GetUser()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(user)
}

func TestGetAddress(t *testing.T) {
	list, err := session.GetAddress()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(list)
}

func TestAllCheck(t *testing.T) {
	err := ios_service.AllCheck()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("All check success")
}

func TestGetCart(t *testing.T) {
	cartMap, err := ios_service.GetCart()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s", json.MustEncodeToString(cartMap))
	out := make(map[string]any)
	json.MustTransform(cartMap, &out)
	t.Logf("%#v", out)
	// t.Logf("%#v", out["products"].([]map[string]any)) // panic
	t.Logf("%#v", out["products"].([]any))
}

func TestGetMultiReserveTime(t *testing.T) {
	cartMap := ios_service.MockCartMap()
	now := time.Now()
	reserveTimes, err := ios_service.GetMultiReserveTime(cartMap)
	t.Logf("Millisecond => %d ms", time.Now().Sub(now).Milliseconds())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(json.MustEncodePrettyString(reserveTimes))
}

// TestMockMultiReserveTime 模拟运力数据
func TestMockMultiReserveTime(t *testing.T) {
	reserveTimes := ios_service.MockMultiReserveTime()
	t.Log(json.MustEncodePrettyString(reserveTimes))
}

func TestCheckOrder(t *testing.T) {
	err := ios_service.AllCheck()
	if err != nil {
		t.Error(err)
		return
	}
	cartMap, err := ios_service.GetCart()
	if err != nil {
		t.Error(err)
		return
	}
	_, err = ios_service.GetMultiReserveTime(cartMap)
	if err != nil {
		t.Error(err)
		return
	}
	orderMap, err := ios_service.CheckOrder(cartMap)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", orderMap)
}

func TestAddNewOrder(t *testing.T) {
	err := service.AddOrder()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Create order success")
}

func TestSnapUpOnce(t *testing.T) {
	service.SnapUpOnce(service.PickUpMode)
}
