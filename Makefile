.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

APP = dingdong
MAIN = ./cmd/dingdong
CONFIG = config.json
RELEASE_WIN = release/windows
RELEASE_WIN_APP = ${RELEASE_WIN}/${APP}
RELEASE_LINUX = release/linux
RELEASE_LINUX_APP = ${RELEASE_LINUX}/${APP}

all: start

generate:
	go generate ./...

build:
	go build -ldflags "-s -w" -a -o ${APP} ${MAIN} && upx ./${APP}

build-win: generate
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${RELEASE_WIN_APP} ${MAIN} && upx ${RELEASE_WIN_APP}

build-linux: generate
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${RELEASE_LINUX_APP} ${MAIN} && upx ${RELEASE_LINUX_APP}

start: generate
	go run -race ${MAIN}

test:
	go test -v ./...
