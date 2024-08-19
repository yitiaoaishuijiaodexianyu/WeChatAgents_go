package plug_in

import (
	"WeChatAgents_go/common"
	_struct "WeChatAgents_go/struct"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"strings"
)

func RequestSong(c *gin.Context) {
	var message _struct.Message
	if ok := c.ShouldBindJSON(&message); ok != nil {
		return
	}
	var result = _struct.PlugInResult{}
	if message.CurrentPacket.Data.AddMsg.Content[0:6] != "点歌" {
		return
	}
	MusicName := strings.Replace(message.CurrentPacket.Data.AddMsg.Content, "点歌", "", -1)
	MusicName = strings.Replace(MusicName, " ", "", -1)
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
	result.Xml = "<appmsg appid=\"wx79f2c4418704b4f8\" sdkver=\"0\"><title>" + t.PrimitiveData.Data.Title + "</title><des>" + t.PrimitiveData.Data.Singer + "</des><action>view</action><type>3</type><showtype>0</showtype><content /><url>" + t.PrimitiveData.Data.Link + "</url><dataurl>" + t.PrimitiveData.Data.MusicUrl + "</dataurl><lowurl>" + t.PrimitiveData.Data.Link + "</lowurl><lowdataurl>" + t.PrimitiveData.Data.MusicUrl + "</lowdataurl><recorditem /><thumburl>" + t.PrimitiveData.Data.Cover + "</thumburl><messageaction /><laninfo /><extinfo /><sourceusername /><sourcedisplayname /><commenturl /><appattach><totallen>0</totallen><attachid /><emoticonmd5></emoticonmd5><fileext /><aeskey></aeskey></appattach><webviewshared><publisherId /><publisherReqId>0</publisherReqId></webviewshared><weappinfo><pagepath /><username /><appid /><appservicetype>0</appservicetype></weappinfo><websearch /><songalbumurl>" + t.PrimitiveData.Data.Cover + "</songalbumurl></appmsg><fromusername></fromusername><scene>0</scene><appinfo><version>57</version><appname>酷狗音乐</appname></appinfo><commenturl />"
	result.Type = "appMsg"
	result.ReceiverId = message.CurrentPacket.Data.AddMsg.FromUserName
	result.BotId = message.CurrentWxid
	c.JSON(200, common.ResultCommon(0, result, "点歌成功"))
	return
}
