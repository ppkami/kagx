// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package client

import (
	"fmt"

	"gopkg.in/ini.v1"
)

//客户端配置项
type Conf struct {
	IP     string           //服务器IP
	Token  []byte           //验证码
	Port   uint16           //远程连接监管服务器端口
	Proxys map[uint16]Proxy //代理服务器组
}

//代理服务器配置
type Proxy struct {
	Name       string //服务器配置自定义名称
	RemoteIP   string //远程代理服务器Ip
	RemotePort uint16 //远程代理服务器端口
	LocalIP    string //本地被代理IP
	LocalPort  uint16 //本地被代理端口
}

var c *Conf //配置

func init() {
	//初始化配置信息
	c = &Conf{
		IP:     "127.0.0.1",
		Token:  []byte("kagx"),
		Port:   9000,
		Proxys: make(map[uint16]Proxy),
	}
}

//通过配置文件设置服务器配置
func LoadConfFile(filePath string) *Conf {
	//----加载配置文件，获取配置信息----//
	cfg, err := ini.Load(filePath)
	if err != nil {
		panic(fmt.Errorf("读取服务器配置文件错误: %v\n", err))
	}

	//配置监管服务器
	if cfg.Section("").HasKey("ip") {
		c.IP = cfg.Section("").Key("ip").String()
	}
	if cfg.Section("").HasKey("token") {
		c.Token = []byte(cfg.Section("").Key("token").String())
	}
	if cfg.Section("").HasKey("supervise_port") {
		rSPortInt, _ := cfg.Section("").Key("supervise_port").Int()
		c.Port = uint16(rSPortInt)
	}

	//配置代理服务器
	var (
		proxyName     string //代理简称
		remotePortInt int    //代理端口
		remoteIp      string //代理IP
		remotePort    uint16 //代理端口
		localPortInt  int    //本地端口
		localPort     uint16 //本地端口
		localIp       string //本地ip
		rPConfs       map[uint16]Proxy
	)
	rPConfs = make(map[uint16]Proxy)
	for _, section := range cfg.Sections() {
		proxyName = section.Name()
		if proxyName == "DEFAULT" {
			continue
		}

		remoteIp = cfg.Section("").Key("ip").String()

		remotePortInt, _ = section.Key("remote_port").Int()
		remotePort = uint16(remotePortInt)

		localPortInt, _ = section.Key("local_port").Int()
		localPort = uint16(localPortInt)

		localIp = section.Key("local_ip").String()

		rPConfs[remotePort] = Proxy{
			Name:       proxyName,
			RemoteIP:   remoteIp,
			RemotePort: remotePort,
			LocalIP:    localIp,
			LocalPort:  localPort,
		}
	}

	c.Proxys = rPConfs

	return c
}

//获取配置信息
func GetConf() *Conf {
	return c
}
