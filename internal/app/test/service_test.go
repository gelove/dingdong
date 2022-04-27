package test

import (
	"testing"
	"time"

	_ "dingdong/internal/app/config"
	"dingdong/internal/app/dto/reserve_time"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/service"
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

func TestGetHomeFlowDetail(t *testing.T) {
	list, err := service.GetHomeFlowDetail()
	if err != nil {
		t.Error(err)
	}
	t.Log(json.MustEncodePrettyString(list))
}

func TestGetUser(t *testing.T) {
	user, err := session.GetUser()
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
}

func TestGetAddress(t *testing.T) {
	list, err := session.GetAddress()
	if err != nil {
		t.Error(err)
	}
	t.Log(list)
}

func TestAllCheck(t *testing.T) {
	err := service.AllCheck()
	if err != nil {
		t.Error(err)
	}
}

func TestGetCart(t *testing.T) {
	cartMap, err := service.GetCart()
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", cartMap)
}

func TestGetMultiReserveTime(t *testing.T) {
	cartMap := service.MockCartMap()
	now := time.Now()
	_, err := service.GetMultiReserveTime(cartMap)
	t.Logf("Millisecond => %d ms", time.Now().Sub(now).Milliseconds())
	if err != nil {
		t.Error(err)
	}
}

// TestMockMultiReserveTime 模拟运力数据
func TestMockMultiReserveTime(t *testing.T) {
	task := service.NewTask()
	defer task.Finished()
	task.MockMultiReserveTime()
	times := task.ReserveTime()
	t.Log(json.MustEncodePrettyString(times))
}

func TestCheckOrder(t *testing.T) {
	reserveTimes := &reserve_time.GoTimes{}
	cartMap, err := service.GetCart()
	if err != nil {
		t.Error(err)
	}
	orderMap, err := service.CheckOrder(cartMap, reserveTimes)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v", orderMap)
}

func TestAddNewOrder(t *testing.T) {
	cartMap, err := service.GetCart()
	if err != nil {
		t.Error(err)
	}
	reserveTimes, err := service.GetMultiReserveTime(cartMap)
	if err != nil {
		t.Error(err)
	}
	orderMap, err := service.CheckOrder(cartMap, reserveTimes)
	if err != nil {
		t.Error(err)
	}
	err = service.AddNewOrder(cartMap, reserveTimes, orderMap)
	if err != nil {
		t.Error(err)
	}
}

func TestSnapUpOnce(t *testing.T) {
	service.SnapUpOnce()
}
