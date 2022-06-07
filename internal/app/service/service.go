package service

import (
	"context"
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
	"dingdong/internal/app/service/ios_service"
	"dingdong/internal/app/service/notify"
	"dingdong/pkg/json"
)

const (
	FirstSnapUp = 1 << iota
	SecondSnapUp
	PickUpMode
)

type errCode uintptr

const (
	durationMinMillis = 450
	durationMaxMillis = 550
	durationGapMillis = durationMaxMillis - durationMinMillis
)

type Task struct {
	sync.RWMutex
	sync.WaitGroup
	timeOut       context.Context
	Finished      context.CancelFunc
	completed     bool
	reserveTime   *reserve_time.GoTimes
	cartMap       map[string]any
	checkOrderMap map[string]string
}

func NewTask() *Task {
	timeOut, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	return &Task{timeOut: timeOut, Finished: cancel}
}

func (t *Task) Completed() bool {
	t.RLock()
	defer t.RUnlock()
	return t.completed
}

func (t *Task) SetCompleted(completed bool) *Task {
	t.Lock()
	defer t.Unlock()
	t.completed = completed
	return t
}

func (t *Task) ReserveTime() *reserve_time.GoTimes {
	t.RLock()
	defer t.RUnlock()
	if t.reserveTime == nil {
		return nil
	}
	reserveTime := *t.reserveTime
	return &reserveTime
}

func (t *Task) SetReserveTime(reserveTime *reserve_time.GoTimes) *Task {
	t.Lock()
	defer t.Unlock()
	if reserveTime != nil {
		t.reserveTime = reserveTime
	}
	return t
}

func (t *Task) CartMap() map[string]any {
	t.RLock()
	defer t.RUnlock()
	if t.cartMap == nil {
		return nil
	}
	copyCart := make(map[string]any)
	json.MustTransform(t.cartMap, &copyCart)
	return copyCart
}

func (t *Task) SetCartMap(cartMap map[string]any) *Task {
	t.Lock()
	defer t.Unlock()
	if cartMap != nil {
		t.cartMap = cartMap
	}
	return t
}

func (t *Task) CheckOrderMap() map[string]string {
	t.RLock()
	defer t.RUnlock()
	if t.checkOrderMap == nil {
		return nil
	}
	orderMap := make(map[string]string)
	json.MustTransform(t.checkOrderMap, &orderMap)
	return orderMap
}

func (t *Task) SetCheckOrderMap(checkOrderMap map[string]string) *Task {
	t.Lock()
	defer t.Unlock()
	if checkOrderMap != nil {
		t.checkOrderMap = checkOrderMap
	}
	return t
}

// AllCheck 不一定需要, 只起补充作用
func (t *Task) AllCheck() {
	defer t.Done()
	for {
		select {
		case <-t.timeOut.Done():
			log.Println("AllCheck finished")
			return
		default:
			err := ios_service.AllCheck()
			if err != nil {
				log.Println(err)
				duration := durationMinMillis + rand.Intn(durationGapMillis)
				<-time.After(time.Duration(duration) * time.Millisecond)
				continue
			}
			log.Println("===== 购物车全选成功 =====")
			<-time.After(time.Second * 5)
		}
	}
}

func (t *Task) GetCart() {
	defer t.Done()
	for {
		select {
		case <-t.timeOut.Done():
			log.Println("GetCart finished")
			return
		default:
			duration := durationMinMillis + rand.Intn(durationGapMillis)
			cartMap, err := ios_service.GetCart()
			if err != nil {
				log.Println(err)
				if errs.Is(err, errs.NoValidProduct) {
					t.Finished()
					return
				}
				<-time.After(time.Duration(duration) * time.Millisecond)
				continue
			}
			t.SetCartMap(cartMap)
			<-time.After(time.Duration(duration) * time.Millisecond)
		}
	}
}

func (t *Task) GetMultiReserveTime() {
	defer t.Done()
	for {
		select {
		case <-t.timeOut.Done():
			log.Println("GetMultiReserveTime finished")
			return
		default:
			cartMap := t.CartMap()
			if cartMap == nil {
				<-time.After(10 * time.Millisecond)
				continue
			}
			duration := durationMinMillis + rand.Intn(durationGapMillis)
			reserveTime, err := ios_service.GetMultiReserveTime(cartMap)
			if err != nil {
				log.Println(err)
				if errs.Is(err, errs.NoReserveTimeAndRetry) {
					t.Finished()
					return
				}
				<-time.After(time.Duration(duration) * time.Millisecond)
				continue
			}
			t.SetReserveTime(reserveTime)
			log.Printf("[叮咚]发现可用运力[%s-%s], 请尽快下单!", reserveTime.StartTime, reserveTime.EndTime)
			// log.Println("reserveTime => ", json.MustEncodeToString(reserveTime))
			<-time.After(time.Duration(duration) * time.Millisecond)
		}
	}
}

func (t *Task) MockMultiReserveTime() {
	reserveTime := ios_service.MockMultiReserveTime()
	t.SetReserveTime(reserveTime)
}

