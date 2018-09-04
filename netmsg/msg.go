// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

//请求服务器数据处理
//访问监管服务器数据、代理服务数据转发相关处理
package netmsg

import (
	"bytes"
	"crypto/md5"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/vmihailenco/msgpack"
)

//消息监听
type Msg struct {
	token  []byte                        //验证码
	conn   *net.UDPConn                  //与远程UDP服务连接
	routes map[int]func(*Respone, error) //UDP请求路由
}

//远程服务响应信息
type Respone struct {
	Token      []byte       //验证码
	Conn       *net.UDPConn //与远程UDP服务连接
	RemoteAddr *net.UDPAddr //远程UDP服务地址
	Msg        []byte       //消息
}

//创建新消息
func New(token []byte, conn *net.UDPConn) *Msg {
	return &Msg{
		token:  token,
		conn:   conn,
		routes: make(map[int]func(*Respone, error)),
	}
}

//设置路由
func (m *Msg) GET(code int, callback func(*Respone, error)) {
	m.routes[code] = callback
}

//启动消息监听
func (m *Msg) Run() {
	data := make([]byte, 1024)
	for {
		n, raddr, err := m.conn.ReadFromUDP(data)
		if err != nil {
			continue
		}
		respone := &Respone{
			Token:      m.token,
			Conn:       m.conn,
			RemoteAddr: raddr,
			Msg:        nil,
		}

		//获取路由和加密数据
		code, _ := strconv.Atoi(string(data[:1]))
		msg := data[1:n]
		callback := m.routes[code]
		//过期验证
		rTime, err := strconv.ParseInt(string(msg[:8]), 16, 64)
		if err != nil {
			callback(respone, err)
			continue
		}
		if time.Now().UTC().Unix()-rTime > 60 {
			callback(respone, errors.New("the token is expire"))
			continue
		}

		//解密验证
		tokenMd5 := md5.Sum(append(m.token, msg[:8]...))
		if !bytes.Equal(tokenMd5[:], msg[8:24]) {
			callback(respone, errors.New("the token is illegal"))
			continue
		}
		//解密数据
		decode := xor(msg[24:], m.token)
		respone.Msg = decode

		callback(respone, nil)
	}
}

//封装消息格式
func MountMsg(code int, token []byte, claim interface{}) ([]byte, error) {
	now := time.Now().UTC().Unix()
	nowStr := strconv.FormatInt(now, 16)
	nowBuf := []byte(nowStr)

	tokenMd5 := md5.Sum(append(token, nowBuf...))
	data, err := msgpack.Marshal(claim)
	if err != nil {
		return nil, err
	}
	encrypt := xor(data, token)
	msg := append(nowBuf, tokenMd5[:]...)
	msg = append(msg, encrypt...)

	pre := []byte(strconv.Itoa(code))
	msg = append(pre, msg...)

	return msg, err
}

func xor(msg []byte, key []byte) []byte {
	en := make([]byte, len(msg))
	kl := len(key)
	for i, v := range msg {
		en[i] = v ^ key[i%kl]
	}

	return en
}
