package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// SQLServerDatabase SQL Serverデータベース接続
type SQLServerDatabase struct {
	DB *gorm.DB
}

// NewSQLServerDatabase SQL Serverデータベース接続を初期化
func NewSQLServerDatabase() (*SQLServerDatabase, error) {
	host := os.Getenv("SQLSERVER_HOST")
	user := os.Getenv("SQLSERVER_USER")
	password := os.Getenv("SQLSERVER_PASSWORD")
	database := os.Getenv("SQLSERVER_DATABASE")

	if host == "" || user == "" || password == "" || database == "" {
		return nil, fmt.Errorf("SQL Server connection info not provided in environment variables")
	}

	// SQL Server接続文字列
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", user, password, host, database)

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
