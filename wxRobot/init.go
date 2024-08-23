package wxRobot

import (
	_struct "WeChatAgents_go/struct"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func init() {

	getIdiom()
}

func configInit() {

}

// getIdiom 获取成语数据
func getIdiom() {
	//// 打开JSON文件
	//file, err := os.Open("./resource/idiom.json")
	//if err != nil {
	//	return
	//}
	//defer file.Close()
	//
	//// 读取文件内容
	//content, err := ioutil.ReadAll(file)
	//if err != nil {
	//	return
	//}

	resp, _ := resty.New().R().Get("https://frz.fan/resource/idiom.json")

	var idiom []_struct.Idiom
	// 解析JSON数据
	if err := json.Unmarshal(resp.Body(), &idiom); err != nil {
		return
	}
	var idiomMap = make(map[string]_struct.Idiom)
	var idiomFirstMaps = make(map[string][]_struct.Idiom)
	var idiomStrings []string
	for _, v := range idiom {
		idiomMap[v.Word] = v
		idiomStrings = append(idiomStrings, v.Word)
		// 如果存在
		if _, ok := idiomFirstMaps[v.First]; ok {
			// 追加进去
			idiomFirstMaps[v.First] = append(idiomFirstMaps[v.First], v)
		} else {
			idiomFirstMaps[v.First] = append(idiomFirstMaps[v.First], v)
		}
	}
	_struct.IdiomMap = idiomMap
	_struct.IdiomStrings = idiomStrings
	_struct.IdiomFirstMap = idiomFirstMaps
	fmt.Println("成语词典初始化成功")
}
