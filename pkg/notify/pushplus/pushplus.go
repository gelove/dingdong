package pushplus

import (
	"fmt"
	"net/http"
	"sync"

	"dingdong/internal/app/pkg/ddmc/session"
	"dingdong/internal/app/pkg/errs"
	"dingdong/internal/app/pkg/errs/code"
	"dingdong/pkg/json"
	"dingdong/pkg/notify"
)

const URL = "http://www.pushplus.plus/send"

var cache sync.Map

type data struct {
	Token   string `json:"token"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type pusher struct {
	token string
}

func New(token string) notify.Notifier {
	if v, ok := cache.Load(token); ok {
		return v.(notify.Notifier)
	}
	instance := &pusher{token: token}
	cache.Store(token, instance)
	return instance
}

func (b *pusher) Name() string {
	return "PushPlus"
}

func (b *pusher) Send(title, content string) error {
	d := &data{
		Token:   b.token,
		Title:   title,
		Content: content,
	}

	resp, err := session.Client().R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(json.MustEncode(d)).
		Send(http.MethodPost, URL)
	if err != nil {
		return errs.Wrap(code.RequestFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		return errs.WithMessage(code.InvalidResponse, fmt.Sprintf("%d %s", resp.StatusCode, resp.String()))
	}
	return nil
}
