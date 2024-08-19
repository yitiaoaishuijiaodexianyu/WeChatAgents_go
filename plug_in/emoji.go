package plug_in

import (
	"WeChatAgents_go/common"
	_struct "WeChatAgents_go/struct"
	"github.com/gin-gonic/gin"
)

func TestEmoji(c *gin.Context) {
	var message _struct.Message
	if ok := c.ShouldBindJSON(&message); ok != nil {
		return
	}
	var result = _struct.PlugInResult{}
	result.Type = "emoji"
	result.ReceiverId = message.CurrentPacket.Data.AddMsg.FromUserName
	result.BotId = message.CurrentWxid
	result.EmojiMd5 = "2ad578fcfecda0f58e90e701b49348aa"
	result.EmojiLength = 81258
	c.JSON(200, common.ResultCommon(0, result, "发送emoji成功"))
	return
}
