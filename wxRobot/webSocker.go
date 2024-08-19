package wxRobot

import (
	_struct "WeChatAgents_go/struct"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"strings"
	"time"
)

func websocketConn() {
	var webSocketClientUrl = "wss://" + _struct.Config.Robot[0].ServiceHose + "/ws/" + _struct.Config.Robot[0].BotWxid + "/" + _struct.Config.Robot[0].SecurityCode
	var err error
	_struct.WebSocketConn, _, err = websocket.DefaultDialer.Dial(webSocketClientUrl, nil)
	if err != nil {
		fmt.Println("正在重连....")
		time.Sleep(time.Second * 5)
		// 进行重连
		websocketConn()
	}
	fmt.Printf("WechatAgents_go_client启动成功\n")
}

func checkWebsocketConn() {
	for true {
		err := _struct.WebSocketConn.WriteMessage(1, []byte("你还在不"))
		if err != nil {
			fmt.Println("链接死了重连")
			websocketConn()
		}
		time.Sleep(time.Second * 5)
	}
}

func WebSocketClientStart() {
	websocketConn()
	defer func() {
		if err := _struct.WebSocketConn.Close(); err != nil {
			websocketConn()
		}
	}()
	// 获取已知群的群成员信息
	go GetKnownGroupInfo()
	// 主动检查websocket是不是死了
	go checkWebsocketConn()
	for {
		_, message, err := _struct.WebSocketConn.ReadMessage()
		// 放开这个注释就能看到原始消息
		//fmt.Printf(string(message))
		if err != nil {
			websocketConn()
			time.Sleep(time.Second * 5)
			continue
		}

		// 这里是调用一些别的东西异步返回的处理
		if strings.Contains(string(message), "CgiBaseResponse") {
			go CgiResponseProcess(message)
			continue
		}
		// 这里是消息的处理
		var messages _struct.Message
		if err := json.Unmarshal(message, &messages); err != nil {
			continue
		}
		go MessageProcess(messages)
	}
}
