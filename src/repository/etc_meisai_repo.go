package repository

import (
	"fmt"
	"time"

	"github.com/yhonda-ohishi/db_service/src/models/mysql"
	"gorm.io/gorm"
)

// ETCMeisaiRepository リポジトリインターフェース
type ETCMeisaiRepository interface {
	Create(data *mysql.ETCMeisai) error
	GetByID(id int64) (*mysql.ETCMeisai, error)
	Update(data *mysql.ETCMeisai) error
	DeleteByID(id int64) error
	List(params *ETCMeisaiListParams) ([]*mysql.ETCMeisai, int64, error)
	ListByHash(hash string) ([]*mysql.ETCMeisai, error)
	ListByDateRange(start, end time.Time) ([]*mysql.ETCMeisai, error)
}

// ETCMeisaiListParams リスト取得用パラメータ
type ETCMeisaiListParams struct {
	Hash      *string
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
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
func (r *etcMeisaiRepo) Create(data *mysql.ETCMeisai) error {
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
func (r *etcMeisaiRepo) GetByID(id int64) (*mysql.ETCMeisai, error) {
	var data mysql.ETCMeisai

	result := r.db.First(&data, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, mysql.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get record: %w", result.Error)
	}

	return &data, nil
}

// Update データ更新
func (r *etcMeisaiRepo) Update(data *mysql.ETCMeisai) error {
	if err := data.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// 既存レコードを確認
	existing, err := r.GetByID(data.ID)
	if err != nil {
		if err == mysql.ErrRecordNotFound {
			return mysql.ErrRecordNotFound
		}
		return err
	}

	// 更新実行
	result := r.db.Model(existing).Updates(data)
	if result.Error != nil {
		return fmt.Errorf("failed to update record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return mysql.ErrRecordNotFound
	}

	return nil
}

// DeleteByID IDでデータ削除
func (r *etcMeisaiRepo) DeleteByID(id int64) error {
	result := r.db.Delete(&mysql.ETCMeisai{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete record: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return mysql.ErrRecordNotFound
	}

	return nil
}

// List 条件付きリスト取得
func (r *etcMeisaiRepo) List(params *ETCMeisaiListParams) ([]*mysql.ETCMeisai, int64, error) {
	var data []*mysql.ETCMeisai
	var totalCount int64

	query := r.db.Model(&mysql.ETCMeisai{})

	// 条件の適用
	if params.Hash != nil && *params.Hash != "" {
		query = query.Where("hash = ?", *params.Hash)
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

// ListByHash hashでリスト取得
func (r *etcMeisaiRepo) ListByHash(hash string) ([]*mysql.ETCMeisai, error) {
	var data []*mysql.ETCMeisai

	if err := r.db.Where("hash = ?", hash).
		Order("date_to DESC").
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to list by hash: %w", err)
	}

	return data, nil
}

// ListByDateRange 日付範囲でリスト取得
func (r *etcMeisaiRepo) ListByDateRange(start, end time.Time) ([]*mysql.ETCMeisai, error) {
	var data []*mysql.ETCMeisai

	if err := r.db.Where("date_to BETWEEN ? AND ?", start, end).
		Order("date_to DESC").
		Find(&data).Error; err != nil {
		return nil, fmt.Errorf("failed to list by date range: %w", err)
	}

	return data, nil
}
