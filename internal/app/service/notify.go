package service

import (
	"log"

	"dingdong/pkg/notify/bark"
)

func Push(key, detail string) {
	icon := "https://images.liqucn.com/img/h1/h968/img201709191609410_info300X300.jpg"
	b := bark.New(key, icon, "minuet")
	err := b.Send("买菜助手", detail)
	if err != nil {
		log.Printf("Bark 推送失败 => %+v", err)
		return
	}
}
