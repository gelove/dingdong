package test

import (
	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/session"
)

const configFile = "../../../config.json"

func init() {
	config.Initialize(configFile)
	// session.Initialize()
	session.InitializeMock()
}
