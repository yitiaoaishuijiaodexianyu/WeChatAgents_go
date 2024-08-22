package _struct

import (
	"encoding/json"
	"math/rand"
	"time"
)

var ReqIdMap = make(map[int]struct {
	Type     string `json:"type"`
	Status   int    `json:"status"`
	Data     string `json:"data"`
	NewReqId int    `json:"new_req_id"`
	//Chatroom string `json:"chatroom"`
})

// 可能会出现重复 重复了就会出现意想不到的问题 自行修复
func randReqId() int {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 生成一个10位的随机数
	minNumber := int64(1000000000) // 10位数的最小值
	maxNumber := int64(9999999999) // 10位数的最大值

	randomNumber := rand.Int63n(maxNumber-minNumber+1) + minNumber
	// 保证没有重复的 有重复的就重复运行 -- 要保证前面的处理完后及时删除这里面的
	if _, ok := ReqIdMap[int(randomNumber)]; ok {
		randReqId()
	}
	return int(randomNumber)
}

// SendText 发送文本
func SendText(botWxId string, ToUserName string, content string, atUserList string) ([]byte, int) {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			ToUserName string `json:"ToUserName"`
			Content    string `json:"Content"`
			MsgType    int    `json:"MsgType"`
			AtUsers    string `json:"AtUsers"`
		} `json:"CgiRequest"`
	}
	var t T
	t.ReqId = randReqId()
	t.BotWxid = botWxId
	t.CgiCmd = 522
	t.CgiRequest.MsgType = 1
	t.CgiRequest.ToUserName = ToUserName
	t.CgiRequest.Content = content
	t.CgiRequest.AtUsers = atUserList
	result, _ := json.Marshal(t)

	ReqIdMap[t.ReqId] = struct {
		Type     string `json:"type"`
		Status   int    `json:"status"`
		Data     string `json:"data"`
		NewReqId int    `json:"new_req_id"`
	}(struct {
		Type     string
		Status   int
		Data     string
		NewReqId int
	}{Type: "text", Status: 1, Data: string(result), NewReqId: 0})

	return result, t.ReqId
}

// SendEmoji 发送emoji
func SendEmoji(botWxId string, ToUserName string, EmojiMd5 string, EmojiLen int) ([]byte, int) {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			ToUserName string `json:"ToUserName"`
			EmojiMd5   string `json:"EmojiMd5"`
			EmojiLen   int    `json:"EmojiLen"`
		} `json:"CgiRequest"`
	}
	var t T
	t.ReqId = randReqId()
	t.BotWxid = botWxId
	t.CgiCmd = 175
	t.CgiRequest.ToUserName = ToUserName
	t.CgiRequest.EmojiMd5 = EmojiMd5
	t.CgiRequest.EmojiLen = EmojiLen
	result, _ := json.Marshal(t)

	ReqIdMap[t.ReqId] = struct {
		Type     string `json:"type"`
		Status   int    `json:"status"`
		Data     string `json:"data"`
		NewReqId int    `json:"new_req_id"`
	}(struct {
		Type     string
		Status   int
		Data     string
		NewReqId int
	}{Type: "emoji", Status: 1, Data: string(result), NewReqId: 0})

	return result, t.ReqId
}

// SendVoice 发送语音条
func SendVoice(botWxId string, ToUserName string, VoiceUrl string, VoiceTime int) ([]byte, int) {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			ToUserName string `json:"ToUserName"`
			VoiceUrl   string `json:"VoiceUrl"`
			VoiceTime  int    `json:"VoiceTime"`
		} `json:"CgiRequest"`
	}
	var t T
	t.ReqId = randReqId()
	t.BotWxid = botWxId
	t.CgiCmd = 127
	t.CgiRequest.ToUserName = ToUserName
	t.CgiRequest.VoiceUrl = VoiceUrl
	t.CgiRequest.VoiceTime = VoiceTime
	result, _ := json.Marshal(t)

	ReqIdMap[t.ReqId] = struct {
		Type     string `json:"type"`
		Status   int    `json:"status"`
		Data     string `json:"data"`
		NewReqId int    `json:"new_req_id"`
	}(struct {
		Type     string
		Status   int
		Data     string
		NewReqId int
	}{Type: "voice", Status: 1, Data: string(result), NewReqId: 0})

	return result, t.ReqId
}

