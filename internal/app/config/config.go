package config

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"dingdong/pkg/json"
)

type Config struct {
	SnapUp         uint8             `json:"snap_up"`         // 抢购开关 0: 关 1: 6点抢 2: 8点半抢 3: 6点和8点半都抢
	PickUpNeeded   bool              `json:"pick_up_needed"`  // 闲时捡漏开关
	MonitorNeeded  bool              `json:"monitor_needed"`  // 监视器开关 监视是否有可配送时段
	NotifyNeeded   bool              `json:"notify_needed"`   // 通知开关 发现有可配送时段时通知大家有可购商品
	NotifyInterval int               `json:"notify_interval"` // 通知间隔 单位: 分钟
	Headers        map[string]string `json:"headers"`
	Params         map[string]string `json:"params"`
	Users          []string          `json:"users"`
}

type conf struct {
	sync.RWMutex
	Pid      int
	FilePath string
	Config   Config
}

func NewConf(pid int, path string) *conf {
	return &conf{
		Pid:      pid,
		FilePath: path,
	}
}

var c *conf

func Initialize(path string) {
	pid := os.Getpid()
	log.Println("当前程序 PID => ", pid)
	c = NewConf(pid, path)

	if !load() {
		os.Exit(1)
	}
	// 热更新配置可能有多种触发方式，这里使用系统信号量sigusr1实现
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR1)
	go func() {
		for {
			<-s
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
	json.MustDecode(f, &temp)
	c.Lock()
	c.Config = temp
	c.Unlock()
	log.Println("Load conf =>", json.MustEncodeToString(c))
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

func Reload() error {
	return syscall.Kill(Pid(), syscall.SIGUSR1)
}

func Pid() int {
	return c.Pid
}

func FilePath() string {
	return c.FilePath
}
