export PATH := $(GOPATH)/bin:$(PATH)

build: pkg kagxs kagxc

#服务器
kagxc:
	go build -o bin/kagxc ./cmd/kagxc

#客户端
kagxs:
	go build -o bin/kagxs ./cmd/kagxs

#第三方包
pkg:
	go get -d gopkg.in/urfave/cli.v1
	go get -d gopkg.in/ini.v1
	go get -u github.com/sirupsen/logrus
	go get -d golang.org/x/crypto/ssh/terminal
	go get -u github.com/vmihailenco/msgpack
	go get -d github.com/satori/go.uuid
