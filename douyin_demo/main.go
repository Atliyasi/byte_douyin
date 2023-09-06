package main

import (
	"douyin_demo/dao"
	"douyin_demo/log"
	"douyin_demo/router"
)

func main() {
	// 初始化日志信息
	log.Init()
	// 初始化数据库信息
	dao.InitDB()
	// 初始化路由信息，并启动服务
	router.InitRouter()
}
