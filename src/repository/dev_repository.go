package repository

import (
	"time"

	"github.com/yhonda-ohishi/db_service/src/models/mysql"
	"gorm.io/gorm"
)

// DevRepository ローカルDB用のリポジトリ（読み書き可能）
type DevRepository struct {
	db *gorm.DB
}

// NewDevRepository DevRepositoryのコンストラクタ
func NewDevRepository(db *gorm.DB) *DevRepository {
	return &DevRepository{
		db: db,
	}
}

// TimeCardDevRepository インターフェース
type TimeCardDevRepository interface {
	Create(timeCard *mysql.TimeCard) error
	Update(timeCard *mysql.TimeCard) error
	GetByCompositeKey(datetime time.Time, id int) (*mysql.TimeCard, error)
	GetAll(limit, offset int, orderBy string) ([]*mysql.TimeCard, int64, error)
	Delete(datetime time.Time, id int) error
}

// TimeCardDevRepositoryImpl 実装
type TimeCardDevRepositoryImpl struct {
	*DevRepository
}

// NewTimeCardDevRepository TimeCardDevRepositoryのコンストラクタ
func NewTimeCardDevRepository(db *gorm.DB) TimeCardDevRepository {
	return &TimeCardDevRepositoryImpl{
		DevRepository: NewDevRepository(db),
	}
}

// Create タイムカードデータ作成
func (r *TimeCardDevRepositoryImpl) Create(timeCard *mysql.TimeCard) error {
	return r.db.Create(timeCard).Error
}

// Update タイムカードデータ更新
func (r *TimeCardDevRepositoryImpl) Update(timeCard *mysql.TimeCard) error {
	return r.db.Save(timeCard).Error
}

// GetByCompositeKey 複合主キー（datetime + id）でタイムカードデータを取得
func (r *TimeCardDevRepositoryImpl) GetByCompositeKey(datetime time.Time, id int) (*mysql.TimeCard, error) {
	var timeCard mysql.TimeCard
	if err := r.db.Where("datetime = ? AND id = ?", datetime, id).First(&timeCard).Error; err != nil {
		return nil, err
	}
	return &timeCard, nil
}

// GetAll 全タイムカードデータを取得
func (r *TimeCardDevRepositoryImpl) GetAll(limit, offset int, orderBy string) ([]*mysql.TimeCard, int64, error) {
	var timeCards []*mysql.TimeCard
	var totalCount int64

	// 総件数を取得
	if err := r.db.Model(&mysql.TimeCard{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// データを取得
	query := r.db.Limit(limit).Offset(offset)
	if orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("datetime DESC")
	}

	if err := query.Find(&timeCards).Error; err != nil {
		return nil, 0, err
	}

	return timeCards, totalCount, nil
}

// Delete タイムカードデータ削除
func (r *TimeCardDevRepositoryImpl) Delete(datetime time.Time, id int) error {
	return r.db.Where("datetime = ? AND id = ?", datetime, id).Delete(&mysql.TimeCard{}).Error
}
