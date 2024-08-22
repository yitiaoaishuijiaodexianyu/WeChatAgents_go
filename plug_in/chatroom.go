package plug_in

import (
	"WeChatAgents_go/common"
	_struct "WeChatAgents_go/struct"
	"github.com/gin-gonic/gin"
)

func DelChatroomMember(c *gin.Context) {
	var message _struct.Message
	if ok := c.ShouldBindJSON(&message); ok != nil {
		return
	}
	var result = _struct.PlugInResult{}
	if message.CurrentPacket.Data.AddMsg.AtId == "" {
		return
	}
	result.Type = "delChatroomMember"
	result.ReceiverId = message.CurrentPacket.Data.AddMsg.FromUserName
	result.UserWxId = message.CurrentPacket.Data.AddMsg.AtId
	result.BotId = message.CurrentWxid
	c.JSON(200, common.ResultCommon(0, result, "踢了他成功"))
	return
}
