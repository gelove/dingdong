package service

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"dingdong/internal/app/config"
	"dingdong/internal/app/dto/reserve_time"
	"dingdong/internal/app/pkg/date"
	"dingdong/pkg/json"
)

var task *Task

const (
	FirstSnapUp uint8 = iota + 1
	SecondSnapUp
)

type Task struct {
	notifyCh      chan struct{}
	pickUpCh      chan struct{}
	Completed     bool
	reserveTime   *reserve_time.GoTimes
	cartMap       map[string]interface{}
	checkOrderMap map[string]interface{}
}

func init() {
	task = NewTask()
}

func NewTask() *Task {
	return &Task{
		notifyCh: make(chan struct{}),
		pickUpCh: make(chan struct{}),
	}
}

func (t *Task) ReserveTime() *reserve_time.GoTimes {
	return t.reserveTime
}

func (t *Task) SetReserveTime(reserveTime *reserve_time.GoTimes) *Task {
	if reserveTime != nil {
		t.reserveTime = reserveTime
	}
	return t
}

func (t *Task) CartMap() map[string]interface{} {
	return t.cartMap
}

func (t *Task) SetCartMap(cartMap map[string]interface{}) *Task {
	if cartMap != nil {
		t.cartMap = cartMap
	}
	return t
}

func (t *Task) CheckOrderMap() map[string]interface{} {
	return t.checkOrderMap
}

func (t *Task) SetCheckOrderMap(checkOrderMap map[string]interface{}) *Task {
	if checkOrderMap != nil {
		t.checkOrderMap = checkOrderMap
	}
	return t
}

// AllCheck 不一定需要, 只起补充作用
func (t *Task) AllCheck() {
	for {
		if t.Completed {
			return
		}
		err := AllCheck()
		if err != nil {
			log.Println(err)
		} else {
			return
		}
		duration := 3000 + rand.Intn(2000)
		<-time.After(time.Duration(duration) * time.Millisecond)
	}
}

func (t *Task) GetCart() {
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
			return
		}
		duration := 50 + rand.Intn(50)
		<-time.After(time.Duration(duration) * time.Millisecond)
	}
}

func (t *Task) GetMultiReserveTime() {
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
			log.Println("reserveTime => ", json.MustEncodeToString(reserveTime))
			return
		}
		duration := 50 + rand.Intn(50)
		<-time.After(time.Duration(duration) * time.Millisecond)
	}
}

func (t *Task) CheckOrder() {
	for {
		if t.Completed {
			return
		}
		if t.CartMap() == nil || t.ReserveTime() == nil {
			<-time.After(20 * time.Millisecond)
			continue
		}
		log.Println("===== 生成订单信息 =====")
		checkOrderMap, err := CheckOrder(t.CartMap(), t.ReserveTime())
		if err != nil {
			log.Println(err)
		} else {
			t.SetCheckOrderMap(checkOrderMap)
			return
		}
		duration := 50 + rand.Intn(50)
		<-time.After(time.Duration(duration) * time.Millisecond)
	}
}

func (t *Task) AddNewOrder() {
	for {
		if t.Completed {
			return
		}
		if t.CartMap() == nil || t.ReserveTime() == nil || t.CheckOrderMap() == nil {
			<-time.After(5 * time.Millisecond)
			continue
		}
		log.Println("===== 提交订单 =====")
		err := AddNewOrder(t.CartMap(), t.ReserveTime(), t.CheckOrderMap())
		if err != nil {
			log.Println(err)
			<-time.After(5 * time.Millisecond)
			continue
		}
		t.Completed = true
		conf := config.Get()
		Push(conf.Users[0], "已成功下单, 请尽快完成支付")
		return
	}
}

func timeTrigger() bool {
	conf := config.Get()
	now := time.Now()
	firstTime := date.FirstSnapUpTime()
	// log.Println(conf.SnapUp&FirstSnapUp == FirstSnapUp, now, firstTime, now.Unix(), firstTime.Unix())
	if conf.SnapUp&FirstSnapUp == FirstSnapUp && now.Unix() == firstTime.Unix()-3 {
		log.Println("===== 6点抢购开始 =====")
		return true
	}
	secondTime := date.SecondSnapUpTime()
	// log.Println(conf.SnapUp&SecondSnapUp == SecondSnapUp, now, secondTime, now.Unix(), secondTime.Unix())
	if conf.SnapUp&SecondSnapUp == SecondSnapUp && now.Unix() == secondTime.Unix()-3 {
		log.Println("===== 8点半抢购开始 =====")
		return true
	}
	return false
}

func SnapUpOnce() {
	go task.AllCheck()

	go task.GetCart()

	go task.GetMultiReserveTime()

	go task.CheckOrder()

	go task.AddNewOrder()

	timer := time.NewTimer(time.Minute * 3)
	<-timer.C
}

func SnapUp() {
	timer := time.Tick(time.Second)
	for {
		select {
		case <-timer:
			if timeTrigger() {
				go SnapUpOnce()
			}
		}
	}
}

func GetMultiReserveTimeAndNotify(cartMap map[string]interface{}) {
	_, err := GetMultiReserveTime(cartMap)
	if err != nil {
		log.Println(err)
		return
	}

	conf := config.Get()
	if conf.PickUpNeeded {
		task.pickUpCh <- struct{}{}
	}

	if conf.NotifyNeeded {
		task.notifyCh <- struct{}{}
	}
}

func PickUp() {
	for {
		<-task.pickUpCh
		SnapUpOnce()
	}
}

func Notify() {
	for {
		<-task.notifyCh
		now := time.Now()
		conf := config.Get()
		interval := time.Duration(conf.NotifyInterval) * time.Minute
		if now.Before(lastNotify.Add(interval)) {
			continue
		}
		lastNotify = now
		list, err := GetHomeFlowDetail()
		if err != nil {
			log.Printf("获取首页产品失败 => %+v", err)
		}
		productNames := make([]string, 0, 10)
		for i, item := range list {
			if i >= 10 {
				continue
			}
			letter := []rune(item.Name)
			if len(letter) > 5 {
				productNames = append(productNames, string(letter[:5]))
				continue
			}
			productNames = append(productNames, string(letter))
		}
		ellipsis := ""
		if len(list) >= 10 {
			ellipsis = "..."
		}
		products := strings.Join(productNames, " ")
		wg := new(sync.WaitGroup)
		for k, v := range conf.Users {
			if k > 0 {
				continue
			}
			wg.Add(1)
			go func(key string) {
				defer wg.Done()
				Push(key, fmt.Sprintf("叮咚买菜当前可配送请尽快下单[%s%s]", products, ellipsis))
			}(v)
		}
		wg.Wait()
	}
}