// SendImg 发送图片
func SendImg(botWxId string, ToUserName string, AppMsgXml string) ([]byte, int) {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			ToUserName string `json:"ToUserName"`
			AppMsgXml  string `json:"ImageXml"`
		} `json:"CgiRequest"`
	}
	var t T
	t.ReqId = randReqId()
	t.BotWxid = botWxId
	t.CgiCmd = 110
	t.CgiRequest.ToUserName = ToUserName
	t.CgiRequest.AppMsgXml = AppMsgXml
	result, _ := json.Marshal(t)

	ReqIdMap[t.ReqId] = struct {
		Type     string `json:"type"`
		Status   int    `json:"status"`
		Data     string `json:"data"`
		NewReqId int    `json:"new_req_id"`
	}(struct {
		Type     string
		Status   int
		Data     string
		NewReqId int
	}{Type: "image", Status: 1, Data: string(result), NewReqId: 0})

	return result, t.ReqId
}

// SendAppMessage 发送App信息
func SendAppMessage(botWxId string, ToUserName string, AppMsgXml string, MsgType int) ([]byte, int) {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			ToUserName string `json:"ToUserName"`
			AppMsgXml  string `json:"AppMsgXml"`
			MsgType    int    `json:"MsgType"`
		} `json:"CgiRequest"`
	}
	var t T
	t.ReqId = randReqId()
	t.BotWxid = botWxId
	t.CgiCmd = 222
	t.CgiRequest.ToUserName = ToUserName
	t.CgiRequest.AppMsgXml = AppMsgXml
	t.CgiRequest.MsgType = MsgType
	result, _ := json.Marshal(t)

	ReqIdMap[t.ReqId] = struct {
		Type     string `json:"type"`
		Status   int    `json:"status"`
		Data     string `json:"data"`
		NewReqId int    `json:"new_req_id"`
	}(struct {
		Type     string
		Status   int
		Data     string
		NewReqId int
	}{Type: "app_message", Status: 1, Data: string(result), NewReqId: 0})

	return result, t.ReqId
}

func SendPatMessage(botWxId string, ChatUserName string, PattedUsername string, Scene int) ([]byte, int) {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			ChatUserName   string `json:"ChatUserName"`
			PattedUsername string `json:"PattedUsername"`
			Scene          int    `json:"Scene"`
		} `json:"CgiRequest"`
	}

	var t T
	t.ReqId = randReqId()
	t.BotWxid = botWxId
	t.CgiCmd = 849
	t.CgiRequest.ChatUserName = ChatUserName
	t.CgiRequest.PattedUsername = PattedUsername
	t.CgiRequest.Scene = Scene
	result, _ := json.Marshal(t)

	ReqIdMap[t.ReqId] = struct {
		Type     string `json:"type"`
		Status   int    `json:"status"`
		Data     string `json:"data"`
		NewReqId int    `json:"new_req_id"`
	}(struct {
		Type     string
		Status   int
		Data     string
		NewReqId int
	}{Type: "pat_message", Status: 1, Data: string(result), NewReqId: 0})

	return result, t.ReqId
}

// GetWxIdInfo 根据wxid获取信息  wxId为群号时 获取群信息-包含群员列表(群员详情需传入wxId)
func GetWxIdInfo(botWxId string, wxId string) ([]byte, int) {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			Wxid []string `json:"Wxid"`
		} `json:"CgiRequest"`
	}
	var t T
	t.ReqId = randReqId()
	t.CgiCmd = 182
	t.BotWxid = botWxId
	t.CgiRequest.Wxid = []string{wxId}
	result, _ := json.Marshal(t)

	ReqIdMap[t.ReqId] = struct {
		Type     string `json:"type"`
		Status   int    `json:"status"`
		Data     string `json:"data"`
		NewReqId int    `json:"new_req_id"`
	}(struct {
		Type     string
		Status   int
		Data     string
		NewReqId int
	}{Type: "get_info", Status: 1, Data: string(result), NewReqId: 0})

	return result, t.ReqId
}

