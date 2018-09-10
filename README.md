# kagx

kagx 一个可用于内网穿透的反向代理应用，目前支持 tcp, http 协议。

## 开发状态

目前为学习项目，供研究学习，慎用于生产环境中

## 配置文件

* 服务端kagxs.ini

    ```
    port=40000##服务器监管服务端口
    token=kagx##验证码，用于数据安全校验
    ```

* 客户端kagxc.ini

    ```
    ip=x.x.x.x##服务器ip
    supervise_port=40000##服务器监管服务端口
    token=kagx##验证码，用于数据安全校验

    ##外网用户访问的x.x.x.x:30000，数据将会转发到客户端本地127.0.0.1:22
    [ssh]
    remote_port=30000##服务器代理端口
    local_ip=127.0.0.1##客户端本地代理ip，可以是客户端本机，或者局域网的某台服务器
    local_port=22##客户端本地代理端口
    ```

## 使用

1. 将项目git clone到本地$GOPATH环境中

    ```
    $ git clone https://github.com/ppkami/kagx.git
    ```
2. 编译

    ```
    $ cd $GOPATH/src/github.com/ppkami/kagx
    $ make build
    ```

    `$GOPATH/src/github.com/ppkami/kagx/bin`目录下会生成两个执行文件, 分别是服务端应用`kags_*_*`和客户端应用`kagc_*_*`, 根据服务器操作系统使用对应的执行文件。你可重命名文件，比如你是linux 64操作系统，则可以执行

    ```
    $ mv kags_linux_amd64 kags
    $ mv kagc_linux_amd64 kagc
    ```

3. 服务端

    将文件`bin/kagxs`和`conf/kagxs.ini`传到服务器x.x.x.x, 登陆服务器x.x.x.x后执行

    ```
    $ /usr/local/kagx/kagxs -c /usr/local/kagx/conf/kagxs.ini
    ```

4. 客户端

    将文件`bin/kagxc`和`conf/kagxc.ini`放到你本地主机上，并在本地主机上执行

    ```
    $ /usr/local/kagx/kagxc -c /usr/local/kagx/conf/kagxs.ini
    ```

5. 外网访问

    本实例为ssh，所以执行

    ```
    $ ssh -p 30000 username@x.x.x.x
    ```

    其中username为你本地主机的登陆用户名，x.x.x.x为线上服务器ip，30000为线上服务器代理端口并转发本地主机ssh端口22

## 使用Docker

1. 服务器

    假设你的线上服务器IP是`x.x.x.x`， 登陆服务器后，创建配置文件`/usr/local/kagx/conf/kagxs.ini`，参照[服务端配置文件](#配置文件)

    ```
    $ docker run --name kagxs -d -p 40000:40000/udp -v /usr/local/kagx/conf/kagxs.ini:/usr/local/kagx/conf/kagxs.ini pjy20050506/kagx:server-0.0.1
    ```

2. 客户端

    在本地客户端创建配置文件`/usr/local/kagx/conf/kagxc.ini`，参照[客户端配置文件](#配置文件)

    ```
    $ docker run --name kagxc -d -P -v /usr/local/kagx/conf/kagxc.ini:/usr/local/kagx/conf/kagxc.ini pjy20050506/kagx:client-0.0.1
    ```

    ##注意：##kagxc.ini配置文件中local_ip不要写成127.0.0.1，因为执行kagx客户端是在docker容器中，应将127.0.0.1改成客户端主机在局域网中的ip地址，比如192.168.1.105

6. 外网访问

    本实例为ssh，所以执行

    ```
    $ ssh -p 30000 username@x.x.x.x
    ```

    其中username为你本地主机的登陆用户名，x.x.x.x为线上服务器ip，30000为线上服务器代理端口并转发本地主机ssh端口22
