// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package client

import (
	"fmt"
	"net"
	"time"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"

	"github.com/ppkami/kagx/netmsg"
)

//远程主机
type App struct {
	*Conf                 //客户单配置信息
	identity []byte       //身份唯一标识
	conn     *net.UDPConn //与远程监管主机连接
}

//远程管理器申请服务
func ApplyService(c *Conf) *App {
	a := new(App)
	a.identity = uuid.Must(uuid.NewV4()).Bytes()
	a.Conf = c

	return a
}

//连接远程监管服务器
func (a *App) connectSuperise() {
	//获取配置信息
	ip := a.IP
	port := a.Port

	//连接监管服务器
	raddr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	conn, _ := net.DialUDP("udp", nil, raddr)

	a.conn = conn
}

//请求监管服务器，获取代理服务器资源
func (a *App) dispatchProxy() {
	//获取配置信息
	token := a.Token
	proxys := a.Proxys
	conn := a.conn

	//定时请求远程服务器开启代理服务
	//目的：
	//1.远程服务挂起，重启后接收到客户端的心跳请求，重开代理服务
	//2.如果客户端挂起，服务端在规定时间内没接收到客户单心跳包，则关闭对应端口代理服务器
	for _, p := range proxys {
		var send = func(identity []byte, port uint16) {
			msg, _ := netmsg.MountMsg(netmsg.PROXY_SERVER, token, &netmsg.Proxy{
				Identity: identity,
				Port:     port,
			})
			conn.Write(msg)
		}
		send(a.identity, p.RemotePort)
		go func(identity []byte, port uint16) {
			for {
				<-time.Tick(10 * time.Second)
				send(identity, port)
			}
		}(a.identity, p.RemotePort)
	}
}

//接收监管服务器下发通知
func (a *App) accept() {
	//获取配置信息
	proxys := a.Proxys
	token := a.Token
	conn := a.conn
	//监听监管服务器下发信息
	m := netmsg.New(token, conn)

	//服务端成功开启代理服务通知
	m.GET(netmsg.SUCCESS_START_PROXY, func(r *netmsg.Respone, err error) {
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Info("启动代理服务失败...")
			return
		}
		var proxy netmsg.Proxy
		msg := r.Msg
		err = msgpack.Unmarshal(msg, &proxy)
		if err != nil {
			return
		}

		port := proxy.Port
		log.WithFields(log.Fields{
			"port": port,
		}).Info("成功启动代理服务...")
	})
	//验证错误通知
	m.GET(netmsg.PROXY_VALIDATE_FAIL, func(r *netmsg.Respone, err error) {
		log.Warn("启动代理请求验证错误...请检查token配置是否与服务器一致...")
	})
	//代理服务器端口已经被占用
	m.GET(netmsg.PROXY_PORT_EXIST, func(r *netmsg.Respone, err error) {
		var proxy netmsg.Proxy
		msg := r.Msg
		err = msgpack.Unmarshal(msg, &proxy)
		if err != nil {
			return
		}

		port := proxy.Port
		log.WithFields(log.Fields{
			"port": port,
		}).Warn("代理服务器端口已经被占用")
	})
	//外网用户访问代理服务
	m.GET(netmsg.VISITOR_REQUEST_PROXY, func(r *netmsg.Respone, err error) {
		var vistor netmsg.Proxy
		msg := r.Msg
		token := r.Token
		err = msgpack.Unmarshal(msg, &vistor)
		if err != nil {
			return
		}
		proxyPort := vistor.Port

		t := NewTunnel(proxys[proxyPort], token)
		t.Foward()
	})
	m.Run()
}

func (a *App) Run() {
	//连接监管服务器
	a.connectSuperise()

	//请求分发
	a.dispatchProxy()

	//处理服务器下发通知
	a.accept()
}
