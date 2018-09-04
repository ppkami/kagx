// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

//客户端管理
package client

import (
	log "github.com/sirupsen/logrus"
)

//启动客户端
func Start() {
	log.Info("客户端启动...")
	//获取配置信息
	conf := GetConf()
	//----调度代理服务器----//
	//请求远程监管服务器
	app := ApplyService(conf)
	//请求分发代理服务器
	app.Run()
}
