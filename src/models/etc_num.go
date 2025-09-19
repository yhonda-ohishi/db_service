package models

import "time"

// ETCNum etc_numテーブルのモデル（本番DB）
type ETCNum struct {
	ETCCardNum    string     `gorm:"column:etc_card_num;primaryKey;size:20" json:"etc_card_num"`
	CarID         string     `gorm:"column:car_id;primaryKey;size:6" json:"car_id"`
	StartDateTime *time.Time `gorm:"column:start_date_time;type:datetime" json:"start_date_time"`
	DueDateTime   *time.Time `gorm:"column:due_date_time;type:datetime" json:"due_date_time"`
	ToChange      *bool      `gorm:"column:to_change;type:tinyint(1)" json:"to_change"`
}

// TableName テーブル名を指定
func (ETCNum) TableName() string {
	return "etc_num"
}