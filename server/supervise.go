// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package server

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/ppkami/kagx/netmsg"
	"github.com/vmihailenco/msgpack"
)

//监管服务器
type Supervise struct {
	SuperviseServer //监管服务器配置
}

//创建监管服务器
func GenSupervise(c *Conf) (s Supervise) {
	s.Token = c.SuperviseServer.Token
	s.Port = c.SuperviseServer.Port

	return
}

//启动监管服务器
func (s *Supervise) Run() {
	token := s.Token
	port := s.Port

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", net.IPv4zero, port))
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("打洞服务器启动失败")
		return
	}

	log.Info("启动监管服务器...")

	m := netmsg.New(token, conn)
	//监听启动代理服务器请求
	//客户端发送心跳包，用于检测代理服务的健康度，并执行相关发代理服务器操作
	m.GET(netmsg.PROXY_SERVER, func(r *netmsg.Respone, err error) {
		msg := r.Msg
		token := r.Token
		conn := r.Conn
		raddr := r.RemoteAddr

		if err != nil {
			msg, _ := netmsg.MountMsg(netmsg.PROXY_VALIDATE_FAIL, token, &netmsg.Proxy{
				Port: 0,
			})
			conn.WriteToUDP(msg, raddr)
			return
		}
		var proxyMsg netmsg.Proxy
		msgpack.Unmarshal(msg, &proxyMsg)
		//验证代理服务是否已经开启
		p, c := GetProxy(proxyMsg.Identity, proxyMsg.Port)
		//代理服务端口已经打开，则写入心跳验证包，用来判断客户端是否关闭，从而关闭代理端口
		if c == PROXY_OWNER {
			p.UpdateHeartbeatTime()
		}
		//代理端口已经被其他客户端占用
		if c == PROXY_EXIST {
			msg, _ := netmsg.MountMsg(netmsg.PROXY_PORT_EXIST, token, &netmsg.Proxy{
				Port: proxyMsg.Port,
			})
			conn.WriteToUDP(msg, raddr)
		}
		//启动代理服务器
		if c == PROXY_FREE {
			go doRunProxyServer(conn, raddr, proxyMsg.Identity, proxyMsg.Port, token)
		}
	})
	m.Run()
}

//开启代理服务器
func doRunProxyServer(conn *net.UDPConn, caddr *net.UDPAddr, identity []byte, port uint16, token []byte) {
	proxy := NewProxy(token, identity, port)
	//TODO 端口占用导致启动服务失败暂不处理
	err := proxy.Run()

	if err != nil {
		msg, _ := netmsg.MountMsg(netmsg.PROXY_PORT_EXIST, token, &netmsg.Proxy{
			Port: port,
		})
		conn.WriteToUDP(msg, caddr)
		return
	}

	msg, _ := netmsg.MountMsg(netmsg.SUCCESS_START_PROXY, token, &netmsg.Proxy{
		Port: port,
	})
	conn.WriteToUDP(msg, caddr)

	//监听用户访问客户端行为
	v, err := proxy.Visitor() //获取外网访客通知
	if err != nil {
		return
	}
	//通知客户端连接代理服务器
	for {
		<-v
		msg, _ := netmsg.MountMsg(netmsg.VISITOR_REQUEST_PROXY, token, &netmsg.Proxy{
			Port: port,
		})
		conn.WriteToUDP(msg, caddr)
	}

}
