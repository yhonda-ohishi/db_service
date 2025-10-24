package repository

import (
	"time"

	"github.com/yhonda-ohishi/db_service/src/config"
	"github.com/yhonda-ohishi/db_service/src/models/mysql"
)

// ProdRepository 本番DB用のリポジトリ（読み取り専用）
type ProdRepository struct {
	prodDB *config.ProdDatabase
}

// NewProdRepository 本番DB用リポジトリのコンストラクタ
func NewProdRepository(prodDB *config.ProdDatabase) *ProdRepository {
	return &ProdRepository{
		prodDB: prodDB,
	}
}

// DTakoCarsRepository インターフェース
type DTakoCarsRepository interface {
	GetAll(limit, offset int) ([]*mysql.DTakoCars, int64, error)
	GetByID(id int) (*mysql.DTakoCars, error)
	GetByCarCode(carCode string) (*mysql.DTakoCars, error)
}

// DTakoEventsRepository インターフェース
type DTakoEventsRepository interface {
	GetAll(limit, offset int, orderBy string) ([]*mysql.DTakoEvents, int64, error)
	GetByID(id int64) (*mysql.DTakoEvents, error)
	GetByOperationNo(operationNo string, eventTypes []string, startTime, endTime *time.Time) ([]*mysql.DTakoEvents, error)
}

// DTakoCarsRepositoryImpl 実装
type DTakoCarsRepositoryImpl struct {
	*ProdRepository
}

// NewDTakoCarsRepository DTakoCarsリポジトリのコンストラクタ
func NewDTakoCarsRepository(prodDB *config.ProdDatabase) DTakoCarsRepository {
	return &DTakoCarsRepositoryImpl{
		ProdRepository: NewProdRepository(prodDB),
	}
}

// GetAll 全車輌情報を取得
func (r *DTakoCarsRepositoryImpl) GetAll(limit, offset int) ([]*mysql.DTakoCars, int64, error) {
	var cars []*mysql.DTakoCars
	var totalCount int64

	// 総数取得
	if err := r.prodDB.DB.Model(&mysql.DTakoCars{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// データ取得
	if err := r.prodDB.DB.Limit(limit).Offset(offset).Find(&cars).Error; err != nil {
		return nil, 0, err
	}

	return cars, totalCount, nil
}

// GetByID IDで車輌情報を取得
func (r *DTakoCarsRepositoryImpl) GetByID(id int) (*mysql.DTakoCars, error) {
	var car mysql.DTakoCars
	if err := r.prodDB.DB.Where("id = ?", id).First(&car).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

// GetByCarCode 車輌CDで車輌情報を取得
func (r *DTakoCarsRepositoryImpl) GetByCarCode(carCode string) (*mysql.DTakoCars, error) {
	var car mysql.DTakoCars
	if err := r.prodDB.DB.Where("車輌CD = ?", carCode).First(&car).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

// DTakoEventsRepositoryImpl 実装
type DTakoEventsRepositoryImpl struct {
	*ProdRepository
}

// NewDTakoEventsRepository DTakoEventsリポジトリのコンストラクタ
func NewDTakoEventsRepository(prodDB *config.ProdDatabase) DTakoEventsRepository {
	return &DTakoEventsRepositoryImpl{
		ProdRepository: NewProdRepository(prodDB),
	}
}

// GetAll 全イベント情報を取得
func (r *DTakoEventsRepositoryImpl) GetAll(limit, offset int, orderBy string) ([]*mysql.DTakoEvents, int64, error) {
	var events []*mysql.DTakoEvents
	var totalCount int64

	// 総数取得
	if err := r.prodDB.DB.Model(&mysql.DTakoEvents{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// デフォルトのorder byを設定
	if orderBy == "" {
		orderBy = "開始日時 DESC"
	}

	// データ取得
	if err := r.prodDB.DB.Limit(limit).Offset(offset).Order(orderBy).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, totalCount, nil
}

// GetByID IDでイベント情報を取得
func (r *DTakoEventsRepositoryImpl) GetByID(id int64) (*mysql.DTakoEvents, error) {
	var event mysql.DTakoEvents
	if err := r.prodDB.DB.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// GetByOperationNo 運行NOでイベント情報を取得（フィルタ付き）
func (r *DTakoEventsRepositoryImpl) GetByOperationNo(operationNo string, eventTypes []string, startTime, endTime *time.Time) ([]*mysql.DTakoEvents, error) {
	var events []*mysql.DTakoEvents
	query := r.prodDB.DB.Where("運行NO = ?", operationNo)

	// イベントタイプでフィルタ
	if len(eventTypes) > 0 {
		query = query.Where("イベント種類 IN ?", eventTypes)
	}

	// 開始時刻でフィルタ
	if startTime != nil {
		query = query.Where("開始日時 >= ?", startTime)
	}

	// 終了時刻でフィルタ
	if endTime != nil {
		query = query.Where("開始日時 <= ?", endTime)
	}

	if err := query.Order("開始日時 ASC").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