// UploadCdnImg 上传cdn图片
func UploadCdnImg(botWxId string, toUserName string, path string) ([]byte, int) {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			ToUserName string `json:"ToUserName"`
			FileType   int    `json:"FileType"`
			FileUrl    string `json:"FileUrl"`
		} `json:"CgiRequest"`
	}
	var t T
	t.ReqId = randReqId()
	t.BotWxid = botWxId
	t.CgiCmd = 0
	t.CgiRequest.ToUserName = toUserName
	t.CgiRequest.FileType = 2
	t.CgiRequest.FileUrl = path
	result, _ := json.Marshal(t)

	ReqIdMap[t.ReqId] = struct {
		Type     string `json:"type"`
		Status   int    `json:"status"`
		Data     string `json:"data"`
		NewReqId int    `json:"new_req_id"`
	}(struct {
		Type     string
		Status   int
		Data     string
		NewReqId int
	}{Type: "upload_image", Status: 1, Data: string(result), NewReqId: 0})

	return result, t.ReqId
}

// UploadCdnFile 上传cdn文件
func UploadCdnFile(botWxId string, toUserName string, path string) ([]byte, int) {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			ToUserName string `json:"ToUserName"`
			FileType   int    `json:"FileType"`
			FileName   string `json:"FileName"`
			FileUrl    string `json:"FileUrl"`
		} `json:"CgiRequest"`
	}
	var t T
	t.ReqId = randReqId()
	t.BotWxid = botWxId
	t.CgiCmd = 0
	t.CgiRequest.ToUserName = toUserName
	t.CgiRequest.FileType = 5
	t.CgiRequest.FileName = "test.mp4"
	t.CgiRequest.FileUrl = path
	result, _ := json.Marshal(t)

	ReqIdMap[t.ReqId] = struct {
		Type     string `json:"type"`
		Status   int    `json:"status"`
		Data     string `json:"data"`
		NewReqId int    `json:"new_req_id"`
	}(struct {
		Type     string
		Status   int
		Data     string
		NewReqId int
	}{Type: "upload_file", Status: 1, Data: string(result), NewReqId: 0})

	return result, t.ReqId
}

// DownloadCdnImg 下载cdn图片
func DownloadCdnImg() {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			AesKey   string `json:"AesKey"`
			FileType int    `json:"FileType"`
			FileId   string `json:"FileId"`
		} `json:"CgiRequest"`
	}
}

// DownloadCdnFile 下载cdn文件
func DownloadCdnFile() {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			AesKey   string `json:"AesKey"`
			FileName string `json:"FileName"`
			FileType int    `json:"FileType"`
			FileId   string `json:"FileId"`
		} `json:"CgiRequest"`
	}
}

// DelChatroomMember 删除群成员
func DelChatroomMember(botWxId string, toUserName string, chatRoomId string) ([]byte, int) {
	type T struct {
		ReqId      int    `json:"ReqId"`
		BotWxid    string `json:"BotWxid"`
		CgiCmd     int    `json:"CgiCmd"`
		CgiRequest struct {
			Wxid       []string `json:"Wxid"`
			ChatroomID string   `json:"ChatroomID"`
		} `json:"CgiRequest"`
	}
	var t T
	t.ReqId = randReqId()
	t.BotWxid = botWxId
	t.CgiCmd = 179
	t.CgiRequest.Wxid = []string{toUserName}
	t.CgiRequest.ChatroomID = chatRoomId
	result, _ := json.Marshal(t)

	ReqIdMap[t.ReqId] = struct {
		Type     string `json:"type"`
		Status   int    `json:"status"`
		Data     string `json:"data"`
		NewReqId int    `json:"new_req_id"`
	}(struct {
		Type     string
		Status   int
		Data     string
		NewReqId int
	}{Type: "del_chatroom_member", Status: 1, Data: string(result), NewReqId: 0})

	return result, t.ReqId
}

