//go:build integration
// +build integration

package integration

import (
	"testing"
	"time"

	"github.com/yhonda-ohishi/db_service/src/config"
	"github.com/yhonda-ohishi/db_service/src/models"
	"github.com/yhonda-ohishi/db_service/src/repository"
)

func TestCompositeKeyOperations(t *testing.T) {
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

	// テストデータ
	testData := &models.DTakoUriageKeihi{
		SrchID:      "COMPOSITE_TEST001",
		Datetime:    time.Now(),
		KeihiC:      1,
		Price:       1000.0,
		Km:          50.5,
		DtakoRowID:  "DTAKO_COMP001",
		DtakoRowIDR: "DTAKO_COMP001R",
	}

	// 複合主キーでの作成
	t.Run("Create with composite key", func(t *testing.T) {
		err := repo.Create(testData)
		if err != nil {
			t.Errorf("Failed to create record: %v", err)
		}
	})

	// 複合主キーでの取得
	t.Run("Get by composite key", func(t *testing.T) {
		result, err := repo.GetByCompositeKey(testData.SrchID, testData.Datetime, testData.KeihiC)
		if err != nil {
			t.Errorf("Failed to get record: %v", err)
		}
		if result == nil {
			t.Error("Expected non-nil result")
		}
		if result != nil && result.Price != testData.Price {
			t.Errorf("Price mismatch: got %v, want %v", result.Price, testData.Price)
		}
	})

	// 複合主キーの重複エラーテスト
	t.Run("Duplicate composite key error", func(t *testing.T) {
		duplicate := &models.DTakoUriageKeihi{
			SrchID:      testData.SrchID,
			Datetime:    testData.Datetime,
			KeihiC:      testData.KeihiC,
			Price:       2000.0, // 異なる値
			DtakoRowID:  "DTAKO_COMP002",
			DtakoRowIDR: "DTAKO_COMP002R",
		}
		err := repo.Create(duplicate)
		if err == nil {
			t.Error("Expected duplicate key error, but got nil")
		}
	})

	// 部分キーでの検索
	t.Run("List by partial key", func(t *testing.T) {
		results, err := repo.ListBySrchID(testData.SrchID)
		if err != nil {
			t.Errorf("Failed to list by srch_id: %v", err)
		}
		if len(results) == 0 {
			t.Error("Expected at least one result")
		}
	})

	// 複合主キーでの削除
	t.Run("Delete by composite key", func(t *testing.T) {
		err := repo.DeleteByCompositeKey(testData.SrchID, testData.Datetime, testData.KeihiC)
		if err != nil {
			t.Errorf("Failed to delete record: %v", err)
		}

		// 削除確認
		result, err := repo.GetByCompositeKey(testData.SrchID, testData.Datetime, testData.KeihiC)
		if err == nil && result != nil {
			t.Error("Record should have been deleted")
		}
	})
}
