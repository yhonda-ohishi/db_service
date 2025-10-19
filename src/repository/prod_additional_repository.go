package repository

import (
	"github.com/yhonda-ohishi/db_service/src/config"
	"github.com/yhonda-ohishi/db_service/src/models"
)

// DTakoFerryRowsRepository インターフェース（本番DB用）
type DTakoFerryRowsProdRepository interface {
	GetAll(limit, offset int) ([]*models.DTakoFerryRows, int64, error)
	GetByID(id int32) (*models.DTakoFerryRows, error)
	GetByUnkoNo(unkoNo string) ([]*models.DTakoFerryRows, error)
}

// DTakoRowsRepository インターフェース
type DTakoRowsRepository interface {
	GetAll(limit, offset int, orderBy string) ([]*models.DTakoRows, int64, error)
	GetByID(id string) (*models.DTakoRows, error)
	GetByOperationNo(operationNo string) ([]*models.DTakoRows, error)
}

// ETCNumRepository インターフェース
type ETCNumRepository interface {
	GetAll(limit, offset int) ([]*models.ETCNum, int64, error)
	GetByETCCardNum(etcCardNum string) ([]*models.ETCNum, error)
	GetByCarID(carID string) ([]*models.ETCNum, error)
}

// DTakoFerryRowsProdRepositoryImpl 本番DB用実装
type DTakoFerryRowsProdRepositoryImpl struct {
	*ProdRepository
}

// NewDTakoFerryRowsProdRepository DTakoFerryRowsリポジトリのコンストラクタ（本番DB用）
func NewDTakoFerryRowsProdRepository(prodDB *config.ProdDatabase) DTakoFerryRowsProdRepository {
	return &DTakoFerryRowsProdRepositoryImpl{
		ProdRepository: NewProdRepository(prodDB),
	}
}

// GetAll 全フェリー運行データを取得
func (r *DTakoFerryRowsProdRepositoryImpl) GetAll(limit, offset int) ([]*models.DTakoFerryRows, int64, error) {
	var rows []*models.DTakoFerryRows
	var totalCount int64

	// 総数取得
	if err := r.prodDB.DB.Model(&models.DTakoFerryRows{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// データ取得
	if err := r.prodDB.DB.Limit(limit).Offset(offset).Order("運行日 DESC").Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, totalCount, nil
}

// GetByID IDでフェリー運行データを取得
func (r *DTakoFerryRowsProdRepositoryImpl) GetByID(id int32) (*models.DTakoFerryRows, error) {
	var row models.DTakoFerryRows
	if err := r.prodDB.DB.Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// GetByUnkoNo 運行NOでフェリー運行データを取得
func (r *DTakoFerryRowsProdRepositoryImpl) GetByUnkoNo(unkoNo string) ([]*models.DTakoFerryRows, error) {
	var rows []*models.DTakoFerryRows
	if err := r.prodDB.DB.Where("運行NO = ?", unkoNo).Order("運行日 ASC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
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
func (r *DTakoRowsRepositoryImpl) GetAll(limit, offset int, orderBy string) ([]*models.DTakoRows, int64, error) {
	var rows []*models.DTakoRows
	var totalCount int64

	// 総数取得
	if err := r.prodDB.DB.Model(&models.DTakoRows{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// デフォルトのorder byを設定
	if orderBy == "" {
		orderBy = "読取日 DESC"
	}

	// データ取得
	if err := r.prodDB.DB.Limit(limit).Offset(offset).Order(orderBy).Find(&rows).Error; err != nil {
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

// CarsRepository インターフェース
type CarsRepository interface {
	GetAll(limit, offset int, orderBy string) ([]*models.Cars, int64, error)
	GetByID(id string) (*models.Cars, error)
	GetByBumonCodeID(bumonCodeID string) ([]*models.Cars, error)
}

// DriversRepository インターフェース
type DriversRepository interface {
	GetAll(limit, offset int, orderBy string) ([]*models.Drivers, int64, error)
	GetByID(id int) (*models.Drivers, error)
	GetByBumon(bumon string) ([]*models.Drivers, error)
}

// CarsRepositoryImpl 実装
type CarsRepositoryImpl struct {
	*ProdRepository
}

// NewCarsRepository Carsリポジトリのコンストラクタ
func NewCarsRepository(prodDB *config.ProdDatabase) CarsRepository {
	return &CarsRepositoryImpl{
		ProdRepository: NewProdRepository(prodDB),
	}
}

// GetAll 全車両情報を取得
func (r *CarsRepositoryImpl) GetAll(limit, offset int, orderBy string) ([]*models.Cars, int64, error) {
	var cars []*models.Cars
	var totalCount int64

	// 総数取得
	if err := r.prodDB.DB.Model(&models.Cars{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// デフォルトのorder byを設定
	if orderBy == "" {
		orderBy = "id ASC"
	}

	// データ取得
	if err := r.prodDB.DB.Limit(limit).Offset(offset).Order(orderBy).Find(&cars).Error; err != nil {
		return nil, 0, err
	}

	return cars, totalCount, nil
}

// GetByID IDで車両情報を取得
func (r *CarsRepositoryImpl) GetByID(id string) (*models.Cars, error) {
	var car models.Cars
	if err := r.prodDB.DB.Where("id = ?", id).First(&car).Error; err != nil {
		return nil, err
	}
	return &car, nil
}

// GetByBumonCodeID 部門コードで車両情報を取得
func (r *CarsRepositoryImpl) GetByBumonCodeID(bumonCodeID string) ([]*models.Cars, error) {
	var cars []*models.Cars
	if err := r.prodDB.DB.Where("bumon_code_id = ?", bumonCodeID).Order("id ASC").Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}

// DriversRepositoryImpl 実装
type DriversRepositoryImpl struct {
	*ProdRepository
}

// NewDriversRepository Driversリポジトリのコンストラクタ
func NewDriversRepository(prodDB *config.ProdDatabase) DriversRepository {
	return &DriversRepositoryImpl{
		ProdRepository: NewProdRepository(prodDB),
	}
}

// GetAll 全ドライバー情報を取得
func (r *DriversRepositoryImpl) GetAll(limit, offset int, orderBy string) ([]*models.Drivers, int64, error) {
	var drivers []*models.Drivers
	var totalCount int64

	// 総数取得
	if err := r.prodDB.DB.Model(&models.Drivers{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// デフォルトのorder byを設定
	if orderBy == "" {
		orderBy = "id ASC"
	}

	// データ取得
	if err := r.prodDB.DB.Limit(limit).Offset(offset).Order(orderBy).Find(&drivers).Error; err != nil {
		return nil, 0, err
	}

	return drivers, totalCount, nil
}

// GetByID IDでドライバー情報を取得
func (r *DriversRepositoryImpl) GetByID(id int) (*models.Drivers, error) {
	var driver models.Drivers
	if err := r.prodDB.DB.Where("id = ?", id).First(&driver).Error; err != nil {
		return nil, err
	}
	return &driver, nil
}

// GetByBumon 部門コードでドライバー情報を取得
func (r *DriversRepositoryImpl) GetByBumon(bumon string) ([]*models.Drivers, error) {
	var drivers []*models.Drivers
	if err := r.prodDB.DB.Where("bumon = ?", bumon).Order("id ASC").Find(&drivers).Error; err != nil {
		return nil, err
	}
	return drivers, nil
}
