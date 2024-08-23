package plug_in

import (
	"WeChatAgents_go/common"
	_struct "WeChatAgents_go/struct"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"time"
)

func StartIdiomSolitaire(c *gin.Context) {
	var message _struct.Message
	if ok := c.ShouldBindJSON(&message); ok != nil {
	}
	var result = _struct.PlugInResult{}
	// 拿出来全部的成语  用成语作为key  内容作为value
	idiomMap := _struct.IdiomMap
	// 取出来所有成语[]string  只取出来成语 不要内容
	idiomStrings := _struct.IdiomStrings
	rand.Seed(int64(time.Now().Nanosecond()))
	randomNum := rand.Intn(len(idiomStrings))
	answer := ""
	correctTips := ""
	for _, v := range _struct.IdiomFirstMap[idiomMap[idiomStrings[randomNum]].Last] {
		answer += v.Word + "|"
		correctTips += "成语：[" + v.Word + "]\n" + "拼音：[" + v.Pinyin + "]\n" + "典故：[" + v.Derivation + "]|"
	}
	// 这里是取出来答案
	result.Type = "game"
	str := ""
	str += "成语接龙开始:" + "\n"
	// 成语
	str += "我先来出题: [" + idiomStrings[randomNum] + "]\n"
	// 最后一个汉字的拼音
	str += "拼音: [" + idiomMap[idiomStrings[randomNum]].Last + "]\n"
	str += "你来接..."

	if len(answer) > 0 {
		result.Answer = answer[0 : len(answer)-1]
		result.CorrectTips = correctTips[0 : len(correctTips)-1]
	} else {
		result.Answer = answer
		result.CorrectTips = correctTips
	}
	result.ReceiverId = message.CurrentPacket.Data.AddMsg.FromUserName
	result.BotId = message.CurrentWxid
	result.IsGame = 1
	result.GameStartName = "开始成语接龙"
	result.GameEndTime = int(common.GetCurrentTimestamp()) + 60
	resp, _ := resty.New().R().SetBody(map[string]interface{}{
		"receiver_id": message.CurrentPacket.Data.AddMsg.FromUserName,
		"text":        str,
		"at_ids":      "",
		"bot_wx_id":   message.CurrentWxid,
	}).Post("http://127.0.0.1:6636/api/SendText")
	var res struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	json.Unmarshal(resp.Body(), &res)
	if res.Code == 0 {
		c.JSON(200, common.ResultCommon(0, result, "开始成语接龙成功"))
	}
	return
}
