package bark

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/imroc/req/v3"

	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/json"
	"dingdong/pkg/notify"
)

const URL = "https://api.day.app/push"

var cache sync.Map

type bark struct {
	DeviceKey string `json:"device_key"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Sound     string `json:"sound,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Badge     int    `json:"badge,omitempty"`
	Group     string `json:"group,omitempty"`
	Url       string `json:"url,omitempty"`
}

func New(key, title, body, icon, sound string) notify.Notifier {
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
	instance := &bark{DeviceKey: key, Title: title, Body: body, Icon: icon, Sound: sound}
	cache.Store(key, instance)
	return instance
}

func (b bark) Name() string {
	return "Bark"
}

func (b bark) Send() error {
	resp, err := req.C().R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(json.MustEncode(b)).
		Send(http.MethodPost, URL)
	if err != nil {
		return errs.Wrap(code.RequestFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		return errs.WithMessage(code.InvalidResponse, fmt.Sprintf("%d %s", resp.StatusCode, resp.String()))
	}
	return nil
}
