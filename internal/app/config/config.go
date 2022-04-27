package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"dingdong/pkg/yaml"
)

type Config struct {
	Name               string            `yaml:"name"`                 // 程序名称
	Addr               string            `yaml:"addr"`                 // web服务地址
	BaseConcurrency    int               `yaml:"base_concurrency"`     // 基础并发数(除了提交订单的其他请求, 默认为1)
	SubmitConcurrency  int               `yaml:"submit_concurrency"`   // 提交订单并发数(默认为2)
	SnapUp             int               `yaml:"snap_up"`              // 抢购开关 0: 关 1: 6点抢 2: 8点半抢 3: 6点和8点半都抢
	AdvanceTime        int64             `yaml:"advance_time"`         // 抢购提前进入时间 单位:秒
	PickUpNeeded       bool              `yaml:"pick_up_needed"`       // 闲时捡漏开关
	MonitorNeeded      bool              `yaml:"monitor_needed"`       // 监视器开关 监视是否有可配送时段
	MonitorSuccessWait int               `yaml:"monitor_success_wait"` // 成功监听(发起捡漏或通知)之后的休息时间 单位:分钟
	MonitorIntervalMin int               `yaml:"monitor_interval_min"` // 监视器调用接口的最小时间间隔 单位:秒
	MonitorIntervalMax int               `yaml:"monitor_interval_max"` // 监视器调用接口的最大时间间隔 单位:秒
	NotifyNeeded       bool              `yaml:"notify_needed"`        // 通知开关 发现有可配送时段时通知大家有可购商品
	AudioNeeded        bool              `yaml:"audio_needed"`         // 播放音频开关 在下单成功后播放音频
	Headers            map[string]string `yaml:"headers"`              // 请求头
	Mock               map[string]string `yaml:"mock,omitempty"`       // 模拟参数测试用
	Bark               []string          `yaml:"bark"`                 // bark 用户
	PushPlus           []string          `yaml:"push_plus"`            // push_plus 用户
}

type conf struct {
	sync.RWMutex
	refresh  chan struct{}
	Pid      int
	FilePath string
	Config   Config
}

func NewConf(pid int, path string) *conf {
	return &conf{
		refresh:  make(chan struct{}),
		Pid:      pid,
		FilePath: path,
	}
}

var c *conf

func Initialize(path string) {
	pid := os.Getpid()
	log.Println("当前程序进程 PID => ", pid)
	c = NewConf(pid, path)

	if !load() {
		os.Exit(1)
	}
	// 热更新配置
	go func() {
		for {
			<-c.refresh
			if load() {
				log.Println("热更新配置成功")
			} else {
				log.Println("热更新配置失败")
			}
		}
	}()
}

// Exists 判断文件是否存在 存在返回true 不存在返回false
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func load() bool {
	f, err := ioutil.ReadFile(c.FilePath)
	if err != nil {
		log.Println("Load config error =>", err)
		return false
	}
	var temp Config
	yaml.MustDecode(f, &temp)
	c.Lock()
	c.Config = temp
	c.Unlock()
	bs := yaml.MustEncode(c.Config)
	log.Printf("已加载配置 => \n%s\n\n", string(bs))
	return true
}

func Get() Config {
	c.RLock()
	defer c.RUnlock()
	return c.Config
}

func Set(bytes []byte) error {
	path := FilePath()
	if !Exists(path) {
		return errors.New("配置文件不存在")
	}
	return ioutil.WriteFile(path, bytes, 0666)
}

func Reload() {
	c.refresh <- struct{}{}
}

func Pid() int {
	return c.Pid
}

func FilePath() string {
	return c.FilePath
}
