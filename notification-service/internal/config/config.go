package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
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
	SSLmode  string `yaml:"sslmode" env-default:"require"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}
type Config struct {
	Env     string     `yaml:"env" env-default:"local"`
	GRPC    GRPCConfig `yaml:"grpc_server" env-required:"true"`
	Storage `yaml:"storage"`
}

func MustLoad() *Config {
	err := godotenv.Load("notification-service/config/config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Error empty config path")
	}

	//configPath := "./config/local.yaml"

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Panicf("CONFIG_PATH file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	cfg.UserDb = os.Getenv("POSTGRES_USER")
	cfg.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Dbname = os.Getenv("POSTGRES_DB")

	return &cfg
}
