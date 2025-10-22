package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// SQLServerDatabase SQL Serverデータベース接続
type SQLServerDatabase struct {
	DB *gorm.DB
}

// NewSQLServerDatabase SQL Serverデータベース接続を初期化
func NewSQLServerDatabase() (*SQLServerDatabase, error) {
	// .envファイルの読み込み
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	host := os.Getenv("SQLSERVER_HOST")
	instance := os.Getenv("SQLSERVER_INSTANCE")
	user := os.Getenv("SQLSERVER_USER")
	password := os.Getenv("SQLSERVER_PASSWORD")
	database := os.Getenv("SQLSERVER_DATABASE")

	if host == "" || user == "" || password == "" || database == "" {
		return nil, fmt.Errorf("SQL Server connection info not provided in environment variables")
	}

	// SQL Server接続文字列の構築
	// インスタンス名がある場合は "server\instance" 形式にする
	serverAddress := host
	if instance != "" {
		serverAddress = fmt.Sprintf("%s\\%s", host, instance)
	}

	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s",
		user,
		password,
		serverAddress,
		database)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQL Server database: %v", err)
	}

	log.Printf("SQL Server database connected: %s", database)

	return &SQLServerDatabase{DB: db}, nil
}

// HealthCheck SQL Serverデータベース接続確認
func SQLServerHealthCheck(db *SQLServerDatabase) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
