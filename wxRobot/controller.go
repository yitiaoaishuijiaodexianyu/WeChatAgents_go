package wxRobot

import (
	"WeChatAgents_go/common"
	_struct "WeChatAgents_go/struct"
	"github.com/gin-gonic/gin"
	"time"
)

type Queue struct {
	items []_struct.Message // 使用切片来存储队列元素
}

var q Queue

func Enqueue(item _struct.Message) {
	// 在访问共享资源前加锁
	mu.Lock()
	defer mu.Unlock()
	if len(q.items) > 200 {
		q.items = nil
	}
	q.items = append(q.items, item) // 将元素追加到切片的末尾
}

// GetMessage http方式获取消息
func GetMessage(c *gin.Context) {
	// 在访问共享资源前加锁
	mu.Lock()
	defer mu.Unlock()
	c.JSON(200, common.ResultCommon(0, q.items, "消息获取成功"))
	q.items = nil
}

// GetChatroomList 获取已知群的列表
func GetChatroomList(c *gin.Context) {
	chatroomList := _struct.KnownGroupConfig
	c.JSON(200, common.ResultCommon(0, chatroomList, "获取已知群列表成功"))
	return
}

// GetChatroomUserList 获取群成员列表
func GetChatroomUserList(c *gin.Context) {
	type Data struct {
		ChatroomId string `json:"chatroom_id"`
	}
	var data Data
	if ok := c.ShouldBindJSON(&data); ok != nil {
		c.JSON(200, common.ResultCommon(0, ChatroomUserInfo, "参数解析失败,返回全部已知的信息"))
		return
	}
	c.JSON(200, common.ResultCommon(0, ChatroomUserInfo[data.ChatroomId], "获取已知群列表成功"))
	return
}

// GetUserInfo 获取用户详细信息
func GetUserInfo(c *gin.Context) {

}

// SendText 发送文本信息
func SendText(c *gin.Context) {
	type Data struct {
		ReceiverId string `json:"receiver_id"`
		Text       string `json:"text"`
		AtIds      string `json:"at_ids"`
		BotWxId    string `json:"bot_wx_id"`
	}
	var data Data
	if ok := c.ShouldBindJSON(&data); ok != nil {
		c.JSON(200, common.ResultCommon(400, "", "参数解析失败"))
		return
	}
	var text []byte
	var reqId int
	if data.AtIds != "" {
		text, reqId = _struct.SendText(data.BotWxId, data.ReceiverId, data.Text, data.AtIds)
	} else {
		text, reqId = _struct.SendText(data.BotWxId, data.ReceiverId, data.Text, "")
	}
	_struct.WebSocketConn.WriteMessage(1, text)
	for true {
		if _struct.ReqIdMap[reqId].Status == 0 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "消息发送成功"))
			return
		}
		if _struct.ReqIdMap[reqId].Status == -2 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "消息发送失败"))
			return
		}
		time.Sleep(time.Second * 1)
	}
}

// SendImage 发送图片
func SendImage(c *gin.Context) {
	type Data struct {
		ReceiverId string `json:"receiver_id"`
		ImageUrl   string `json:"image_url"`
		BotWxId    string `json:"bot_wx_id"`
	}
	var data Data
	if ok := c.ShouldBindJSON(&data); ok != nil {
		c.JSON(200, common.ResultCommon(400, "", "参数解析失败"))
		return
	}
	image, reqId := _struct.UploadCdnImg(data.BotWxId, data.ReceiverId, data.ImageUrl)
	ResponseImgMap[reqId] = _struct.ImgInfo{
		CurrentWxid:  data.BotWxId,
		FromUserName: data.ReceiverId,
		Type:         1,
	}
	_struct.WebSocketConn.WriteMessage(1, image)
	for true {
		if _struct.ReqIdMap[reqId].Status == 0 {
			// 这里检测下新的reqId是否存在 存在的话说明是还没成功 去检测新的 然后一次性删除俩
			if _struct.ReqIdMap[reqId].NewReqId != 0 {
				for true {
					if _struct.ReqIdMap[_struct.ReqIdMap[reqId].NewReqId].Status == 1 {
						delete(_struct.ReqIdMap, reqId)
						delete(_struct.ReqIdMap, _struct.ReqIdMap[reqId].NewReqId)
						c.JSON(200, common.ResultCommon(0, "", "图片发送成功"))
						return
					}
					if _struct.ReqIdMap[reqId].Status == -2 {
						delete(_struct.ReqIdMap, reqId)
						c.JSON(200, common.ResultCommon(0, "", "图片发送失败"))
						return
					}
					time.Sleep(time.Second * 1)
				}
			}
		}
		if _struct.ReqIdMap[reqId].Status == -2 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "上传图片失败"))
			return
		}
		time.Sleep(time.Second * 1)
	}
}

