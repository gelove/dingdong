package app

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dingdong/assets"
	"dingdong/internal/app/api"
	"dingdong/internal/app/config"
	"dingdong/internal/app/service"
	"dingdong/internal/app/service/ios_service"
	"dingdong/internal/app/service/meituan"
)

var (
	pickUpCh         = make(chan struct{})
	dingDongNotifyCh = make(chan string)
	merTuanNotifyCh  = make(chan string)
)

func Run() {
	done := make(chan struct{}, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())

	conf := config.Get()
	server := newWebServer(conf.Addr)

	go func() {
		s := <-quit
		log.Printf("Got signal: %+v. Server is shutting down...", s)
		defer cancel()

		// gracefulShutdown 给服务器30秒时间优雅关闭
		// graceCtx, graceCancel := context.WithTimeout(context.Background(), 30*time.Second)
		// defer graceCancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	go SnapUp(ctx)
	go PickUp(ctx, pickUpCh)
	go DingDongNotify(ctx, dingDongNotifyCh)
	go MeiTuanNotify(ctx, merTuanNotifyCh)
	go DingDongMonitor(ctx, dingDongNotifyCh, pickUpCh)
	go MeiTuanMonitor(ctx, merTuanNotifyCh)

	log.Println("Server is ready to handle requests at", conf.Addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", conf.Addr, err)
	}

	<-done
	log.Println("Server stopped")
}

func newWebServer(addr string) *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/", api.ConfigView)
	router.HandleFunc("/set", api.SetConfig)
	router.HandleFunc("/config", api.ConfigView)
	router.HandleFunc("/notify", api.Notify)
	router.HandleFunc("/address", api.GetAddress)
	router.HandleFunc("/addOrder", api.AddOrder)
	router.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets.FS))))

	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

// SnapUp 抢购
func SnapUp(ctx context.Context) {
	timer := time.Tick(time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Println("SnapUp stopped")
			return
		case <-timer:
			mode := service.TimeTrigger()
			if mode > 0 {
				go service.SnapUpOnce(mode)
			}
		}
	}
}

// PickUp 捡漏
func PickUp(ctx context.Context, pickUpCh <-chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			log.Println("PickUp stopped")
			return
		case <-pickUpCh:
			go service.SnapUpOnce(service.PickUpMode)
		}
	}
}

// DingDongMonitor 监视器 监听运力
func DingDongMonitor(ctx context.Context, notifyCh chan<- string, pickUpCh chan<- struct{}) {
	for {
		select {
		case <-ctx.Done():
			log.Println("DingDongMonitor stopped")
			return
		default:
			<-time.After(time.Second)
			conf := config.GetDingDong()
			if !conf.MonitorNeeded && !conf.PickUpNeeded {
				continue
			}
			now := time.Now()
			// 每分钟运行一次
			if now.Second() != 0 {
				continue
			}
			// 每5分钟运行一次
			// if now.Minute()%5 != 0 || now.Second() != 0 {
			// 	continue
			// }
			// 随机在第1-5秒运行
			random := rand.Intn(3)
			<-time.After(time.Second * time.Duration(random))
			if dingDongIsPeak() {
				log.Println("[叮咚]当前高峰期或暂未营业")
				continue
			}
			MonitorAndPickUp(notifyCh, pickUpCh)
		}
	}
}

func dingDongIsPeak() bool {
	now := time.Now()
	if now.Hour() >= 0 && now.Hour() < 9 {
		return true
	}
	if now.Hour() >= 22 {
		return true
	}
	return false
}

func MonitorAndPickUp(notifyCh chan<- string, pickUpCh chan<- struct{}) {
	cartMap := ios_service.MockCartMap()
	_, err := ios_service.GetMultiReserveTime(cartMap)
	if err != nil {
		log.Println("[叮咚]", err)
		return
	}
	// detail := getDetails()
	detail := "[叮咚买菜]当前有配送运力请尽快下单"
	log.Println(detail)

	conf := config.GetDingDong()
	if conf.PickUpNeeded {
		pickUpCh <- struct{}{}
	}

	if conf.NotifyNeeded {
		notifyCh <- detail
	}

	<-time.After(time.Duration(conf.MonitorSuccessWait) * time.Minute)
}

func DingDongNotify(ctx context.Context, notifyCh <-chan string) {
	for {
		select {
		case <-ctx.Done():
			log.Println("DingDongNotify stopped")
			return
		case detail := <-notifyCh:
			service.Notify(detail)
		}
	}
}

func MeiTuanMonitor(ctx context.Context, notifyCh chan<- string) {
	for {
		select {
		case <-ctx.Done():
			log.Println("MeiTuanMonitor stopped")
			return
		default:
			<-time.After(time.Second)
			conf := config.GetMeiTuan()
			if !conf.MonitorNeeded {
				continue
			}
			now := time.Now()
			if now.Second() != 0 {
				continue
			}
			// 每分钟在第1-3秒运行一次
			random := rand.Intn(5) + 1
			<-time.After(time.Second * time.Duration(random))
			if meiTuanIsPeak() {
				log.Println("[美团]当前高峰期或暂未营业")
				continue
			}
			MeiTuanMonitorAndNotify(notifyCh)
		}
	}
}

func meiTuanIsPeak() bool {
	now := time.Now()
	if now.Hour() >= 0 && now.Hour() < 6 {
		return true
	}
	if now.Hour() == 6 && now.Minute() < 30 {
		return true
	}
	if now.Hour() >= 22 {
		return true
	}
	return false
}

func MeiTuanMonitorAndNotify(notifyCh chan<- string) {
	result, err := meituan.GetMultiReserveTime()
	if err != nil {
		log.Println(err)
		return
	}
	if result.Data.CycleType != 0 {
		log.Println("[美团]", result.Data.Msg)
		return
	}
	detail := "[美团买菜]当前有配送运力请尽快下单"
	log.Println(detail)

	conf := config.GetMeiTuan()
	if conf.NotifyNeeded {
		notifyCh <- detail
	}

	<-time.After(time.Duration(conf.MonitorSuccessWait) * time.Minute)
}

func MeiTuanNotify(ctx context.Context, notifyCh chan string) {
	for {
		select {
		case <-ctx.Done():
			log.Println("MeiTuanNotify stopped")
			return
		case detail := <-notifyCh:
			service.Notify(detail)
		}
	}
}
