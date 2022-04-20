package app

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"dingdong/internal/app/api"
	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/date"
	"dingdong/internal/app/service"
)

const (
	FirstSnapUp uint8 = iota + 1
	SecondSnapUp
)

func Run() {
	http.HandleFunc("/", api.SayWelcome)
	http.HandleFunc("/set", api.SetConfig)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe => ", err)
	}
}

// Monitor 监视器 每8-15秒调用一次接口
func Monitor() {
	cartMap := service.MockCartMap()
	for {
		conf := config.Get()
		if conf.MonitorNeeded {
			if isPeak() {
				log.Println("当前高峰期或暂未营业")
			} else {
				service.GetMultiReserveTimeAndNotify(cartMap)
			}
		}
		duration := 8 + rand.Intn(8)
		<-time.After(time.Duration(duration) * time.Second)
	}
}

func isPeak() bool {
	now := time.Now()
	if now.Hour() >= 0 && now.Hour() < 6 {
		return true
	}
	if now.Hour() == 6 && now.Minute() < 30 {
		return true
	}
	if now.Hour() == 8 {
		return true
	}
	return false
}

func Shopping() {
	timer := time.Tick(time.Second)
	for {
		select {
		case <-timer:
			if timeTrigger() {
				go service.Run()
			}
		}
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
