export GO111MODULE=on
export GOPROXY=https://goproxy.io,direct
export GOOS=linux
export GIN_MODE=release
bUild: export GOARCH=amd64

BINARY=app

all: help

build: clean bCommand

run:	build
	./${BINARY}

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

bCommand:
	go build -o ${BINARY} .

help:
	@echo "make build - 编译 Go 代码"
	@echo "make run - 编译并运行 Go 代码"
	@echo "make clean - 移除二进制文件"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
