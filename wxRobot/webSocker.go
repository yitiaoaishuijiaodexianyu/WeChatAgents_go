package wxRobot

import (
	"WeChatAgents_go/config"
	_struct "WeChatAgents_go/struct"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

var configInfo = config.GetConfigInfo()

var webSocketClientUrl = "wss://" + configInfo.SerciveHose + "/ws/" + configInfo.BotWxid + "/" + configInfo.SecurityCode

var conn *websocket.Conn

func websocketConn() {
	var err error
	conn, _, err = websocket.DefaultDialer.Dial(webSocketClientUrl, nil)
	if err != nil {
		fmt.Println("正在重连....")
		time.Sleep(time.Second * 5)
		// 进行重连
		websocketConn()
	}
	fmt.Printf("WechatAgents_go_client启动成功\n")
}

func WebSocketClientStart() {
	websocketConn()
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for {
		_, message, err := conn.ReadMessage()
		// 放开这个注释就能看到原始消息
		//fmt.Printf(string(message))
		if err != nil {
			websocketConn()
			return
		}
		// 这里是调用一些别的东西异步返回的处理
		if strings.Contains(string(message), "CgiBaseResponse") {
			go CgiResponseProcess(message, conn)
			continue
		}
		// 这里是消息的处理
		var messages _struct.Message
		if err := json.Unmarshal(message, &messages); err != nil {
			continue
		}
		go MessageProcess(messages, conn)
	}
}
