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
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	log.Println(dir)
	config.Initialize(dir + "/config.json")
	session.Initialize(dir + "/sign.js")

	app.Run()
}
