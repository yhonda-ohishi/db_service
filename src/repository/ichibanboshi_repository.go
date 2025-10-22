package repository

import (
	"github.com/yhonda-ohishi/db_service/src/config"
	"github.com/yhonda-ohishi/db_service/src/models/ichibanboshi"
)

// IchibanboshiRepository SQL Server用のリポジトリ（読み取り専用）
type IchibanboshiRepository struct {
	sqlServerDB *config.SQLServerDatabase
}

// NewIchibanboshiRepository SQL Server用リポジトリのコンストラクタ
func NewIchibanboshiRepository(sqlServerDB *config.SQLServerDatabase) *IchibanboshiRepository {
	return &IchibanboshiRepository{
		sqlServerDB: sqlServerDB,
	}
}

// UntenNippoMeisaiRepository 運転日報明細リポジトリインターフェース
type UntenNippoMeisaiRepository interface {
	GetAll(limit, offset int, orderBy string) ([]*ichibanboshi.UntenNippoMeisai, int64, error)
	GetByNippoK(nippoK, haishaK, sharyoC string) (*ichibanboshi.UntenNippoMeisai, error)
	GetBySharyoC(sharyoC string, limit int) ([]*ichibanboshi.UntenNippoMeisai, error)
	GetByDateRange(startDate, endDate string, limit, offset int) ([]*ichibanboshi.UntenNippoMeisai, int64, error)
}

// ShainMasterRepository 社員マスタリポジトリインターフェース
type ShainMasterRepository interface {
	GetAll(limit, offset int, orderBy string) ([]*ichibanboshi.ShainMaster, int64, error)
	GetByShainC(shainC string) (*ichibanboshi.ShainMaster, error)
	GetByBumonC(bumonC string) ([]*ichibanboshi.ShainMaster, error)
}

// ChiikiMasterRepository 地域マスタリポジトリインターフェース
type ChiikiMasterRepository interface {
	GetAll(limit, offset int, orderBy string) ([]*ichibanboshi.ChiikiMaster, int64, error)
	GetByChiikiC(chiikiC string) (*ichibanboshi.ChiikiMaster, error)
}

// ChikuMasterRepository 地区マスタリポジトリインターフェース
type ChikuMasterRepository interface {
	GetAll(limit, offset int, orderBy string) ([]*ichibanboshi.ChikuMaster, int64, error)
	GetByChikuC(chikuC string) (*ichibanboshi.ChikuMaster, error)
	GetByChiikiC(chiikiC string) ([]*ichibanboshi.ChikuMaster, error)
}

// UntenNippoMeisaiRepositoryImpl 運転日報明細リポジトリ実装
type UntenNippoMeisaiRepositoryImpl struct {
	*IchibanboshiRepository
}

// NewUntenNippoMeisaiRepository 運転日報明細リポジトリのコンストラクタ
func NewUntenNippoMeisaiRepository(sqlServerDB *config.SQLServerDatabase) UntenNippoMeisaiRepository {
	return &UntenNippoMeisaiRepositoryImpl{
		IchibanboshiRepository: NewIchibanboshiRepository(sqlServerDB),
	}
}

