// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

//服务器管理
package server

//启动服务器
func Start() {
	//获取配置信息
	conf := GetConf()
	//----启动监管服务器----//
	supervise := GenSupervise(conf)
	supervise.Run()
}
