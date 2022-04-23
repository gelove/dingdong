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
)

const (
	FirstSnapUp uint8 = iota + 1
	SecondSnapUp
)

const (
	durationMinMillis = 450
	durationMaxMillis = 550
	durationGapMillis = durationMaxMillis - durationMinMillis
)

var task *Task

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
			duration := durationMinMillis + rand.Intn(durationGapMillis)
			<-time.After(time.Duration(duration) * time.Millisecond)
		}
		return
	}
}

func (t *Task) GetCart() {
	for {
		if t.Completed {
			return
		}
		cartMap, err := GetCart()
		if err != nil {
			log.Println(err)
		} else {
			t.SetCartMap(cartMap)
			log.Println("===== 购物车商品已更新 =====")
		}
		duration := durationMinMillis + rand.Intn(durationGapMillis)
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
		reserveTime, err := GetMultiReserveTime(t.CartMap())
		if err != nil {
			log.Println(err)
		} else {
			t.SetReserveTime(reserveTime)
			log.Println("===== 有效配送时段已更新 =====")
			// log.Println("reserveTime => ", json.MustEncodeToString(reserveTime))
		}
		duration := durationMinMillis + rand.Intn(durationGapMillis)
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
		checkOrderMap, err := CheckOrder(t.CartMap(), t.ReserveTime())
		if err != nil {
			log.Println(err)
		} else {
			t.SetCheckOrderMap(checkOrderMap)
			log.Println("===== 订单信息已更新 =====")
		}
		duration := durationMinMillis + rand.Intn(durationGapMillis)
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
		err := AddNewOrder(t.CartMap(), t.ReserveTime(), t.CheckOrderMap())
		if err != nil {
			log.Println(err)
			duration := 100 + rand.Intn(50)
			<-time.After(time.Duration(duration) * time.Millisecond)
			continue
		}
		t.Completed = true
		detail := "已成功下单, 请尽快完成支付"
		log.Println(detail)
		conf := config.Get()
		Push(conf.Users[0], detail)
		return
	}
}

func timeTrigger() bool {
	conf := config.Get()
	now := time.Now()
	firstTime := date.FirstSnapUpUnix()
	// log.Println(conf.SnapUp&FirstSnapUp == FirstSnapUp, now, firstTime, now.Unix(), firstTime)
	if conf.SnapUp&FirstSnapUp == FirstSnapUp && now.Unix() == firstTime-20 {
		log.Println("===== 6点抢购开始 =====")
		return true
	}
	secondTime := date.SecondSnapUpUnix()
	// log.Println(conf.SnapUp&SecondSnapUp == SecondSnapUp, now, secondTime, now.Unix(), secondTime)
	if conf.SnapUp&SecondSnapUp == SecondSnapUp && now.Unix() == secondTime-20 {
		log.Println("===== 8点半抢购开始 =====")
		return true
	}
	return false
}

func SnapUpOnce() {
	conf := config.Get()

	for i := 0; i < conf.BaseConcurrency; i++ {
		go task.AllCheck()

		go task.GetCart()

		go task.GetMultiReserveTime()

		go task.CheckOrder()
	}

	for i := 0; i < conf.SubmitConcurrency; i++ {
		go task.AddNewOrder()
	}

	timer := time.NewTimer(time.Minute * 2)
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
