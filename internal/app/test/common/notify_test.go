package common

import (
	"testing"

	"dingdong/internal/app/config"
	"dingdong/internal/app/service/notify"
	"dingdong/pkg/notify/player"
)

func init() {
	config.Initialize("../../../../config.yml")
}

func TestBark(t *testing.T) {
	conf := config.Get()
	if len(conf.Bark) > 0 {
		notify.Push(conf.Bark[0], "Bark 测试")
	}
}

func TestPushPlus(t *testing.T) {
	conf := config.Get()
	if len(conf.PushPlus) > 0 {
		notify.PushPlus(conf.PushPlus[0], "PushPlus 测试")
	}
}

func TestPlayMp3(t *testing.T) {
	p := player.Default()
	err := p.Send()
	if err != nil {
		t.Error(err)
	}
}
