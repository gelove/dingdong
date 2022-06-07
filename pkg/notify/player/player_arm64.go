package player

import (
	"dingdong/pkg/notify"
)

type Player struct {
	Audio string
}

func New(audioName string) notify.Notifier {
	return &Player{
		Audio: audioName,
	}
}

func Default() notify.Notifier {
	return New("audio/order.mp3")
}

func (p *Player) Name() string {
	return "Mp3 Player"
}

func (p *Player) Send() error {
	return nil
}
