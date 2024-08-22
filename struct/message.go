package _struct

// Message 接收消息的结构体
type Message struct {
	CurrentPacket struct {
		WebConnId string `json:"WebConnId"` // 不知道干啥的一直双空
		Data      struct {
			AddMsg struct {
				MsgId          int         `json:"MsgId"`
				FromUserName   string      `json:"FromUserName"`   //消息来自哪里 群id或者用户id
				ToUserName     string      `json:"ToUserName"`     //目前看到是自己的wx_id
				MsgType        int         `json:"MsgType"`        // 消息类型
				Content        string      `json:"Content"`        //消息内容
				Status         int         `json:"Status"`         //不知道
				ImgStatus      int         `json:"ImgStatus"`      //不知道
				ImgBuf         interface{} `json:"ImgBuf"`         //不知道
				CreateTime     int         `json:"CreateTime"`     //一个时间戳
				MsgSource      string      `json:"MsgSource"`      // 原消息
				PushContent    string      `json:"PushContent"`    // 未知
				NewMsgId       int64       `json:"NewMsgId"`       // 消息的id 应该是
				NewMsgIdExt    string      `json:"NewMsgIdExt"`    // 也是消息的id
				ActionUserName string      `json:"ActionUserName"` // 发消息的人的wx_id
				ActionNickName string      `json:"ActionNickName"` // 发消息的人的微信昵称--为空需要调用查询接口
				ChatroomName   string      `json:"chatroom_name"`  // 群名称
				AtId           string      `json:"at_id"`
				RawContent     string      `json:"raw_content"`
			} `json:"AddMsg"`
			Contact struct {
				MsgType         int         `json:"MsgType"`
				UserName        string      `json:"UserName"`
				NickName        string      `json:"NickName"`
				Signature       string      `json:"Signature"`
				SmallHeadImgUrl string      `json:"SmallHeadImgUrl"`
				BigHeadImgUrl   string      `json:"BigHeadImgUrl"`
				Province        string      `json:"Province"`
				City            string      `json:"City"`
				Remark          string      `json:"Remark"`
				Alias           string      `json:"Alias"`
				Sex             int         `json:"Sex"`
				ContactType     int         `json:"ContactType"`
				VerifyFlag      int         `json:"VerifyFlag"`
				LabelLists      string      `json:"LabelLists"`
				ChatRoomOwner   string      `json:"ChatRoomOwner"`
				EncryptUsername string      `json:"EncryptUsername"`
				ExtInfo         string      `json:"ExtInfo"`
				ExtInfoExt      string      `json:"ExtInfoExt"`
				ChatRoomMember  interface{} `json:"ChatRoomMember"`
				Ticket          string      `json:"Ticket"`
				ChatroomVersion int         `json:"ChatroomVersion"`
			} `json:"Contact"`
			EventName      string `json:"EventName"` //事件名称
			ChatUserName   string `json:"ChatUserName"`
			FromUserName   string `json:"FromUserName"`
			PattedUserName string `json:"PattedUserName"`
			Template       string `json:"Template"`
		} `json:"Data"`
	} `json:"CurrentPacket"`
	CurrentWxid string `json:"CurrentWxid"` // 挂机的机器人的wxid 可作用与多开
	UUid        string `json:"UUid"`        //不知道干啥的- -
}

type Response struct {
	CgiBaseResponse struct {
		ErrMsg string `json:"ErrMsg"`
		Ret    int    `json:"Ret"`
	} `json:"CgiBaseResponse"`
	ReqId        int64  `json:"ReqId"`
	ResponseData string `json:"ResponseData"`
}

type SearchChatroomInfo struct {
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

type PlugInResult struct {
	ReceiverId    string `json:"receiver_id"`
	Type          string `json:"type"`
	Message       string `json:"message"`
	AtIds         string `json:"at_ids"`
	Url           string `json:"url"`
	EmojiMd5      string `json:"emoji_md5"`
	EmojiLength   int    `json:"emoji_length"`
	PatId         string `json:"pat_id"`
	Xml           string `json:"xml"`
	BotId         string `json:"bot_id"`
	UserWxId      string `json:"user_wx_id"`
	IsGame        int    `json:"is_game"`
	GameStartName string `json:"game_start_name"`
	Answer        string `json:"answer"`
	GameEndTime   int    `json:"game_end_time"`
}

type ChatroomUser struct {
	WxId     string `json:"wx_id"`
	Username string `json:"username"`
}

type ImgInfo struct {
	CurrentWxid  string
	FromUserName string
	Type         int
}

type GetUserInfo struct {
	Type int
}

type GameInfo struct {
	Answer        string
	GameStartName string
	GameEndTime   int
	Status        int
}
