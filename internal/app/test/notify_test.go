package test

import (
	"testing"

	"dingdong/internal/app/config"
	"dingdong/internal/app/service/notify"
	"dingdong/pkg/notify/player"
)

func TestPush(t *testing.T) {
	conf := config.Get()
	if len(conf.Users) > 0 {
		notify.Push(conf.Users[0], "测试")
	}
}

func TestPushToAndroid(t *testing.T) {
	conf := config.Get()
	if len(conf.AndroidUsers) > 0 {
		notify.PushToAndroid(conf.AndroidUsers[0], "测试")
	}
}

func TestPlayMp3(t *testing.T) {
	p := player.New("../../../audio.mp3")
	err := p.Send()
	if err != nil {
		t.Error(err)
	}
}
