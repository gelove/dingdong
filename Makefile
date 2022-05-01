.PHONY: start generate build compress build-upx test pack

# 当前年月日时分秒
NOW = $(shell date '+%Y%m%d%H%M%S')
# 目录名 通过shell获取
#DIR = $(shell basename ${CURDIR})
# 目录名 通过makefile内置函数获取
DIR = $(notdir ${CURDIR})

# 系统变量 darwin(默认macOS),linux,windows
OS ?= darwin
# 体系架构 amd64(默认),arm64
ARCH ?= amd64
CGO ?= 0
# 可执行文件扩展名 如果是windows系统则为.exe, 否则为空
# ${OS:windows=} 表示将变量中的字符串windows替换为空, 如果此时OS值为windows则表达式返回false
EXT = $(if ${OS:windows=},,.exe)
# 目录名作为应用名称
APP = dingdong
# 可执行文件名称
APP_EXE = ${APP}${EXT}
MAIN = ./cmd/${APP}
CONFIG = config.yml
RELEASE_DIR = release
RELEASE_OS = ${RELEASE_DIR}/${OS}-${ARCH}
RELEASE_APP = ${RELEASE_OS}/${APP_EXE}

all: start

.PHONY echo:
	@echo ${CURDIR} $(dir ${CURDIR}) $(notdir ${CURDIR}) ${EXT}

generate:
	go generate ./...

build:
	CGO_ENABLE=${CGO} GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-s -w" -a -o ${RELEASE_APP} ${MAIN}

compress:
	upx ${RELEASE_APP}

build-upx: generate build compress

start:
	go run -race ${MAIN}

test:
	go test -v ./...

pack: build-upx
	cp -r ${CONFIG} $(RELEASE_OS)
	cd $(RELEASE_DIR) && zip -r ${APP}-${OS}-${ARCH}-$(NOW).zip ${OS}-${ARCH}
