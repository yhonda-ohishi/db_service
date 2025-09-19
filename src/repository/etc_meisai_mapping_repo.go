package repository

import (
	"fmt"

	"github.com/yhonda-ohishi/db_service/src/models"
	"gorm.io/gorm"
)

// ETCMeisaiMappingRepository リポジトリインターフェース
type ETCMeisaiMappingRepository interface {
	Create(data *models.ETCMeisaiMapping) error
	GetByID(id int64) (*models.ETCMeisaiMapping, error)
	Update(data *models.ETCMeisaiMapping) error
	DeleteByID(id int64) error
	List(params *ETCMeisaiMappingListParams) ([]*models.ETCMeisaiMapping, int64, error)
	GetDTakoRowIDsByHash(hash string) ([]string, error)
}

// ETCMeisaiMappingListParams リスト取得用パラメータ
type ETCMeisaiMappingListParams struct {
	ETCMeisaiHash *string
	DTakoRowID    *string
	Limit         int
	Offset        int
}

// etcMeisaiMappingRepo リポジトリ実装
type etcMeisaiMappingRepo struct {
	db *gorm.DB
}

// NewETCMeisaiMappingRepository リポジトリのコンストラクタ
func NewETCMeisaiMappingRepository(db *gorm.DB) ETCMeisaiMappingRepository {
	return &etcMeisaiMappingRepo{db: db}
}

// Create マッピング作成
func (r *etcMeisaiMappingRepo) Create(data *models.ETCMeisaiMapping) error {
	if err := data.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := r.db.Create(data).Error; err != nil {
		return fmt.Errorf("failed to create mapping: %w", err)
	}

	return nil
}

// GetByID ID指定でマッピング取得
func (r *etcMeisaiMappingRepo) GetByID(id int64) (*models.ETCMeisaiMapping, error) {
	var data models.ETCMeisaiMapping
	if err := r.db.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("mapping not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get mapping: %w", err)
	}

	return &data, nil
}

// Update マッピング更新
func (r *etcMeisaiMappingRepo) Update(data *models.ETCMeisaiMapping) error {
	if err := data.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := r.db.Save(data).Error; err != nil {
		return fmt.Errorf("failed to update mapping: %w", err)
	}

	return nil
}

// DeleteByID ID指定でマッピング削除
func (r *etcMeisaiMappingRepo) DeleteByID(id int64) error {
	result := r.db.Delete(&models.ETCMeisaiMapping{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete mapping: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("mapping not found")
	}

	return nil
}

// List マッピング一覧取得
func (r *etcMeisaiMappingRepo) List(params *ETCMeisaiMappingListParams) ([]*models.ETCMeisaiMapping, int64, error) {
	var data []*models.ETCMeisaiMapping
	var totalCount int64

	query := r.db.Model(&models.ETCMeisaiMapping{})

	// 条件の適用
	if params.ETCMeisaiHash != nil && *params.ETCMeisaiHash != "" {
		query = query.Where("etc_meisai_hash = ?", *params.ETCMeisaiHash)
	}
	if params.DTakoRowID != nil && *params.DTakoRowID != "" {
		query = query.Where("dtako_row_id = ?", *params.DTakoRowID)
	}

	// 総数取得
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count mappings: %w", err)
	}

	// データ取得
	if err := query.Order("created_at DESC").
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&data).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list mappings: %w", err)
	}

	return data, totalCount, nil
}

// GetDTakoRowIDsByHash ハッシュからDTakoRowIDのリストを取得
func (r *etcMeisaiMappingRepo) GetDTakoRowIDsByHash(hash string) ([]string, error) {
	var mappings []*models.ETCMeisaiMapping

	if err := r.db.Where("etc_meisai_hash = ?", hash).
		Find(&mappings).Error; err != nil {
		return nil, fmt.Errorf("failed to get mappings by hash: %w", err)
	}

	dtakoRowIDs := make([]string, len(mappings))
	for i, mapping := range mappings {
		dtakoRowIDs[i] = mapping.DTakoRowID
	}

	return dtakoRowIDs, nil
}