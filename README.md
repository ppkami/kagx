# kagx

kagx 一个可用于内网穿透的反向代理应用，目前支持 tcp, http 协议。

## 开发状态

目前为学习项目，供研究学习，慎用于生产环境中

## 创建

1. 将项目git clone到本地$GOPATH环境中

```
$ git clone https://github.com/ppkami/kagx.git
```
2. 编译

```
$ cd $GOPATH/src/github.com/ppkami/kagx
$ make build
```

`$GOPATH/src/github.com/ppkami/kagx/bin`目录下会生成两个执行文件, 分别是服务端应用`kags`和客户端应用`kagc`

3. 配置文件

文件位于`$GOPATH/src/github.com/ppkami/kagx/conf`目录下

服务器配置文件`kagxs.ini`

```
port=9595##服务器监管服务端口
token=kagx##验证码，用于数据安全校验
```

客户端配置文件`kagxc.ini`

```
ip=x.x.x.x##服务器ip
supervise_port=9595##服务器监管服务端口
token=kagx##验证码，用于数据安全校验

##外网用户访问的x.x.x.x:6000，数据将会转发到客户端本地127.0.0.1:22
[ssh]
remote_port=6000##服务器代理端口
local_ip=127.0.0.1##客户端本地代理ip，可以是客户端本机，或者局域网的某台服务器
local_port=22##客户端本地代理端口
```

4. 执行服务端

将文件`bin/kagxs`和`conf/kagxs.ini`传到服务器x.x.x.x, 登陆服务器x.x.x.x后执行

```
$ YOUR_BIN_PATH/kagxs -c YOUR_CONFIGURE_PATH/kagxs.ini
```

5. 执行客户端

将文件`bin/kagxc`和`conf/kagxc.ini`放到你本地主机上，并在本地主机上执行

```
$ YOUR_BIN_PATH/kagxc -c YOUR_CONFIGURE_PATH/kagxc.ini
```

6. 外网访问

本实例为ssh，所以执行

```
$ ssh -p 6000 username@x.x.x.x
```

其中username为你本地主机的登陆用户名，x.x.x.x为线上服务器ip，6000为线上服务器代理端口并转发本地主机ssh端口22
