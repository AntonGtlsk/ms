package config

import (
	config_entity "chat-ms/internal/entity/config"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustLoad() *config_entity.Config {
	configPath := "config/local.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg config_entity.Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
