package repository

import (
	"github.com/yhonda-ohishi/db_service/src/config"
	"github.com/yhonda-ohishi/db_service/src/models"
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
	GetAll(limit, offset int) ([]*models.DTakoCars, int64, error)
	GetByID(id int) (*models.DTakoCars, error)
	GetByCarCode(carCode string) (*models.DTakoCars, error)
}

// DTakoEventsRepository インターフェース
type DTakoEventsRepository interface {
	GetAll(limit, offset int) ([]*models.DTakoEvents, int64, error)
	GetByID(id int64) (*models.DTakoEvents, error)
	GetByOperationNo(operationNo string) ([]*models.DTakoEvents, error)
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
func (r *DTakoCarsRepositoryImpl) GetAll(limit, offset int) ([]*models.DTakoCars, int64, error) {
	var cars []*models.DTakoCars
	var totalCount int64

	// 総数取得
	if err := r.prodDB.DB.Model(&models.DTakoCars{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// データ取得
	if err := r.prodDB.DB.Limit(limit).Offset(offset).Find(&cars).Error; err != nil {
		return nil, 0, err
	}

	return cars, totalCount, nil
}

// GetByID IDで車輌情報を取得
func (r *DTakoCarsRepositoryImpl) GetByID(id int) (*models.DTakoCars, error) {
	var car models.DTakoCars
	if err := r.prodDB.DB.Where("id = ?", id).First(&car).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

// GetByCarCode 車輌CDで車輌情報を取得
func (r *DTakoCarsRepositoryImpl) GetByCarCode(carCode string) (*models.DTakoCars, error) {
	var car models.DTakoCars
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
func (r *DTakoEventsRepositoryImpl) GetAll(limit, offset int) ([]*models.DTakoEvents, int64, error) {
	var events []*models.DTakoEvents
	var totalCount int64

	// 総数取得
	if err := r.prodDB.DB.Model(&models.DTakoEvents{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// データ取得
	if err := r.prodDB.DB.Limit(limit).Offset(offset).Order("開始日時 DESC").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, totalCount, nil
}

// GetByID IDでイベント情報を取得
func (r *DTakoEventsRepositoryImpl) GetByID(id int64) (*models.DTakoEvents, error) {
	var event models.DTakoEvents
	if err := r.prodDB.DB.Where("id = ?", id).First(&event).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

// GetByOperationNo 運行NOでイベント情報を取得
func (r *DTakoEventsRepositoryImpl) GetByOperationNo(operationNo string) ([]*models.DTakoEvents, error) {
	var events []*models.DTakoEvents
	if err := r.prodDB.DB.Where("運行NO = ?", operationNo).Order("開始日時 ASC").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
