export PATH := $(GOPATH)/bin:$(PATH)

LDFLAGS := -s -w

all: build

build: pkg kagxs kagxc

#服务器
kagxc:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/kagxc_darwin_amd64 ./cmd/kagxc
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/kagxc_linux_amd64 ./cmd/kagxc
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/kagxc_windows_amd64 ./cmd/kagxc
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/kagxc_linux_arm64 ./cmd/kagxc

#客户端
kagxs:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/kagxs_darwin_amd64 ./cmd/kagxs
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/kagxs_linux_amd64 ./cmd/kagxs
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/kagxs_windows_amd64 ./cmd/kagxs
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/kagxs_linux_arm64 ./cmd/kagxs

clean:
	rm -rf ./bin/*
#第三方包
pkg:
	go get -d gopkg.in/urfave/cli.v1
	go get -d gopkg.in/ini.v1
	go get -d github.com/sirupsen/logrus
	go get -d github.com/vmihailenco/msgpack
	go get -d github.com/satori/go.uuid

tar:
	tar -zcvf bin/kagxc_darwin_amd64.tar.gz -C bin kagxc_darwin_amd64
	tar -zcvf bin/kagxc_linux_amd64.tar.gz -C bin kagxc_linux_amd64
	tar -zcvf bin/kagxc_linux_arm64.tar.gz -C bin kagxc_linux_arm64
	tar -zcvf bin/kagxc_windows_amd64.tar.gz -C bin kagxc_windows_amd64
	tar -zcvf bin/kagxs_darwin_amd64.tar.gz -C bin kagxs_darwin_amd64
	tar -zcvf bin/kagxs_linux_amd64.tar.gz -C bin kagxs_linux_amd64
	tar -zcvf bin/kagxs_linux_arm64.tar.gz -C bin kagxs_linux_arm64
	tar -zcvf bin/kagxs_windows_amd64.tar.gz -C bin kagxs_windows_amd64
