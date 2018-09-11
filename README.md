# kagx

kagx 一个可用于内网穿透的反向代理应用，目前支持 tcp 协议。

## 开发状态

目前为学习项目，供研究学习，慎用于生产环境中

## 配置文件

* 服务端kagxs.ini

    ```
    # 服务器监管服务端口
    port=40000
    # 验证码，用于数据安全校验
    token=kagx
    ```

* 客户端kagxc.ini

    ```
    # 服务器ip
    ip=x.x.x.x
    # 服务器监管服务端口
    supervise_port=40000
    # 验证码，用于数据安全校验
    token=kagx

    # 外网用户访问的x.x.x.x:30000，数据将会转发到客户端本地127.0.0.1:22
    [ssh]
    # 服务器代理端口
    remote_port=30000
    # 客户端本地代理ip，可以是客户端本机，或者局域网的某台服务器
    local_ip=127.0.0.1
    # 客户端本地代理端口
    local_port=22
    ```

## 快速搭建

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
    $ docker run --name kagxs -d -p 40000:40000/udp -p 30000-30010:30000-30010 -v /usr/local/kagx/conf/kagxs.ini:/usr/local/kagx/conf/kagxs.ini pjy20050506/kagx:server-0.0.1
    ```

2. 客户端

    在本地客户端创建配置文件`/usr/local/kagx/conf/kagxc.ini`，参照[客户端配置文件](#配置文件)

    ```
    $ docker run --name kagxc -d -P -v /usr/local/kagx/conf/kagxc.ini:/usr/local/kagx/conf/kagxc.ini pjy20050506/kagx:client-0.0.1
    ```

    _注意：_ kagxc.ini配置文件中local_ip不要写成127.0.0.1，因为执行kagx客户端是在docker容器中，应将127.0.0.1改成客户端主机在局域网中的ip地址，比如192.168.1.105

6. 外网访问

    本实例为ssh，所以执行

    ```
    $ ssh -p 30000 username@x.x.x.x
    ```

    其中username为你本地主机的登陆用户名，x.x.x.x为线上服务器ip，30000为线上服务器代理端口并转发本地主机ssh端口22

## http web站点搭建示例

基于kagx搭建http内网穿透，示例以[ghost](https://github.com/TryGhost/Ghost)为web项目，配合nginx，搭建外网能访问的博客网站

1. 客户端搭建ghost

    参考[ghost](https://github.com/TryGhost/Ghost)，快速搭建`ghost install local`，完成后确保在客户端能访问http://127.0.0.1:2368

2. 客户端kagxc.ini配置

    ```
    # 服务器ip
    ip=x.x.x.x
    # 服务器监管服务端口
    supervise_port=40000
    # 验证码，用于数据安全校验
    token=kagx

    # 外网用户访问的x.x.x.x:30001，数据将会转发到客户端本地127.0.0.1:2368
    [ghost]
    # 服务器代理端口
    remote_port=30001
    # ghost站点ip
    local_ip=127.0.0.1
    # ghost服务端口
    local_port=2368

    ```

    执行kagxc

    ```
    $ /usr/local/kagx/kagxc -c /usr/local/kagx/conf/kagxc.ini
    ```

3. 服务器端kagxs.ini配置

    ```
    # 服务器监管服务端口
    port=40000
    # 验证码，用于数据安全校验
    token=kagx
    ```

    执行kagxs

    ```
    $ /usr/local/kagx/kagxs -c /usr/local/kagx/conf/kagxs.ini
    ```

    假设你的线上服务器ip是x.x.x.x, 这是你访问http://x.x.x.x:30001, 将会是转发客户端ghost站点

4. 服务器nginx配置

    配合nginx使用，外网则能通过域名访问站点

    ```
    server {
        listen 80;
        server_name ghost.domain.com;


        location / {
            proxy_pass http://127.0.0.1:30001;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        }
    }
    ```

    重置nginx后，访问`http://ghost.domain.com/`将会是ghost站点
