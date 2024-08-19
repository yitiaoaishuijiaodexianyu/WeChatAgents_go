package main

import (
	"WeChatAgents_go/config"
	"WeChatAgents_go/wxRobot"
)

func main() {
	// 初始化配置文件
	config.InitConfig()
	// 运行websocket
	go wxRobot.WebSocketClientStart()
	// 运行http管理端
	wxRobot.HttpRun()
}
