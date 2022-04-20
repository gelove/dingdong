package test

import (
	"log"
	"testing"

	"dingdong/internal/app/config"
	_ "dingdong/internal/app/config"
	"dingdong/internal/app/service"
	"dingdong/pkg/json"
)

func TestPush(t *testing.T) {
	conf := config.Get()
	service.Push(conf.Users[0], "测试", "")
}

func TestGetHomeFlowDetail(t *testing.T) {
	list, err := service.GetHomeFlowDetail()
	if err != nil {
		t.Error(err)
	}
	t.Log(list)
}

func TestGetMultiReserveTime(t *testing.T) {
	cartMap := service.MockCartMap()
	_, err := service.GetMultiReserveTime(cartMap)
	if err != nil {
		t.Error(err)
	}
}

func TestGetAddress(t *testing.T) {
	list, err := service.GetAddress()
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
	m, err := service.GetCart()
	if err != nil {
		t.Error(err)
	}
	t.Log(json.MustEncodeToString(m))
}

func TestRun(t *testing.T) {
	service.Run()
}

// TestRunOnce 此为单次执行模式 用于在非高峰期测试下单 也必须满足3个前提条件 1.有收货地址 2.购物车有商品 3.有配送时间段
func TestRunOnce(t *testing.T) {
	conf := config.Get()
	addressId := conf.Params["address_id"]
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
	checkOrderMap, err := service.GetCheckOrder(cartMap, reserveTime)
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
