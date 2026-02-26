package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	OneAPIType        string   `json:"oneapi_type" yaml:"oneapi_type"`
	ExcludeChannel    []int    `json:"exclude_channel" yaml:"exclude_channel"`
	ExcludeModel      []string `json:"exclude_model" yaml:"exclude_model"`
	Models            []string `json:"models" yaml:"models"`
	ForceModels       bool     `json:"force_models" yaml:"force_models"`
	ForceInsideModels bool     `json:"force_inside_models" yaml:"force_inside_models"`
	TimePeriod        string   `json:"time_period" yaml:"time_period"`
	MaxConcurrent     int      `json:"max_concurrent" yaml:"max_concurrent"`
	RPS               int      `json:"rps" yaml:"rps"`
	RPM               int      `json:"rpm" yaml:"rpm"`
	Timeout           int      `json:"timeout" yaml:"timeout"`
	DbType            string   `json:"db_type" yaml:"db_type"`
	DbDsn             string   `json:"db_dsn" yaml:"db_dsn"`
	DoNotModifyDb     bool     `json:"do_not_modify_db" yaml:"do_not_modify_db"`
	BaseURL           string   `json:"base_url" yaml:"base_url"`
	SystemToken       string   `json:"system_token" yaml:"system_token"`
	UptimeKuma        struct {
		Status     string            `json:"status" yaml:"status"`
		ModelURL   map[string]string `json:"model_url" yaml:"model_url"`
		ChannelURL map[string]string `json:"channel_url" yaml:"channel_url"`
	} `json:"uptime-kuma" yaml:"uptime-kuma"`
	Notification struct {
		SMTP struct {
			Enabled  bool   `json:"enabled" yaml:"enabled"`
			Host     string `json:"host" yaml:"host"`
			Port     int    `json:"port" yaml:"port"`
			Username string `json:"username" yaml:"username"`
			Password string `json:"password" yaml:"password"`
			From     string `json:"from" yaml:"from"`
			To       string `json:"to" yaml:"to"`
		} `json:"smtp" yaml:"smtp"`
		Webhook struct {
			Enabled  bool   `json:"enabled" yaml:"enabled"`
			Type     string `json:"type" yaml:"type"`
			Telegram struct {
				ChatID string `json:"chat_id" yaml:"chat_id"`
				Retry  int    `json:"retry" yaml:"retry"`
			} `json:"telegram" yaml:"telegram"`
			Secret string `json:"secret" yaml:"secret"`
		} `json:"webhook" yaml:"webhook"`
	} `json:"notification" yaml:"notification"`
}

func loadConfig() (*Config, error) {
	// 尝试加载不同格式的配置文件
	possibleConfigs := []string{"config.yaml", "config.yml", "config.json"}

	var configFile string
	var configData []byte

	for _, filename := range possibleConfigs {
		if _, err := os.Stat(filename); err == nil {
			configFile = filename
			configData, err = ioutil.ReadFile(filename)
			if err != nil {
				return nil, fmt.Errorf("读取配置文件 %s 失败: %v", filename, err)
			}
			break
		}
	}

	if configFile == "" {
		return nil, fmt.Errorf("找不到配置文件，请创建 config.yaml、config.yml 或 config.json")
	}

	var config Config
	ext := strings.ToLower(filepath.Ext(configFile))

	if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(configData, &config); err != nil {
			return nil, fmt.Errorf("解析YAML配置文件失败: %v", err)
		}
		fmt.Printf("使用YAML格式配置文件: %s\n", configFile)
	} else {
		if err := json.Unmarshal(configData, &config); err != nil {
			return nil, fmt.Errorf("解析JSON配置文件失败: %v", err)
		}
		fmt.Printf("使用JSON格式配置文件: %s\n", configFile)
	}

	// 设置默认值
	if config.OneAPIType == "" {
		config.OneAPIType = "oneapi"
	}

	if config.DbType == "" {
		config.DbType = "mysql"
	}

	if config.BaseURL != "" && config.BaseURL[len(config.BaseURL)-1] == '/' {
		config.BaseURL = config.BaseURL[:len(config.BaseURL)-1]
	}

	if config.UptimeKuma.Status == "" {
		config.UptimeKuma.Status = "disabled"
	}

	if config.MaxConcurrent == 0 {
		config.MaxConcurrent = 5
	}

	// 仅当 rps 和 rpm 都未配置时，才给 rps 设默认值，以便用户可以通过只配置 rpm 来禁用 rps
	if config.RPS == 0 && config.RPM == 0 {
		config.RPS = 5
	}

	if config.Timeout == 0 {
		config.Timeout = 10
	}

	return &config, nil
}
