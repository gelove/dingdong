package bark

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/imroc/req/v3"

	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/notify"
)

const API = "https://api.day.app"

var cache sync.Map

type bark struct {
	Api       string `json:"-"`
	DeviceKey string `json:"device_key"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Sound     string `json:"sound,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Group     string `json:"group,omitempty"`
	Url       string `json:"url,omitempty"`
	Badge     int    `json:"badge,omitempty"`
}

func New(key, title, body, icon, sound string) notify.Notifier {
	api := API
	if strings.HasPrefix(key, "http") {
		list := strings.Split(key, "/")
		key = list[3]
		api = strings.Join(list[:3], "/")
	}
	if v, ok := cache.Load(key); ok {
		instance := v.(*bark)
		if instance.Title != title {
			instance.Title = title
		}
		if instance.Body != body {
			instance.Body = body
		}
		if instance.Icon != icon {
			instance.Icon = icon
		}
		if instance.Sound != sound {
			instance.Sound = sound
		}
		return instance
	}
	instance := &bark{DeviceKey: key, Api: api, Title: title, Body: body, Icon: icon, Sound: sound}
	cache.Store(key, instance)
	return instance
}

func (b bark) Name() string {
	return "Bark"
}

func (b bark) Send() error {
	url := fmt.Sprintf("%s/%s/%s/%s?sound=%s&icon=%s", b.Api, b.DeviceKey, b.Title, b.Body, b.Sound, b.Icon)
	resp, err := req.C().R().Send(http.MethodGet, url)
	if err != nil {
		return errs.Wrap(code.RequestFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		return errs.WithMessage(code.ResponseError, fmt.Sprintf("%d %s", resp.StatusCode, resp.String()))
	}
	return nil
}
