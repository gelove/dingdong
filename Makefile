.PHONY: start build

NOW = $(shell date -u '+%Y%m%d%I%M%S')

APP = dingdong
MAIN = ./cmd/dingdong
CONFIG = config.json
RELEASE_MAC = release/macOS
RELEASE_MAC_APP = ${RELEASE_MAC}/${APP}
RELEASE_WIN = release/windows
RELEASE_WIN_APP = ${RELEASE_WIN}/${APP}.exe
RELEASE_LINUX = release/linux
RELEASE_LINUX_APP = ${RELEASE_LINUX}/${APP}

all: start

generate:
	go generate ./...

build-mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -a -o ${RELEASE_MAC_APP} ${MAIN} && upx ./${RELEASE_MAC_APP}

build-win: generate
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o ${RELEASE_WIN_APP} ${MAIN} && upx ${RELEASE_WIN_APP}

build-linux: generate
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ${RELEASE_LINUX_APP} ${MAIN} && upx ${RELEASE_LINUX_APP}

start: generate
	go run -race ${MAIN}

test:
	go test -v ./...
