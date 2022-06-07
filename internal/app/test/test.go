package test

import (
	"dingdong/internal/app/config"
	"dingdong/internal/app/pkg/ddmc/ios_session"
)

func init() {
	config.Initialize("../../../config.yml")
	ios_session.InitializeMock("../../../session", "金山")
}
