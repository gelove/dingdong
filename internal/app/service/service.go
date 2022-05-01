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
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/internal/app/service/meituan"
	"dingdong/internal/app/service/notify"
)

const (
	FirstSnapUp = iota + 1
	SecondSnapUp
	PickUpMode
)

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
	cartMap       map[string]interface{}
	checkOrderMap map[string]interface{}
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
	return t.reserveTime
}

func (t *Task) SetReserveTime(reserveTime *reserve_time.GoTimes) *Task {
	t.Lock()
	defer t.Unlock()
	if reserveTime != nil {
		t.reserveTime = reserveTime
	}
	return t
}

func (t *Task) CartMap() map[string]interface{} {
	t.RLock()
	defer t.RUnlock()
	return t.cartMap
}

func (t *Task) SetCartMap(cartMap map[string]interface{}) *Task {
	t.Lock()
	defer t.Unlock()
	if cartMap != nil {
		t.cartMap = cartMap
	}
	return t
}

func (t *Task) CheckOrderMap() map[string]interface{} {
	t.RLock()
	defer t.RUnlock()
	return t.checkOrderMap
}

func (t *Task) SetCheckOrderMap(checkOrderMap map[string]interface{}) *Task {
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
			err := AllCheck()
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
			cartMap, err := GetCart()
			if err != nil {
				log.Println(err)
				if e, ok := err.(errs.Error); ok {
					if e.CodeEqual(code.NoValidProduct) {
						t.Finished()
						return
					}
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
			if t.CartMap() == nil {
				<-time.After(10 * time.Millisecond)
				continue
			}
			duration := durationMinMillis + rand.Intn(durationGapMillis)
			reserveTime, err := GetMultiReserveTime(t.CartMap())
			if err != nil {
				log.Println(err)
				<-time.After(time.Duration(duration) * time.Millisecond)
				continue
			}
			t.SetReserveTime(reserveTime)
			log.Println("===== 有效运力时段已更新 =====")
			// log.Println("reserveTime => ", json.MustEncodeToString(reserveTime))
			<-time.After(time.Duration(duration) * time.Millisecond)
		}
	}
}

func (t *Task) MockMultiReserveTime() {
	reserveTime := MockMultiReserveTime()
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
			if t.CartMap() == nil || t.ReserveTime() == nil {
				<-time.After(10 * time.Millisecond)
				continue
			}
			duration := durationMinMillis + rand.Intn(durationGapMillis)
			checkOrderMap, err := CheckOrder(t.CartMap(), t.ReserveTime())
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
			if t.CartMap() == nil || t.ReserveTime() == nil || t.CheckOrderMap() == nil {
				<-time.After(10 * time.Millisecond)
				continue
			}
			log.Println("===== 准备提交订单 =====")
			err := AddNewOrder(t.CartMap(), t.ReserveTime(), t.CheckOrderMap())
			if err != nil {
				_err := errs.New(code.ReserveTimeIsDisabled)
				if !errs.As(err, &_err) {
					t.Finished()
					return
				}
				log.Println(err)
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

func timeTrigger() int {
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

	if mode == FirstSnapUp || mode == SecondSnapUp {
		task.MockMultiReserveTime() // 更新运力
	}

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

// SnapUp 抢购
func SnapUp(ctx context.Context) {
	timer := time.Tick(time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Println("SnapUp force stopped")
			return
		case <-timer:
			mode := timeTrigger()
			if mode > 0 {
				go SnapUpOnce(mode)
			}
		}
	}
}

// PickUp 捡漏
func PickUp(ctx context.Context, pickUpCh <-chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			log.Println("PickUp force stopped")
			return
		case <-pickUpCh:
			go SnapUpOnce(PickUpMode)
		}

	}
}

func MonitorAndPickUp(cartMap map[string]interface{}, notifyCh chan<- struct{}, pickUpCh chan<- struct{}) {
	_, err := GetMultiReserveTime(cartMap)
	if err != nil {
		log.Println(err)
		return
	}

	conf := config.GetDingDong()
	if conf.PickUpNeeded {
		pickUpCh <- struct{}{}
	}

	if conf.NotifyNeeded {
		notifyCh <- struct{}{}
	}

	<-time.After(time.Duration(conf.MonitorSuccessWait) * time.Minute)
}

func DingDongNotify(notifyCh <-chan struct{}) {
	for {
		<-notifyCh
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

		Notify(fmt.Sprintf("叮咚买菜当前有运力请尽快下单[%s%s]", products, ellipsis))
	}
}

func MeiTuanMonitorAndNotify(notifyCh chan<- struct{}) {
	result, err := meituan.GetMultiReserveTime()
	if err != nil {
		log.Println(err)
		return
	}
	if result.Data.CycleType != 0 {
		log.Println("美团 =>", result.Data.Msg)
		return
	}

	conf := config.GetMeiTuan()
	if conf.NotifyNeeded {
		notifyCh <- struct{}{}
	}

	<-time.After(time.Duration(conf.MonitorSuccessWait) * time.Minute)
}

func MeiTuanNotify(notifyCh <-chan struct{}) {
	for {
		<-notifyCh
		Notify("美团买菜当前有运力请尽快下单")
	}
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
	err := AllCheck()
	if err != nil {
		return err
	}
	cartMap, err := GetCart()
	if err != nil {
		return err
	}
	reserveTimes, err := GetMultiReserveTime(cartMap)
	if err != nil {
		return err
	}
	orderMap, err := CheckOrder(cartMap, reserveTimes)
	if err != nil {
		return err
	}
	return AddNewOrder(cartMap, reserveTimes, orderMap)
}