// GetAll 全運転日報明細を取得
func (r *UntenNippoMeisaiRepositoryImpl) GetAll(limit, offset int, orderBy string) ([]*ichibanboshi.UntenNippoMeisai, int64, error) {
	var meisai []*ichibanboshi.UntenNippoMeisai
	var totalCount int64

	// 総数取得
	if err := r.sqlServerDB.DB.Model(&ichibanboshi.UntenNippoMeisai{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// デフォルトのorder byを設定
	if orderBy == "" {
		orderBy = "管理年月日 DESC"
	}

	// データ取得
	if err := r.sqlServerDB.DB.Limit(limit).Offset(offset).Order(orderBy).Find(&meisai).Error; err != nil {
		return nil, 0, err
	}

	return meisai, totalCount, nil
}

// GetByNippoK 日報K、配車K、車輌Cで運転日報明細を取得（複合主キー）
func (r *UntenNippoMeisaiRepositoryImpl) GetByNippoK(nippoK, haishaK, sharyoC string) (*ichibanboshi.UntenNippoMeisai, error) {
	var meisai ichibanboshi.UntenNippoMeisai
	if err := r.sqlServerDB.DB.Where("日報K = ? AND 配車K = ? AND 車輌C = ?", nippoK, haishaK, sharyoC).First(&meisai).Error; err != nil {
		return nil, err
	}
	return &meisai, nil
}

// GetBySharyoC 車輌Cで運転日報明細を取得
func (r *UntenNippoMeisaiRepositoryImpl) GetBySharyoC(sharyoC string, limit int) ([]*ichibanboshi.UntenNippoMeisai, error) {
	var meisai []*ichibanboshi.UntenNippoMeisai
	if err := r.sqlServerDB.DB.Where("車輌C = ?", sharyoC).Limit(limit).Order("管理年月日 DESC").Find(&meisai).Error; err != nil {
		return nil, err
	}
	return meisai, nil
}

// GetByDateRange 日付範囲で運転日報明細を取得
func (r *UntenNippoMeisaiRepositoryImpl) GetByDateRange(startDate, endDate string, limit, offset int) ([]*ichibanboshi.UntenNippoMeisai, int64, error) {
	var meisai []*ichibanboshi.UntenNippoMeisai
	var totalCount int64

	query := r.sqlServerDB.DB.Where("管理年月日 BETWEEN ? AND ?", startDate, endDate)

	// 総数取得
	if err := query.Model(&ichibanboshi.UntenNippoMeisai{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// データ取得
	if err := query.Limit(limit).Offset(offset).Order("管理年月日 DESC").Find(&meisai).Error; err != nil {
		return nil, 0, err
	}

	return meisai, totalCount, nil
}

// ShainMasterRepositoryImpl 社員マスタリポジトリ実装
type ShainMasterRepositoryImpl struct {
	*IchibanboshiRepository
}

// NewShainMasterRepository 社員マスタリポジトリのコンストラクタ
func NewShainMasterRepository(sqlServerDB *config.SQLServerDatabase) ShainMasterRepository {
	return &ShainMasterRepositoryImpl{
		IchibanboshiRepository: NewIchibanboshiRepository(sqlServerDB),
	}
}

// GetAll 全社員マスタを取得
func (r *ShainMasterRepositoryImpl) GetAll(limit, offset int, orderBy string) ([]*ichibanboshi.ShainMaster, int64, error) {
	var shain []*ichibanboshi.ShainMaster
	var totalCount int64

	// 総数取得
	if err := r.sqlServerDB.DB.Model(&ichibanboshi.ShainMaster{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// デフォルトのorder byを設定
	if orderBy == "" {
		orderBy = "社員C ASC"
	}

	// データ取得
	if err := r.sqlServerDB.DB.Limit(limit).Offset(offset).Order(orderBy).Find(&shain).Error; err != nil {
		return nil, 0, err
	}

	return shain, totalCount, nil
}

// GetByShainC 社員Cで社員マスタを取得
func (r *ShainMasterRepositoryImpl) GetByShainC(shainC string) (*ichibanboshi.ShainMaster, error) {
	var shain ichibanboshi.ShainMaster
	if err := r.sqlServerDB.DB.Where("社員C = ?", shainC).First(&shain).Error; err != nil {
		return nil, err
	}
	return &shain, nil
}

// GetByBumonC 部門Cで社員マスタを取得
func (r *ShainMasterRepositoryImpl) GetByBumonC(bumonC string) ([]*ichibanboshi.ShainMaster, error) {
	var shain []*ichibanboshi.ShainMaster
	if err := r.sqlServerDB.DB.Where("部門C = ?", bumonC).Order("社員C ASC").Find(&shain).Error; err != nil {
		return nil, err
	}
	return shain, nil
}

// ChiikiMasterRepositoryImpl 地域マスタリポジトリ実装
type ChiikiMasterRepositoryImpl struct {
	*IchibanboshiRepository
}

// NewChiikiMasterRepository 地域マスタリポジトリのコンストラクタ
func NewChiikiMasterRepository(sqlServerDB *config.SQLServerDatabase) ChiikiMasterRepository {
	return &ChiikiMasterRepositoryImpl{
		IchibanboshiRepository: NewIchibanboshiRepository(sqlServerDB),
	}
}

// GetAll 全地域マスタを取得
func (r *ChiikiMasterRepositoryImpl) GetAll(limit, offset int, orderBy string) ([]*ichibanboshi.ChiikiMaster, int64, error) {
	var chiiki []*ichibanboshi.ChiikiMaster
	var totalCount int64

	// 総数取得
	if err := r.sqlServerDB.DB.Model(&ichibanboshi.ChiikiMaster{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// デフォルトのorder byを設定
	if orderBy == "" {
		orderBy = "地域C ASC"
	}

	// データ取得
	if err := r.sqlServerDB.DB.Limit(limit).Offset(offset).Order(orderBy).Find(&chiiki).Error; err != nil {
		return nil, 0, err
	}

	return chiiki, totalCount, nil
}

// GetByChiikiC 地域Cで地域マスタを取得
func (r *ChiikiMasterRepositoryImpl) GetByChiikiC(chiikiC string) (*ichibanboshi.ChiikiMaster, error) {
	var chiiki ichibanboshi.ChiikiMaster
	if err := r.sqlServerDB.DB.Where("地域C = ?", chiikiC).First(&chiiki).Error; err != nil {
		return nil, err
	}
	return &chiiki, nil
}

// ChikuMasterRepositoryImpl 地区マスタリポジトリ実装
type ChikuMasterRepositoryImpl struct {
	*IchibanboshiRepository
}

// NewChikuMasterRepository 地区マスタリポジトリのコンストラクタ
func NewChikuMasterRepository(sqlServerDB *config.SQLServerDatabase) ChikuMasterRepository {
	return &ChikuMasterRepositoryImpl{
		IchibanboshiRepository: NewIchibanboshiRepository(sqlServerDB),
	}
}

// GetAll 全地区マスタを取得
func (r *ChikuMasterRepositoryImpl) GetAll(limit, offset int, orderBy string) ([]*ichibanboshi.ChikuMaster, int64, error) {
	var chiku []*ichibanboshi.ChikuMaster
	var totalCount int64

	// 総数取得
	if err := r.sqlServerDB.DB.Model(&ichibanboshi.ChikuMaster{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// デフォルトのorder byを設定
	if orderBy == "" {
		orderBy = "地区C ASC"
	}

	// データ取得
	if err := r.sqlServerDB.DB.Limit(limit).Offset(offset).Order(orderBy).Find(&chiku).Error; err != nil {
		return nil, 0, err
	}

	return chiku, totalCount, nil
}

// GetByChikuC 地区Cで地区マスタを取得
func (r *ChikuMasterRepositoryImpl) GetByChikuC(chikuC string) (*ichibanboshi.ChikuMaster, error) {
	var chiku ichibanboshi.ChikuMaster
	if err := r.sqlServerDB.DB.Where("地区C = ?", chikuC).First(&chiku).Error; err != nil {
		return nil, err
	}
	return &chiku, nil
}

// GetByChiikiC 地域Cで地区マスタを取得
func (r *ChikuMasterRepositoryImpl) GetByChiikiC(chiikiC string) ([]*ichibanboshi.ChikuMaster, error) {
	var chiku []*ichibanboshi.ChikuMaster
	if err := r.sqlServerDB.DB.Where("地域C = ?", chiikiC).Order("地区C ASC").Find(&chiku).Error; err != nil {
		return nil, err
	}
	return chiku, nil
}