// SendAppMsg 发送app消息
func SendAppMsg(c *gin.Context) {
	type Data struct {
		ReceiverId string `json:"receiver_id"`
		Xml        string `json:"xml"`
		BotWxId    string `json:"bot_wx_id"`
	}
	var data Data
	if ok := c.ShouldBindJSON(&data); ok != nil {
		c.JSON(200, common.ResultCommon(400, "", "参数解析失败"))
		return
	}
	appMsg, reqId := _struct.SendAppMessage(data.BotWxId, data.ReceiverId, data.Xml, 49)
	_struct.WebSocketConn.WriteMessage(1, appMsg)
	for true {
		if _struct.ReqIdMap[reqId].Status == 0 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "app消息发送成功"))
			return
		}
		if _struct.ReqIdMap[reqId].Status == -2 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "app消息发送失败，请检查xml信息"))
			return
		}
		time.Sleep(time.Second * 1)
	}
	//{
	//	"receiver_id":"39139856094@chatroom",
	//	"xml":"<appmsg appid=\"wx79f2c4418704b4f8\" sdkver=\"0\"><title>七里香</title><des>周杰伦</des><action>view</action><type>3</type><showtype>0</howtype><content /<url>https://www.kugou.com/song/#hash=2C7CEB6CC2340ECC8948E0ACE62F0CF8</url><dataurl>http://fsandroid.tx.kugou.com/202408220905/2f5740debd60c24149719b23565169b/3/2c7ceb6cc2340ecc8948e0ace62f0cf8/yp/full/ap1005_us776295431_mi336d5ebc5436534e61d16e63ddfca327_pi2_mx0_qu128_s799382533.mp3?info=cache?from=longzhu_api</ataurl><lowurl>https://www.kugou.com/song/#hash=2C7CEB6CC2340ECC8948E0ACE62F0CF8</lowurl><lowdataurl>http://fsandroid.tx.kugou.com/202408220905/2f5740debd60c24149719b23565169b/v3/2c7ceb6cc2340ecc8948e0ace62f0cf8/yp/full/ap1005_us776295431_mi336d5ebc5436534e61d16e63ddfca327_pi2_mx0_qu128_s799382533.mp3?nfo=cache?from=longzhu_api</lowdataurl><recorditem /><thumburl>https://singerimg.kugou.com/uploadpic/softhead/400/20230510/20230510173043311.jpg</humburl><messageaction /><laninfo /><extinfo /><sourceusername /><sourcedisplayname /><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></moticonmd5><fileext /><aeskey></aeskey></appattach><webviewshared><publisherId /><publisherReqId>0</publisherReqId></webviewshared><weappinfo><pagepath /<username /><appid /><appservicetype>0</appservicetype></weappinfo><websearch /><songalbumurl>https://singerimg.kugou.com/uploadpic/softhead/400/20230510/0230510173043311.jpg</songalbumurl></appmsg><fromusername></fromusername><scene>0</scene><appinfo><version>57</version><appname>酷狗音乐</appname></ppinfo><commenturl />"
	//}
}

