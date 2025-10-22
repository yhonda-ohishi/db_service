package mysql

import "time"

// Cars carsテーブルのモデル（本番DB）
type Cars struct {
	ID               string     `gorm:"column:id;primaryKey;size:6" json:"id"`
	ID4              int        `gorm:"column:id4" json:"id4"`
	Name             *string    `gorm:"column:name;size:40" json:"name,omitempty"`
	NameR            *string    `gorm:"column:name_R;size:10" json:"name_r,omitempty"`
	Shashu           *string    `gorm:"column:shashu;size:2" json:"shashu,omitempty"`
	Sekisai          *float64   `gorm:"column:sekisai" json:"sekisai,omitempty"`
	Youseki          *float64   `gorm:"column:youseki" json:"youseki,omitempty"`
	RegDate          *time.Time `gorm:"column:reg_date" json:"reg_date,omitempty"`
	NextInspectDate  *time.Time `gorm:"column:next_inspect_date" json:"next_inspect_date,omitempty"`
	ParchDate        *time.Time `gorm:"column:parch_date" json:"parch_date,omitempty"`
	ScrapDate        *time.Time `gorm:"column:scrap_date" json:"scrap_date,omitempty"`
	BumonCodeID      *string    `gorm:"column:bumon_code_id;size:3" json:"bumon_code_id,omitempty"`
	DriverID         *int       `gorm:"column:driver_id" json:"driver_id,omitempty"`
	ETC              *int       `gorm:"column:etc" json:"etc,omitempty"`
	Dai1             int        `gorm:"column:大①" json:"dai1"`
	Chu1             int        `gorm:"column:中①" json:"chu1"`
	Sho1             int        `gorm:"column:小①" json:"sho1"`
	Dai2             int        `gorm:"column:大②" json:"dai2"`
	Chu2             int        `gorm:"column:中②" json:"chu2"`
	Sho2             int        `gorm:"column:小②" json:"sho2"`
	DaiChuSho1       *string    `gorm:"column:大中小①;size:8" json:"daichusho1,omitempty"`
	DaiChuSho2       *string    `gorm:"column:大中小②;size:8" json:"daichusho2,omitempty"`
}

// TableName テーブル名を指定
func (Cars) TableName() string {
	return "cars"
}
