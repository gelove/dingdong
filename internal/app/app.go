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
	go SnapUp()
	go PickUp()
	go Notify()
	http.HandleFunc("/", api.SayWelcome)
	http.HandleFunc("/set", api.SetConfig)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		log.Fatal("ListenAndServe => ", err)
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

func SnapUp() {
	service.SnapUp()
}

func PickUp() {
	service.PickUp()
}

func Notify() {
	service.Notify()
}
