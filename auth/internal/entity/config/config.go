package config_entity

import (
	"fmt"
	"time"
)

type Config struct {
	Env           string `yaml:"env" env-default:"local"`
	HTTPServer    `yaml:"http_server"`
	SQLConnection `yaml:"sql_connection"`
	Logger        Logger `yaml:"logger"`
	JWTKey        string `yaml:"jwt_key"`
	BotToken      string `yaml:"bot_token"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type SQLConnection struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"dbName"`
}

func (c SQLConnection) StoragePath() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Host, c.Port, c.DbName)
}

type Logger struct {
	Folder           string                `yaml:"folder"`
	Filenames        []map[string][]string `yaml:"filenames"`
	WebhookUrls      []map[string][]string `yaml:"webhooks"`
	UndefinedWebhook string                `yaml:"undefinedWebhook"`
}
