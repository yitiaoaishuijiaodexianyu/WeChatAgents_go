package wxRobot

import (
	"WeChatAgents_go/common"
	"WeChatAgents_go/config"
	_struct "WeChatAgents_go/struct"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
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

var ChatroomUserInfo = make(map[string][]_struct.ChatroomUser)

// 存请求的reqId来判断 [reqid]:[类型(自定义如何处理)]
var reqType = make(map[int]int)

// GameStatus 存游戏的开始状态[群id]:[[status]:[1],[timestamp]:[时间戳]]
var GameStatus = make(map[string]map[string]int)

// 存游戏的答案 [群id]:[答案]
var gameAnswer = make(map[string]string)

// getChatRoomInfo 获取群的信息
func getChatRoomInfo(botWxId string, chatRoomId string) {
	result, reqId := _struct.GetWxIdInfo(botWxId, chatRoomId)
	ResponseUserInfoMap[reqId] = _struct.GetUserInfo{Type: 2}
	reqType[reqId] = 2
	_struct.WebSocketConn.WriteMessage(1, result)
}

// checkChatroom 检查这个群是否已知 不知道这个群的话就写入到yml中 保证 程序重启时 能提前去获取一下群成员
func checkChatroom(chatroomId string, chatroomName string) {
	for _, v := range _struct.KnownGroupConfig.KnownGroup {
		if v.ChatroomId == chatroomId {
			// 这里就不处理了 存在的话
			return
		}
	}
	// 循环结束表示不存在 这里处理一下
	config.WriteChatroomConfig(chatroomId, chatroomName)
}

// searchAtId 查找被at的人的id 目前发现有三种不同的情况
func searchAtId(xml string) string {
	atId := ""
	// 定义正则表达式模式 这是一种情况
	pattern := `<atuserlist><!\[CDATA\[,([^\]]+)\]\]></atuserlist>`
	// 使用re.FindStringSubmatch进行匹配
	match := regexp.MustCompile(pattern).FindStringSubmatch(xml)
	// 检查是否匹配成功
	if match != nil && len(match) > 1 {
		// 被at的人的id
		return match[1]
	}
	// 如果第一个模式没有匹配成功，尝试第二个模式 这是一种情况
	pattern = `<atuserlist>(.*?)</atuserlist>`
	match = regexp.MustCompile(pattern).FindStringSubmatch(xml)
	if match != nil && len(match) > 1 {
		atID := match[1]
		if atID[0] == '<' {
			// 如果atID以"<"开头，尝试使用第三个模式 这又是一种情况
			pattern = `<!\[CDATA\[([^\]]+)\]\]>`
			match = regexp.MustCompile(pattern).FindStringSubmatch(atID)
			if match != nil && len(match) > 1 {
				// 被at的人的id
				return match[1]
			}
		} else {
			// 被at的人的id
			return atID
		}
	}
	return atId
}

// MessageProcess 消息处理
func MessageProcess(message _struct.Message) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	// 收到事件(去检查是否有人退群)
	if message.CurrentPacket.Data.EventName == "ON_EVENT_CONTACT_CHANGE" {
		result, reqId := _struct.GetWxIdInfo(_struct.Config.Robot[0].BotWxid, message.CurrentPacket.Data.Contact.UserName)
		reqType[reqId] = 3
		_struct.WebSocketConn.WriteMessage(1, result)
		return
	}

	if message.CurrentPacket.Data.EventName == "ON_EVENT_PAT_MSG" {
		if strings.Contains(message.CurrentPacket.Data.Template, "拍了拍我") {
			// 使用正则表达式匹配 ${} 之间的内容
			wxidre := regexp.MustCompile(`\$\{([^}]+)\}`)
			// 查找匹配的内容
			match := wxidre.FindStringSubmatch(message.CurrentPacket.Data.Template)
			wxid := ""
			// match[1] 是提取出的 id
			if len(match) > 1 {
				wxid = match[1]
				var patArr = []string{
					"再拍就打你呦[敲打]",
					"拍我干嘛，是不是想我啦[疑问]",
					"呜呜呜，别拍啦，再拍就要傻掉啦[流泪]",
					"再拍我信不信发恐怖片吓你😱",
					"再拍我就把你绑起来打屁屁[发怒]",
					"别拍啦，别拍啦，这就给你唱首歌听",
					"表情包",
				}

				var silkArr = []string{
					"https://fanruizhecn.serv00.net/silk/2420.silk",
					"https://fanruizhecn.serv00.net/silk/2430.silk",
					"https://fanruizhecn.serv00.net/silk/2440.silk",
					"https://fanruizhecn.serv00.net/silk/2495.silk",
					"https://fanruizhecn.serv00.net/silk/1017.silk",
				}

				var emoji = []map[string]int{
					{"2ad578fcfecda0f58e90e701b49348aa": 81258},
					{"b4afd3d68a5e6ab7cc95f70be3334eb4": 31106},
					{"042779e61ac4c8dd171d1212adb2b3e5": 80537},
					{"574fd2bb851c38a66c4a6354826cc3b5": 770386},
					{"1a2042cdc6a50da28c5843d6db86a8a9": 134388},
				}
				count := len(patArr)
				rand.Seed(int64(time.Now().Nanosecond()))
				randomNum := rand.Intn(count)

				if patArr[randomNum] == "表情包" {
					ecount := len(emoji)
					rand.Seed(int64(time.Now().Nanosecond()))
					erandomNum := rand.Intn(ecount)
					randomEmoji := emoji[erandomNum]
					// 遍历 map，取出 key 和 value
					for key, value := range randomEmoji {
						result, _ := _struct.SendEmoji(message.CurrentWxid, message.CurrentPacket.Data.FromUserName, key, value)
						_struct.WebSocketConn.WriteMessage(1, result)
					}
					return
				}

				result, _ := _struct.SendText(message.CurrentWxid, message.CurrentPacket.Data.FromUserName, "@"+UserList[wxid]+" "+patArr[randomNum], wxid)
				_struct.WebSocketConn.WriteMessage(1, result)

				if patArr[randomNum] == "别拍啦，别拍啦，这就给你唱首歌听" {
					mcount := len(silkArr)
					rand.Seed(int64(time.Now().Nanosecond()))
					mRandomNum := rand.Intn(mcount)
					results, _ := _struct.SendVoice(message.CurrentWxid, message.CurrentPacket.Data.FromUserName, silkArr[mRandomNum], 10)
					_struct.WebSocketConn.WriteMessage(1, results)
				}
				// https://fanruizhecn.serv00.net/silk/2420.silk
			}
		}
		return
	}

	// 入群欢迎
	if message.CurrentPacket.Data.AddMsg.MsgType == 10000 {
		joinGroup(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.Content, message.CurrentPacket.Data.AddMsg.FromUserName)
		return
	}

	// 如果检测到不存在已知的群中 获取一次用户消息
	if _, ok := ChatroomInfo[message.CurrentPacket.Data.AddMsg.FromUserName]; !ok {
		if strings.Contains(message.CurrentPacket.Data.AddMsg.FromUserName, "@chatroom") {
			getChatRoomInfo(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName)
		}
		time.Sleep(time.Second * 1)
	}

	content := "===============消息块==================\n"
	content += "时间：" + common.GetCurrentTime() + "\n"
	strLength := len(message.CurrentPacket.Data.AddMsg.Content)
	if len(message.CurrentPacket.Data.AddMsg.Content) > 99 {
		strLength = 99
	}
	if message.CurrentPacket.Data.AddMsg.FromUserName == _struct.Config.Robot[0].BotWxid {
		content += "机器人发言：\n"
		if strings.Contains(message.CurrentPacket.Data.AddMsg.ToUserName, "@chatroom") {
			content += "群名：[" + ChatroomInfo[message.CurrentPacket.Data.AddMsg.ToUserName] + "] 群id：[" + message.CurrentPacket.Data.AddMsg.ToUserName + "]\n"
		}
		content += "用户名：[" + UserList[message.CurrentPacket.Data.AddMsg.FromUserName] + "] 用户id：[" + message.CurrentPacket.Data.AddMsg.FromUserName + "]\n"

		content += "我的发言：[" + message.CurrentPacket.Data.AddMsg.Content[0:strLength] + "] 消息Id：" + strconv.Itoa(int(message.CurrentPacket.Data.AddMsg.NewMsgId)) + "\n"
	} else {
		if strings.Contains(message.CurrentPacket.Data.AddMsg.FromUserName, "@chatroom") {
			content += "群名：[" + ChatroomInfo[message.CurrentPacket.Data.AddMsg.FromUserName] + "] 群id：[" + message.CurrentPacket.Data.AddMsg.FromUserName + "]\n"
		}
		content += "用户名：[" + UserList[message.CurrentPacket.Data.AddMsg.ActionUserName] + "] 用户id：[" + message.CurrentPacket.Data.AddMsg.ActionUserName + "]\n"
		content += "群友发言：[" + message.CurrentPacket.Data.AddMsg.Content[0:strLength] + "] 消息Id：" + strconv.Itoa(int(message.CurrentPacket.Data.AddMsg.NewMsgId)) + "\n"
	}
	content += "===============消息块=================="
	if int(message.CurrentPacket.Data.AddMsg.NewMsgId) != 0 {
		fmt.Println(content)
	}
	if message.CurrentPacket.Data.AddMsg.Content == "清空运行缓存" && message.CurrentPacket.Data.AddMsg.ActionUserName == _struct.Config.Robot[0].AdminWxId {
		fmt.Println("清空缓存成功")
	}

	if _, ok := UserList[message.CurrentPacket.Data.AddMsg.ActionUserName]; ok {
		message.CurrentPacket.Data.AddMsg.ActionNickName = UserList[message.CurrentPacket.Data.AddMsg.ActionUserName]
	}

	var plugIn = _struct.PlugInConfig
	for _, v := range plugIn.PlugIn {
		if v.MatchingMode == 1 {
			if message.CurrentPacket.Data.AddMsg.Content != v.PlugInName {
				continue
			}
			requestData, _ := json.Marshal(message)
			response, _ := resty.New().R().SetBody(requestData).Post(v.Url)
			resultHandle(response.Body())
			break
		}
		if v.MatchingMode == 2 {
			plugInLength := len(v.PlugInName)
			userMessage := message.CurrentPacket.Data.AddMsg.Content
			if len(userMessage) < plugInLength {
				continue
			}
			if userMessage[0:plugInLength] != v.PlugInName {
				continue
			}

			if v.PlugInName == "阿呆" {
				message.CurrentPacket.Data.AddMsg.Content = message.CurrentPacket.Data.AddMsg.Content[6:]
			}

			requestData, _ := json.Marshal(message)
			response, _ := resty.New().R().SetBody(requestData).Post(v.Url)
			resultHandle(response.Body())
			break
		}
		if v.MatchingMode == 3 {
			if strings.Contains(message.CurrentPacket.Data.AddMsg.Content, v.PlugInName) {
				requestData, _ := json.Marshal(message)
				response, _ := resty.New().R().SetBody(requestData).Post(v.Url)
				resultHandle(response.Body())
				break
			}
		}
		if v.MatchingMode == 4 {
			plugInNameArr := strings.Split(v.PlugInName, "|")
			for _, vv := range plugInNameArr {
				if message.CurrentPacket.Data.AddMsg.Content == vv {
					requestData, _ := json.Marshal(message)
					response, _ := resty.New().R().SetBody(requestData).Post(v.Url)
					resultHandle(response.Body())
					break
				}
			}
		}
	}

	// 判断猜歌名
	if _, ok := gameAnswer[message.CurrentPacket.Data.AddMsg.FromUserName]; ok {
		if message.CurrentPacket.Data.AddMsg.Content == gameAnswer[message.CurrentPacket.Data.AddMsg.FromUserName] {
			delete(gameAnswer, message.CurrentPacket.Data.AddMsg.FromUserName)
			// 在访问共享资源前加锁
			mu.Lock()
			result, _ := _struct.SendText(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "@"+UserList[message.CurrentPacket.Data.AddMsg.ActionUserName]+" 恭喜回答正确："+message.CurrentPacket.Data.AddMsg.Content, message.CurrentPacket.Data.AddMsg.ActionUserName)
			_struct.WebSocketConn.WriteMessage(1, result)
			message.CurrentPacket.Data.AddMsg.Content = "开始猜歌名"
			delete(GameStatus, message.CurrentPacket.Data.AddMsg.FromUserName)
			time.Sleep(time.Second * 1)
			// 自己去调用一次开始猜歌名
			MessageProcess(message)
			// 释放锁
			mu.Unlock()
		}
	}

	// 踢人
	if message.CurrentPacket.Data.AddMsg.Content == "踢了他" && message.CurrentPacket.Data.AddMsg.ActionUserName == "wxid_za7ku9u4uu5q21" && message.CurrentPacket.Data.AddMsg.AtId != "" {
		result, _ := _struct.DelChatroomMember(_struct.Config.Robot[0].BotWxid, message.CurrentPacket.Data.AddMsg.AtId, message.CurrentPacket.Data.AddMsg.FromUserName)
		_struct.WebSocketConn.WriteMessage(1, result)
	}

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
		//resp, _ := resty.New().R().Get("https://fanruizhecn.serv00.net/radio.json")
		resp, _ := resty.New().R().Get("https://frz.fan/resource/radio.json")
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
		//musicGameContent += "地址：[" + "https://fanruizhecn.serv00.net/silk/" + t[key[randomNum]].Id + ".silk" + "]\n"
		musicGameContent += "地址：[" + "https://frz.fan/resource/silk/" + t[key[randomNum]].Id + ".silk" + "]\n"
		musicGameContent += "===============开始猜歌名消息块=================="
		fmt.Println(musicGameContent)

		gameAnswer[message.CurrentPacket.Data.AddMsg.FromUserName] = t[key[randomNum]].Answer
		//result, _ := _struct.SendVoice(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "https://fanruizhecn.serv00.net/silk/"+t[key[randomNum]].Id+".silk", 10)
		result, _ := _struct.SendVoice(message.CurrentWxid, message.CurrentPacket.Data.AddMsg.FromUserName, "https://frz.fan/resource/silk/"+t[key[randomNum]].Id+".silk", 10)
		_struct.WebSocketConn.WriteMessage(1, result)
	}
}

