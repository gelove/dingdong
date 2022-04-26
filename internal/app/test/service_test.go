package test

import (
	"log"
	"testing"

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
	res, err := js.Call(jsFile, "sign", json.MustEncodeToString(params))
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
	t.Log(list)
}

func TestGetAddress(t *testing.T) {
	list, err := session.GetAddress()
	if err != nil {
		t.Error(err)
	}
	t.Log(list)
}

func TestGetUser(t *testing.T) {
	user, err := session.GetUser()
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
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
	_, err := service.GetMultiReserveTime(cartMap)
	if err != nil {
		t.Error(err)
	}
}

func TestMockMultiReserveTime(t *testing.T) {
	task := service.NewTask()
	task.MockMultiReserveTime()
	t.Log(task.ReserveTime())
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

// TestRunOnce 此为单次执行模式 用于在非高峰期测试下单 也必须满足3个前提条件 1.有收货地址 2.购物车有商品 3.有配送时间段
func TestRunOnce(t *testing.T) {
	addressId := session.Address().Id
	if addressId == "" {
		t.Error("address_id is empty")
		return
	}

	log.Println("===== 获取有效的商品 =====")
	err := service.AllCheck()
	if err != nil {
		t.Error(err)
		return
	}

	cartMap, err := service.GetCart()
	if err != nil {
		t.Error(err)
		return
	}
	if len(cartMap) == 0 {
		t.Error("cart is empty")
		return
	}

	log.Println("===== 获取有效的配送时段 =====")
	reserveTime, err := service.GetMultiReserveTime(cartMap)
	if err != nil {
		t.Error(err)
		return
	}
	if reserveTime == nil {
		t.Error("reserveTime is empty")
		return
	}

	log.Println("===== 生成订单信息 =====")
	checkOrderMap, err := service.CheckOrder(cartMap, reserveTime)
	if err != nil {
		t.Error(err)
		return
	}
	if len(checkOrderMap) == 0 {
		t.Error("checkOrderMap is empty")
		return
	}
	log.Println("订单总金额 =>", checkOrderMap["price"])

	log.Println("===== 提交订单 =====")
	err = service.AddNewOrder(cartMap, reserveTime, checkOrderMap)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("提交订单成功")
}
