package mysql

import "time"

// TimeCard タイムカードモデル（本番DB、読み取り専用）
type TimeCard struct {
	Datetime     time.Time `gorm:"column:datetime;primaryKey;type:datetime;not null"`
	ID           int       `gorm:"column:id;primaryKey;type:int(11);not null"`
	MachineIP    string    `gorm:"column:machine_ip;type:varchar(20);not null"`
	State        string    `gorm:"column:state;type:varchar(20);not null"`
	StateDetail  *string   `gorm:"column:state_detail;type:varchar(20)"`
	Created      time.Time `gorm:"column:created;type:datetime;not null"`
	Modified     time.Time `gorm:"column:modified;type:datetime;not null"`
}

// TableName テーブル名を指定
func (TimeCard) TableName() string {
	return "time_card"
}
