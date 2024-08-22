package wxRobot

import (
	"WeChatAgents_go/plug_in"
	"WeChatAgents_go/plug_in/ai"
	_struct "WeChatAgents_go/struct"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HttpRun() {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})

	// 这里是程序自带的功能
	r.POST("/PlugIn/RequestSong", plug_in.RequestSong)
	r.POST("/PlugIn/DelChatroomMember", plug_in.DelChatroomMember)
	r.POST("/PlugIn/Pat", plug_in.Pat)
	r.POST("/PlugIn/TestEmoji", plug_in.TestEmoji)
	r.POST("/PlugIn/TestImage", plug_in.TestImage)
	r.POST("/PlugIn/XhAi", ai.XhAi)
	r.POST("/PlugIn/StarSign", plug_in.StarSign)
	r.POST("/PlugIn/StartGuessMusic", plug_in.StartGuessMusic)

	// 下面是主动发送
	// 如果你不会写websocket 可以使用下面的方法
	// 反正在你本地运行，写个for循环来请求 /api/GetMessage 获取消息就行了
	r.POST("/api/GetMessage", GetMessage)                   // 获取消息
	r.POST("/api/GetChatroomList", GetChatroomList)         // 获取已知群列表
	r.POST("/api/GetChatroomUserList", GetChatroomUserList) // 获取群内成员
	r.POST("/api/GetUserInfo", GetUserInfo)                 // 获取成员信息
	r.POST("/api/SendText", SendText)                       // 发送文本
	r.POST("/api/SendImage", SendImage)                     // 发送图片
	r.POST("/api/SendAppMsg", SendAppMsg)                   // 发送xml消息
	r.POST("/api/SendPat", SendPat)                         // 发送拍一拍消息
	r.POST("/api/SendEmoji", SendEmoji)                     // 发送表情包消息
	r.POST("/api/SendVoice", SendVoice)                     // 发送语音条
	r.POST("/api/DelChatroomMember", DelChatroomMember)     // 删除群成员

	r.Run(_struct.Config.HttpServer.Host + ":" + _struct.Config.HttpServer.Port)
}