// CgiResponseProcess 发送消息后要在这里处理回来的信息 reqId 为处理的标识
func CgiResponseProcess(info []byte) {
	re := regexp.MustCompile("\"ReqId\":(.*?),")
	reqInfo := re.FindStringSubmatch(string(info))
	if len(reqInfo) < 1 {
		return
	}
	// 根据这个reqId 去找对应的处理吧 不然结构体虽然相同但是类型不同
	reqId, _ := strconv.Atoi(reqInfo[1])
	//fmt.Printf("请求的ID：%d\n", reqId)
	var response _struct.Response
	json.Unmarshal(info, &response)

	// 是图片的
	if response.ReqId != 0 && ResponseImgMap[int(response.ReqId)].Type == 1 {
		result, reqIds := _struct.SendImg(ResponseImgMap[int(response.ReqId)].CurrentWxid, ResponseImgMap[int(response.ReqId)].FromUserName, response.ResponseData)
		reqStruct := _struct.ReqIdMap[reqId]            // 获取结构体副本
		reqStruct.Status = response.CgiBaseResponse.Ret // 修改副本
		reqStruct.NewReqId = reqIds                     //  要将第二个放进去 不然检测不到图片发送成功的回调
		_struct.ReqIdMap[reqId] = reqStruct             //  重新赋值回 map
		_struct.WebSocketConn.WriteMessage(1, result)
	}

	// 是文件的
	if response.ReqId != 0 && ResponseImgMap[int(response.ReqId)].Type == 2 {
		result, reqIds := _struct.SendAppMessage(ResponseImgMap[int(response.ReqId)].CurrentWxid, ResponseImgMap[int(response.ReqId)].FromUserName, response.ResponseData, 49)
		reqStruct := _struct.ReqIdMap[reqId]            // 获取结构体副本
		reqStruct.Status = response.CgiBaseResponse.Ret // 修改副本
		reqStruct.NewReqId = reqIds                     //  要将第二个放进去 不然检测不到图片发送成功的回调
		_struct.ReqIdMap[reqId] = reqStruct             //  重新赋值回 map
		_struct.WebSocketConn.WriteMessage(1, result)
	}

	// 获取群成员信息
	if response.ReqId != 0 && ResponseUserInfoMap[int(response.ReqId)].Type == 2 {
		var t _struct.SearchChatroomInfo
		json.Unmarshal(info, &t)
		if len(t.ResponseData) < 1 {
			return
		}
		// 将用户信息写入
		var Userinfo []_struct.ChatroomUser
		for _, v := range t.ResponseData[0].ChatRoomMember {
			UserList[v.Wxid] = v.NickName
			Userinfo = append(Userinfo, _struct.ChatroomUser{
				WxId:     v.Wxid,
				Username: v.NickName,
			})
		}
		// 写入检测退群检测
		ChatroomUserInfo[t.ResponseData[0].UserName] = Userinfo
		// 将群信息写入
		ChatroomInfo[t.ResponseData[0].UserName] = t.ResponseData[0].NickName
		go checkChatroom(t.ResponseData[0].UserName, t.ResponseData[0].NickName)
	}

	// 这里是拿到reqType 为3 的 这里可以判断谁退群了
	if v, ok := reqType[reqId]; ok && v == 3 {
		var t _struct.SearchChatroomInfo
		json.Unmarshal(info, &t)
		if len(t.ResponseData) < 1 {
			return
		}
		if len(ChatroomUserInfo[t.ResponseData[0].UserName]) != len(t.ResponseData[0].ChatRoomMember) {
			//fmt.Println(t.ResponseData[0].ChatRoomMember)
			// 之前还有多少人
			var oldUser []string
			for _, vv := range ChatroomUserInfo[t.ResponseData[0].UserName] {
				oldUser = append(oldUser, vv.WxId)
			}
			// 当前还有多少人
			var newUser []string
			for _, vv := range t.ResponseData[0].ChatRoomMember {
				newUser = append(newUser, vv.Wxid)
			}
			// 存放旧的切片中多出来的用户
			var leftUsers []string

			// 遍历旧的切片，找出那些不在新的切片中的用户
			for _, old := range oldUser {
				found := false
				for _, newU := range newUser {
					if old == newU {
						found = true
						break
					}
				}
				// 如果没有在新切片中找到，说明该用户已离开
				if !found {
					leftUsers = append(leftUsers, old)
				}
			}
			if len(leftUsers) != 0 {
				for _, v := range leftUsers {
					for _, vv := range ChatroomUserInfo[t.ResponseData[0].UserName] {
						if v == vv.WxId {
							str := "<appmsg appid=\"\" sdkver=\"0\"><title>[" + vv.Username + "]退出了群聊</title><des>" + v + "\n" + common.GetCurrentTime() + "</des><action>view</action><type>5</type><showtype>0</showtype><content /><url>https://apifox.com/apidoc/shared-edbfcebc-6263-4e87-9813-54520c1b3c19</url><dataurl /><lowurl /><lowdataurl /><recorditem /><thumburl>https://wx.qlogo.cn/mmopen/r48cSSlr7jgFutEJFpmolCux6WWZsm92KLTOmWITDvqPVIO5kLpTblfqsxuGzaZvGkgHsBOohkWuZlZuF48hRVEIcjRu1wVF/64</thumburl><messageaction /><laninfo /><md5></md5><extinfo /><sourceusername>gh_0c617dab0f5f</sourceusername><sourcedisplayname>关注公众号: 一条爱睡觉的咸鱼</sourcedisplayname><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></emoticonmd5><fileext>jpg</fileext><filekey></filekey><cdnthumburl></cdnthumburl><aeskey></aeskey><cdnthumbaeskey></cdnthumbaeskey><cdnthumbmd5></cdnthumbmd5><encryver>1</encryver><cdnthumblength>1830</cdnthumblength><cdnthumbheight>100</cdnthumbheight><cdnthumbwidth>100</cdnthumbwidth></appattach><weappinfo><pagepath /><username /><appid /><appservicetype>0</appservicetype></weappinfo><websearch /></appmsg><fromusername>wxid_k9i0ws42v8bt12</fromusername><scene>0</scene><appinfo><version>1</version><appname /></appinfo><commenturl />"
							result, _ := _struct.SendAppMessage(_struct.Config.Robot[0].BotWxid, t.ResponseData[0].UserName, str, 49)
							_struct.WebSocketConn.WriteMessage(1, result)
							getChatRoomInfo(_struct.Config.Robot[0].BotWxid, t.ResponseData[0].UserName)
						}
					}
				}
			}
		}
	}
	reqStruct := _struct.ReqIdMap[reqId]            // 获取结构体副本
	reqStruct.Status = response.CgiBaseResponse.Ret // 修改副本
	_struct.ReqIdMap[reqId] = reqStruct             // 重新赋值回 map
}

