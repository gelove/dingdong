package service

import (
	"log"
	"math/rand"
	"time"

	"dingdong/internal/app/dto/reserve_time"
)

type task struct {
	Completed     bool
	reserveTime   *reserve_time.GoTimes
	cartMap       map[string]interface{}
	checkOrderMap map[string]interface{}
}

func NewTask() *task {
	return new(task)
}

func (t *task) ReserveTime() *reserve_time.GoTimes {
	return t.reserveTime
}

func (t *task) SetReserveTime(reserveTime *reserve_time.GoTimes) *task {
	if reserveTime != nil {
		t.reserveTime = reserveTime
	}
	return t
}

func (t *task) CartMap() map[string]interface{} {
	return t.cartMap
}

func (t *task) SetCartMap(cartMap map[string]interface{}) *task {
	if cartMap != nil {
		t.cartMap = cartMap
	}
	return t
}

func (t *task) CheckOrderMap() map[string]interface{} {
	return t.checkOrderMap
}

func (t *task) SetCheckOrderMap(checkOrderMap map[string]interface{}) *task {
	if checkOrderMap != nil {
		t.checkOrderMap = checkOrderMap
	}
	return t
}

// AllCheck 不一定需要, 只起补充作用
func (t *task) AllCheck() {
	for {
		if t.Completed {
			return
		}
		err := AllCheck()
		if err != nil {
			log.Println(err)
		}
		duration := 3000 + rand.Intn(2000)
		<-time.After(time.Duration(duration) * time.Millisecond)
	}
}

func (t *task) GetCart() {
	for {
		if t.Completed {
			return
		}
		log.Println("===== 获取有效的商品 =====")
		cartMap, err := GetCart()
		if err != nil {
			log.Println(err)
		} else {
			t.SetCartMap(cartMap)
		}
		duration := 150 + rand.Intn(50)
		<-time.After(time.Duration(duration) * time.Millisecond)
	}
}

func (t *task) GetMultiReserveTime() {
	for {
		if t.Completed {
			return
		}
		if t.CartMap() == nil {
			<-time.After(20 * time.Millisecond)
			continue
		}
		log.Println("===== 获取有效的配送时段 =====")
		reserveTime, err := GetMultiReserveTime(t.CartMap())
		if err != nil {
			log.Println(err)
		} else {
			t.SetReserveTime(reserveTime)
		}
		duration := 150 + rand.Intn(50)
		<-time.After(time.Duration(duration) * time.Millisecond)
	}
}

func (t *task) GetCheckOrder() {
	for {
		if t.Completed {
			return
		}
		if t.CartMap() == nil || t.ReserveTime() == nil {
			<-time.After(20 * time.Millisecond)
			continue
		}
		log.Println("===== 生成订单信息 =====")
		checkOrderMap, err := GetCheckOrder(t.CartMap(), t.ReserveTime())
		if err != nil {
			log.Println(err)
		} else {
			t.SetCheckOrderMap(checkOrderMap)
		}
		duration := 150 + rand.Intn(50)
		<-time.After(time.Duration(duration) * time.Millisecond)
	}
}

func (t *task) AddNewOrder() {
	for {
		if t.Completed {
			return
		}
		if t.CartMap() == nil || t.ReserveTime() == nil || t.CheckOrderMap() == nil {
			<-time.After(20 * time.Millisecond)
			continue
		}
		log.Println("===== 提交订单 =====")
		err := AddNewOrder(t.CartMap(), t.ReserveTime(), t.CheckOrderMap())
		if err != nil {
			log.Println(err)
			duration := 150 + rand.Intn(50)
			<-time.After(time.Duration(duration) * time.Millisecond)
			continue
		}
		t.Completed = true
	}
}

func Run() {
	t := NewTask()

	go t.AllCheck()

	go t.GetCart()

	go t.GetMultiReserveTime()

	go t.GetCheckOrder()

	go t.AddNewOrder()

	timer := time.NewTimer(time.Minute * 3)
	<-timer.C
}