// SendPat 拍一拍
func SendPat(c *gin.Context) {
	type Data struct {
		ReceiverId string `json:"receiver_id"`
		PatId      string `json:"pat_id"`
		BotWxId    string `json:"bot_wx_id"`
	}
	var data Data
	if ok := c.ShouldBindJSON(&data); ok != nil {
		c.JSON(200, common.ResultCommon(400, "", "参数解析失败"))
		return
	}
	pat, reqId := _struct.SendPatMessage(data.BotWxId, data.ReceiverId, data.PatId, 0)
	_struct.WebSocketConn.WriteMessage(1, pat)
	for true {
		if _struct.ReqIdMap[reqId].Status == 0 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "拍一拍消息发送成功"))
			return
		}
		if _struct.ReqIdMap[reqId].Status == -2 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "拍一拍消息发送失败"))
			return
		}
		time.Sleep(time.Second * 1)
	}
}

// SendEmoji 发送表情包
func SendEmoji(c *gin.Context) {
	type Data struct {
		ReceiverId  string `json:"receiver_id"`
		EmoJiMd5    string `json:"emoji_md5"`
		EmoJiLength int    `json:"emoji_length"`
		BotWxId     string `json:"bot_wx_id"`
	}
	var data Data
	if ok := c.ShouldBindJSON(&data); ok != nil {
		c.JSON(200, common.ResultCommon(400, "", "参数解析失败"))
		return
	}
	emoji, reqId := _struct.SendEmoji(data.BotWxId, data.ReceiverId, data.EmoJiMd5, data.EmoJiLength)
	_struct.WebSocketConn.WriteMessage(1, emoji)
	for true {
		if _struct.ReqIdMap[reqId].Status == 0 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "emoji消息发送成功"))
			return
		}
		if _struct.ReqIdMap[reqId].Status == -2 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "emoji消息发送失败"))
			return
		}
		time.Sleep(time.Second * 1)
	}
}

// SendVoice 发送语音条
func SendVoice(c *gin.Context) {
	type Data struct {
		ReceiverId  string `json:"receiver_id"`
		VoiceUrl    string `json:"voice_url"`
		VoiceLength int    `json:"voice_length"`
		BotWxId     string `json:"bot_wx_id"`
	}
	var data Data
	if ok := c.ShouldBindJSON(&data); ok != nil {
		c.JSON(200, common.ResultCommon(400, "", "参数解析失败"))
		return
	}
	voice, reqId := _struct.SendVoice(data.BotWxId, data.ReceiverId, data.VoiceUrl, data.VoiceLength)
	_struct.WebSocketConn.WriteMessage(1, voice)
	for true {
		if _struct.ReqIdMap[reqId].Status == 0 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "语音消息发送成功"))
			return
		}
		if _struct.ReqIdMap[reqId].Status == -2 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "语音消息发送失败"))
			return
		}
		time.Sleep(time.Second * 1)
	}
}

// DelChatroomMember 删除群成员
func DelChatroomMember(c *gin.Context) {
	type Data struct {
		ReceiverId string `json:"receiver_id"`
		UserWxId   string `json:"user_wx_id"`
		BotWxId    string `json:"bot_wx_id"`
	}
	var data Data
	if ok := c.ShouldBindJSON(&data); ok != nil {
		c.JSON(200, common.ResultCommon(400, "", "参数解析失败"))
		return
	}
	voice, reqId := _struct.DelChatroomMember(data.BotWxId, data.ReceiverId, data.UserWxId)
	_struct.WebSocketConn.WriteMessage(1, voice)
	for true {
		if _struct.ReqIdMap[reqId].Status == 0 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "删除群成员发送成功"))
			return
		}
		if _struct.ReqIdMap[reqId].Status == -2 {
			delete(_struct.ReqIdMap, reqId)
			c.JSON(200, common.ResultCommon(0, "", "删除群成员发送失败"))
			return
		}
		time.Sleep(time.Second * 1)
	}
}
