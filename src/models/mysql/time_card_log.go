package mysql

import "time"

// TimeCardLog タイムカードログテーブル
type TimeCardLog struct {
	Datetime    string    `gorm:"column:datetime;primaryKey;type:varchar(30);not null"` // RFC3339形式のタイムスタンプ
	ID          int       `gorm:"column:id;primaryKey;type:int(11);not null;default:0"` // ユーザーID（0はゲスト/不明）
	CardID      string    `gorm:"column:card_id;type:varchar(50);not null"`             // カードID（FeliCa UIDなど）
	MachineIP   string    `gorm:"column:machine_ip;type:varchar(100);not null"`         // マシンIP/Reader ID
	State       string    `gorm:"column:state;type:varchar(20);not null"`               // 状態: in/out
	StateDetail *string   `gorm:"column:state_detail;type:varchar(50)"`                 // 状態詳細（オプション）
	Created     time.Time `gorm:"column:created;type:datetime;not null"`                // 作成日時
	Modified    time.Time `gorm:"column:modified;type:datetime;not null"`               // 更新日時
}

func (TimeCardLog) TableName() string {
	return "timecard_logs"
}
