package repository

import (
	"fmt"
	"time"

	"github.com/yhonda-ohishi/db_service/src/models"
	"gorm.io/gorm"
)

// DTakoFerryRowsRepository リポジトリインターフェース
type DTakoFerryRowsRepository interface {
	Create(data *models.DTakoFerryRows) error
	GetByID(id int32) (*models.DTakoFerryRows, error)
	Update(data *models.DTakoFerryRows) error
	DeleteByID(id int32) error
	List(params *DTakoFerryRowsListParams) ([]*models.DTakoFerryRows, int64, error)
	ListByUnkoNo(unkoNo string) ([]*models.DTakoFerryRows, error)
	ListByDateRange(start, end time.Time) ([]*models.DTakoFerryRows, error)
}

// DTakoFerryRowsListParams リスト取得用パラメータ
type DTakoFerryRowsListParams struct {
	UnkoNo    *string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

// dtakoFerryRowsRepo リポジトリ実装
type dtakoFerryRowsRepo struct {
	db *gorm.DB
}

// NewDTakoFerryRowsRepository リポジトリのコンストラクタ
func NewDTakoFerryRowsRepository(db *gorm.DB) DTakoFerryRowsRepository {
	return &dtakoFerryRowsRepo{db: db}
}

// Create データ作成
func (r *dtakoFerryRowsRepo) Create(data *models.DTakoFerryRows) error {
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
func (r *dtakoFerryRowsRepo) GetByID(id int32) (*models.DTakoFerryRows, error) {
	var data models.DTakoFerryRows

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
func (r *dtakoFerryRowsRepo) Update(data *models.DTakoFerryRows) error {
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
func (r *dtakoFerryRowsRepo) DeleteByID(id int32) error {
	result := r.db.Delete(&models.DTakoFerryRows{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return models.ErrRecordNotFound
	}

	return nil
}

// List 条件付きリスト取得
func (r *dtakoFerryRowsRepo) List(params *DTakoFerryRowsListParams) ([]*models.DTakoFerryRows, int64, error) {
	var data []*models.DTakoFerryRows
	var totalCount int64

	query := r.db.Model(&models.DTakoFerryRows{})

	// 条件の適用
	if params.UnkoNo != nil && *params.UnkoNo != "" {
		query = query.Where("運行NO = ?", *params.UnkoNo)
	}
	if params.StartDate != nil {
		query = query.Where("運行日 >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("運行日 <= ?", *params.EndDate)
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
	if err := query.Order("運行日 DESC").Find(&data).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list records: %w", err)
	}

	return data, totalCount, nil
}

// ListByUnkoNo 運行NOでリスト取得
func (r *dtakoFerryRowsRepo) ListByUnkoNo(unkoNo string) ([]*models.DTakoFerryRows, error) {
	var data []*models.DTakoFerryRows

	if err := r.db.Where("運行NO = ?", unkoNo).
		Order("運行日 DESC").
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to list by unko_no: %w", err)
	}

	return data, nil
}

// ListByDateRange 日付範囲でリスト取得
func (r *dtakoFerryRowsRepo) ListByDateRange(start, end time.Time) ([]*models.DTakoFerryRows, error) {
	var data []*models.DTakoFerryRows

	if err := r.db.Where("運行日 BETWEEN ? AND ?", start, end).
		Order("運行日 DESC").
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to list by date range: %w", err)
	}

	return data, nil
}