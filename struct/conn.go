package _struct

import "github.com/gorilla/websocket"

// WebSocketConn 方便让http接口来调用
var WebSocketConn *websocket.Conn

var WebsocketConnMap = make(map[string]*websocket.Conn)
