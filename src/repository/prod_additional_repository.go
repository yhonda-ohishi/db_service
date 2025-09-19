package repository

import (
	"github.com/yhonda-ohishi/db_service/src/config"
	"github.com/yhonda-ohishi/db_service/src/models"
)

// DTakoRowsRepository インターフェース
type DTakoRowsRepository interface {
	GetAll(limit, offset int) ([]*models.DTakoRows, int64, error)
	GetByID(id string) (*models.DTakoRows, error)
	GetByOperationNo(operationNo string) ([]*models.DTakoRows, error)
}

// ETCNumRepository インターフェース
type ETCNumRepository interface {
	GetAll(limit, offset int) ([]*models.ETCNum, int64, error)
	GetByETCCardNum(etcCardNum string) ([]*models.ETCNum, error)
	GetByCarID(carID string) ([]*models.ETCNum, error)
}

// DTakoRowsRepositoryImpl 実装
type DTakoRowsRepositoryImpl struct {
	*ProdRepository
}

// NewDTakoRowsRepository DTakoRowsリポジトリのコンストラクタ
func NewDTakoRowsRepository(prodDB *config.ProdDatabase) DTakoRowsRepository {
	return &DTakoRowsRepositoryImpl{
		ProdRepository: NewProdRepository(prodDB),
	}
}

// GetAll 全運行データを取得
func (r *DTakoRowsRepositoryImpl) GetAll(limit, offset int) ([]*models.DTakoRows, int64, error) {
	var rows []*models.DTakoRows
	var totalCount int64

	// 総数取得
	if err := r.prodDB.DB.Model(&models.DTakoRows{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// データ取得
	if err := r.prodDB.DB.Limit(limit).Offset(offset).Order("読取日 DESC").Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, totalCount, nil
}

// GetByID IDで運行データを取得
func (r *DTakoRowsRepositoryImpl) GetByID(id string) (*models.DTakoRows, error) {
	var row models.DTakoRows
	if err := r.prodDB.DB.Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// GetByOperationNo 運行NOで運行データを取得
func (r *DTakoRowsRepositoryImpl) GetByOperationNo(operationNo string) ([]*models.DTakoRows, error) {
	var rows []*models.DTakoRows
	if err := r.prodDB.DB.Where("運行NO = ?", operationNo).Order("読取日 ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ETCNumRepositoryImpl 実装
type ETCNumRepositoryImpl struct {
	*ProdRepository
}

// NewETCNumRepository ETCNumリポジトリのコンストラクタ
func NewETCNumRepository(prodDB *config.ProdDatabase) ETCNumRepository {
	return &ETCNumRepositoryImpl{
		ProdRepository: NewProdRepository(prodDB),
	}
}

// GetAll 全ETCカード番号を取得
func (r *ETCNumRepositoryImpl) GetAll(limit, offset int) ([]*models.ETCNum, int64, error) {
	var etcNums []*models.ETCNum
	var totalCount int64

	// 総数取得
	if err := r.prodDB.DB.Model(&models.ETCNum{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// データ取得
	if err := r.prodDB.DB.Limit(limit).Offset(offset).Find(&etcNums).Error; err != nil {
		return nil, 0, err
	}

	return etcNums, totalCount, nil
}

// GetByETCCardNum ETCカード番号でデータを取得
func (r *ETCNumRepositoryImpl) GetByETCCardNum(etcCardNum string) ([]*models.ETCNum, error) {
	var etcNums []*models.ETCNum
	if err := r.prodDB.DB.Where("etc_card_num = ?", etcCardNum).Find(&etcNums).Error; err != nil {
		return nil, err
	}
	return etcNums, nil
}

// GetByCarID 車輌IDでETCカード番号を取得
func (r *ETCNumRepositoryImpl) GetByCarID(carID string) ([]*models.ETCNum, error) {
	var etcNums []*models.ETCNum
	if err := r.prodDB.DB.Where("car_id = ?", carID).Find(&etcNums).Error; err != nil {
		return nil, err
	}
	return etcNums, nil
}