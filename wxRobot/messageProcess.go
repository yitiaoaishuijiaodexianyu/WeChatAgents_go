package wxRobot

import (
	"WeChatAgents_go/common"
	"WeChatAgents_go/config"
	_struct "WeChatAgents_go/struct"
	"encoding/json"
	"fmt"
	resty "github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ResponseImgMap = make(map[int]_struct.ImgInfo)
var ResponseUserInfoMap = make(map[int]_struct.GetUserInfo)
var mu sync.Mutex

// UserList [用户的wx_id]:[微信名]
var UserList = make(map[string]string)

// ChatroomInfo [群id]:[群名]
var ChatroomInfo = make(map[string]string)

// 存请求的reqId来判断 [reqid]:[类型(自定义如何处理)]
var reqType = make(map[int]int)

// GameStatus 存游戏的开始状态[群id]:[[status]:[1],[timestamp]:[时间戳]]
var GameStatus = make(map[string]map[string]int)

// 存游戏的答案 [群id]:[答案]
var gameAnswer = make(map[string]string)

// getChatRoomInfo 获取群的信息
func getChatRoomInfo(botWxId string, chatRoomId string, c *websocket.Conn) {
	result, reqId := GetWxIdInfo(botWxId, chatRoomId)
	ResponseUserInfoMap[reqId] = _struct.GetUserInfo{Type: 2}
	reqType[reqId] = 2
	c.WriteMessage(1, result)
}

// MessageProcess 消息处理
func MessageProcess(message _struct.Message, c *websocket.Conn) {

	config := config.GetConfigInfo()
	fmt.Println(config)

	// 入群欢迎
	if message.CurrentPacket.Data.AddMsg.MsgType == 10000 {
		joinGroup(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.Content, message.CurrentPacket.Data.AddMsg.FromUserName, c)
	}

	if _, ok := ChatroomInfo[message.CurrentPacket.Data.AddMsg.FromUserName]; !ok {
		if strings.Contains(message.CurrentPacket.Data.AddMsg.FromUserName, "@chatroom") {
			getChatRoomInfo(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, c)
		}
		time.Sleep(time.Second * 1)
	}

	content := "===============消息块==================\n"
	content += "时间：" + common.GetCurrentTime() + "\n"
	if strings.Contains(message.CurrentPacket.Data.AddMsg.FromUserName, "@chatroom") {
		content += "群名：[" + ChatroomInfo[message.CurrentPacket.Data.AddMsg.FromUserName] + "] 群id：[" + message.CurrentPacket.Data.AddMsg.FromUserName + "]\n"
	}
	content += "用户名：[" + UserList[message.CurrentPacket.Data.AddMsg.ActionUserName] + "] 用户id：[" + message.CurrentPacket.Data.AddMsg.ActionUserName + "]\n"
	content += "消息：[" + message.CurrentPacket.Data.AddMsg.Content + "]\n"
	content += "===============消息块=================="
	fmt.Println(content)

	// 判断猜歌名
	if _, ok := gameAnswer[message.CurrentPacket.Data.AddMsg.FromUserName]; ok {
		if message.CurrentPacket.Data.AddMsg.Content == gameAnswer[message.CurrentPacket.Data.AddMsg.FromUserName] {
			delete(gameAnswer, message.CurrentPacket.Data.AddMsg.FromUserName)
			// 在访问共享资源前加锁
			mu.Lock()
			result := SendText(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "@"+UserList[message.CurrentPacket.Data.AddMsg.ActionUserName]+" 恭喜回答正确："+message.CurrentPacket.Data.AddMsg.Content, message.CurrentPacket.Data.AddMsg.ActionUserName)
			c.WriteMessage(1, result)
			message.CurrentPacket.Data.AddMsg.Content = "开始猜歌名"
			delete(GameStatus, message.CurrentPacket.Data.AddMsg.FromUserName)
			time.Sleep(time.Second * 1)
			// 自己去调用一次开始猜歌名
			MessageProcess(message, c)
			// 释放锁
			mu.Unlock()
		}
	}

	if message.CurrentPacket.Data.AddMsg.Content == "发个表情" {
		result := SendEmoji(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "2ad578fcfecda0f58e90e701b49348aa", 81258)
		c.WriteMessage(1, result)
	}

	//if message.CurrentPacket.Data.AddMsg.Content == "发个视频" {
	//	result, reqId := UploadCdnFile(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "https://dayu.qqsuu.cn/moyuribaoshipin/file/20240618.mp4")
	//	ResponseImgMap[reqId] = _struct.ImgInfo{
	//		CurrentWxid:  message.CurrentWxid,
	//		FromUserName: message.CurrentPacket.Data.AddMsg.FromUserName,
	//		Type:         2,
	//	}
	//	fmt.Println(reqId)
	//	fmt.Println(string(result))
	//	c.WriteMessage(1, result)
	//}

	// 发微信收藏的表情包的
	if message.CurrentPacket.Data.AddMsg.Content == "后入鸭子" {
		result := SendEmoji(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "2ad578fcfecda0f58e90e701b49348aa", 81258)
		c.WriteMessage(1, result)
	}

	// 点歌功能
	if strings.Contains(message.CurrentPacket.Data.AddMsg.Content, "点歌") {
		MusicName := strings.Replace(message.CurrentPacket.Data.AddMsg.Content, "点歌", "", -1)
		resp, _ := resty.New().R().Get("https://api.frz379.com/go_api/api/GetMusicMp3?musicName=" + MusicName)
		type T struct {
			Code          int         `json:"code"`
			Message       string      `json:"message"`
			DataItem      interface{} `json:"data_item"`
			PrimitiveData struct {
				Data struct {
					Code     int    `json:"code"`
					Cover    string `json:"cover"`  //图片
					Link     string `json:"link"`   //歌曲网页链接
					Lyrics   string `json:"lyrics"` //歌词
					Msg      string `json:"msg"`
					MusicUrl string `json:"music_url"` //歌的链接
					Singer   string `json:"singer"`    //作者
					Title    string `json:"title"`     //歌名
				} `json:"data"`
				Source string `json:"source"`
			} `json:"primitive_data"`
		}
		var t T
		json.Unmarshal(resp.Body(), &t)
		result := SendAppMessage(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "<appmsg appid=\"wx79f2c4418704b4f8\" sdkver=\"0\"><title>"+t.PrimitiveData.Data.Title+"</title><des>"+t.PrimitiveData.Data.Singer+"</des><action>view</action><type>3</type><showtype>0</showtype><content /><url>"+t.PrimitiveData.Data.Link+"</url><dataurl>"+t.PrimitiveData.Data.MusicUrl+"</dataurl><lowurl>"+t.PrimitiveData.Data.Link+"</lowurl><lowdataurl>"+t.PrimitiveData.Data.MusicUrl+"</lowdataurl><recorditem /><thumburl>"+t.PrimitiveData.Data.Cover+"</thumburl><messageaction /><laninfo /><extinfo /><sourceusername /><sourcedisplayname /><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></emoticonmd5><fileext /><aeskey></aeskey></appattach><webviewshared><publisherId /><publisherReqId>0</publisherReqId></webviewshared><weappinfo><pagepath /><username /><appid /><appservicetype>0</appservicetype></weappinfo><websearch /><songalbumurl>"+t.PrimitiveData.Data.Cover+"</songalbumurl></appmsg><fromusername></fromusername><scene>0</scene><appinfo><version>57</version><appname>酷狗音乐</appname></appinfo><commenturl />", 49)
		c.WriteMessage(1, result)
	}

	// 发小程序
	if message.CurrentPacket.Data.AddMsg.Content == "抓大鹅" {
		result := SendAppMessage(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "<appmsg appid=\"wxb98ac240fd74b0e3\" sdkver=\"\"><title>抓大鹅</title><des /><action>view</action><type>33</type><showtype>0</showtype><content /><url>https://mp.weixin.qq.com/mp/waerrpage?appid=wxb98ac240fd74b0e3&amp;amp;type=upgrade&amp;amp;upgradetype=3#wechat_redirect</url><dataurl /><lowurl /><lowdataurl /><recorditem /><thumburl /><messageaction /><laninfo /><md5>a5dfa1ea225c278d8377802e315296fa</md5><extinfo /><sourceusername /><sourcedisplayname>抓大鹅</sourcedisplayname><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></emoticonmd5><fileext>jpg</fileext><filekey>ad8d5f0d5ca82658ab6774d7a8f0b3ec</filekey><cdnthumburl>3057020100044b30490201000204946eb63d02032df08e0204aab4bc770204669f1f1c042432626663653161352d333566362d343261612d623165332d6234376365323635613234360204052808030201000405004c543e00</cdnthumburl><aeskey>8f51cc4fbcd32970914d2aa1657dbd54</aeskey><cdnthumbaeskey>8f51cc4fbcd32970914d2aa1657dbd54</cdnthumbaeskey><cdnthumbmd5>a5dfa1ea225c278d8377802e315296fa</cdnthumbmd5><encryver>1</encryver><cdnthumblength>61763</cdnthumblength><cdnthumbheight>100</cdnthumbheight><cdnthumbwidth>100</cdnthumbwidth></appattach><webviewshared><publisherId /><publisherReqId>0</publisherReqId></webviewshared><weappinfo><pagepath /><username>gh_730ace0831c4@app</username><appid>wxb98ac240fd74b0e3</appid><type>2</type><weappiconurl>http://mmbiz.qpic.cn/sz_mmbiz_png/OFye4rdPwyccR5UXGz9z5X74I2ghib0yT0pU4aFufUDy11cR8IiccLf69rba9ecRuUwAX3WtMGNNPE1nhDnBzialA/640?wx_fmt=png&amp;wxfrom=200</weappiconurl><appservicetype>4</appservicetype><shareId>2_wxb98ac240fd74b0e3_1837218910_1721704180_1</shareId></weappinfo><websearch /></appmsg>", 49)
		c.WriteMessage(1, result)
	}
	// 发小程序
	if message.CurrentPacket.Data.AddMsg.Content == "羊了个羊" {
		result := SendAppMessage(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "<appmsg appid=\"wx141bfb9b73c970a9\" sdkver=\"\"><title>什么能阻止我睡觉？就是这头羊！！</title><des /><action>view</action><type>33</type><showtype>0</showtype><content /><url>https://mp.weixin.qq.com/mp/waerrpage?appid=wx141bfb9b73c970a9&amp;amp;type=upgrade&amp;amp;upgradetype=3#wechat_redirect</url><dataurl /><lowurl /><lowdataurl /><recorditem /><thumburl /><messageaction /><laninfo /><md5>43692fcd974d406c33aad13fff56079f</md5><extinfo /><sourceusername /><sourcedisplayname>羊了个羊</sourcedisplayname><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></emoticonmd5><fileext>jpg</fileext><filekey>3e77abc20b8bc03644274f1113aede03</filekey><cdnthumburl>3057020100044b30490201000204946eb63d02032df08e0204a8b4bc77020466b2de3c042438636464613930342d623235622d343037662d383166332d3639663532363638366538380204052408030201000405004c505600</cdnthumburl><aeskey>15e6af19df3980c4e983761e85e065bf</aeskey><cdnthumbaeskey>15e6af19df3980c4e983761e85e065bf</cdnthumbaeskey><cdnthumbmd5>43692fcd974d406c33aad13fff56079f</cdnthumbmd5><encryver>1</encryver><cdnthumblength>30882</cdnthumblength><cdnthumbheight>100</cdnthumbheight><cdnthumbwidth>100</cdnthumbwidth></appattach><webviewshared><publisherId /><publisherReqId>0</publisherReqId></webviewshared><weappinfo><pagepath>?1</pagepath><username>gh_6f4e5ea45e95@app</username><appid>wx141bfb9b73c970a9</appid><type>3</type><weappiconurl>http://mmbiz.qpic.cn/mmbiz_png/jLpM6Xfk87kvrEQhHQ19tPF2QqP38GRR375xiaRJcibsrHPHKndpL2Y8lfWU2aXiamsFE3zPeGWbjRzgdP1RR0slg/640?wx_fmt=png&amp;wxfrom=200</weappiconurl><appservicetype>4</appservicetype><sharekey>GCM_EYIanFwK-4WHiKnUbj90yVR2Kbib0zpU_CRtnpwakXIkAgwICoFFIfJC3B30tw-IrLZo1yKq2T1Wk7pwv0BwpSJxl3Sgwk8OgA82uRY~</sharekey><shareId>2_wx141bfb9b73c970a9_1837218910_1722998330_1</shareId></weappinfo><websearch /></appmsg><fromusername></fromusername><scene>0</scene><appinfo><version>3</version><appname>羊了个羊</appname></appinfo><commenturl />", 49)
		c.WriteMessage(1, result)
	}

	if message.CurrentPacket.Data.AddMsg.Content == "发个图片" {
		result := SendImg(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "<msg><img aeskey=\"cd8ccfc701d1bb8d41d2dafc1809aaa8\" encryver=\"1\" cdnthumbaeskey=\"cd8ccfc701d1bb8d41d2dafc1809aaa8\" cdnthumburl=\"3057020100044b304902010002040c2ac9e502032f5081020415eff98c02046694a573042432316339303161312d376666302d343539332d623830622d3639346336646466623530610204051418020201000405004c54a100\" cdnthumblength=\"21104\" cdnthumbheight=\"200\" cdnthumbwidth=\"200\" cdnmidheight=\"0\" cdnmidwidth=\"0\" cdnhdheight=\"0\" cdnhdwidth=\"0\" cdnmidimgurl=\"3057020100044b304902010002040c2ac9e502032f5081020415eff98c02046694a573042432316339303161312d376666302d343539332d623830622d3639346336646466623530610204051418020201000405004c54a100\" length=\"85326\" md5=\"cd8ccfc701d1bb8d41d2dafc1809aaa8\" /></msg>")
		c.WriteMessage(1, result)
	}

	// 发图片的
	if message.CurrentPacket.Data.AddMsg.Content == "刺激刺激" {
		if message.CurrentPacket.Data.AddMsg.ActionUserName != "wxid_za7ku9u4uu5q21" {
			return
		}
		result, reqId := UploadCdnImg(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "https://fanruizhecn.serv00.net/fl/")
		ResponseImgMap[reqId] = _struct.ImgInfo{
			CurrentWxid:  message.CurrentWxid,
			FromUserName: message.CurrentPacket.Data.AddMsg.FromUserName,
			Type:         1,
		}
		c.WriteMessage(1, result)
	}

	// 拍一拍
	if message.CurrentPacket.Data.AddMsg.Content == "拍拍我" {
		result := SendPatMessage(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, message.CurrentPacket.Data.AddMsg.ActionUserName, 0)
		c.WriteMessage(1, result)
	}

	//if message.CurrentPacket.Data.AddMsg.Content == "听歌" {
	//	result := SendVoice(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "https://fanruizhecn.serv00.net/98.silk", 10)
	//	fmt.Println(string(result))
	//	c.WriteMessage(1, result)
	//}

	// 猜歌名游戏
	if message.CurrentPacket.Data.AddMsg.Content == "开始猜歌名" {
		if _, ok := GameStatus[message.CurrentPacket.Data.AddMsg.FromUserName]; ok {
			return
		}
		GameStatus[message.CurrentPacket.Data.AddMsg.FromUserName] = map[string]int{"status": 1, "timestamp": int(common.GetCurrentTimestamp()) + 60}
		var t map[string]struct {
			Id       string `json:"id"`
			Aid      string `json:"aid"`
			LogId    string `json:"log_id"`
			RadioUrl string `json:"radio_url"`
			Answer   string `json:"answer,omitempty"`
		}
		resp, _ := resty.New().R().Get("https://fanruizhecn.serv00.net/radio.json")
		json.Unmarshal(resp.Body(), &t)
		var key []string
		for k, _ := range t {
			key = append(key, k)
		}
		count := len(key)
		rand.Seed(int64(time.Now().Nanosecond()))
		randomNum := rand.Intn(count)

		musicGameContent := "===============开始猜歌名消息块==================\n"
		musicGameContent += "时间：" + common.GetCurrentTime() + "\n"
		if strings.Contains(message.CurrentPacket.Data.AddMsg.FromUserName, "@chatroom") {
			musicGameContent += "群名：[" + ChatroomInfo[message.CurrentPacket.Data.AddMsg.FromUserName] + "] 群id：[" + message.CurrentPacket.Data.AddMsg.FromUserName + "]\n"
		}
		musicGameContent += "用户名：[" + UserList[message.CurrentPacket.Data.AddMsg.ActionUserName] + "] 用户id：[" + message.CurrentPacket.Data.AddMsg.ActionUserName + "]\n"
		musicGameContent += "答案：[" + t[key[randomNum]].Answer + "]\n"
		musicGameContent += "地址：[" + "https://fanruizhecn.serv00.net/silk/" + t[key[randomNum]].Id + ".silk" + "]\n"
		musicGameContent += "===============开始猜歌名消息块=================="
		fmt.Println(musicGameContent)

		gameAnswer[message.CurrentPacket.Data.AddMsg.FromUserName] = t[key[randomNum]].Answer
		result := SendVoice(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "https://fanruizhecn.serv00.net/silk/"+t[key[randomNum]].Id+".silk", 10)
		c.WriteMessage(1, result)
	}

	// 没玩明白
	if message.CurrentPacket.Data.AddMsg.Content == "上传文件" {
		result, _ := UploadCdnFile(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "")
		c.WriteMessage(1, result)
	}
}

func CgiResponseProcess(info []byte, c *websocket.Conn) {
	re := regexp.MustCompile("\"ReqId\":(.*?),")
	reqInfo := re.FindStringSubmatch(string(info))
	if len(reqInfo) < 1 {
		return
	}

	// 根据这个reqId 去找对应的处理吧 不然结构体虽然相同但是类型不同
	reqId, _ := strconv.Atoi(reqInfo[1])

	fmt.Printf("请求的ID：%d\n", reqId)

	var response _struct.Response
	if err := json.Unmarshal(info, &response); err != nil {
	}
	// 是图片的
	if response.ReqId != 0 && ResponseImgMap[int(response.ReqId)].Type == 1 {
		result := SendImg(ResponseImgMap[int(response.ReqId)].CurrentWxid, ResponseImgMap[int(response.ReqId)].FromUserName, response.ResponseData)
		c.WriteMessage(1, result)
	}

	// 是文件的
	if response.ReqId != 0 && ResponseImgMap[int(response.ReqId)].Type == 2 {
		result := SendAppMessage(ResponseImgMap[int(response.ReqId)].CurrentWxid, ResponseImgMap[int(response.ReqId)].FromUserName, response.ResponseData, 49)
		c.WriteMessage(1, result)
	}

	// 获取群成员信息
	if response.ReqId != 0 && ResponseUserInfoMap[int(response.ReqId)].Type == 2 {
		type T struct {
			CgiBaseResponse struct {
				ErrMsg string `json:"ErrMsg"`
				Ret    int    `json:"Ret"`
			} `json:"CgiBaseResponse"`
			ReqId        int64 `json:"ReqId"`
			ResponseData []struct {
				MsgType         int    `json:"MsgType"`
				UserName        string `json:"UserName"`
				NickName        string `json:"NickName"`
				Signature       string `json:"Signature"`
				SmallHeadImgUrl string `json:"SmallHeadImgUrl"`
				BigHeadImgUrl   string `json:"BigHeadImgUrl"`
				Province        string `json:"Province"`
				City            string `json:"City"`
				Remark          string `json:"Remark"`
				Alias           string `json:"Alias"`
				Sex             int    `json:"Sex"`
				ContactType     int    `json:"ContactType"`
				VerifyFlag      int    `json:"VerifyFlag"`
				LabelLists      string `json:"LabelLists"`
				ChatRoomOwner   string `json:"ChatRoomOwner"`
				EncryptUsername string `json:"EncryptUsername"`
				ExtInfo         string `json:"ExtInfo"`
				ExtInfoExt      string `json:"ExtInfoExt"`
				ChatRoomMember  []struct {
					Wxid               string `json:"Wxid"`
					NickName           string `json:"NickName"`
					ChatroomMemberFlag int    `json:"ChatroomMemberFlag"`
				} `json:"ChatRoomMember"`
				Ticket          string `json:"Ticket"`
				ChatroomVersion int    `json:"ChatroomVersion"`
			} `json:"ResponseData"`
		}
		var t T
		json.Unmarshal(info, &t)
		if len(t.ResponseData) < 1 {
			return
		}
		// 将用户信息写入
		for _, v := range t.ResponseData[0].ChatRoomMember {
			UserList[v.Wxid] = v.NickName
		}
		// 将群信息写入
		ChatroomInfo[t.ResponseData[0].UserName] = t.ResponseData[0].NickName
	}
}

// 退出群聊
func exitGroup() {

}

// 加入群聊
func joinGroup(CurrentWxid string, content string, roomId string, c *websocket.Conn) {
	if strings.Contains(content, "加入了群聊") && strings.Contains(content, "邀请你加入了群聊") == false {
		re := regexp.MustCompile("\"(.*?)\"邀请\"(.*?)\"加入了群聊")
		matches := re.FindAllStringSubmatch(content, -1)
		if len(matches) < 1 {
			return
		}
		if len(matches[0]) >= 3 {
			str := "<appmsg appid=\"\" sdkver=\"0\"><title>欢迎新人[" + matches[0][2] + "]进群</title><des>邀请人 :" + matches[0][1] + "\n发送[功能]获取玩法</des><action>view</action><type>5</type><showtype>0</showtype><content /><url>https://apifox.com/apidoc/shared-edbfcebc-6263-4e87-9813-54520c1b3c19</url><dataurl /><lowurl /><lowdataurl /><recorditem /><thumburl>https://wx.qlogo.cn/mmopen/r48cSSlr7jgFutEJFpmolCux6WWZsm92KLTOmWITDvqPVIO5kLpTblfqsxuGzaZvGkgHsBOohkWuZlZuF48hRVEIcjRu1wVF/64</thumburl><messageaction /><laninfo /><md5></md5><extinfo /><sourceusername>gh_0c617dab0f5f</sourceusername><sourcedisplayname>关注公众号: 一条爱睡觉的咸鱼</sourcedisplayname><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></emoticonmd5><fileext>jpg</fileext><filekey></filekey><cdnthumburl></cdnthumburl><aeskey></aeskey><cdnthumbaeskey></cdnthumbaeskey><cdnthumbmd5></cdnthumbmd5><encryver>1</encryver><cdnthumblength>1830</cdnthumblength><cdnthumbheight>100</cdnthumbheight><cdnthumbwidth>100</cdnthumbwidth></appattach><weappinfo><pagepath /><username /><appid /><appservicetype>0</appservicetype></weappinfo><websearch /></appmsg><fromusername>wxid_k9i0ws42v8bt12</fromusername><scene>0</scene><appinfo><version>1</version><appname /></appinfo><commenturl />"
			result := SendAppMessage(CurrentWxid, roomId, str, 49)
			e := c.WriteMessage(1, result)
			fmt.Println(e)
		}
	}

	if strings.Contains(content, "分享的二维码加入群聊") {
		re := regexp.MustCompile("\"(.*?)\"通过扫描\"(.*?)\"分享的二维码加入群聊")
		matches := re.FindAllStringSubmatch(content, -1)
		if len(matches) < 1 {
			return
		}
		if len(matches[0]) >= 3 {
			str := "<appmsg appid=\"\" sdkver=\"0\"><title>欢迎新人[" + matches[0][2] + "]进群</title><des>邀请人 :" + matches[0][1] + "\n发送[功能]获取玩法</des><action>view</action><type>5</type><showtype>0</showtype><content /><url>https://apifox.com/apidoc/shared-edbfcebc-6263-4e87-9813-54520c1b3c19</url><dataurl /><lowurl /><lowdataurl /><recorditem /><thumburl>https://wx.qlogo.cn/mmopen/r48cSSlr7jgFutEJFpmolCux6WWZsm92KLTOmWITDvqPVIO5kLpTblfqsxuGzaZvGkgHsBOohkWuZlZuF48hRVEIcjRu1wVF/64</thumburl><messageaction /><laninfo /><md5></md5><extinfo /><sourceusername>gh_0c617dab0f5f</sourceusername><sourcedisplayname>关注公众号: 一条爱睡觉的咸鱼</sourcedisplayname><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></emoticonmd5><fileext>jpg</fileext><filekey></filekey><cdnthumburl></cdnthumburl><aeskey></aeskey><cdnthumbaeskey></cdnthumbaeskey><cdnthumbmd5></cdnthumbmd5><encryver>1</encryver><cdnthumblength>1830</cdnthumblength><cdnthumbheight>100</cdnthumbheight><cdnthumbwidth>100</cdnthumbwidth></appattach><weappinfo><pagepath /><username /><appid /><appservicetype>0</appservicetype></weappinfo><websearch /></appmsg><fromusername>wxid_k9i0ws42v8bt12</fromusername><scene>0</scene><appinfo><version>1</version><appname /></appinfo><commenturl />"
			result := SendAppMessage(CurrentWxid, roomId, str, 49)
			e := c.WriteMessage(1, result)
			fmt.Println(e)
		}
	}
}

//"<appmsg appid=\"\" sdkver=\"0\"><title>欢迎新人["+matches[0][1]+"]进群</title><des>邀请人 :"+matches[0][2]+"\n发送[功能]获取玩法</des><action>view</action><type>5</type><showtype>0</showtype><content /><url>https://apifox.com/apidoc/shared-edbfcebc-6263-4e87-9813-54520c1b3c19</url><dataurl /><lowurl /><lowdataurl /><recorditem /><thumburl>https://wx.qlogo.cn/mmopen/r48cSSlr7jgFutEJFpmolCux6WWZsm92KLTOmWITDvqPVIO5kLpTblfqsxuGzaZvGkgHsBOohkWuZlZuF48hRVEIcjRu1wVF/64</thumburl><messageaction /><laninfo /><md5></md5><extinfo /><sourceusername>gh_0c617dab0f5f</sourceusername><sourcedisplayname>关注公众号: 一条爱睡觉的咸鱼</sourcedisplayname><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></emoticonmd5><fileext>jpg</fileext><filekey></filekey><cdnthumburl></cdnthumburl><aeskey></aeskey><cdnthumbaeskey></cdnthumbaeskey><cdnthumbmd5></cdnthumbmd5><encryver>1</encryver><cdnthumblength>1830</cdnthumblength><cdnthumbheight>100</cdnthumbheight><cdnthumbwidth>100</cdnthumbwidth></appattach><weappinfo><pagepath /><username /><appid /><appservicetype>0</appservicetype></weappinfo><websearch /></appmsg><fromusername>wxid_k9i0ws42v8bt12</fromusername><scene>0</scene><appinfo><version>1</version><appname /></appinfo><commenturl />"
