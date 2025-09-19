package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config アプリケーション設定
type Config struct {
	// データベース設定
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// gRPCサーバー設定
	GRPCPort int

	// 接続プール設定
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int
}

// LoadConfig 環境変数から設定を読み込み
func LoadConfig() (*Config, error) {
	// .envファイルの読み込み（存在する場合）
	_ = godotenv.Load()

	config := &Config{}

	// データベース設定の読み込み
	config.DBHost = getEnv("DB_HOST", "localhost")
	config.DBPort = getEnvAsInt("DB_PORT", 3306)
	config.DBUser = getEnv("DB_USER", "")
	config.DBPassword = getEnv("DB_PASSWORD", "")
	config.DBName = getEnv("DB_NAME", "db1")

	// 必須フィールドのチェック
	if config.DBUser == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if config.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	// gRPC設定
	config.GRPCPort = getEnvAsInt("GRPC_PORT", 50051)

	// 接続プール設定
	config.MaxOpenConns = getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	config.MaxIdleConns = getEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	config.ConnMaxLifetime = getEnvAsInt("DB_CONN_MAX_LIFETIME", 3600) // 秒単位
	config.ConnMaxIdleTime = getEnvAsInt("DB_CONN_MAX_IDLE_TIME", 300)  // 秒単位

	return config, nil
}

// GetDSN データソース名を生成
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}

// GetGRPCAddress gRPCサーバーのアドレスを取得
func (c *Config) GetGRPCAddress() string {
	return fmt.Sprintf(":%d", c.GRPCPort)
}

// getEnv 環境変数を取得（デフォルト値付き）
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 環境変数を整数として取得
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// Validate 設定の妥当性を検証
func (c *Config) Validate() error {
	if c.DBHost == "" {
		return fmt.Errorf("DBHost cannot be empty")
	}
	if c.DBPort <= 0 || c.DBPort > 65535 {
		return fmt.Errorf("invalid DBPort: %d", c.DBPort)
	}
	if c.DBName == "" {
		return fmt.Errorf("DBName cannot be empty")
	}
	if c.GRPCPort <= 0 || c.GRPCPort > 65535 {
		return fmt.Errorf("invalid GRPCPort: %d", c.GRPCPort)
	}
	if c.MaxOpenConns <= 0 {
		return fmt.Errorf("MaxOpenConns must be positive")
	}
	if c.MaxIdleConns <= 0 || c.MaxIdleConns > c.MaxOpenConns {
		return fmt.Errorf("invalid MaxIdleConns: %d", c.MaxIdleConns)
	}
	return nil
}