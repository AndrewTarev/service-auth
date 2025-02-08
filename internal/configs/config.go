package configs

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Конфигурация сервера
type ServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes int           `mapstructure:"max_header_bytes"`
}

// Конфигурация логирования
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputFile string `mapstructure:"output_file"`
}

// Конфигурация базы данных
type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// Конфигурация Auth
type AuthConfig struct {
	SigningKey      string        `mapstructure:"signing_key"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// Полная конфигурация
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Logging  LoggerConfig   `mapstructure:"logging"`
	Database PostgresConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

// LoadConfig загружает конфигурацию из файлов и переменных окружения
func LoadConfig(path string) (*Config, error) {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Инициализация Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Чтение конфигурации
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not load YAML config file: %v", err)
	}

	// Маппинг данных в структуру Config
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Валидация конфигурации
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return nil, fmt.Errorf("invalid server port: %d", config.Server.Port)
	}
	if config.Server.ReadTimeout <= 0 {
		config.Server.ReadTimeout = 5 * time.Second
	}
	if config.Server.WriteTimeout <= 0 {
		config.Server.WriteTimeout = 10 * time.Second
	}

	return &config, nil
}
