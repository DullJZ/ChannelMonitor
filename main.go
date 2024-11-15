package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Channel struct {
	ID      int
	Type    int
	Name    string
	BaseURL string
	Key     string
	Status  int
}

var (
	db     *sql.DB
	config *Config
)

func initDB(dbType, dsn string) error {
	var err error
	
	switch dbType {
	case "mysql":
		db, err = sql.Open("mysql", dsn)
	case "sqlite":
		db, err = sql.Open("sqlite3", dsn)
		if err != nil {
			return err
		}
		db.SetMaxOpenConns(1)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", dbType)
	}
	
	if err != nil {
		return err
	}
	
	return db.Ping()
}

func fetchChannels() ([]Channel, error) {
	query := "SELECT id, type, name, base_url, `key`, status FROM channels"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []Channel
	for rows.Next() {
		var c Channel
		if err := rows.Scan(&c.ID, &c.Type, &c.Name, &c.BaseURL, &c.Key, &c.Status); err != nil {
			return nil, err
		}
		
		switch c.Type {
		case 40:
			c.BaseURL = "https://api.siliconflow.cn"
		case 999:
			c.BaseURL = "https://api.siliconflow.cn"
		}
		
		if contains(config.ExcludeChannel, c.ID) {
			log.Printf("渠道 %s(ID:%d) 在排除列表中，跳过\n", c.Name, c.ID)
			continue
		}
		channels = append(channels, c)
	}
	log.Printf("获取到 %d 个渠道\n", len(channels))
	log.Printf("准备测试的渠道：%v\n", channels)
	return channels, nil
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func testModels(channel Channel) ([]string, error) {
	var availableModels []string
	// 从/v1/models接口获取模型列表
	req, err := http.NewRequest("GET", channel.BaseURL+"/v1/models", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败：%v", err)
	}
	req.Header.Set("Authorization", "Bearer "+channel.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	modelList := []string{}
	if err != nil {
		log.Println("获取模型列表失败：", err, "尝试自定义模型列表")
		modelList = config.Models
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("获取模型列表失败，状态码：%d，响应：%s", resp.StatusCode, string(body))
		}

		// 解析响应JSON
		var response struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("解析模型列表失败：%v", err)
		}
		// 提取模型ID列表
		for _, model := range response.Data {
			modelList = append(modelList, model.ID)
		}
	}
	// 测试模型
	for _, model := range modelList {
		url := channel.BaseURL
		if !strings.Contains(channel.BaseURL, "/v1/chat/completions") {
			if !strings.HasSuffix(channel.BaseURL, "/chat") {
				if !strings.HasSuffix(channel.BaseURL, "/v1") {
					url += "/v1"
				}
				url += "/chat"
			}
			url += "/completions"
		}

		// 构造请求
		reqBody := map[string]interface{}{
			"model": model,
			"messages": []map[string]string{
				{"role": "user", "content": "Hello! Reply in short"},
			},
		}
		jsonData, _ := json.Marshal(reqBody)

		log.Printf("测试渠道 %s(ID:%d) 的模型 %s\n", channel.Name, channel.ID, model)

		req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
		if err != nil {
			log.Println("创建请求失败：", err)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+channel.Key)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("\033[31m请求失败：%v\033[0m\n", err)
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusOK {
			// 根据返回内容判断是否成功
			availableModels = append(availableModels, model)
			log.Printf("\033[32m渠道 %s(ID:%d) 的模型 %s 测试成功\033[0m\n", channel.Name, channel.ID, model)
		} else {
			log.Printf("\033[31m渠道 %s(ID:%d) 的模型 %s 测试失败，状态码：%d，响应：%s\033[0m\n", channel.Name, channel.ID, model, resp.StatusCode, string(body))
		}
	}
	return availableModels, nil
}

func updateModels(channelID int, models []string) error {
	modelsStr := strings.Join(models, ",")
	query := "UPDATE channels SET models = ? WHERE id = ?"
	_, err := db.Exec(query, modelsStr, channelID)
	return err
}

func main() {
	var err error
	config, err = loadConfig()
	if err != nil {
		log.Fatal("加载配置失败：", err)
	}

	// 解析时间周期
	duration, err := time.ParseDuration(config.TimePeriod)
	if err != nil {
		log.Fatal("解析时间周期失败：", err)
	}

	err = initDB(config.DbType, config.DbDsn)
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}
	defer db.Close()

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		log.Println("开始检测...")
		channels, err := fetchChannels()
		if err != nil {
			log.Printf("\033[31m获取渠道失败：%v\033[0m\n", err)
			continue
		}

		for _, channel := range channels {
			log.Printf("开始测试渠道 %s(ID:%d) 的模型\n", channel.Name, channel.ID)
			models, err := testModels(channel)
			if err != nil {
				log.Printf("\033[31m渠道 %s(ID:%d) 测试模型失败：%v\033[0m\n", channel.Name, channel.ID, err)
				continue
			}
			err = updateModels(channel.ID, models)
			if err != nil {
				log.Printf("\033[31m更新渠道 %s(ID:%d) 的模型失败：%v\033[0m\n", channel.Name, channel.ID, err)
			} else {
				log.Printf("渠道 %s(ID:%d) 可用模型：%v\n", channel.Name, channel.ID, models)
			}
		}

		// 等待下一个周期
		<-ticker.C
	}
}
