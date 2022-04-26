package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/pkg/json"
)

func SayWelcome(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome to this website")
}

// GetAddress 获取地址
func GetAddress(w http.ResponseWriter, _ *http.Request) {
	list, err := session.GetAddress()
	if err != nil {
		_, _ = io.WriteString(w, err.Error()+"\n")
		return
	}
	_, err = io.WriteString(w, json.MustEncodeToString(list)+"\n")
	if err != nil {
		log.Println("io.WriteString error =>", err)
	}
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
	if r.Form.Get("base_concurrency") != "" {
		baseConcurrency, err := strconv.Atoi(r.Form.Get("base_concurrency"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.BaseConcurrency = baseConcurrency
	}
	if r.Form.Get("submit_concurrency") != "" {
		submitConcurrency, err := strconv.Atoi(r.Form.Get("submit_concurrency"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.SubmitConcurrency = submitConcurrency
	}
	if r.Form.Get("snap_up") != "" {
		snapUp, err := strconv.Atoi(r.Form.Get("snap_up"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.SnapUp = snapUp
	}
	if r.Form.Get("advance_time") != "" {
		advanceTime, err := strconv.Atoi(r.Form.Get("advance_time"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.AdvanceTime = int64(advanceTime)
	}
	if r.Form.Get("pick_up_needed") != "" {
		pickUpNeeded := r.Form.Get("pick_up_needed")
		conf.PickUpNeeded = pickUpNeeded != "0"
	}
	if r.Form.Get("monitor_needed") != "" {
		monitorNeeded := r.Form.Get("monitor_needed")
		conf.MonitorNeeded = monitorNeeded != "0"
	}
	if r.Form.Get("monitor_success_wait") != "" {
		monitorSuccessWait, err := strconv.Atoi(r.Form.Get("monitor_success_wait"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.MonitorSuccessWait = monitorSuccessWait
	}
	if r.Form.Get("monitor_interval_min") != "" {
		monitorIntervalMin, err := strconv.Atoi(r.Form.Get("monitor_interval_min"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.MonitorIntervalMin = monitorIntervalMin
	}
	if r.Form.Get("monitor_interval_max") != "" {
		monitorIntervalMax, err := strconv.Atoi(r.Form.Get("monitor_interval_max"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		if monitorIntervalMax < conf.MonitorIntervalMin {
			_, _ = io.WriteString(w, "监控间隔最大值不能小于最小值\n")
			return
		}
		conf.MonitorIntervalMax = monitorIntervalMax
	}
	if r.Form.Get("notify_needed") != "" {
		notifyNeeded := r.Form.Get("notify_needed")
		conf.NotifyNeeded = notifyNeeded != "0"
	}
	if r.Form.Get("audio_needed") != "" {
		audioNeeded := r.Form.Get("audio_needed")
		conf.AudioNeeded = audioNeeded != "0"
	}
	if r.Form.Get("users") != "" {
		list := strings.Split(r.Form.Get("users"), ",")
		conf.Users = append(conf.Users, list...)
	}
	if r.Form.Get("an_users") != "" {
		list := strings.Split(r.Form.Get("an_users"), ",")
		conf.AndroidUsers = append(conf.AndroidUsers, list...)
	}

	// 更新配置并写入文件
	bs := json.MustEncodePretty(conf)
	err = config.Set(bs)
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
		return
	}
	config.Reload()
	_, err = io.WriteString(w, "重载配置文件成功\n"+string(bs))
	if err != nil {
		log.Println("io.WriteString error =>", err)
	}
}

func validParams(r *http.Request) bool {
	if r.Form.Get("base_concurrency") != "" {
		return true
	}
	if r.Form.Get("submit_concurrency") != "" {
		return true
	}
	if r.Form.Get("snap_up") != "" {
		return true
	}
	if r.Form.Get("advance_time") != "" {
		return true
	}
	if r.Form.Get("pick_up_needed") != "" {
		return true
	}
	if r.Form.Get("monitor_needed") != "" {
		return true
	}
	if r.Form.Get("monitor_success_wait") != "" {
		return true
	}
	if r.Form.Get("monitor_interval_min") != "" {
		return true
	}
	if r.Form.Get("monitor_interval_max") != "" {
		return true
	}
	if r.Form.Get("notify_needed") != "" {
		return true
	}
	if r.Form.Get("audio_needed") != "" {
		return true
	}
	if r.Form.Get("users") != "" {
		return true
	}
	if r.Form.Get("an_users") != "" {
		return true
	}
	return false
}
