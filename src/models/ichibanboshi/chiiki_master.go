package ichibanboshi

// ChiikiMaster 地域マスタテーブルのモデル（SQL Server）
type ChiikiMaster struct {
	ChiikiC string  `gorm:"column:地域C;primaryKey;size:6" json:"chiiki_c"`
	ChiikiN *string `gorm:"column:地域N;size:32" json:"chiiki_n,omitempty"`
	ChiikiR *string `gorm:"column:地域R;size:16" json:"chiiki_r,omitempty"`
	ChiikiF *string `gorm:"column:地域F;size:12" json:"chiiki_f,omitempty"`
}

// TableName テーブル名を指定
func (ChiikiMaster) TableName() string {
	return "地域ﾏｽﾀ"
}
