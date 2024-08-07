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
			} `json:"AddMsg"`
			EventName string `json:"EventName"` //事件名称
		} `json:"Data"`
	} `json:"CurrentPacket"`
	CurrentWxid string `json:"CurrentWxid"` // 挂机的机器人的wxid
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

type ImgInfo struct {
	CurrentWxid  string
	FromUserName string
	Type         int
}

type GetUserInfo struct {
	Type int
}
