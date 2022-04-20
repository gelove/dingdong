package bark

import (
	"fmt"
	"net/http"

	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/json"
	"dingdong/pkg/notify"
)

const barkURL = "https://api.day.app/push"

var m = make(map[string]notify.Notifier)

type data struct {
	Badge     int    `json:"badge,omitempty"`
	DeviceKey string `json:"device_key"`
	Title     string `json:"title"`
	Body      string `json:"body,omitempty"`
	Sound     string `json:"sound,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Group     string `json:"group,omitempty"`
	Url       string `json:"url,omitempty"`
}

type bark struct {
	key   string
	icon  string
	sound string
}

func New(key, icon, sound string) notify.Notifier {
	if _, ok := m[key]; !ok {
		m[key] = &bark{key: key, icon: icon, sound: sound}
	}
	return m[key]
}

func Reset(key, icon, sound string) notify.Notifier {
	m[key] = &bark{key: key, icon: icon, sound: sound}
	return m[key]
}

func (b *bark) Name() string {
	return "Bark"
}

func (b *bark) Send(title, body string) error {
	d := &data{
		DeviceKey: b.key,
		Title:     title,
		Body:      body,
		Icon:      b.icon,
		Sound:     b.sound,
	}

	resp, err := session.Client().R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(json.MustEncode(d)).
		Send(http.MethodPost, barkURL)
	if err != nil {
		return errs.Wrap(code.RequestFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		return errs.WithMessage(code.InvalidResponse, fmt.Sprintf("%d %s", resp.StatusCode, resp.String()))
	}
	return nil
}
