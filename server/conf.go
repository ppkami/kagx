// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package server

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type Conf struct {
	*SuperviseServer
}

//监管服务器配置
type SuperviseServer struct {
	Token []byte
	Port  uint16 //端口
}

var c *Conf

func init() {
	//初始化配置信息
	c = &Conf{
		&SuperviseServer{
			Token: []byte("kagx"),
			Port:  9000,
		},
	}
}

//通过配置文件设置服务器
func LoadConfFile(filePath string) *Conf {
	//----加载配置文件，获取配置信息----//
	cfg, err := ini.Load(filePath)
	if err != nil {
		panic(fmt.Errorf("读取服务器配置文件错误: %v\n", err))
	}

	//配置监管服务器
	sPort := c.Port
	sToken := c.Token

	if cfg.Section("").HasKey("port") {
		sPortInt, _ := cfg.Section("").Key("port").Int()
		sPort = uint16(sPortInt)
	}

	if cfg.Section("").HasKey("token") {
		sTokenStr := cfg.Section("").Key("token").String()
		sToken = []byte(sTokenStr)
	}

	c.Port = sPort
	c.Token = sToken

	return c
}

//获取配置信息
func GetConf() *Conf {
	return c
}
