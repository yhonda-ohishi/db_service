package models

import (
	"time"
)

// DTakoFerryRows フェリー運行データモデル
type DTakoFerryRows struct {
	// 主キー（自動インクリメント）
	ID int32 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 運行情報
	UnkoNo       string    `gorm:"column:運行NO;size:23;not null" json:"unko_no"`
	UnkoDate     time.Time `gorm:"column:運行日;type:date;not null" json:"unko_date"`
	YomitoriDate time.Time `gorm:"column:読取日;type:date;not null" json:"yomitori_date"`

	// 事業所情報
	JigyoshoCD   int32  `gorm:"column:事業所CD;not null" json:"jigyosho_cd"`
	JigyoshoName string `gorm:"column:事業所名;size:20;not null" json:"jigyosho_name"`

	// 車両情報
	SharyoCD   int32  `gorm:"column:車輌CD;not null" json:"sharyo_cd"`
	SharyoName string `gorm:"column:車輌名;size:20;not null" json:"sharyo_name"`

	// 乗務員情報
	JomuinCD1         int32  `gorm:"column:乗務員CD1;not null" json:"jomuin_cd1"`
	JomuinName1       string `gorm:"column:乗務員名１;size:20;not null" json:"jomuin_name1"`
	TaishoJomuinKbn   int32  `gorm:"column:対象乗務員区分;not null" json:"taisho_jomuin_kbn"`

	// 運行時間
	KaishiDatetime time.Time `gorm:"column:開始日時;not null" json:"kaishi_datetime"`
	ShuryoDatetime time.Time `gorm:"column:終了日時;not null" json:"shuryo_datetime"`

	// フェリー会社情報
	FerryCompanyCD   int32  `gorm:"column:フェリー会社CD;not null" json:"ferry_company_cd"`
	FerryCompanyName string `gorm:"column:フェリー会社名;size:20;not null" json:"ferry_company_name"`

	// 乗降場情報
	NoribaCD   int32  `gorm:"column:乗場CD;not null" json:"noriba_cd"`
	NoribaName string `gorm:"column:乗場名;size:20;not null" json:"noriba_name"`
	Bin        string `gorm:"column:便;size:10;not null" json:"bin"`
	OribaCD    int32  `gorm:"column:降場CD;not null" json:"oriba_cd"`
	OribaName  string `gorm:"column:降場名;size:20;not null" json:"oriba_name"`

	// 精算情報
	SeisanKbn     int32  `gorm:"column:精算区分;not null" json:"seisan_kbn"`
	SeisanKbnName string `gorm:"column:精算区分名;size:20;not null" json:"seisan_kbn_name"`

	// 料金情報
	HyojunRyokin   int32 `gorm:"column:標準料金;not null" json:"hyojun_ryokin"`
	KeiyakuRyokin  int32 `gorm:"column:契約料金;not null" json:"keiyaku_ryokin"`

	// 車種・距離情報
	KosoShashuKbn     int32  `gorm:"column:航送車種区分;not null" json:"koso_shashu_kbn"`
	KosoShashuKbnName string `gorm:"column:航送車種区分名;size:20;not null" json:"koso_shashu_kbn_name"`
	MinashiKyori      int32  `gorm:"column:見なし距離;not null" json:"minashi_kyori"`

	// 検索用フィールド
	FerrySrch *string `gorm:"column:ferry_srch;size:60" json:"ferry_srch,omitempty"`
}

// TableName テーブル名を指定
func (DTakoFerryRows) TableName() string {
	return "dtako_ferry_rows"
}

// Validate バリデーション
func (d *DTakoFerryRows) Validate() error {
	if d.UnkoNo == "" {
		return ErrInvalidUnkoNo
	}
	if d.UnkoDate.IsZero() {
		return ErrInvalidUnkoDate
	}
	if d.YomitoriDate.IsZero() {
		return ErrInvalidYomitoriDate
	}
	if d.JigyoshoCD <= 0 {
		return ErrInvalidJigyoshoCD
	}
	if d.JigyoshoName == "" {
		return ErrInvalidJigyoshoName
	}
	if d.SharyoCD <= 0 {
		return ErrInvalidSharyoCD
	}
	if d.SharyoName == "" {
		return ErrInvalidSharyoName
	}
	if d.HyojunRyokin < 0 {
		return ErrInvalidHyojunRyokin
	}
	if d.KeiyakuRyokin < 0 {
		return ErrInvalidKeiyakuRyokin
	}
	if d.MinashiKyori < 0 {
		return ErrInvalidMinashiKyori
	}
	return nil
}

// CalculateDiscountAmount 割引額を計算
func (d *DTakoFerryRows) CalculateDiscountAmount() int32 {
	if d.HyojunRyokin > d.KeiyakuRyokin {
		return d.HyojunRyokin - d.KeiyakuRyokin
	}
	return 0
}

// CalculateDiscountRate 割引率を計算
func (d *DTakoFerryRows) CalculateDiscountRate() float64 {
	if d.HyojunRyokin == 0 {
		return 0
	}
	discount := d.CalculateDiscountAmount()
	return float64(discount) / float64(d.HyojunRyokin) * 100
}

// GetOperationDuration 運行時間を取得
func (d *DTakoFerryRows) GetOperationDuration() time.Duration {
	return d.ShuryoDatetime.Sub(d.KaishiDatetime)
}