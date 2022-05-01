package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		s := <-c
		log.Printf("Got signal: %+v", s)
		cancel()
	}()

	app.Run(ctx)
}
