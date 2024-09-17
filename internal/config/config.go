package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Storage struct {
	Host     string `yaml:"host" env-default:"localhost"`
	UserDb   string `yaml:"userdb"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port" env-default:"5432"`
	Dbname   string `yaml:"dbname"`
	SSLmode  string `yaml:"sslmode"`
}

type HTTPServer struct {
	Address  string `yaml:"address" env-default:"localhost:8080"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env:"HTTP-SERVER-PASSWORD"`
}
type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server" env-required:"true"`
	Storage    `yaml:"storage" env-required:"true"`
}

func MustLoad() *Config {
	err := godotenv.Load("config/config.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Error empty config path")
	}
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	return &cfg
}