// 加入群聊
func joinGroup(CurrentWxid string, content string, roomId string) {
	if strings.Contains(content, "加入了群聊") && strings.Contains(content, "邀请你加入了群聊") == false {
		re := regexp.MustCompile("\"(.*?)\"邀请\"(.*?)\"加入了群聊")
		matches := re.FindAllStringSubmatch(content, -1)
		if len(matches) < 1 {
			return
		}
		if len(matches[0]) >= 3 {
			str := "<appmsg appid=\"\" sdkver=\"0\"><title>欢迎新人[" + matches[0][2] + "]进群</title><des>邀请人 :" + matches[0][1] + "\n发送[功能]获取玩法</des><action>view</action><type>5</type><showtype>0</showtype><content /><url>https://apifox.com/apidoc/shared-edbfcebc-6263-4e87-9813-54520c1b3c19</url><dataurl /><lowurl /><lowdataurl /><recorditem /><thumburl>https://wx.qlogo.cn/mmopen/r48cSSlr7jgFutEJFpmolCux6WWZsm92KLTOmWITDvqPVIO5kLpTblfqsxuGzaZvGkgHsBOohkWuZlZuF48hRVEIcjRu1wVF/64</thumburl><messageaction /><laninfo /><md5></md5><extinfo /><sourceusername>gh_0c617dab0f5f</sourceusername><sourcedisplayname>关注公众号: 一条爱睡觉的咸鱼</sourcedisplayname><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></emoticonmd5><fileext>jpg</fileext><filekey></filekey><cdnthumburl></cdnthumburl><aeskey></aeskey><cdnthumbaeskey></cdnthumbaeskey><cdnthumbmd5></cdnthumbmd5><encryver>1</encryver><cdnthumblength>1830</cdnthumblength><cdnthumbheight>100</cdnthumbheight><cdnthumbwidth>100</cdnthumbwidth></appattach><weappinfo><pagepath /><username /><appid /><appservicetype>0</appservicetype></weappinfo><websearch /></appmsg><fromusername>wxid_k9i0ws42v8bt12</fromusername><scene>0</scene><appinfo><version>1</version><appname /></appinfo><commenturl />"
			result, _ := _struct.SendAppMessage(CurrentWxid, roomId, str, 49)
			_struct.WebSocketConn.WriteMessage(1, result)
			result, _ = _struct.SendVoice(CurrentWxid, roomId, "https://frz.fan/resource/rqhy.silk", 8)
			_struct.WebSocketConn.WriteMessage(1, result)
			getChatRoomInfo(_struct.Config.Robot[0].BotWxid, roomId)
		}
	}

	if strings.Contains(content, "分享的二维码加入群聊") {
		re := regexp.MustCompile("\"(.*?)\"通过扫描\"(.*?)\"分享的二维码加入群聊")
		matches := re.FindAllStringSubmatch(content, -1)
		if len(matches) < 1 {
			return
		}
		if len(matches[0]) >= 3 {
			str := "<appmsg appid=\"\" sdkver=\"0\"><title>欢迎新人[" + matches[0][1] + "]进群</title><des>邀请人 :" + matches[0][2] + "\n发送[功能]获取玩法</des><action>view</action><type>5</type><showtype>0</showtype><content /><url>https://apifox.com/apidoc/shared-edbfcebc-6263-4e87-9813-54520c1b3c19</url><dataurl /><lowurl /><lowdataurl /><recorditem /><thumburl>https://wx.qlogo.cn/mmopen/r48cSSlr7jgFutEJFpmolCux6WWZsm92KLTOmWITDvqPVIO5kLpTblfqsxuGzaZvGkgHsBOohkWuZlZuF48hRVEIcjRu1wVF/64</thumburl><messageaction /><laninfo /><md5></md5><extinfo /><sourceusername>gh_0c617dab0f5f</sourceusername><sourcedisplayname>关注公众号: 一条爱睡觉的咸鱼</sourcedisplayname><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></emoticonmd5><fileext>jpg</fileext><filekey></filekey><cdnthumburl></cdnthumburl><aeskey></aeskey><cdnthumbaeskey></cdnthumbaeskey><cdnthumbmd5></cdnthumbmd5><encryver>1</encryver><cdnthumblength>1830</cdnthumblength><cdnthumbheight>100</cdnthumbheight><cdnthumbwidth>100</cdnthumbwidth></appattach><weappinfo><pagepath /><username /><appid /><appservicetype>0</appservicetype></weappinfo><websearch /></appmsg><fromusername>wxid_k9i0ws42v8bt12</fromusername><scene>0</scene><appinfo><version>1</version><appname /></appinfo><commenturl />"
			result, _ := _struct.SendAppMessage(CurrentWxid, roomId, str, 49)
			_struct.WebSocketConn.WriteMessage(1, result)
			result, _ = _struct.SendVoice(CurrentWxid, roomId, "https://frz.fan/resource/rqhy.silk", 8)
			_struct.WebSocketConn.WriteMessage(1, result)
			getChatRoomInfo(_struct.Config.Robot[0].BotWxid, roomId)
		}
	}
}

