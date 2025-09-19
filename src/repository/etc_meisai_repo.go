package repository

import (
	"fmt"
	"time"

	"github.com/yhonda-ohishi/db_service/src/models"
	"gorm.io/gorm"
)

// ETCMeisaiRepository リポジトリインターフェース
type ETCMeisaiRepository interface {
	Create(data *models.ETCMeisai) error
	GetByID(id int64) (*models.ETCMeisai, error)
	Update(data *models.ETCMeisai) error
	DeleteByID(id int64) error
	List(params *ETCMeisaiListParams) ([]*models.ETCMeisai, int64, error)
	ListByDtakoRowID(dtakoRowID string) ([]*models.ETCMeisai, error)
	ListByDateRange(start, end time.Time) ([]*models.ETCMeisai, error)
}

// ETCMeisaiListParams リスト取得用パラメータ
type ETCMeisaiListParams struct {
	DtakoRowID *string
	StartDate  *time.Time
	EndDate    *time.Time
	Limit      int
	Offset     int
}

// etcMeisaiRepo リポジトリ実装
type etcMeisaiRepo struct {
	db *gorm.DB
}

// NewETCMeisaiRepository リポジトリのコンストラクタ
func NewETCMeisaiRepository(db *gorm.DB) ETCMeisaiRepository {
	return &etcMeisaiRepo{db: db}
}

// Create データ作成
func (r *etcMeisaiRepo) Create(data *models.ETCMeisai) error {
	if err := data.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	result := r.db.Create(data)
	if result.Error != nil {
		return fmt.Errorf("failed to create record: %w", result.Error)
	}

	return nil
}

// GetByID IDでデータ取得
func (r *etcMeisaiRepo) GetByID(id int64) (*models.ETCMeisai, error) {
	var data models.ETCMeisai

	result := r.db.First(&data, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, models.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get record: %w", result.Error)
	}

	return &data, nil
}

// Update データ更新
func (r *etcMeisaiRepo) Update(data *models.ETCMeisai) error {
	if err := data.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// 既存レコードを確認
	existing, err := r.GetByID(data.ID)
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

// DeleteByID IDでデータ削除
func (r *etcMeisaiRepo) DeleteByID(id int64) error {
	result := r.db.Delete(&models.ETCMeisai{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.ErrRecordNotFound
	}

	return nil
}

// List 条件付きリスト取得
func (r *etcMeisaiRepo) List(params *ETCMeisaiListParams) ([]*models.ETCMeisai, int64, error) {
	var data []*models.ETCMeisai
	var totalCount int64

	query := r.db.Model(&models.ETCMeisai{})

	// 条件の適用
	if params.DtakoRowID != nil && *params.DtakoRowID != "" {
		query = query.Where("dtako_row_id = ?", *params.DtakoRowID)
	}
	if params.StartDate != nil {
		query = query.Where("date_to >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("date_to <= ?", *params.EndDate)
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
	if err := query.Order("date_to DESC").Find(&data).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list records: %w", err)
	}

	return data, totalCount, nil
}

// ListByDtakoRowID dtako_row_idでリスト取得
func (r *etcMeisaiRepo) ListByDtakoRowID(dtakoRowID string) ([]*models.ETCMeisai, error) {
	var data []*models.ETCMeisai

	if err := r.db.Where("dtako_row_id = ?", dtakoRowID).
		Order("date_to DESC").
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to list by dtako_row_id: %w", err)
	}

	return data, nil
}

// ListByDateRange 日付範囲でリスト取得
func (r *etcMeisaiRepo) ListByDateRange(start, end time.Time) ([]*models.ETCMeisai, error) {
	var data []*models.ETCMeisai

	if err := r.db.Where("date_to BETWEEN ? AND ?", start, end).
		Order("date_to DESC").
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to list by date range: %w", err)
	}

	return data, nil
}