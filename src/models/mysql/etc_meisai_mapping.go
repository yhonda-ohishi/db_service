package mysql

import "time"

// ETCMeisaiMapping ETC明細とDTakoRowsの関連付けテーブル
type ETCMeisaiMapping struct {
	// 主キー（自動インクリメント）
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// ETC明細のハッシュ値
	ETCMeisaiHash string `gorm:"column:etc_meisai_hash;size:64;not null;index" json:"etc_meisai_hash"`

	// 運行データID
	DTakoRowID string `gorm:"column:dtako_row_id;size:24;not null;index" json:"dtako_row_id"`

	// マッピング作成日時
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`

	// マッピング更新日時
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`

	// マッピング作成者/システム
	CreatedBy string `gorm:"column:created_by;size:50;not null" json:"created_by"`

	// 備考
	Notes *string `gorm:"column:notes;size:200" json:"notes,omitempty"`
}

// TableName テーブル名を指定
func (ETCMeisaiMapping) TableName() string {
	return "etc_meisai_mapping"
}

// Validate バリデーション
func (m *ETCMeisaiMapping) Validate() error {
	if m.ETCMeisaiHash == "" {
		return ErrInvalidHash
	}
	if m.DTakoRowID == "" {
		return ErrInvalidDtakoRowID
	}
	if m.CreatedBy == "" {
		return ErrInvalidCreatedBy
	}
	if m.CreatedAt.IsZero() {
		return ErrInvalidCreatedAt
	}
	if m.UpdatedAt.IsZero() {
		return ErrInvalidUpdatedAt
	}
	return nil
}

// BeforeCreate GORM作成前フック
func (m *ETCMeisaiMapping) BeforeCreate() {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
}

// BeforeUpdate GORM更新前フック
func (m *ETCMeisaiMapping) BeforeUpdate() {
	m.UpdatedAt = time.Now()
}

// SetDefaults デフォルト値を設定
func (m *ETCMeisaiMapping) SetDefaults() {
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	m.UpdatedAt = now
	if m.CreatedBy == "" {
		m.CreatedBy = "system"
	}
}