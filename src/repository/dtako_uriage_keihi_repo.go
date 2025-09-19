package repository

import (
	"fmt"
	"time"

	"github.com/yhonda-ohishi/db_service/src/models"
	"gorm.io/gorm"
)

// DTakoUriageKeihiRepository リポジトリインターフェース
type DTakoUriageKeihiRepository interface {
	Create(data *models.DTakoUriageKeihi) error
	GetByCompositeKey(srchID string, datetime time.Time, keihiC int32) (*models.DTakoUriageKeihi, error)
	Update(data *models.DTakoUriageKeihi) error
	DeleteByCompositeKey(srchID string, datetime time.Time, keihiC int32) error
	List(params *ListParams) ([]*models.DTakoUriageKeihi, int64, error)
	ListBySrchID(srchID string) ([]*models.DTakoUriageKeihi, error)
	ListByDtakoRowID(dtakoRowID string) ([]*models.DTakoUriageKeihi, error)
	ListByDateRange(start, end time.Time) ([]*models.DTakoUriageKeihi, error)
}

// ListParams リスト取得用パラメータ
type ListParams struct {
	DtakoRowID *string
	StartDate  *time.Time
	EndDate    *time.Time
	Limit      int
	Offset     int
}

// dtakoUriageKeihiRepo リポジトリ実装
type dtakoUriageKeihiRepo struct {
	db *gorm.DB
}

// NewDTakoUriageKeihiRepository リポジトリのコンストラクタ
func NewDTakoUriageKeihiRepository(db *gorm.DB) DTakoUriageKeihiRepository {
	return &dtakoUriageKeihiRepo{db: db}
}

// Create データ作成
func (r *dtakoUriageKeihiRepo) Create(data *models.DTakoUriageKeihi) error {
	if err := data.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	result := r.db.Create(data)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return models.ErrDuplicateKey
		}
		return fmt.Errorf("failed to create record: %w", result.Error)
	}

	return nil
}

// GetByCompositeKey 複合キーでデータ取得
func (r *dtakoUriageKeihiRepo) GetByCompositeKey(srchID string, datetime time.Time, keihiC int32) (*models.DTakoUriageKeihi, error) {
	var data models.DTakoUriageKeihi

	result := r.db.Where("srch_id = ? AND datetime = ? AND keihi_c = ?",
		srchID, datetime, keihiC).First(&data)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, models.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get record: %w", result.Error)
	}

	return &data, nil
}

// Update データ更新
func (r *dtakoUriageKeihiRepo) Update(data *models.DTakoUriageKeihi) error {
	if err := data.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// 複合キーで既存レコードを確認
	existing, err := r.GetByCompositeKey(data.SrchID, data.Datetime, data.KeihiC)
	if err != nil {
		if err == models.ErrRecordNotFound {
			return models.ErrRecordNotFound
		}
		return err
	}

	// 更新実行
	result := r.db.Model(existing).Updates(data)
	if result.Error != nil {
		return fmt.Errorf("failed to update record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.ErrRecordNotFound
	}

	return nil
}

// DeleteByCompositeKey 複合キーでデータ削除
func (r *dtakoUriageKeihiRepo) DeleteByCompositeKey(srchID string, datetime time.Time, keihiC int32) error {
	result := r.db.Where("srch_id = ? AND datetime = ? AND keihi_c = ?",
		srchID, datetime, keihiC).Delete(&models.DTakoUriageKeihi{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.ErrRecordNotFound
	}

	return nil
}

// List 条件付きリスト取得
func (r *dtakoUriageKeihiRepo) List(params *ListParams) ([]*models.DTakoUriageKeihi, int64, error) {
	var data []*models.DTakoUriageKeihi
	var totalCount int64

	query := r.db.Model(&models.DTakoUriageKeihi{})

	// 条件の適用
	if params.DtakoRowID != nil && *params.DtakoRowID != "" {
		query = query.Where("dtako_row_id = ?", *params.DtakoRowID)
	}
	if params.StartDate != nil {
		query = query.Where("datetime >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("datetime <= ?", *params.EndDate)
	}

	// 総件数取得
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count records: %w", err)
	}

	// ページネーション
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	// データ取得
	if err := query.Order("datetime DESC").Find(&data).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list records: %w", err)
	}

	return data, totalCount, nil
}

// ListBySrchID srch_idでリスト取得
func (r *dtakoUriageKeihiRepo) ListBySrchID(srchID string) ([]*models.DTakoUriageKeihi, error) {
	var data []*models.DTakoUriageKeihi

	if err := r.db.Where("srch_id = ?", srchID).
		Order("datetime DESC").
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to list by srch_id: %w", err)
	}

	return data, nil
}

// ListByDtakoRowID dtako_row_idでリスト取得
func (r *dtakoUriageKeihiRepo) ListByDtakoRowID(dtakoRowID string) ([]*models.DTakoUriageKeihi, error) {
	var data []*models.DTakoUriageKeihi

	if err := r.db.Where("dtako_row_id = ?", dtakoRowID).
		Order("datetime DESC").
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to list by dtako_row_id: %w", err)
	}

	return data, nil
}

// ListByDateRange 日付範囲でリスト取得
func (r *dtakoUriageKeihiRepo) ListByDateRange(start, end time.Time) ([]*models.DTakoUriageKeihi, error) {
	var data []*models.DTakoUriageKeihi

	if err := r.db.Where("datetime BETWEEN ? AND ?", start, end).
		Order("datetime DESC").
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to list by date range: %w", err)
	}

	return data, nil
}

// isDuplicateKeyError 重複キーエラーの判定
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	// MySQLの重複キーエラーコード: 1062
	return contains(err.Error(), "1062") || contains(err.Error(), "Duplicate entry")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		   len(s) >= len(substr) && contains(s[1:], substr)
}