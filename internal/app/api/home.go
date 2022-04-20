package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"dingdong/internal/app/config"
	"dingdong/pkg/json"
)

func SayWelcome(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome to this website")
}

// SetConfig 通过接口更新配置并重载配置
func SetConfig(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // 解析参数，默认是不会解析的，表单提交或者url传值
	if err != nil {
		_, _ = io.WriteString(w, err.Error()+"\n")
		return
	}
	if r.Method != http.MethodGet {
		_, _ = io.WriteString(w, "只支持GET方法\n")
		return
	}
	// log.Println("SetConfig Form => ", r.Form)
	// 获取旧配置
	conf := config.Get()
	if !validParams(r) {
		_, _ = io.WriteString(w, "没有提交任何有效的更新\n")
		return
	}
	if r.Form.Get("snap_up") != "" {
		snapUp, err := strconv.Atoi(r.Form.Get("snap_up"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.SnapUp = uint8(snapUp)
	}
	if r.Form.Get("pick_up_needed") != "" {
		pickUpNeeded := r.Form.Get("pick_up_needed")
		conf.MonitorNeeded = pickUpNeeded != "0"
	}
	if r.Form.Get("monitor_needed") != "" {
		monitorNeeded := r.Form.Get("monitor_needed")
		conf.MonitorNeeded = monitorNeeded != "0"
	}
	if r.Form.Get("notify_needed") != "" {
		notifyNeeded := r.Form.Get("notify_needed")
		conf.NotifyNeeded = notifyNeeded != "0"
	}
	if r.Form.Get("notify_interval") != "" {
		interval, err := strconv.Atoi(r.Form.Get("notify_interval"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.NotifyInterval = interval
	}
	if r.Form.Get("users") != "" {
		list := strings.Split(r.Form.Get("users"), ",")
		conf.Users = append(conf.Users, list...)
	}

	// 更新配置并写入文件
	bs := json.MustEncodePretty(conf)
	err = config.Set(bs)
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
		return
	}
	err = config.Reload()
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
		return
	}
	_, err = io.WriteString(w, "重载配置文件成功\n"+string(bs))
	if err != nil {
		log.Println("io.WriteString error =>", err)
	}
}

func validParams(r *http.Request) bool {
	if r.Form.Get("snap_up") != "" {
		return true
	}
	if r.Form.Get("pick_up_needed") != "" {
		return true
	}
	if r.Form.Get("monitor_needed") != "" {
		return true
	}
	if r.Form.Get("notify_needed") != "" {
		return true
	}
	if r.Form.Get("notify_interval") != "" {
		return true
	}
	if r.Form.Get("users") != "" {
		return true
	}
	return false
}
