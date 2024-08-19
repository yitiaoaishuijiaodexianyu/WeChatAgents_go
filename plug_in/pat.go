package plug_in

import (
	"WeChatAgents_go/common"
	_struct "WeChatAgents_go/struct"
	"github.com/gin-gonic/gin"
)

func Pat(c *gin.Context) {
	var message _struct.Message
	if ok := c.ShouldBindJSON(&message); ok != nil {
		return
	}
	var result = _struct.PlugInResult{}
	result.PatId = message.CurrentPacket.Data.AddMsg.ActionUserName
	if message.CurrentPacket.Data.AddMsg.AtId != "" {
		result.PatId = message.CurrentPacket.Data.AddMsg.AtId
	}
	result.Type = "pat"
	result.ReceiverId = message.CurrentPacket.Data.AddMsg.FromUserName
	result.BotId = message.CurrentWxid
	c.JSON(200, common.ResultCommon(0, result, "拍一拍成功"))
	return
}
