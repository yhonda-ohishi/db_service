package models

import (
	"time"
)

// ETCMeisai ETC明細データモデル
type ETCMeisai struct {
	// 主キー（自動インクリメント）
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// 日時フィールド
	DateFr     *time.Time `gorm:"column:date_fr" json:"date_fr,omitempty"`
	DateTo     time.Time  `gorm:"column:date_to;not null" json:"date_to"`
	DateToDate time.Time  `gorm:"column:date_to_date;type:date;not null" json:"date_to_date"`

	// IC情報
	IcFr string `gorm:"column:IC_fr;size:30;not null" json:"ic_fr"`
	IcTo string `gorm:"column:IC_to;size:30;not null" json:"ic_to"`

	// 料金情報
	PriceBf  *int32 `gorm:"column:price_bf" json:"price_bf,omitempty"`
	Descount *int32 `gorm:"column:descount" json:"descount,omitempty"`
	Price    int32  `gorm:"column:price;not null" json:"price"`

	// 車両情報
	Shashu    int32  `gorm:"column:shashu;not null" json:"shashu"`
	CarIDNum  *int32 `gorm:"column:car_id_num" json:"car_id_num,omitempty"`
	EtcNum    string `gorm:"column:etc_num;size:20;not null" json:"etc_num"`
	Detail    *string `gorm:"column:detail;size:40" json:"detail,omitempty"`

	// 関連情報
	DtakoRowID *string `gorm:"column:dtako_row_id;size:24" json:"dtako_row_id,omitempty"`
}

// TableName テーブル名を指定
func (ETCMeisai) TableName() string {
	return "etc_meisai"
}

// Validate バリデーション
func (e *ETCMeisai) Validate() error {
	if e.DateTo.IsZero() {
		return ErrInvalidDateTo
	}
	if e.DateToDate.IsZero() {
		return ErrInvalidDateToDate
	}
	if e.IcFr == "" {
		return ErrInvalidIcFr
	}
	if e.IcTo == "" {
		return ErrInvalidIcTo
	}
	if e.Price < 0 {
		return ErrInvalidPrice
	}
	if e.PriceBf != nil && *e.PriceBf < 0 {
		return ErrInvalidPriceBf
	}
	if e.Descount != nil && *e.Descount < 0 {
		return ErrInvalidDescount
	}
	if e.Shashu <= 0 {
		return ErrInvalidShashu
	}
	if e.EtcNum == "" {
		return ErrInvalidEtcNum
	}
	return nil
}

// CalculateDiscountRate 割引率を計算
func (e *ETCMeisai) CalculateDiscountRate() float64 {
	if e.PriceBf == nil || *e.PriceBf == 0 {
		return 0
	}
	if e.Descount == nil || *e.Descount == 0 {
		return 0
	}
	return float64(*e.Descount) / float64(*e.PriceBf) * 100
}

// GetFinalPrice 最終価格を取得
func (e *ETCMeisai) GetFinalPrice() int32 {
	return e.Price
}

// HasDtakoRowID 運行NOの有無を確認
func (e *ETCMeisai) HasDtakoRowID() bool {
	return e.DtakoRowID != nil && *e.DtakoRowID != ""
}