package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProdDatabase 本番データベース接続（読み取り専用）
type ProdDatabase struct {
	DB *gorm.DB
}

// NewProdDatabase 本番データベース接続の初期化
func NewProdDatabase() (*ProdDatabase, error) {
	host := os.Getenv("PROD_DB_HOST")
	port := os.Getenv("PROD_DB_PORT")
	user := os.Getenv("PROD_DB_USER")
	password := os.Getenv("PROD_DB_PASSWORD")
	dbname := os.Getenv("PROD_DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to production database: %v", err)
	}

	// 接続テスト
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping production database: %v", err)
	}

	// 読み取り専用の接続プール設定
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	log.Printf("Production database connection established: %s:%s/%s", host, port, dbname)

	return &ProdDatabase{DB: db}, nil
}

// Close 接続を閉じる
func (pdb *ProdDatabase) Close() error {
	sqlDB, err := pdb.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
