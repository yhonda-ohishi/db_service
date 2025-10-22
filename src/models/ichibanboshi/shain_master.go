package ichibanboshi

import "time"

// ShainMaster 社員マスタテーブルのモデル（SQL Server）
type ShainMaster struct {
	ShainC          string     `gorm:"column:社員C;primaryKey;size:4" json:"shain_c"`
	ShainN          *string    `gorm:"column:社員N;size:16" json:"shain_n,omitempty"`
	ShainR          *string    `gorm:"column:社員R;size:8" json:"shain_r,omitempty"`
	ShainF          *string    `gorm:"column:社員F;size:30" json:"shain_f,omitempty"`
	YubinBango      *string    `gorm:"column:郵便番号;size:8" json:"yubin_bango,omitempty"`
	Jusho1          *string    `gorm:"column:住所1;size:40" json:"jusho1,omitempty"`
	Jusho2          *string    `gorm:"column:住所2;size:40" json:"jusho2,omitempty"`
	DenwaBango      *string    `gorm:"column:電話番号;size:13" json:"denwa_bango,omitempty"`
	KeitaiBango     *string    `gorm:"column:携帯番号;size:13" json:"keitai_bango,omitempty"`
	ShainK          string     `gorm:"column:社員K;size:1" json:"shain_k"`
	Seibetsu        string     `gorm:"column:性別;size:1" json:"seibetsu"`
	Ketsuekigata    string     `gorm:"column:血液型;size:1" json:"ketsuekigata"`
	Seinengappi     *time.Time `gorm:"column:生年月日" json:"seinengappi,omitempty"`
	NyushaNengappi  *time.Time `gorm:"column:入社年月日" json:"nyusha_nengappi,omitempty"`
	TaishokuNengappi *time.Time `gorm:"column:退職年月日" json:"taishoku_nengappi,omitempty"`
	DaiBunrui1      string     `gorm:"column:大分類1;size:1" json:"dai_bunrui1"`
	ChuBunrui1      string     `gorm:"column:中分類1;size:2" json:"chu_bunrui1"`
	ShoBunrui1      string     `gorm:"column:小分類1;size:2" json:"sho_bunrui1"`
	DaiBunrui2      string     `gorm:"column:大分類2;size:1" json:"dai_bunrui2"`
	ChuBunrui2      string     `gorm:"column:中分類2;size:2" json:"chu_bunrui2"`
	ShoBunrui2      string     `gorm:"column:小分類2;size:2" json:"sho_bunrui2"`
	KodanPlate      *string    `gorm:"column:公団ﾌﾟﾚｰﾄ;size:4" json:"kodan_plate,omitempty"`
	UriageMokuhyogaku int      `gorm:"column:売上目標額" json:"uriage_mokuhyogaku"`
	UntenMenkyoK    string     `gorm:"column:運転免許K;size:1" json:"unten_menkyo_k"`
	MenkyoshoBango  *string    `gorm:"column:免許証番号;size:12" json:"menkyosho_bango,omitempty"`
	JikaiKoshinbi   *time.Time `gorm:"column:次回更新日" json:"jikai_koshinbi,omitempty"`
	JishaYoshaK     string     `gorm:"column:自車傭車K;size:1" json:"jisha_yosha_k"`
	KeisanK         string     `gorm:"column:計算K;size:1" json:"keisan_k"`
	ShiharaiRitsu   float64    `gorm:"column:支払率;type:decimal" json:"shiharai_ritsu"`
	HasuK           string     `gorm:"column:端数K;size:1" json:"hasu_k"`
	BumonC          string     `gorm:"column:部門C;size:3" json:"bumon_c"`
	UnchinPatternC  string     `gorm:"column:運賃ﾊﾟﾀｰﾝC;size:3" json:"unchin_pattern_c"`
	Kiji1           *string    `gorm:"column:記事1;size:40" json:"kiji1,omitempty"`
	Kiji2           *string    `gorm:"column:記事2;size:40" json:"kiji2,omitempty"`
	Kiji3           *string    `gorm:"column:記事3;size:40" json:"kiji3,omitempty"`
	Kiji4           *string    `gorm:"column:記事4;size:40" json:"kiji4,omitempty"`
	Kiji5           *string    `gorm:"column:記事5;size:40" json:"kiji5,omitempty"`
	Yobi1           *string    `gorm:"column:予備1;size:64" json:"yobi1,omitempty"`
	Yobi2           *string    `gorm:"column:予備2;size:64" json:"yobi2,omitempty"`
	Yobi3           *string    `gorm:"column:予備3;size:64" json:"yobi3,omitempty"`
	Yobi4           *string    `gorm:"column:予備4;size:64" json:"yobi4,omitempty"`
	Yobi5           *string    `gorm:"column:予備5;size:64" json:"yobi5,omitempty"`
	KinmuTaikeiC    string     `gorm:"column:勤務体系C;size:2" json:"kinmu_taikei_c"`
}

// TableName テーブル名を指定
func (ShainMaster) TableName() string {
	return "社員ﾏｽﾀ"
}
