package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Storage struct {
	Host     string `yaml:"host" env-default:"localhost"`
	UserDb   string `yaml:"userdb"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port" env-default:"5432"`
	Dbname   string `yaml:"dbname"`
	SSLmode  string `yaml:"sslmode"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}
type Config struct {
	Env     string     `yaml:"env" env-default:"local"`
	GRPC    GRPCConfig `yaml:"grpc" env-required:"true"`
	Storage `yaml:"storage" env-required:"true"`
}

func MustLoad() *Config {
	//err := godotenv.Load("config/config.env")
	//if err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//}
	//
	//configPath := os.Getenv("CONFIG_PATH")
	//if configPath == "" {
	//	log.Fatal("Error empty config path")
	//}

	configPath := "./config/local.yaml"

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Panicf("CONFIG_PATH file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	return &cfg
}