//{"CurrentPacket":{"WebConnId":"","Data":{"AddMsg":{"MsgId":1404310834,"FromUserName":"39139856094@chatroom","ToUserName":"wxid_tj1hdj6zuh3b12","MsgType":49,"Content":"\u003c?xml version=\"1.0\"?\u003e\n\u003cmsg\u003e\n\t\u003cappmsg appid=\"wxb98ac240fd74b0e3\" sdkver=\"\"\u003e\n\t\t\u003ctitle\u003e抓大鹅\u003c/title\u003e\n\t\t\u003cdes /\u003e\n\t\t\u003caction\u003eview\u003c/action\u003e\n\t\t\u003ctype\u003e33\u003c/type\u003e\n\t\t\u003cshowtype\u003e0\u003c/showtype\u003e\n\t\t\u003ccontent /\u003e\n\t\t\u003curl\u003ehttps://mp.weixin.qq.com/mp/waerrpage?appid=wxb98ac240fd74b0e3\u0026amp;amp;type=upgrade\u0026amp;amp;upgradetype=3#wechat_redirect\u003c/url\u003e\n\t\t\u003cdataurl /\u003e\n\t\t\u003clowurl /\u003e\n\t\t\u003clowdataurl /\u003e\n\t\t\u003crecorditem /\u003e\n\t\t\u003cthumburl /\u003e\n\t\t\u003cmessageaction /\u003e\n\t\t\u003claninfo /\u003e\n\t\t\u003cmd5\u003ea5dfa1ea225c278d8377802e315296fa\u003c/md5\u003e\n\t\t\u003cextinfo /\u003e\n\t\t\u003csourceusername /\u003e\n\t\t\u003csourcedisplayname\u003e抓大鹅\u003c/sourcedisplayname\u003e\n\t\t\u003ccommenturl /\u003e\n\t\t\u003cappattach\u003e\n\t\t\t\u003ctotallen\u003e0\u003c/totallen\u003e\n\t\t\t\u003cattachid /\u003e\n\t\t\t\u003cemoticonmd5\u003e\u003c/emoticonmd5\u003e\n\t\t\t\u003cfileext\u003ejpg\u003c/fileext\u003e\n\t\t\t\u003cfilekey\u003ead8d5f0d5ca82658ab6774d7a8f0b3ec\u003c/filekey\u003e\n\t\t\t\u003ccdnthumburl\u003e3057020100044b30490201000204946eb63d02032df08e0204aab4bc770204669f1f1c042432626663653161352d333566362d343261612d623165332d6234376365323635613234360204052808030201000405004c543e00\u003c/cdnthumburl\u003e\n\t\t\t\u003caeskey\u003e8f51cc4fbcd32970914d2aa1657dbd54\u003c/aeskey\u003e\n\t\t\t\u003ccdnthumbaeskey\u003e8f51cc4fbcd32970914d2aa1657dbd54\u003c/cdnthumbaeskey\u003e\n\t\t\t\u003ccdnthumbmd5\u003ea5dfa1ea225c278d8377802e315296fa\u003c/cdnthumbmd5\u003e\n\t\t\t\u003cencryver\u003e1\u003c/encryver\u003e\n\t\t\t\u003ccdnthumblength\u003e61763\u003c/cdnthumblength\u003e\n\t\t\t\u003ccdnthumbheight\u003e100\u003c/cdnthumbheight\u003e\n\t\t\t\u003ccdnthumbwidth\u003e100\u003c/cdnthumbwidth\u003e\n\t\t\u003c/appattach\u003e\n\t\t\u003cwebviewshared\u003e\n\t\t\t\u003cpublisherId /\u003e\n\t\t\t\u003cpublisherReqId\u003e0\u003c/publisherReqId\u003e\n\t\t\u003c/webviewshared\u003e\n\t\t\u003cweappinfo\u003e\n\t\t\t\u003cpagepath /\u003e\n\t\t\t\u003cusername\u003egh_730ace0831c4@app\u003c/username\u003e\n\t\t\t\u003cappid\u003ewxb98ac240fd74b0e3\u003c/appid\u003e\n\t\t\t\u003ctype\u003e2\u003c/type\u003e\n\t\t\t\u003cweappiconurl\u003ehttp://mmbiz.qpic.cn/sz_mmbiz_png/OFye4rdPwyccR5UXGz9z5X74I2ghib0yT0pU4aFufUDy11cR8IiccLf69rba9ecRuUwAX3WtMGNNPE1nhDnBzialA/640?wx_fmt=png\u0026amp;wxfrom=200\u003c/weappiconurl\u003e\n\t\t\t\u003cappservicetype\u003e4\u003c/appservicetype\u003e\n\t\t\t\u003cshareId\u003e2_wxb98ac240fd74b0e3_1837218910_1721704180_1\u003c/shareId\u003e\n\t\t\u003c/weappinfo\u003e\n\t\t\u003cwebsearch /\u003e\n\t\u003c/appmsg\u003e\n\t\u003cfromusername\u003ewxid_za7ku9u4uu5q21\u003c/fromusername\u003e\n\t\u003cscene\u003e0\u003c/scene\u003e\n\t\u003cappinfo\u003e\n\t\t\u003cversion\u003e1\u003c/version\u003e\n\t\t\u003cappname\u003e抓大鹅\u003c/appname\u003e\n\t\u003c/appinfo\u003e\n\t\u003ccommenturl /\u003e\n\u003c/msg\u003e\n","Status":3,"ImgStatus":2,"ImgBuf":null,"CreateTime":1721704220,"MsgSource":"\u003cmsgsource\u003e\n\t\u003ctmp_node\u003e\n\t\t\u003cpublisher-id\u003e\u003c/publisher-id\u003e\n\t\u003c/tmp_node\u003e\n\t\u003csec_msg_node\u003e\n\t\t\u003cuuid\u003e07daf9c012e6743d8ac47b14b4306237_\u003c/uuid\u003e\n\t\t\u003crisk-file-flag /\u003e\n\t\t\u003crisk-file-md5-list /\u003e\n\t\u003c/sec_msg_node\u003e\n\t\u003csilence\u003e1\u003c/silence\u003e\n\t\u003cmembercount\u003e23\u003c/membercount\u003e\n\t\u003csignature\u003eV1_OOG3J2hM|v1_Isf5+E3A\u003c/signature\u003e\n\u003c/msgsource\u003e\n","PushContent":"","NewMsgId":7400409546359626802,"NewMsgIdExt":"7400409546359626802","ActionUserName":"wxid_za7ku9u4uu5q21","ActionNickName":""},"EventName":"ON_EVENT_MSG_NEW"}},"CurrentWxid":"wxid_tj1hdj6zuh3b12","UUid":"guKFsdH6XXGamvw_9Qlo"}

