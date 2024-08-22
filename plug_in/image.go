package plug_in

import (
	"WeChatAgents_go/common"
	_struct "WeChatAgents_go/struct"
	"github.com/gin-gonic/gin"
)

func TestImage(c *gin.Context) {
	var message _struct.Message
	if ok := c.ShouldBindJSON(&message); ok != nil {
		return
	}
	if message.CurrentPacket.Data.AddMsg.ActionUserName != "wxid_za7ku9u4uu5q21" {
		return
	}
	var result = _struct.PlugInResult{}
	result.Type = "image"
	result.ReceiverId = message.CurrentPacket.Data.AddMsg.FromUserName
	result.BotId = message.CurrentWxid
	result.Url = "https://fanruizhecn.serv00.net/fl"
	c.JSON(200, common.ResultCommon(0, result, "图片发送成功"))
	return
}
