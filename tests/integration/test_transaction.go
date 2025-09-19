//go:build integration
// +build integration

package integration

import (
	"testing"
	"time"

	"github.com/yhonda-ohishi/db_service/src/config"
	"github.com/yhonda-ohishi/db_service/src/models"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"gorm.io/gorm"
)

func TestTransactionRollback(t *testing.T) {
	// データベース接続のセットアップ
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	db, err := config.InitDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to init database: %v", err)
	}

	repo := repository.NewDTakoUriageKeihiRepository(db)

	// トランザクションテスト
	t.Run("Transaction rollback on error", func(t *testing.T) {
		// トランザクション開始
		tx := db.Begin()
		txRepo := repository.NewDTakoUriageKeihiRepository(tx)

		// 有効なデータを挿入
		validData := &models.DTakoUriageKeihi{
			SrchID:      "TRANS_TEST001",
			Datetime:    time.Now(),
			KeihiC:      1,
			Price:       1000.0,
			DtakoRowID:  "DTAKO_TRANS001",
			DtakoRowIDR: "DTAKO_TRANS001R",
		}

		err := txRepo.Create(validData)
		if err != nil {
			t.Errorf("Failed to create valid data: %v", err)
		}

		// 無効なデータで意図的にエラーを発生させる
		invalidData := &models.DTakoUriageKeihi{
			SrchID:      validData.SrchID, // 重複キー
			Datetime:    validData.Datetime,
			KeihiC:      validData.KeihiC,
			Price:       2000.0,
			DtakoRowID:  "DTAKO_TRANS002",
			DtakoRowIDR: "DTAKO_TRANS002R",
		}

		err = txRepo.Create(invalidData)
		if err == nil {
			tx.Commit()
			t.Error("Expected error for duplicate key, but got nil")
		} else {
			// エラー発生時はロールバック
			tx.Rollback()
		}

		// ロールバック後、データが存在しないことを確認
		result, _ := repo.GetByCompositeKey(validData.SrchID, validData.Datetime, validData.KeihiC)
		if result != nil {
			t.Error("Data should not exist after rollback")
		}
	})

	// 複数テーブルのトランザクション
	t.Run("Multi-table transaction", func(t *testing.T) {
		tx := db.Begin()

		// DTakoUriageKeihiの作成
		keihiRepo := repository.NewDTakoUriageKeihiRepository(tx)
		keihiData := &models.DTakoUriageKeihi{
			SrchID:      "MULTI_TRANS001",
			Datetime:    time.Now(),
			KeihiC:      1,
			Price:       3000.0,
			DtakoRowID:  "DTAKO_MULTI001",
			DtakoRowIDR: "DTAKO_MULTI001R",
		}

		err := keihiRepo.Create(keihiData)
		if err != nil {
			tx.Rollback()
			t.Errorf("Failed to create keihi data: %v", err)
			return
		}

		// ETCMeisaiの作成
		etcRepo := repository.NewETCMeisaiRepository(tx)
		etcData := &models.ETCMeisai{
			DateTo:     time.Now(),
			DateToDate: time.Now().Format("2006-01-02"),
			IcFr:       "東京IC",
			IcTo:       "横浜IC",
			Price:      1500,
			Shashu:     1,
			EtcNum:     "1234567890123456",
			DtakoRowID: &keihiData.DtakoRowID,
		}

		err = etcRepo.Create(etcData)
		if err != nil {
			tx.Rollback()
			t.Errorf("Failed to create etc data: %v", err)
			return
		}

		// コミット
		if err := tx.Commit().Error; err != nil {
			t.Errorf("Failed to commit transaction: %v", err)
		}

		// データの存在確認
		result, _ := repo.GetByCompositeKey(keihiData.SrchID, keihiData.Datetime, keihiData.KeihiC)
		if result == nil {
			t.Error("Data should exist after successful commit")
		}
	})
}

// トランザクションの分離レベルテスト
func TestTransactionIsolation(t *testing.T) {
	cfg, _ := config.LoadConfig()
	db, _ := config.InitDatabase(cfg)

	t.Run("Read committed isolation", func(t *testing.T) {
		// トランザクション1を開始
		tx1 := db.Begin(&gorm.Session{PrepareStmt: false})
		repo1 := repository.NewDTakoUriageKeihiRepository(tx1)

		// トランザクション2を開始
		tx2 := db.Begin(&gorm.Session{PrepareStmt: false})
		repo2 := repository.NewDTakoUriageKeihiRepository(tx2)

		// トランザクション1でデータ作成
		testData := &models.DTakoUriageKeihi{
			SrchID:      "ISO_TEST001",
			Datetime:    time.Now(),
			KeihiC:      1,
			Price:       5000.0,
			DtakoRowID:  "DTAKO_ISO001",
			DtakoRowIDR: "DTAKO_ISO001R",
		}

		err := repo1.Create(testData)
		if err != nil {
			t.Errorf("Failed to create in tx1: %v", err)
		}

		// トランザクション2から見えないことを確認（未コミット）
		result, _ := repo2.GetByCompositeKey(testData.SrchID, testData.Datetime, testData.KeihiC)
		if result != nil {
			t.Error("Uncommitted data should not be visible in tx2")
		}

		// トランザクション1をコミット
		tx1.Commit()

		// コミット後はトランザクション2からも見える
		result, _ = repo2.GetByCompositeKey(testData.SrchID, testData.Datetime, testData.KeihiC)
		if result == nil {
			t.Error("Committed data should be visible in tx2")
		}

		tx2.Rollback()
	})
}
