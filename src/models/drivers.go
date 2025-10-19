package models

import "time"

// Drivers driversテーブルのモデル（本番DB）
type Drivers struct {
	ID          int        `gorm:"column:id;primaryKey" json:"id"`
	Name        *string    `gorm:"column:name;size:40" json:"name,omitempty"`
	ShainR      *string    `gorm:"column:社員R;size:20" json:"shain_r,omitempty"`
	Bumon       string     `gorm:"column:bumon;size:3" json:"bumon"`
	JoinDate    *time.Time `gorm:"column:join_date" json:"join_date,omitempty"`
	RetireDate  *time.Time `gorm:"column:retire_date" json:"retire_date,omitempty"`
	Bunrui1     *string    `gorm:"column:bunrui1;size:5" json:"bunrui1,omitempty"`
	Bunrui2     *string    `gorm:"column:bunrui2;size:5" json:"bunrui2,omitempty"`
	Kubun       *int       `gorm:"column:kubun" json:"kubun,omitempty"`
	KinmuTaikei int        `gorm:"column:勤務体系" json:"kinmu_taikei"`
}

// TableName テーブル名を指定
func (Drivers) TableName() string {
	return "drivers"
}
