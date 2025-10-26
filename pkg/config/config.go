package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Config 应用配置结构
type Config struct {
	Rest struct {
		Port int    `json:"port"`
		Base string `json:"base"`
	} `json:"rest"`
	CropId       string `json:"cropId"`
	CropSecret   string `json:"cropSecret"`
	AgentId      string `json:"agentId"`
	Receiver     string `json:"receiver"`
	Token        string `json:"token"`
}

var globalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig() error {
	// 支持通过环境变量 CONFIG_PATH 指定配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.json"
	}

	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("打开配置文件失败: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	globalConfig = config
	return nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if globalConfig == nil {
		panic("配置未初始化，请先调用 LoadConfig()")
	}
	return globalConfig
}

// GetString 获取字符串配置值
func GetString(key string) string {
	cfg := GetConfig()
	switch key {
	case "cropId":
		return cfg.CropId
	case "cropSecret":
		return cfg.CropSecret
	case "agentId":
		return cfg.AgentId
	case "receiver":
		return cfg.Receiver
	case "token":
		return cfg.Token
	default:
		return ""
	}
}

// SetString 覆盖字符串配置值
func SetString(key, value string) {
	if globalConfig == nil {
		globalConfig = &Config{}
	}
	switch key {
	case "cropId":
		globalConfig.CropId = value
	case "cropSecret":
		globalConfig.CropSecret = value
	case "agentId":
		globalConfig.AgentId = value
	case "receiver":
		globalConfig.Receiver = value
	case "token":
		globalConfig.Token = value
	}
}