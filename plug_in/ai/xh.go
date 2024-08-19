package ai

import (
	"WeChatAgents_go/common"
	_struct "WeChatAgents_go/struct"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/**
 *  WebAPI 接口调用示例 接口文档（必看）：https://www.xfyun.cn/doc/spark/Web.html
 * 错误码链接：https://www.xfyun.cn/doc/spark/%E6%8E%A5%E5%8F%A3%E8%AF%B4%E6%98%8E.html（code返回错误码时必看）
 * @author iflytek
 */

//"general":     "wss://spark-api.xf-yun.com/v1.1/chat",
//"generalv2":   "wss://spark-api.xf-yun.com/v2.1/chat",
//"generalv3":   "wss://spark-api.xf-yun.com/v3.1/chat",
//"generalv3.5": "wss://spark-api.xf-yun.com/v3.5/chat",

func XhAi(c *gin.Context) {
	var message _struct.Message
	var result = _struct.PlugInResult{}
	if ok := c.ShouldBindJSON(&message); ok != nil {
		c.JSON(200, common.ResultCommon(1, result, "ai回复失败"))
		return
	}

	result.Type = "text"
	result.ReceiverId = message.CurrentPacket.Data.AddMsg.FromUserName
	result.BotId = message.CurrentWxid

	d := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	//握手并建立websocket 连接
	conn, resp, err := d.Dial(assembleAuthUrl1(_struct.Config.XhConfig.HostUrl, _struct.Config.XhConfig.ApiKey, _struct.Config.XhConfig.ApiSecret), nil)
	if err != nil {
		c.JSON(200, common.ResultCommon(1, result, "ai回复失败"))
		return
	} else if resp.StatusCode != 101 {
		c.JSON(200, common.ResultCommon(1, result, "ai回复失败"))
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
		}
	}()
	go func() {
		data := genParams1(message.CurrentPacket.Data.AddMsg.Content)
		conn.WriteJSON(data)
	}()

	var answer = ""
	//获取返回的数据
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var data map[string]interface{}
		err1 := json.Unmarshal(msg, &data)
		if err1 != nil {
			c.JSON(200, common.ResultCommon(1, result, "ai回复失败"))
			return
		}
		//解析数据
		if _, ok := data["payload"]; ok {
			//fmt.Println("存在aaa")
		} else {
			continue
		}
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			c.JSON(200, common.ResultCommon(1, result, "ai回复失败"))
			return
		}
		status := choices["status"].(float64)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		if status != 2 {
			answer += content
		} else {
			answer += content
			usage := payload["usage"].(map[string]interface{})
			temp := usage["text"].(map[string]interface{})
			_ = temp["total_tokens"].(float64)
			conn.Close()
			break
		}

	}
	//输出返回结果
	if strings.Contains(message.CurrentPacket.Data.AddMsg.FromUserName, "@chatroom") {
		result.Message = "@" + message.CurrentPacket.Data.AddMsg.ActionNickName + " " + answer
		result.AtIds = message.CurrentPacket.Data.AddMsg.ActionUserName
	} else {
		result.Message = answer
	}
	c.JSON(200, common.ResultCommon(0, result, "ai回复成功"))
	return
}

// 生成参数
func genParams1(question string) map[string]interface{} { // 根据实际情况修改返回的数据结构和字段名

	messages := []Message{
		{Role: "user", Content: question},
	}

	data := map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
		"header": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"app_id": _struct.Config.XhConfig.AppId, // 根据实际情况修改返回的数据结构和字段名
		},
		"parameter": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"chat": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"domain":      _struct.Config.XhConfig.Domain, // 根据实际情况修改返回的数据结构和字段名
				"temperature": float64(0.8),                   // 根据实际情况修改返回的数据结构和字段名
				"top_k":       int64(6),                       // 根据实际情况修改返回的数据结构和字段名
				"max_tokens":  int64(2048),                    // 根据实际情况修改返回的数据结构和字段名
				"auditing":    "default",                      // 根据实际情况修改返回的数据结构和字段名
			},
		},
		"payload": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"message": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"text": messages, // 根据实际情况修改返回的数据结构和字段名
			},
		},
	}
	return data // 根据实际情况修改返回的数据结构和字段名
}

// 创建鉴权url  apikey 即 hmac username
func assembleAuthUrl1(hosturl string, apiKey, apiSecret string) string {
	ul, err := url.Parse(hosturl)
	if err != nil {
		fmt.Println(err)
	}
	//签名时间
	date := time.Now().UTC().Format(time.RFC1123)
	//date = "Tue, 28 May 2019 09:10:42 MST"
	//参与签名的字段 host ,date, request-line
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	//拼接签名字符串
	sgin := strings.Join(signString, "\n")
	// fmt.Println(sgin)
	//签名结果
	sha := HmacWithShaTobase64("hmac-sha256", sgin, apiSecret)
	// fmt.Println(sha)
	//构建请求参数 此时不需要urlencoding
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	//将请求参数使用base64编码
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))

	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	//将编码后的字符串url encode后添加到url后面
	callurl := hosturl + "?" + v.Encode()
	return callurl
}

func HmacWithShaTobase64(algorithm, data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

func readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
