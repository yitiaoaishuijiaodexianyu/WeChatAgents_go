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
	r.POST("/PlugIn/Pat", plug_in.Pat)
	r.POST("/PlugIn/TestEmoji", plug_in.TestEmoji)
	r.POST("/PlugIn/TestImage", plug_in.TestImage)
	r.POST("/PlugIn/XhAi", ai.XhAi)
	r.POST("/PlugIn/StarSign", plug_in.StarSign)

	// 下面是主动发送
	r.POST("/api/GetChatroomList", GetChatroomList)
	r.POST("/api/GetChatroomUserList", GetChatroomUserList)
	r.POST("/api/GetUserInfo", GetUserInfo)
	r.POST("/api/SendText", SendText)
	r.POST("/api/SendImage", SendImage)
	r.POST("/api/SendAppMsg", SendAppMsg)
	r.POST("/api/SendPat", SendPat)
	r.POST("/api/SendEmoji", SendEmoji)
	r.POST("/api/SendVoice", SendVoice)

	r.Run(_struct.Config.HttpServer.Host + ":" + _struct.Config.HttpServer.Port)
}
