package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config chứa toàn bộ các biến môi trường của dự án
type Config struct {
	Environment string `mapstructure:"ENV"`
	APIPort     string `mapstructure:"API_PORT"`
	EdgePort    string `mapstructure:"EDGE_PORT"`

	// Database (Postgres)
	PostgresUser string `mapstructure:"POSTGRES_USER"`
	PostgresPass string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDB   string `mapstructure:"POSTGRES_DB"`
	PostgresHost string `mapstructure:"POSTGRES_HOST"`

	// Infra & Message Brokers
	KafkaBrokers string `mapstructure:"KAFKA_BROKERS"` // VD: "localhost:9092"
	RedisURL     string `mapstructure:"REDIS_URL"`
	NatsURL      string `mapstructure:"NATS_URL"`
	DynamoDBURL  string `mapstructure:"DYNAMODB_URL"`
}

// LoadConfig đọc cấu hình từ file .env hoặc từ OS Environment Variables
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Quan trọng: Tự động ghi đè giá trị trong .env nếu hệ điều hành có truyền biến cùng tên
	viper.AutomaticEnv()

	// Thiết lập các giá trị mặc định (Fallback) để code vẫn chạy mượt dưới local
	viper.SetDefault("ENV", "development")
	viper.SetDefault("API_PORT", "8081")
	viper.SetDefault("EDGE_PORT", "8080")
	viper.SetDefault("POSTGRES_HOST", "localhost")
	viper.SetDefault("KAFKA_BROKERS", "localhost:9092")
	viper.SetDefault("REDIS_URL", "localhost:6379")
	viper.SetDefault("NATS_URL", "nats://localhost:4222")
	viper.SetDefault("DYNAMODB_URL", "http://localhost:8000")

	// Đọc file .env
	if err := viper.ReadInConfig(); err != nil {
		// Trên production chúng ta thường không dùng file .env mà xài thẳng env của Docker/OS
		// Nên nếu không thấy file .env thì ta chỉ in ra Warning chứ không crash app.
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning: Không thể đọc file .env: %v", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
