// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package netmsg

//消息监听类型
const (
	PROXY_SERVER            int = iota //启动代理服务器
	SUCCESS_START_PROXY                //成功启动代理服务
	PROXY_VALIDATE_FAIL                //验证错误
	VISITOR_REQUEST_PROXY              //用户访问代理服务
	VISITOR_REQUEST_FORWARD            //用户代理转发
	PROXY_PORT_EXIST                   //代理服务器端口已经占用
)
