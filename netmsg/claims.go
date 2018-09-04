// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package netmsg

//请求服务器启动代理服务器
type Proxy struct {
	Identity []byte
	Port     uint16
}
