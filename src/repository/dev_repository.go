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

// TimeCardLogRepository インターフェース
type TimeCardLogRepository interface {
	Create(log *mysql.TimeCardLog) error
	Update(log *mysql.TimeCardLog) error
	GetByCompositeKey(datetime string, id int) (*mysql.TimeCardLog, error)
	GetAll(limit, offset int, orderBy string) ([]*mysql.TimeCardLog, int64, error)
	GetByCardID(cardID string, limit, offset int) ([]*mysql.TimeCardLog, int64, error)
	Delete(datetime string, id int) error
}

// TimeCardLogRepositoryImpl 実装
type TimeCardLogRepositoryImpl struct {
	*DevRepository
}

// NewTimeCardLogRepository TimeCardLogRepositoryのコンストラクタ
func NewTimeCardLogRepository(db *gorm.DB) TimeCardLogRepository {
	return &TimeCardLogRepositoryImpl{
		DevRepository: NewDevRepository(db),
	}
}

// Create タイムカードログ作成
func (r *TimeCardLogRepositoryImpl) Create(log *mysql.TimeCardLog) error {
	return r.db.Create(log).Error
}

// Update タイムカードログ更新
func (r *TimeCardLogRepositoryImpl) Update(log *mysql.TimeCardLog) error {
	return r.db.Save(log).Error
}

// GetByCompositeKey 複合主キー（datetime + id）でタイムカードログを取得
func (r *TimeCardLogRepositoryImpl) GetByCompositeKey(datetime string, id int) (*mysql.TimeCardLog, error) {
	var log mysql.TimeCardLog
	if err := r.db.Where("datetime = ? AND id = ?", datetime, id).First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// GetAll 全タイムカードログを取得
func (r *TimeCardLogRepositoryImpl) GetAll(limit, offset int, orderBy string) ([]*mysql.TimeCardLog, int64, error) {
	var logs []*mysql.TimeCardLog
	var totalCount int64

	// 総件数を取得
	if err := r.db.Model(&mysql.TimeCardLog{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// データを取得
	query := r.db.Limit(limit).Offset(offset)
	if orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("datetime DESC")
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, totalCount, nil
}

// GetByCardID カードIDでタイムカードログを取得
func (r *TimeCardLogRepositoryImpl) GetByCardID(cardID string, limit, offset int) ([]*mysql.TimeCardLog, int64, error) {
	var logs []*mysql.TimeCardLog
	var totalCount int64

	// 総件数を取得
	if err := r.db.Model(&mysql.TimeCardLog{}).Where("card_id = ?", cardID).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// データを取得
	if err := r.db.Where("card_id = ?", cardID).
		Order("datetime DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, totalCount, nil
}

// Delete タイムカードログ削除
func (r *TimeCardLogRepositoryImpl) Delete(datetime string, id int) error {
	return r.db.Where("datetime = ? AND id = ?", datetime, id).Delete(&mysql.TimeCardLog{}).Error
}
