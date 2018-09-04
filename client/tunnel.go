// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package client

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/ppkami/kagx/netmsg"
	log "github.com/sirupsen/logrus"
)

//隧道
type Tunnel struct {
	token []byte //验证码
	Proxy        //代理服务器
}

//新建隧道
func NewTunnel(p Proxy, token []byte) *Tunnel {
	var t = new(Tunnel)
	t.Proxy = p
	t.token = token

	return t
}

//启动本地与代理数据转发
func (t *Tunnel) Foward() {
	proxyIp := t.Proxy.RemoteIP
	proxyPort := t.Proxy.RemotePort
	localIp := t.Proxy.LocalIP
	localPort := t.Proxy.LocalPort
	token := t.token

	//连接代理服务器
	addr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", proxyIp, proxyPort))
	conn, err := net.DialTCP("tcp4", nil, addr)

	if err != nil {
		log.WithFields(log.Fields{
			"ip":   proxyIp,
			"port": proxyPort,
		}).Warn("连接远程代理服务器失败")
		return
	}

	//向代理服务器发送信息，表示该请求为客户端代理
	request := netmsg.PrefixForwardConn(token)
	conn.Write(request)

	//连接本地服务器
	laddr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", localIp, localPort))
	lconn, err := net.DialTCP("tcp4", nil, laddr)

	if err != nil {
		log.WithFields(log.Fields{
			"ip":   localIp,
			"port": localPort,
		}).Panic("连接本地服务器失败")
		return
	}

	//数据交换
	go func(c1 net.Conn, c2 net.Conn) {
		var wait sync.WaitGroup
		wait.Add(2)
		go pipe(c1, c2, wait)
		go pipe(c2, c1, wait)
		wait.Wait()
	}(conn, lconn)

}

func pipe(dst net.Conn, src net.Conn, wait sync.WaitGroup) {
	defer dst.Close()
	defer src.Close()
	defer wait.Done()

	io.Copy(dst, src)
}
