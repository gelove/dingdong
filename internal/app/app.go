package app

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"dingdong/internal/app/api"
	"dingdong/internal/app/config"
	"dingdong/internal/app/service"
)

func Run() {
	go Monitor()
	go service.SnapUp()
	go service.PickUp()
	go service.Notify()
	http.HandleFunc("/", api.SayWelcome)
	http.HandleFunc("/set", api.SetConfig)
	conf := config.Get()
	err := http.ListenAndServe(conf.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe => ", err)
	}
}

func isPeak() bool {
	now := time.Now()
	if now.Hour() >= 0 && now.Hour() < 8 {
		return true
	}
	if now.Hour() == 8 && now.Minute() < 50 {
		return true
	}
	return false
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
