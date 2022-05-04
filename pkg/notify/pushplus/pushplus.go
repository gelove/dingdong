package pushplus

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

const URL = "http://www.pushplus.plus/send"

var cache sync.Map

type pusher struct {
	Token   string `json:"token"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func New(token, title, content string) notify.Notifier {
	if v, ok := cache.Load(token); ok {
		instance := v.(*pusher)
		if instance.Title != title {
			instance.Title = title
		}
		if instance.Content != content {
			instance.Content = content
		}
		return instance
	}
	instance := &pusher{Token: token, Title: title, Content: content}
	cache.Store(token, instance)
	return instance
}

func (b pusher) Name() string {
	return "PushPlus"
}

func (b pusher) Send() error {
	resp, err := req.C().R().
		SetHeader("Content-Type", "application/json; charset=utf-8").
		SetBody(json.MustEncode(b)).
		Send(http.MethodPost, URL)
	if err != nil {
		return errs.Wrap(code.RequestFailed, err)
	}

	if resp.StatusCode != http.StatusOK {
		return errs.WithMessage(code.ResponseError, fmt.Sprintf("%d %s", resp.StatusCode, resp.String()))
	}
	return nil
}
