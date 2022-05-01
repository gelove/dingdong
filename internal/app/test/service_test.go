package test

import (
	"testing"
	"time"

	_ "dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/service"
	"dingdong/internal/app/service/meituan"
	"dingdong/pkg/js"
	"dingdong/pkg/json"
)

func TestJsCall(t *testing.T) {
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
	err := service.AllCheck()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("All check success")
}

func TestGetCart(t *testing.T) {
	cartMap, err := service.GetCart()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", cartMap)
}

func TestGetMultiReserveTime(t *testing.T) {
	cartMap := service.MockCartMap()
	now := time.Now()
	reserveTimes, err := service.GetMultiReserveTime(cartMap)
	t.Logf("Millisecond => %d ms", time.Now().Sub(now).Milliseconds())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(json.MustEncodePrettyString(reserveTimes))
}

// TestMockMultiReserveTime 模拟运力数据
func TestMockMultiReserveTime(t *testing.T) {
	reserveTimes := service.MockMultiReserveTime()
	t.Log(json.MustEncodePrettyString(reserveTimes))
}

func TestCheckOrder(t *testing.T) {
	cartMap, err := service.GetCart()
	if err != nil {
		t.Error(err)
		return
	}
	reserveTimes, err := service.GetMultiReserveTime(cartMap)
	if err != nil {
		t.Error(err)
		return
	}
	orderMap, err := service.CheckOrder(cartMap, reserveTimes)
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
