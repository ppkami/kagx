// Copyright 2018 kamigx Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

//服务端程序入口
package main

import (
	"os"
	"sort"

	"github.com/ppkami/kagx/server"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	var configFilePath string //配置文件路径

	//----配置命令行----//
	app := cli.NewApp()
	app.Name = "kagxs"
	app.Usage = "a reverse proxy server"
	//命令说明
	app.Flags = []cli.Flag{
		//配置文件路径
		cli.StringFlag{
			Name:        "config, c",
			Value:       "conf/kagxs.ini",
			Usage:       "the server configure",
			Destination: &configFilePath,
		},
	}
	//命令行为
	app.Action = func(c *cli.Context) error {
		//命令行参数输入大于1，则不执行命令，显示帮助信息
		if c.NArg() > 0 {
			cli.ShowAppHelp(c)
			return nil
		}
		//加载配置文件并执行
		server.LoadConfFile(configFilePath)
		server.Start()
		return nil
	}

	//按字母顺序排序自定义命令
	sort.Sort(cli.FlagsByName(app.Flags))
	//执行
	app.Run(os.Args)
}
