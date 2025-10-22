package ichibanboshi

// ChikuMaster 地区マスタテーブルのモデル（SQL Server）
type ChikuMaster struct {
	ChikuC         string  `gorm:"column:地区C;primaryKey;size:6" json:"chiku_c"`
	ChikuN         *string `gorm:"column:地区N;size:40" json:"chiku_n,omitempty"`
	ChikuR         *string `gorm:"column:地区R;size:20" json:"chiku_r,omitempty"`
	ChikuF         *string `gorm:"column:地区F;size:12" json:"chiku_f,omitempty"`
	ChiikiC        string  `gorm:"column:地域C;size:6" json:"chiiki_c"`
	YubinBango     *string `gorm:"column:郵便番号;size:8" json:"yubin_bango,omitempty"`
	Jusho1         *string `gorm:"column:住所1;size:40" json:"jusho1,omitempty"`
	Jusho2         *string `gorm:"column:住所2;size:40" json:"jusho2,omitempty"`
	DenwaBango     *string `gorm:"column:電話番号;size:13" json:"denwa_bango,omitempty"`
	FAXBango       *string `gorm:"column:FAX番号;size:13" json:"fax_bango,omitempty"`
	Tantosha       *string `gorm:"column:担当者;size:16" json:"tantosha,omitempty"`
	Yobi1          *string `gorm:"column:予備1;size:64" json:"yobi1,omitempty"`
	Yobi2          *string `gorm:"column:予備2;size:64" json:"yobi2,omitempty"`
	Yobi3          *string `gorm:"column:予備3;size:64" json:"yobi3,omitempty"`
	Yobi4          *string `gorm:"column:予備4;size:64" json:"yobi4,omitempty"`
	Yobi5          *string `gorm:"column:予備5;size:64" json:"yobi5,omitempty"`
	DGRTokuisakiC string  `gorm:"column:DGR得意先C;size:6" json:"dgr_tokuisaki_c"`
	DGRTokuisakiH string  `gorm:"column:DGR得意先H;size:3" json:"dgr_tokuisaki_h"`
	DGRHinmeiC     string  `gorm:"column:DGR品名C;size:4" json:"dgr_hinmei_c"`
	DGRHinmeiH     string  `gorm:"column:DGR品名H;size:2" json:"dgr_hinmei_h"`
}

// TableName テーブル名を指定
func (ChikuMaster) TableName() string {
	return "地区ﾏｽﾀ"
}
