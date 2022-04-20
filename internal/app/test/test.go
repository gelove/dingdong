package test

import (
	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/session"
)

const configFile = "../../../config.json"
const jsFile = "../../../sign.js"

func init() {
	// dir, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(dir)
	config.Initialize(configFile)
	session.Initialize(jsFile)
}
