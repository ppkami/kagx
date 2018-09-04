// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ppkami/kagx/netmsg"
	log "github.com/sirupsen/logrus"
)

//代理服务器端口使用状态
const (
	PROXY_OWNER int = iota //当前请求客户端匹配的客户单
	PROXY_EXIST            //端口已被其他客户端占用，无法与当前请求的客户端匹配，客户单需要重新更换新的代理端口
	PROXY_FREE             //端口未被使用
)

//代理服务器
type Proxy struct {
	port          uint16        //端口
	token         []byte        //验证码
	identity      []byte        //与客户端绑定唯一
	heartbeatTime time.Time     //代理心跳记录时间
	visitor       chan int      //外网用户访问记录
	forwardConn   chan net.Conn //客户端连接代理
	connMux       int           //允许用户最大请求（防止客户端非法请求，造成无限循环请求代理服务器）
}

var proxys map[uint16]*Proxy

func init() {
	proxys = make(map[uint16]*Proxy)
}

//创建代理服务器
func NewProxy(token []byte, identity []byte, port uint16) *Proxy {
	var p = new(Proxy)
	p.port = port
	p.identity = identity
	p.visitor = make(chan int, 20)
	p.forwardConn = make(chan net.Conn, 20)
	p.token = token
	p.heartbeatTime = time.Now()
	p.connMux = 10

	proxys[port] = p
	return p
}

//获取代理服务器信息
func GetProxy(identity []byte, port uint16) (*Proxy, int) {
	if p, ok := proxys[port]; ok {
		if bytes.Equal(p.identity, identity) {
			return p, PROXY_OWNER
		}

		return p, PROXY_EXIST
	}

	return nil, PROXY_FREE
}

//外网用户请求代理服务器记录
func (p *Proxy) Visitor() (chan int, error) {
	if p, ok := proxys[p.port]; ok {
		return p.visitor, nil
	}

	return nil, errors.New("代理服务器未启动")
}

//启动代理服务器
func (p *Proxy) Run() error {
	port := p.port

	addr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", net.IPv4zero, port))
	listener, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.WithFields(log.Fields{
			"port":  port,
			"error": err,
		}).Error("TCP代理服务器启动失败")
		return err
	}

	log.WithFields(log.Fields{
		"port": port,
	}).Info("TCP代理服务器启动成功")

	//长时间没有收到客户端心跳包，则停止代理服务器
	go func(l *net.TCPListener, p *Proxy) {
		for {
			<-time.Tick(30 * time.Second)
			hb := p.heartbeatTime.Add(10 * time.Second)
			if time.Now().After(hb) {
				l.Close()
				delete(proxys, p.port)

				log.WithFields(log.Fields{
					"port": p.port,
				}).Info("客户端无响应，关闭代理服务")

				break
			} else {
				continue
			}
		}
	}(listener, p)

	//处理用户请求
	go func(l *net.TCPListener) {
		connMuxChan := make(chan int, p.connMux)
		for i := 0; i < p.connMux; i++ {
			connMuxChan <- 1
		}

		for {
			conn, err := listener.AcceptTCP()
			<-connMuxChan
			if err != nil {
				return
			}

			go p.handleConn(conn, connMuxChan)
		}
	}(listener)

	return nil
}

//处理用户请求与客户端服务转发
func (p *Proxy) handleConn(conn net.Conn, connMuxChan chan int) {
	//获取请求数据判断是客户端转发请求，还是外网用户请求
	mark := make([]byte, 25)
	conn.Read(mark)

	//客户端转发连接
	if netmsg.IsForwardConn(mark, p.token) {
		connMuxChan <- 1
		connMuxChan <- 1
		p.forwardConn <- conn
		//外网用户请求
	} else {
		//----标记外网用户访问，通知客户端连接代理服务器----//
		p.visitor <- 1
		//----获取客户端的代理请求，进行数据交换----//
		forwardConn := <-p.forwardConn
		forwardConn.Write(mark) //预读的验证数据发送给客户端，以免用户请求数据发送给客户端缺失
		var wait sync.WaitGroup
		wait.Add(2)
		go pipe(conn, forwardConn, wait)
		go pipe(forwardConn, conn, wait)
		wait.Wait()
	}
}

//更新心跳记录最新时间
func (p *Proxy) UpdateHeartbeatTime() {
	p.heartbeatTime = time.Now()
}

//转发连接通道数据
func pipe(dst net.Conn, src net.Conn, wait sync.WaitGroup) {
	defer dst.Close()
	defer src.Close()
	defer wait.Done()

	io.Copy(dst, src)
}
