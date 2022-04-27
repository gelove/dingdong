package notify

import (
	"log"

	"dingdong/pkg/notify"
	"dingdong/pkg/notify/bark"
	"dingdong/pkg/notify/player"
	"dingdong/pkg/notify/pushplus"
)

func send(n notify.Notifier) {
	err := n.Send()
	if err != nil {
		log.Printf("%s 通知失败 => %+v", n.Name(), err)
	}
}

func Push(token, detail string) {
	if token == "" {
		return
	}
	icon := "https://images.liqucn.com/img/h1/h968/img201709191609410_info300X300.jpg"
	send(bark.New(token, "买菜助手", detail, icon, "minuet"))
}

func PushPlus(token, detail string) {
	if token == "" {
		return
	}
	send(pushplus.New(token, "买菜助手", detail))
}

func PlayMp3() {
	send(player.Default())
}
