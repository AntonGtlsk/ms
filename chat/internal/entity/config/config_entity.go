package config_entity

import "fmt"

type Config struct {
	Env           string    `yaml:"env"`
	WSSServer     WSSServer `yaml:"wss_server"`
	BotToken      string    `yaml:"bot_token"`
	SQLConnection `yaml:"sql_connection"`
	Logger        Logger `yaml:"logger"`
}

type SQLConnection struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"dbName"`
}

type Logger struct {
	Folder           string                `yaml:"folder"`
	Filenames        []map[string][]string `yaml:"filenames"`
	WebhookUrls      []map[string][]string `yaml:"webhooks"`
	UndefinedWebhook string                `yaml:"undefinedWebhook"`
}

func (c SQLConnection) StoragePath() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port, c.DbName)
}

type WSSServer struct {
	Address string `yaml:"address"`
}
