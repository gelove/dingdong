package service

import (
	"fmt"
	"log"

	"dingdong/pkg/notify/bark"
)

func Push(key, products, ellipsis string) {
	b := bark.New(key, "https://pica.zhimg.com/v2-96c09d84051809c3f61797790f9177c9_im.jpg", "minuet")
	err := b.Send("买菜助手", fmt.Sprintf("叮咚买菜当前可配送请尽快下单[%s%s]", products, ellipsis))
	if err != nil {
		log.Printf("Bark 推送失败 => %+v", err)
		return
	}
	log.Printf("Bark 推送成功 => %s", key)
}