//<appmsg appid="wx79f2c4418704b4f8\" sdkver=\"0\"><title>稻香</title><des>周杰伦 · 周杰伦</des><type>3</type><messageext>kugou://start.weixin?{\"cmd\":212,\"jsonStr\":{\"needfav\":0,\"userid\":1443852867,\"wx_music_video\":1,\"fail_process\":4,\"album_audio_id\":\"32042828\",\"bitrate\":128,\"duration\":223503,\"cyt1\":\"wechat\",\"chl2\":\"wx_music_video\",\"pay_type\":3,\"u\":1443852867,\"320privilege\":10,\"type\":\"audio\",\"filename\":\"周杰伦 - 稻香\",\"sqprivilege\":10,\"hash\":\"8909E1809908CD8E3BF6CF85D98B93F0\",\"privilege\":10,\"chl\":\"wechat\",\"old_cpy\":0,\"trans_param\":{\"free_for_ad\":0,\"all_quality_free\":0,\"musicpack_advance\":1,\"free_limited\":0,\"display\":32,\"pay_block_tpl\":0,\"exclusive\":false,\"hash_offset\":{\"end_ms\":60000,\"file_type\":0,\"offset_hash\":\"974381C4C3DAD140FD7CC02466378C93\",\"end_byte\":960115,\"clip_hash\":\"CEDC66EA175AED2ED09E4858BE2C80CD\",\"start_byte\":0,\"start_ms\":0},\"display_rate\":1}},\"type\":1}</messageext><url>https://t3.kugou.com/wc/s/2ylI50eCPV3#wechat_music_url=7B22736F6E675F5761704C69766555524C223A2268747470733A5C2F5C2F6D2E6B75676F752E636F6D5C2F6170695C2F76315C2F7765636861745C2F696E6465783F686173683D3839303945313830393930384344384533424636434638354439384239334630266D69643D323963613233316337616333653037653236353630326135393564303466393539323137373831322676657273696F6E3D313232353526706C61743D32266170697665723D3226616C62756D5F617564696F5F69643D333230343238323826757365725F69643D31343433383532383637266B65793D343135373137643834393731616230616236316166643362313364363166646626636D643D313031266578743D6D3461265F743D3137323136333536343326616C62756D5F69643D3936303339392673686172655F63686C3D776563686174267369676E3D3362613162633438313562333533643665353562613665353466663139353364222C22736F6E675F5769666955524C223A2268747470733A5C2F5C2F6D2E6B75676F752E636F6D5C2F6170695C2F76315C2F7765636861745C2F696E6465783F686173683D3839303945313830393930384344384533424636434638354439384239334630266D69643D323963613233316337616333653037653236353630326135393564303466393539323137373831322676657273696F6E3D313232353526706C61743D32266170697665723D3226616C62756D5F617564696F5F69643D333230343238323826757365725F69643D31343433383532383637266B65793D343135373137643834393731616230616236316166643362313364363166646626636D643D313031266578743D6D3461265F743D3137323136333536343326616C62756D5F69643D3936303339392673686172655F63686C3D776563686174267369676E3D3362613162633438313562333533643665353562613665353466663139353364227D</url><dataurl>https://m.kugou.com/api/v1/wechat/index?hash=8909E1809908CD8E3BF6CF85D98B93F0&amp;mid=29ca231c7ac3e07e265602a595d04f9592177812&amp;version=12255&amp;plat=2&amp;apiver=2&amp;album_audio_id=32042828&amp;user_id=1443852867&amp;key=415717d84971ab0ab61afd3b13d61fdf&amp;cmd=101&amp;ext=m4a&amp;_t=1721635643&amp;album_id=960399&amp;share_chl=wechat&amp;sign=3ba1bc4815b353d6e55ba6e54ff1953d</dataurl><songalbumurl>http://wxapp.tc.qq.com/202/20304/stodownload?filekey=30350201010421301f020200ca0402534804101e81efef5df69efbaf63fb6add59eb5b020305243b040d00000004627466730000000132&amp;hy=SH&amp;storeid=2669e133e0005938e6d81bc5e000000ca00004f50534807d3e031573805212&amp;bizid=1023</songalbumurl><songlyric>[00:09.47]作词：周杰伦[00:16.58]作曲：周杰伦[00:23.69]编曲：黄雨勋[00:30.80]对这个世界如果[00:32.63]你有太多的抱怨[00:34.26]跌倒了就不敢继续往前走[00:37.23]为什么人要这么的脆弱[00:39.74]堕落[00:41.35]请你打开电视看看[00:43.08]多少人为生命在努力[00:45.29]勇敢的走下去[00:47.07]我们是不是该知足[00:49.44]珍惜一切[00:50.72]就算没有拥有[00:54.01]还记得你说家是唯一的城堡[00:57.65]随着稻香河流继续奔跑[01:00.64]微微笑小时候的梦我知道[01:05.38]不要哭让萤火虫带着你逃跑[01:09.37]乡间的歌谣永远的依靠[01:12.33]回家吧回到最初的美好[01:41.01]不要这么容易就想放弃[01:43.23]就像我说的[01:44.68]追不到的梦想[01:45.81]换个梦不就得了[01:47.75]为自己的人生鲜艳上色[01:49.78]先把爱涂上喜欢的颜色[01:52.72]笑一个吧功成名就不是目的[01:55.68]让自己快乐快乐[01:56.97]这才叫做意义[01:58.60]童年的纸飞机[02:00.01]现在终于飞回我手里[02:04.16]所谓的那快乐[02:05.69]赤脚在田里追蜻蜓追到累了[02:08.65]偷摘水果被蜜蜂给叮到怕了[02:11.73]谁在偷笑呢[02:13.21]我靠着稻草人吹着风[02:15.22]唱着歌睡着了[02:17.67]午后吉他在虫鸣中更清脆[02:20.68]阳光洒在路上就不怕心碎[02:23.10]珍惜一切就算没有拥有[02:27.69]还记得你说家是唯一的城堡[02:31.32]随着稻香河流继续奔跑[02:34.28]微微笑小时候的梦我知道[02:39.03]不要哭让萤火虫带着你逃跑[02:43.03]乡间的歌谣永远的依靠[02:45.99]回家吧回到最初的美好[02:51.50]还记得你说家是唯一的城堡[02:54.91]随着稻香河流继续奔跑[02:57.88]微微笑小时候的梦我知道[03:02.80]不要哭让萤火虫带着你逃跑[03:06.44]乡间的歌谣永远的依靠[03:09.40]回家吧回到最初的美好</songlyric><appattach><cdnthumburl>3057020100044b304902010002046d81bc5e02032df08e02047fb4bc770204669e133e042464333831333766302d323135622d343135312d393538322d6335396164666663643461360204051408030201000405004c4f2a00</cdnthumburl><cdnthumbmd5>5959c481fee22790e9d230b5511043e3</cdnthumbmd5><cdnthumblength>37463</cdnthumblength><cdnthumbwidth>720</cdnthumbwidth><cdnthumbheight>720</cdnthumbheight><cdnthumbaeskey>0f0d506a20f722513bf72ea4203c7846</cdnthumbaeskey><aeskey>0f0d506a20f722513bf72ea4203c7846</aeskey><encryver>0</encryver><filekey>wxid_tj1hdj6zuh3b12_391_1721635646</filekey></appattach><md5>5959c481fee22790e9d230b5511043e3</md5><statextstr>GhQKEnd4NzlmMmM0NDE4NzA0YjRmOA==</statextstr><musicShareItem><mid>getlinkclisdkmid_2ylI50eCPV3</mid><mvSingerName>周杰伦 · 周杰伦</mvSingerName><mvExtInfo>kugou://start.weixin?{\"cmd\":212,\"jsonStr\":{\"needfav\":0,\"userid\":1443852867,\"wx_music_video\":1,\"fail_process\":4,\"album_audio_id\":\"32042828\",\"bitrate\":128,\"duration\":223503,\"cyt1\":\"wechat\",\"chl2\":\"wx_music_video\",\"pay_type\":3,\"u\":1443852867,\"320privilege\":10,\"type\":\"audio\",\"filename\":\"周杰伦 - 稻香\",\"sqprivilege\":10,\"hash\":\"8909E1809908CD8E3BF6CF85D98B93F0\",\"privilege\":10,\"chl\":\"wechat\",\"old_cpy\":0,\"trans_param\":{\"free_for_ad\":0,\"all_quality_free\":0,\"musicpack_advance\":1,\"free_limited\":0,\"display\":32,\"pay_block_tpl\":0,\"exclusive\":false,\"hash_offset\":{\"end_ms\":60000,\"file_type\":0,\"offset_hash\":\"974381C4C3DAD140FD7CC02466378C93\",\"end_byte\":960115,\"clip_hash\":\"CEDC66EA175AED2ED09E4858BE2C80CD\",\"start_byte\":0,\"start_ms\":0},\"display_rate\":1}},\"type\":1}</mvExtInfo><mvIdentification>8909E1809908CD8E3BF6CF85D98B93F0</mvIdentification><musicDuration>223503</musicDuration></musicShareItem></appmsg>

