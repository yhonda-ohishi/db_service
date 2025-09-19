package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yhonda-ohishi/db_service/src/config"
	"github.com/yhonda-ohishi/db_service/src/models"
)

func main() {
	fmt.Println("ハッシュ生成とマッピングテーブルテスト開始...")

	// 設定読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("設定読み込みエラー: %v", err)
	}

	// データベース接続
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	defer func() {
		if err := config.CloseDatabase(db); err != nil {
			log.Printf("データベースクローズエラー: %v", err)
		}
	}()

	// マイグレーション実行（新しいテーブル作成）
	err = db.AutoMigrate(
		&models.ETCMeisai{},
		&models.ETCMeisaiMapping{},
	)
	if err != nil {
		log.Fatalf("マイグレーションエラー: %v", err)
	}
	fmt.Println("マイグレーション完了")

	// テストデータ作成とハッシュ生成
	testETCMeisai := &models.ETCMeisai{
		DateTo:     time.Date(2025, 9, 19, 10, 0, 0, 0, time.UTC),
		DateToDate: time.Date(2025, 9, 19, 0, 0, 0, 0, time.UTC),
		IcFr:       "東京IC",
		IcTo:       "横浜IC",
		Price:      1500,
		Shashu:     1,
		EtcNum:     "1234567890123456",
	}

	// ハッシュ生成
	testETCMeisai.SetHash()
	fmt.Printf("生成されたハッシュ: %s\n", testETCMeisai.Hash)

	// バリデーション
	if err := testETCMeisai.Validate(); err != nil {
		log.Fatalf("バリデーションエラー: %v", err)
	}
	fmt.Println("バリデーション成功")

	// データベースに保存
	result := db.Create(testETCMeisai)
	if result.Error != nil {
		log.Fatalf("ETC明細保存エラー: %v", result.Error)
	}
	fmt.Printf("ETC明細保存成功 (ID: %d)\n", testETCMeisai.ID)

	// マッピングテーブルテスト
	mapping := &models.ETCMeisaiMapping{
		ETCMeisaiHash: testETCMeisai.Hash,
		DTakoRowID:    "ROW123456789012345678901",
		CreatedBy:     "test_user",
		Notes:         &[]string{"テストマッピング"}[0],
	}

	// マッピングのバリデーション
	mapping.BeforeCreate()
	if err := mapping.Validate(); err != nil {
		log.Fatalf("マッピングバリデーションエラー: %v", err)
	}
	fmt.Println("マッピングバリデーション成功")

	// マッピング保存
	result = db.Create(mapping)
	if result.Error != nil {
		log.Fatalf("マッピング保存エラー: %v", result.Error)
	}
	fmt.Printf("マッピング保存成功 (ID: %d)\n", mapping.ID)

	// ハッシュでマッピング検索テスト
	var mappings []models.ETCMeisaiMapping
	result = db.Where("etc_meisai_hash = ?", testETCMeisai.Hash).Find(&mappings)
	if result.Error != nil {
		log.Fatalf("マッピング検索エラー: %v", result.Error)
	}

	fmt.Printf("ハッシュ %s に対するマッピング数: %d\n", testETCMeisai.Hash, len(mappings))
	for i, m := range mappings {
		fmt.Printf("  %d. マッピングID: %d, DTakoRowID: %s, 作成者: %s\n",
			i+1, m.ID, m.DTakoRowID, m.CreatedBy)
	}

	// 重複ハッシュテスト
	duplicateETCMeisai := &models.ETCMeisai{
		DateTo:     time.Date(2025, 9, 19, 10, 0, 0, 0, time.UTC),
		DateToDate: time.Date(2025, 9, 19, 0, 0, 0, 0, time.UTC),
		IcFr:       "東京IC",
		IcTo:       "横浜IC",
		Price:      1500,
		Shashu:     1,
		EtcNum:     "1234567890123456",
	}
	duplicateETCMeisai.SetHash()

	if duplicateETCMeisai.Hash == testETCMeisai.Hash {
		fmt.Println("✓ 同じデータから同じハッシュが生成されました")
	} else {
		fmt.Println("✗ ハッシュの一貫性エラー")
	}

	// 異なるデータでのハッシュテスト
	differentETCMeisai := &models.ETCMeisai{
		DateTo:     time.Date(2025, 9, 19, 10, 0, 0, 0, time.UTC),
		DateToDate: time.Date(2025, 9, 19, 0, 0, 0, 0, time.UTC),
		IcFr:       "東京IC",
		IcTo:       "大阪IC", // 異なるIC
		Price:      1500,
		Shashu:     1,
		EtcNum:     "1234567890123456",
	}
	differentETCMeisai.SetHash()

	if differentETCMeisai.Hash != testETCMeisai.Hash {
		fmt.Println("✓ 異なるデータから異なるハッシュが生成されました")
	} else {
		fmt.Println("✗ ハッシュの一意性エラー")
	}

	fmt.Println("ハッシュ生成とマッピングテーブルテスト完了!")
}