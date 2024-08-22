package plug_in

import (
	"WeChatAgents_go/common"
	_struct "WeChatAgents_go/struct"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Horoscope struct {
	Date       int    `json:"date"`
	Name       string `json:"name"`
	QFriend    string `json:"QFriend"`
	Color      string `json:"color"`
	Datetime   string `json:"datetime"`
	Health     string `json:"health"`
	Love       string `json:"love"`
	Work       string `json:"work"`
	Money      string `json:"money"`
	Number     int    `json:"number"`
	Summary    string `json:"summary"`
	All        string `json:"all"`
	Resultcode string `json:"resultcode"`
	ErrorCode  int    `json:"error_code"`
}

func StarSign(c *gin.Context) {
	var HoroscopeList = []string{
		"白羊座",
		"金牛座",
		"双子座",
		"巨蟹座",
		"狮子座",
		"处女座",
		"天秤座",
		"天蝎座",
		"射手座",
		"摩羯座",
		"水瓶座",
		"双鱼座",
	}

	var message _struct.Message
	if ok := c.ShouldBindJSON(&message); ok != nil {
		return
	}

	var result = _struct.PlugInResult{}
	result.Type = "text"
	result.ReceiverId = message.CurrentPacket.Data.AddMsg.FromUserName
	result.BotId = message.CurrentWxid

	key := "63b715bbff8675a30aa0ed1f8cf5b92a"
	HoroscopeName := message.CurrentPacket.Data.AddMsg.Content
	if HoroscopeName == "" {
		return
	}
	found := false
	for _, horoscope := range HoroscopeList {
		if horoscope == HoroscopeName {
			found = true
			break
		}
	}
	if found == false {
		return
	}
	urls := "http://web.juhe.cn/constellation/getAll?key=" + key + "&consName=" + url.QueryEscape(HoroscopeName) + "&type=today"
	results, err := http.Get(urls)
	//接口请求错误
	if err != nil {
		return
	}
	var horoscope Horoscope
	res, _err := io.ReadAll(results.Body)
	//读取数据错误
	if _err != nil {
		return
	}
	// 数据解析错误
	if err = json.Unmarshal(res, &horoscope); err != nil {
		return
	}
	// 接口错误
	if horoscope.ErrorCode != 0 {
		return
	}
	str := "\n"
	str += "星座: " + horoscope.Name + "\n"
	str += "时间: " + horoscope.Datetime + "\n"
	str += "健康指数: " + horoscope.Health + " 财运指数: " + horoscope.Money + "\n"
	str += "爱情指数: " + horoscope.Love + " 工作指数: " + horoscope.Work + "\n"
	str += "综合指数: " + horoscope.All + " 幸运色: " + horoscope.Color + "\n"
	str += "幸运数字: " + strconv.Itoa(horoscope.Number) + " 速配星座: " + horoscope.QFriend + "\n"
	str += "今日概述: " + horoscope.Summary + "\n"
	result.Message = str
	c.JSON(200, common.ResultCommon(0, result, "星座获取成功"))
	return
}
