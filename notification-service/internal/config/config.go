package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type StorageConfig struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	UserDb   string `yaml:"userdb"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	SSLmode  string `yaml:"ssl_mode" env-default:"require"`
}
type SMTPConfig struct {
	Port         string `yaml:"smtp_port"`
	Host         string `yaml:"smtp_host"`
	MailPassword string `yaml:"smtp_password"`
	Mail_from    string `yaml:"smtp_from"`
}
type GRPCConfig struct {
	Host    string        `yaml:"grpc_host" env-default:"localhost"`
	Port    string        `yaml:"grpc_port"`
	Timeout time.Duration `yaml:"timeout"`
}
type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	GRPC     GRPCConfig    `yaml:"grpc_server" env-required:"true"`
	SMTP     SMTPConfig    `yaml:"smtp_server" env-required:"true"`
	Storage  StorageConfig `yaml:"storage"`
	BotToken string        `yaml:"bot_token"`
}

func MustLoad() *Config {
	//err := godotenv.Load(".env")
	//if err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Error empty config path")
	}

	//configPath = "./config/dev.yaml"

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Panicf("CONFIG_PATH file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	cfg.SMTP.MailPassword = os.Getenv("MAIL_PASSWORD")
	cfg.SMTP.Mail_from = os.Getenv("MAIL_FROM")

	cfg.BotToken = os.Getenv("BOT_TOKEN")

	cfg.Storage.UserDb = os.Getenv("POSTGRES_USER")
	cfg.Storage.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Storage.Dbname = os.Getenv("POSTGRES_DB")

	return &cfg
}
