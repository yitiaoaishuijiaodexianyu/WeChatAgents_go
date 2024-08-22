package config

import (
	_struct "WeChatAgents_go/struct"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

var mu sync.Mutex

func InitConfig() {
	GetConfigInfo()
	GetChatroomConfig()
	GetPlugInConfig()
	GetKnownGroupConfig()
}

func GetConfigInfo() {
	// 读取文件
	data, err := ioutil.ReadFile("./config/config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 创建 Config 实例
	var config _struct.ConfigInfo

	// 解析 YAML 数据
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	_struct.Config = config
}

func GetChatroomConfig() {
	// 读取文件
	data, err := ioutil.ReadFile("./config/chatroom_config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 创建 Config 实例
	var config _struct.ChatRoomConfigInfo

	// 解析 YAML 数据
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	_struct.ChatRoomConfig = config
}

func GetKnownGroupConfig() {
	// 读取文件
	data, err := ioutil.ReadFile("./config/KnownGroup.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 创建 Config 实例
	var config _struct.KnownGroupConfigInfo

	// 解析 YAML 数据
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	_struct.KnownGroupConfig = config
}

func GetPlugInConfig() {
	// 读取文件
	data, err := ioutil.ReadFile("./config/plug_in.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 创建 Config 实例
	var config _struct.PlugInConfigInfo

	// 解析 YAML 数据
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	_struct.PlugInConfig = config
}

// WriteChatroomConfig 写入已知群的操作
func WriteChatroomConfig(BotWxId string, chatroomId string, ChatroomName string) {
	// 在访问共享资源前加锁
	mu.Lock()
	_struct.KnownGroupConfig.KnownGroup = append(_struct.KnownGroupConfig.KnownGroup, struct {
		ChatroomId   string `yaml:"chatroom_id"`
		ChatroomName string `yaml:"chatroom_name"`
		BotWxId      string `yaml:"bot_wx_id"`
	}(struct {
		ChatroomId   string
		ChatroomName string
		BotWxId      string
	}{ChatroomId: chatroomId, ChatroomName: ChatroomName, BotWxId: BotWxId}))
	fmt.Println(_struct.KnownGroupConfig)
	// 将结构体编码为 YAML 数据
	data, err := yaml.Marshal(_struct.KnownGroupConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// 写入到文件
	err = ioutil.WriteFile("./config/KnownGroup.yml", data, 0644)
	if err != nil {
	}
	// 释放锁
	mu.Unlock()
}
