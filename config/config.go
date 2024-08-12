package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type ConfigInfo struct {
	Name         string `yaml:"name"`
	Version      string `yaml:"version"`
	SerciveHose  string `yaml:"sercive_hose"`
	SecurityCode string `yaml:"security_code"`
	BotWxid      string `yaml:"BotWxid"`
	Database     struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`
	Chatroom       []string `yaml:"chatroom"`
	RevokeChatroom []string `yaml:"revokeChatroom"`
}

func GetConfigInfo() ConfigInfo {
	// 读取文件
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 创建 Config 实例
	var config ConfigInfo

	// 解析 YAML 数据
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return config
}
