// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package netmsg

import (
	"bytes"
	"crypto/md5"
	"strconv"
	"time"
)

//验证是否为转发连接
func IsForwardConn(r []byte, token []byte) bool {
	//return false
	mark, _ := strconv.Atoi(string(r[:1]))

	if mark == VISITOR_REQUEST_FORWARD && len(r) == 25 {

		//过期验证
		rTime, err := strconv.ParseInt(string(r[1:9]), 16, 64)
		if err != nil {
			return false
		}
		if time.Now().UTC().Unix()-rTime > 60 {
			return false
		}

		tokenMd5 := md5.Sum(append(token, r[1:9]...))
		return bytes.Equal(tokenMd5[:], r[9:25])
	}

	return false
}

//生成转发请求数据前缀
func PrefixForwardConn(token []byte) []byte {
	now := time.Now().UTC().Unix()
	nowStr := strconv.FormatInt(now, 16)
	nowBuf := []byte(nowStr)

	tokenMd5 := md5.Sum(append(token, nowBuf...))

	prefix := []byte(strconv.Itoa(VISITOR_REQUEST_FORWARD))

	prefix = append(prefix, nowBuf...)
	prefix = append(prefix, tokenMd5[:]...)

	return prefix
}
