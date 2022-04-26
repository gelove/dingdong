package notify

import (
	"log"

	"dingdong/pkg/notify/bark"
	"dingdong/pkg/notify/player"
	"dingdong/pkg/notify/pushplus"
)

func Push(token, detail string) {
	if token == "" {
		return
	}
	icon := "https://images.liqucn.com/img/h1/h968/img201709191609410_info300X300.jpg"
	b := bark.New(token, "买菜助手", detail, icon, "minuet")
	err := b.Send()
	if err != nil {
		log.Printf("%s 推送失败 => %+v", b.Name(), err)
	}
}

func PushToAndroid(token, detail string) {
	if token == "" {
		return
	}
	b := pushplus.New(token, "买菜助手", detail)
	err := b.Send()
	if err != nil {
		log.Printf("%s 推送失败 => %+v", b.Name(), err)
	}
}

func PlayMp3() {
	p := player.Default()
	err := p.Send()
	if err != nil {
		log.Printf("%s 播放失败 => %+v", p.Name(), err)
	}
}
