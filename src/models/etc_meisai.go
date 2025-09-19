package models

import (
	"crypto/sha256"
	"fmt"
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
	Shashu   int32   `gorm:"column:shashu;not null" json:"shashu"`
	CarIDNum *int32  `gorm:"column:car_id_num" json:"car_id_num,omitempty"`
	EtcNum   string  `gorm:"column:etc_num;size:20;not null" json:"etc_num"`
	Detail   *string `gorm:"column:detail;size:40" json:"detail,omitempty"`

	// ハッシュ値（データ整合性確認用）
	Hash string `gorm:"column:hash;size:64;not null" json:"hash"`
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
	if e.Hash == "" {
		return ErrInvalidHash
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

// GenerateHash データからハッシュ値を生成（主キー以外の全フィールド）
func (e *ETCMeisai) GenerateHash() string {
	// 主キー以外の全フィールドを組み合わせてハッシュを生成
	dateFrStr := ""
	if e.DateFr != nil {
		dateFrStr = e.DateFr.Format("2006-01-02 15:04:05")
	}

	priceBfStr := ""
	if e.PriceBf != nil {
		priceBfStr = fmt.Sprintf("%d", *e.PriceBf)
	}

	descountStr := ""
	if e.Descount != nil {
		descountStr = fmt.Sprintf("%d", *e.Descount)
	}

	carIDNumStr := ""
	if e.CarIDNum != nil {
		carIDNumStr = fmt.Sprintf("%d", *e.CarIDNum)
	}

	detailStr := ""
	if e.Detail != nil {
		detailStr = *e.Detail
	}

	data := fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%d|%d|%s|%s|%s",
		dateFrStr,
		e.DateTo.Format("2006-01-02 15:04:05"),
		e.DateToDate.Format("2006-01-02"),
		e.IcFr,
		e.IcTo,
		priceBfStr,
		descountStr,
		e.Price,
		e.Shashu,
		carIDNumStr,
		e.EtcNum,
		detailStr)

	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// SetHash ハッシュ値を設定
func (e *ETCMeisai) SetHash() {
	e.Hash = e.GenerateHash()
}