func GetKnownGroupInfo() {
	for _, v := range _struct.KnownGroupConfig.KnownGroup {
		getChatRoomInfo(_struct.Config.Robot[0].BotWxid, v.ChatroomId)
	}
}

// 处理插件返回的结果
func resultHandle(result []byte) {
	type HttpResult struct {
		Code    int                  `json:"code"`
		Data    _struct.PlugInResult `json:"data"`
		Message string               `json:"message"`
	}
	var response HttpResult
	json.Unmarshal(result, &response)
	if response.Code == 0 {
		// 发送文本消息
		if response.Data.Type == "text" {
			var text []byte
			if response.Data.AtIds != "" {
				text, _ = _struct.SendText(_struct.Config.Robot[0].BotWxid, response.Data.ReceiverId, response.Data.Message, response.Data.AtIds)
			} else {
				text, _ = _struct.SendText(_struct.Config.Robot[0].BotWxid, response.Data.ReceiverId, response.Data.Message, "")
			}
			_struct.WebSocketConn.WriteMessage(1, text)
		}

		// 语言消息
		if response.Data.Type == "voice" {

		}

		// 发送图片消息
		if response.Data.Type == "image" {
			image, reqId := _struct.UploadCdnImg(_struct.Config.Robot[0].BotWxid, response.Data.ReceiverId, response.Data.Url)
			ResponseImgMap[reqId] = _struct.ImgInfo{
				CurrentWxid:  _struct.Config.Robot[0].BotWxid,
				FromUserName: response.Data.ReceiverId,
				Type:         1,
			}
			_struct.WebSocketConn.WriteMessage(1, image)
		}

		// 发送拍一拍消息
		if response.Data.Type == "pat" {
			pat, _ := _struct.SendPatMessage(_struct.Config.Robot[0].BotWxid, response.Data.ReceiverId, response.Data.PatId, 0)
			_struct.WebSocketConn.WriteMessage(1, pat)
		}

		// 发送emoji消息
		if response.Data.Type == "emoji" {
			emoji, _ := _struct.SendEmoji(_struct.Config.Robot[0].BotWxid, response.Data.ReceiverId, response.Data.EmojiMd5, response.Data.EmojiLength)
			_struct.WebSocketConn.WriteMessage(1, emoji)
		}

		// 发送app消息
		if response.Data.Type == "appMsg" {
			appMsg, _ := _struct.SendAppMessage(_struct.Config.Robot[0].BotWxid, response.Data.ReceiverId, response.Data.Xml, 49)
			_struct.WebSocketConn.WriteMessage(1, appMsg)
		}
	}
}
