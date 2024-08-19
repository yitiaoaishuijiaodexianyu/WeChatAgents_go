package _struct

var Config ConfigInfo

var PlugInConfig PlugInConfigInfo

var ChatRoomConfig ChatRoomConfigInfo

var KnownGroupConfig KnownGroupConfigInfo

// ConfigInfo 配置信息
type ConfigInfo struct {
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	HttpServer struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"httpServer"`
	Robot []struct {
		ServiceHose  string `yaml:"service_hose"`
		SecurityCode string `yaml:"security_code"`
		BotWxid      string `yaml:"BotWxid"`
		AdminWxId    string `yaml:"admin_wx_id"`
	} `yaml:"robot"`
	XhConfig struct {
		AppId     string `yaml:"appid"`
		ApiSecret string `yaml:"apiSecret"`
		ApiKey    string `yaml:"apiKey"`
		Domain    string `yaml:"domain"`
		HostUrl   string `yaml:"hostUrl"`
	} `yaml:"xh_config"`
}

// PlugInConfigInfo 插件配置信息
type PlugInConfigInfo struct {
	PlugIn []struct {
		Type         string `yaml:"type"`
		PlugInName   string `yaml:"plug_in_name"`
		Method       string `yaml:"method"`
		Url          string `yaml:"url"`
		MatchingMode int    `yaml:"MatchingMode"`
	} `yaml:"PlugIn"`
}

// ChatRoomConfigInfo 群配置信息
type ChatRoomConfigInfo struct {
	Chatroom       []string `yaml:"chatroom"`
	RevokeChatroom []string `yaml:"revokeChatroom"`
}

// KnownGroupConfigInfo 已知群配置信息
type KnownGroupConfigInfo struct {
	KnownGroup []struct {
		ChatroomId   string `yaml:"chatroom_id"`
		ChatroomName string `yaml:"chatroom_name"`
	} `yaml:"knownGroup"`
}
