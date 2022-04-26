package test

import (
	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/session"
)

const configFile = "../../../config.json"
const jsFile = "../../../sign.js"

func init() {
	config.Initialize(configFile)
	session.InitializeMock(jsFile)
}
