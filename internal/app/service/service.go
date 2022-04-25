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
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
)

const (
	FirstSnapUp = iota + 1
	SecondSnapUp
)

const (
	durationMinMillis = 450
	durationMaxMillis = 550
	durationGapMillis = durationMaxMillis - durationMinMillis
)

var (
	notifyCh = make(chan struct{})
	pickUpCh = make(chan struct{})
)

type Task struct {
	Completed     bool
	reserveTime   *reserve_time.GoTimes
	cartMap       map[string]interface{}
	checkOrderMap map[string]interface{}
}

func NewTask() *Task {
	return &Task{}
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
		duration := durationMinMillis + rand.Intn(durationGapMillis)
		cartMap, err := GetCart()
		if err != nil {
			log.Println(err)
			if e, ok := err.(errs.Error); ok {
				if e.Code() == code.NoValidProduct {
					return
				}
			}
			<-time.After(time.Duration(duration) * time.Millisecond)
			continue
		}
		t.SetCartMap(cartMap)
		<-time.After(time.Duration(duration) * time.Millisecond)
		continue
	}
}

func (t *Task) GetMultiReserveTime() {
	for {
		if t.Completed {
			return
		}
		if t.CartMap() == nil {
			<-time.After(10 * time.Millisecond)
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

func (t *Task) MockMultiReserveTime() {
	reserveTime := &reserve_time.GoTimes{}
	halfPastTwoPM := date.TodayUnix(14, 30, 0)
	now := time.Now().Unix()
	conf := config.Get()
	if now >= date.FirstSnapUpUnix()-conf.AdvanceTime && now <= date.FirstSnapUpUnix() {
		reserveTime.StartTimestamp = date.TodayUnix(6, 30, 0)
		reserveTime.EndTimestamp = halfPastTwoPM
		t.SetReserveTime(reserveTime)
		return
	}
	if now >= date.SecondSnapUpUnix()-conf.AdvanceTime && now <= date.SecondSnapUpUnix() {
		reserveTime.StartTimestamp = halfPastTwoPM
		reserveTime.EndTimestamp = date.TodayUnix(22, 30, 0)
		t.SetReserveTime(reserveTime)
		return
	}
	var fiveMinutes int64 = 5 * 60
	if now > halfPastTwoPM-fiveMinutes {
		reserveTime.StartTimestamp = now + fiveMinutes
		// reserveTime.StartTimestamp = (now/fiveMinutes + 1) * fiveMinutes
	} else {
		reserveTime.StartTimestamp = halfPastTwoPM
	}
	reserveTime.EndTimestamp = date.TodayUnix(22, 30, 0)
	t.SetReserveTime(reserveTime)
}

func (t *Task) CheckOrder() {
	for {
		if t.Completed {
			return
		}
		if t.CartMap() == nil || t.ReserveTime() == nil {
			<-time.After(10 * time.Millisecond)
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
			<-time.After(10 * time.Millisecond)
			continue
		}
		err := AddNewOrder(t.CartMap(), t.ReserveTime(), t.CheckOrderMap())
		if err != nil {
			log.Println(err)
			duration := 50 + rand.Intn(50)
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
	if conf.SnapUp&FirstSnapUp == FirstSnapUp && now.Unix() == firstTime-conf.AdvanceTime {
		log.Println("===== 6点抢购开始 =====")
		return true
	}
	secondTime := date.SecondSnapUpUnix()
	// log.Println(conf.SnapUp&SecondSnapUp == SecondSnapUp, now, secondTime, now.Unix(), secondTime)
	if conf.SnapUp&SecondSnapUp == SecondSnapUp && now.Unix() == secondTime-conf.AdvanceTime {
		log.Println("===== 8点半抢购开始 =====")
		return true
	}
	return false
}

func SnapUpOnce() {
	conf := config.Get()
	task := NewTask()
	task.MockMultiReserveTime() // 模拟配送时段

	for i := 0; i < conf.BaseConcurrency; i++ {
		go task.AllCheck()

		go task.GetCart()

		// go task.GetMultiReserveTime()

		go task.CheckOrder()
	}

	for i := 0; i < conf.SubmitConcurrency; i++ {
		go task.AddNewOrder()
	}

	timer := time.NewTimer(time.Minute * 2)
	<-timer.C
}

// SnapUp 抢购
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

// PickUp 捡漏
func PickUp() {
	for {
		<-pickUpCh
		SnapUpOnce()
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
		pickUpCh <- struct{}{}
	}

	if conf.NotifyNeeded {
		notifyCh <- struct{}{}
	}

	<-time.After(10 * time.Minute)
}

func Notify() {
	for {
		<-notifyCh
		now := time.Now()
		conf := config.Get()
		// interval := time.Duration(conf.NotifyInterval) * time.Minute
		// if now.Before(lastNotify.Add(interval)) {
		// 	continue
		// }
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
		for _, v := range conf.Users {
			if v == "" {
				continue
			}
			wg.Add(1)
			go func(token string) {
				defer wg.Done()
				Push(token, fmt.Sprintf("叮咚买菜当前可配送请尽快下单[%s%s]", products, ellipsis))
			}(v)
		}
		for _, v := range conf.AndroidUsers {
			if v == "" {
				continue
			}
			wg.Add(1)
			go func(token string) {
				defer wg.Done()
				PushToAndroid(token, fmt.Sprintf("叮咚买菜当前可配送请尽快下单[%s%s]", products, ellipsis))
			}(v)
		}
		wg.Wait()
	}
}
