package api

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"dingdong/assets"
	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/service"
	"dingdong/internal/app/service/notify"
	"dingdong/pkg/json"
	"dingdong/pkg/yaml"
)

var tmpl *template.Template

func init() {
	t, err := template.ParseFS(assets.FS, "template/*.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl = t
}

func ConfigView(w http.ResponseWriter, r *http.Request) {
	conf := config.GetDingDong()
	err := tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{"title": "叮咚买菜助手", "conf": conf})
	if err != nil {
		_, _ = io.WriteString(w, err.Error()+"\n")
	}
}

// Notify 发送通知测试
func Notify(w http.ResponseWriter, _ *http.Request) {
	conf := config.Get()
	if len(conf.Bark) > 0 {
		for _, v := range conf.Bark {
			notify.Push(v, "Bark 测试")
		}
	}
	if len(conf.PushPlus) > 0 {
		for _, v := range conf.PushPlus {
			notify.PushPlus(v, "PushPlus 测试")
		}
	}
	_, _ = io.WriteString(w, "已发送通知, 如未收到通知请查看配置是否正确\n")
}

// GetAddress 获取地址
func GetAddress(w http.ResponseWriter, _ *http.Request) {
	list, err := session.GetAddress()
	if err != nil {
		_, _ = io.WriteString(w, err.Error()+"\n")
		return
	}
	_, _ = io.WriteString(w, json.MustEncodeToString(list)+"\n")
}

// AddOrder 提交订单
func AddOrder(w http.ResponseWriter, _ *http.Request) {
	err := service.AddOrder()
	if err != nil {
		_, _ = io.WriteString(w, err.Error()+"\n")
		return
	}
	_, _ = io.WriteString(w, "创建订单成功\n")
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
		conf.DingDong.BaseConcurrency = baseConcurrency
	}
	if r.Form.Get("submit_concurrency") != "" {
		submitConcurrency, err := strconv.Atoi(r.Form.Get("submit_concurrency"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.DingDong.SubmitConcurrency = submitConcurrency
	}
	if r.Form.Get("snap_up") != "" {
		snapUp, err := strconv.Atoi(r.Form.Get("snap_up"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.DingDong.SnapUp = snapUp
	}
	if r.Form.Get("advance_time") != "" {
		advanceTime, err := strconv.Atoi(r.Form.Get("advance_time"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.DingDong.AdvanceTime = int64(advanceTime)
	}
	if r.Form.Get("pick_up_needed") != "" {
		pickUpNeeded := r.Form.Get("pick_up_needed")
		conf.DingDong.PickUpNeeded = pickUpNeeded != "0"
	}
	if r.Form.Get("monitor_needed") != "" {
		monitorNeeded := r.Form.Get("monitor_needed")
		conf.DingDong.MonitorNeeded = monitorNeeded != "0"
	}
	if r.Form.Get("monitor_success_wait") != "" {
		monitorSuccessWait, err := strconv.Atoi(r.Form.Get("monitor_success_wait"))
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}
		conf.DingDong.MonitorSuccessWait = monitorSuccessWait
	}
	if r.Form.Get("notify_needed") != "" {
		notifyNeeded := r.Form.Get("notify_needed")
		conf.DingDong.NotifyNeeded = notifyNeeded != "0"
	}
	if r.Form.Get("audio_needed") != "" {
		audioNeeded := r.Form.Get("audio_needed")
		conf.DingDong.AudioNeeded = audioNeeded != "0"
	}
	if r.Form.Get("bark") != "" {
		list := strings.Split(r.Form.Get("bark"), ",")
		conf.Bark = append(conf.Bark, list...)
	}
	if r.Form.Get("push_plus") != "" {
		list := strings.Split(r.Form.Get("push_plus"), ",")
		conf.PushPlus = append(conf.PushPlus, list...)
	}

	// 更新配置并写入文件
	bs := yaml.MustEncode(conf)
	err = config.Set(bs)
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
		return
	}
	config.Reload()
	_, err = io.WriteString(w, "重载配置文件成功\n"+string(bs)+"\n")
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
	if r.Form.Get("notify_needed") != "" {
		return true
	}
	if r.Form.Get("audio_needed") != "" {
		return true
	}
	if r.Form.Get("bark") != "" {
		return true
	}
	if r.Form.Get("push_plus") != "" {
		return true
	}
	return false
}