//{"CgiBaseResponse":{"ErrMsg":"","Ret":0},"ReqId":2008621197,"ResponseData":"\u003cappmsg appid=\"\"  sdkver=\"0\"\u003e\u003ctitle\u003e稻香\u003c/title\u003e\u003cdes\u003e\u003c/des\u003e\u003caction\u003e\u003c/action\u003e\u003ctype\u003e6\u003c/type\u003e\u003cshowtype\u003e0\u003c/showtype\u003e\u003csoundtype\u003e0\u003c/soundtype\u003e\u003cmediatagname\u003e\u003c/mediatagname\u003e\u003cmessageext\u003e\u003c/messageext\u003e\u003cmessageaction\u003e\u003c/messageaction\u003e\u003ccontent\u003e\u003c/content\u003e\u003ccontentattr\u003e0\u003c/contentattr\u003e\u003curl\u003e\u003c/url\u003e\u003clowurl\u003e\u003c/lowurl\u003e\u003cdataurl\u003e\u003c/dataurl\u003e\u003clowdataurl\u003e\u003c/lowdataurl\u003e\u003csongalbumurl\u003e\u003c/songalbumurl\u003e\u003csonglyric\u003e\u003c/songlyric\u003e\u003cappattach\u003e\u003ctotallen\u003e0\u003c/totallen\u003e\u003cattachid\u003e@cdn_3057020100044b3049020100020488b1fbd302032f56c1020494e5e7730204669f5b20042463323163383331612d366165662d343234382d626437652d3736616566626533373738630204011400050201000405004c50b900_d41d8cd98f00b204e9800998ecf8427e_1\u003c/attachid\u003e\u003cemoticonmd5\u003e\u003c/emoticonmd5\u003e\u003cfileext\u003esilk\u003c/fileext\u003e\u003ccdnattachurl\u003e3057020100044b3049020100020488b1fbd302032f56c1020494e5e7730204669f5b20042463323163383331612d366165662d343234382d626437652d3736616566626533373738630204011400050201000405004c50b900\u003c/cdnattachurl\u003e\u003ccdnthumbaeskey\u003e\u003c/cdnthumbaeskey\u003e\u003caeskey\u003ed41d8cd98f00b204e9800998ecf8427e\u003c/aeskey\u003e\u003cencryver\u003e0\u003c/encryver\u003e\u003cfilekey\u003ewxid_elv1n0btkywy000_1633495096\u003c/filekey\u003e\u003coverwrite_newmsgid\u003e0\u003c/overwrite_newmsgid\u003e\u003cfileuploadtoken\u003e\u003c/fileuploadtoken\u003e\u003c/appattach\u003e\u003cextinfo\u003e\u003c/extinfo\u003e\u003csourceusername\u003e\u003c/sourceusername\u003e\u003csourcedisplayname\u003e\u003c/sourcedisplayname\u003e\u003cthumburl\u003e\u003c/thumburl\u003e\u003cmd5\u003ed41d8cd98f00b204e9800998ecf8427e\u003c/md5\u003e\u003cstatextstr\u003e\u003c/statextstr\u003e\u003cdirectshare\u003e0\u003c/directshare\u003e\u003crecorditem\u003e\u003c![CDATA[\u003crecordinfo\u003e\u003cedittime\u003e0\u003c/edittime\u003e\u003cfromscene\u003e0\u003c/fromscene\u003e\u003c/recordinfo\u003e]]\u003e\u003c/recorditem\u003e\u003c/appmsg\u003e"}
