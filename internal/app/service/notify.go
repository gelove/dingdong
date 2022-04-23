package service

import (
	"log"

	"dingdong/pkg/notify/bark"
	"dingdong/pkg/notify/pushplus"
)

func Push(token, detail string) {
	if token == "" {
		return
	}
	icon := "https://images.liqucn.com/img/h1/h968/img201709191609410_info300X300.jpg"
	b := bark.New(token, icon, "minuet")
	err := b.Send("买菜助手", detail)
	if err != nil {
		log.Printf("%s 推送失败 => %+v", b.Name(), err)
		return
	}
}

func PushToAndroid(token, detail string) {
	if token == "" {
		return
	}
	b := pushplus.New(token)
	err := b.Send("买菜助手", detail)
	if err != nil {
		log.Printf("%s 推送失败 => %+v", b.Name(), err)
		return
	}
}
