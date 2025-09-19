package models

// DTakoCars dtako_carsテーブルのモデル（本番DB）
type DTakoCars struct {
	ID                  int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CarCode             string `gorm:"column:車輌CD;size:10" json:"car_code"`
	CarCC               string `gorm:"column:車輌CC;size:6" json:"car_cc"`
	CarName             string `gorm:"column:車輌名;size:40" json:"car_name"`
	BelongOfficeCode    int    `gorm:"column:所属事業所CD" json:"belong_office_code"`
	HighwayCarType      int    `gorm:"column:高速車種区分" json:"highway_car_type"`
	FerryCarType        int    `gorm:"column:航送車種区分" json:"ferry_car_type"`
	EvaluationClassCode int    `gorm:"column:評価分類CD" json:"evaluation_class_code"`
	IdlingType          int    `gorm:"column:アイドリング取扱区分" json:"idling_type"`
	MaxLoadWeightKg     int    `gorm:"column:最大積載重量kg" json:"max_load_weight_kg"`
	CarClass1           int    `gorm:"column:車輌分類1" json:"car_class1"`
	CarClass2           int    `gorm:"column:車輌分類2" json:"car_class2"`
	CarClass3           int    `gorm:"column:車輌分類3" json:"car_class3"`
	CarClass4           int    `gorm:"column:車輌分類4" json:"car_class4"`
	CarClass5           int    `gorm:"column:車輌分類5" json:"car_class5"`
	OperationType       int    `gorm:"column:運用タイプ" json:"operation_type"`
}

// TableName テーブル名を指定
func (DTakoCars) TableName() string {
	return "dtako_cars"
}