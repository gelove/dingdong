package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"dingdong/internal/app"
	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/session"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	defer func() {
		if v := recover(); v != nil {
			log.Printf("[严重错误]: %+v", v)
		}
	}()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Println(dir)
	config.Initialize(dir + "/config.yml")
	session.Initialize()

	app.Run()
}