func (t *Task) CheckOrder() {
	defer t.Done()
	for {
		select {
		case <-t.timeOut.Done():
			log.Println("CheckOrder finished")
			return
		default:
			cartMap := t.CartMap()
			reserveTime := t.ReserveTime()
			if cartMap == nil || reserveTime == nil {
				<-time.After(10 * time.Millisecond)
				continue
			}
			duration := durationMinMillis + rand.Intn(durationGapMillis)
			checkOrderMap, err := ios_service.CheckOrder(cartMap)
			if err != nil {
				log.Println(err)
				<-time.After(time.Duration(duration) * time.Millisecond)
				continue
			}
			t.SetCheckOrderMap(checkOrderMap)
			log.Println("===== 订单信息已更新 =====")
			<-time.After(time.Duration(duration) * time.Millisecond)
		}
	}
}

func (t *Task) AddNewOrder() {
	defer t.Done()

	for {
		select {
		case <-t.timeOut.Done():
			log.Println("AddNewOrder finished")
			return
		default:
			cartMap := t.CartMap()
			reserveTime := t.ReserveTime()
			orderMap := t.CheckOrderMap()
			if cartMap == nil || reserveTime == nil || orderMap == nil {
				<-time.After(10 * time.Millisecond)
				continue
			}
			log.Println("===== 准备提交订单 =====")
			err := ios_service.AddNewOrder(cartMap, reserveTime, orderMap)
			if err != nil {
				log.Println(err)
				if errs.Is(err, errs.ReserveTimeIsDisabled) {
					t.SetReserveTime(nil)
				}
				duration := 20 + rand.Intn(80)
				<-time.After(time.Duration(duration) * time.Millisecond)
				continue
			}
			detail := "已成功下单, 请尽快完成支付"
			log.Println(detail)
			conf := config.Get()
			if conf.DingDong.NotifyNeeded && len(conf.Bark) > 0 {
				go notify.Push(conf.Bark[0], detail)
			}
			if conf.DingDong.NotifyNeeded && len(conf.PushPlus) > 0 {
				go notify.PushPlus(conf.PushPlus[0], detail)
			}
			if conf.DingDong.AudioNeeded {
				go notify.PlayMp3()
			}
			t.Finished()
		}
	}
}

func TimeTrigger() int {
	conf := config.GetDingDong()
	now := time.Now()
	firstTime := date.FirstSnapUpUnix()
	// log.Println(conf.SnapUp&FirstSnapUp == FirstSnapUp, now, firstTime, now.Unix(), firstTime)
	if conf.SnapUp&FirstSnapUp == FirstSnapUp && now.Unix() == firstTime-conf.AdvanceTime {
		log.Println("===== 6点抢购开始 =====")
		return FirstSnapUp
	}
	secondTime := date.SecondSnapUpUnix()
	// log.Println(conf.SnapUp&SecondSnapUp == SecondSnapUp, now, secondTime, now.Unix(), secondTime)
	if conf.SnapUp&SecondSnapUp == SecondSnapUp && now.Unix() == secondTime-conf.AdvanceTime {
		log.Println("===== 8点半抢购开始 =====")
		return SecondSnapUp
	}
	return 0
}

func SnapUpOnce(mode int) {
	conf := config.GetDingDong()
	// wg := new(sync.WaitGroup)
	task := NewTask()
	defer task.Finished()

	for i := 0; i < conf.BaseConcurrency; i++ {
		task.Add(1)
		go task.AllCheck()
	}

	for i := 0; i < conf.BaseConcurrency; i++ {
		task.Add(1)
		go task.GetCart()
	}

	if mode == PickUpMode {
		for i := 0; i < conf.BaseConcurrency; i++ {
			task.Add(1)
			go task.GetMultiReserveTime()
		}
	}

	// if mode == FirstSnapUp || mode == SecondSnapUp {
	// 	task.MockMultiReserveTime() // 更新运力
	// }

	for i := 0; i < conf.BaseConcurrency; i++ {
		task.Add(1)
		go task.CheckOrder()
	}

	submitConcurrency := 1
	// 只提前2秒开始提交订单, 提前太早有可能被风控
	if mode == FirstSnapUp || mode == SecondSnapUp {
		submitConcurrency = conf.SubmitConcurrency
		<-time.After(time.Duration(60-2-time.Now().Second()) * time.Second)
	}
	for i := 0; i < submitConcurrency; i++ {
		task.Add(1)
		go task.AddNewOrder()
	}

	task.Wait()
}

func getDetails() string {
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

	return fmt.Sprintf("[叮咚买菜]当前有配送运力请尽快下单[%s%s]", products, ellipsis)
}

func Notify(content string) {
	conf := config.Get()
	wg := new(sync.WaitGroup)
	for _, v := range conf.Bark {
		if v == "" {
			continue
		}
		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			notify.Push(token, content)
		}(v)
	}
	for _, v := range conf.PushPlus {
		if v == "" {
			continue
		}
		wg.Add(1)
		go func(token string) {
			defer wg.Done()
			notify.PushPlus(token, content)
		}(v)
	}
	wg.Wait()
}

func AddOrder() error {
	err := ios_service.AllCheck()
	if err != nil {
		return err
	}
	cartMap, err := ios_service.GetCart()
	if err != nil {
		return err
	}
	reserveTimes, err := ios_service.GetMultiReserveTime(cartMap)
	if err != nil {
		return err
	}
	orderMap, err := ios_service.CheckOrder(cartMap)
	if err != nil {
		return err
	}
	return ios_service.AddNewOrder(cartMap, reserveTimes, orderMap)
}
