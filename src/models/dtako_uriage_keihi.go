package models

import (
	"time"
)

// DTakoUriageKeihi 経費精算データモデル
type DTakoUriageKeihi struct {
	// 複合主キー
	SrchID   string    `gorm:"column:srch_id;primaryKey;size:44" json:"srch_id"`
	Datetime time.Time `gorm:"column:datetime;primaryKey" json:"datetime"`
	KeihiC   int32     `gorm:"column:keihi_c;primaryKey" json:"keihi_c"`

	// 必須フィールド
	Price       float64 `gorm:"column:price;not null" json:"price"`
	DtakoRowID  string  `gorm:"column:dtako_row_id;size:24;not null" json:"dtako_row_id"`
	DtakoRowIDR string  `gorm:"column:dtako_row_id_r;size:23;not null" json:"dtako_row_id_r"`

	// オプションフィールド
	Km             *float64   `gorm:"column:km" json:"km,omitempty"`
	StartSrchID    *string    `gorm:"column:start_srch_id;size:44" json:"start_srch_id,omitempty"`
	StartSrchTime  *time.Time `gorm:"column:start_srch_time" json:"start_srch_time,omitempty"`
	StartSrchPlace *string    `gorm:"column:start_srch_place;size:50" json:"start_srch_place,omitempty"`
	StartSrchTokui *string    `gorm:"column:start_srch_tokui;size:9" json:"start_srch_tokui,omitempty"`
	EndSrchID      *string    `gorm:"column:end_srch_id;size:44" json:"end_srch_id,omitempty"`
	EndSrchTime    *time.Time `gorm:"column:end_srch_time" json:"end_srch_time,omitempty"`
	EndSrchPlace   *string    `gorm:"column:end_srch_place;size:50" json:"end_srch_place,omitempty"`
	Manual         *bool      `gorm:"column:manual" json:"manual,omitempty"`
}

// TableName テーブル名を指定
func (DTakoUriageKeihi) TableName() string {
	return "dtako_uriage_keihi"
}

// Validate バリデーション
func (d *DTakoUriageKeihi) Validate() error {
	if d.SrchID == "" {
		return ErrInvalidSrchID
	}
	if d.KeihiC < 0 {
		return ErrInvalidKeihiC
	}
	if d.Price < 0 {
		return ErrInvalidPrice
	}
	if d.Km != nil && *d.Km < 0 {
		return ErrInvalidKm
	}
	if d.DtakoRowID == "" {
		return ErrInvalidDtakoRowID
	}
	if d.DtakoRowIDR == "" {
		return ErrInvalidDtakoRowIDR
	}
	return nil
}

// GetCompositeKey 複合キーを取得
func (d *DTakoUriageKeihi) GetCompositeKey() (string, time.Time, int32) {
	return d.SrchID, d.Datetime, d.KeihiC
}

// SetManual 手動フラグを設定
func (d *DTakoUriageKeihi) SetManual(manual bool) {
	d.Manual = &manual
}

// IsManual 手動フラグの確認
func (d *DTakoUriageKeihi) IsManual() bool {
	return d.Manual != nil && *d.Manual
}
