package unit

import (
	"testing"
	"time"

	"github.com/yhonda-ohishi/db_service/src/models"
)

func TestETCMeisaiMapping_Validate(t *testing.T) {
	testTime := time.Now()

	testCases := []struct {
		name    string
		mapping *models.ETCMeisaiMapping
		wantErr bool
	}{
		{
			name: "Valid mapping",
			mapping: &models.ETCMeisaiMapping{
				ETCMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
				DTakoRowID:    "ROW123456789012345678901",
				CreatedAt:     testTime,
				UpdatedAt:     testTime,
				CreatedBy:     "test_user",
			},
			wantErr: false,
		},
		{
			name: "Valid mapping with notes",
			mapping: &models.ETCMeisaiMapping{
				ETCMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
				DTakoRowID:    "ROW123456789012345678901",
				CreatedAt:     testTime,
				UpdatedAt:     testTime,
				CreatedBy:     "test_user",
				Notes:         &[]string{"Test mapping"}[0],
			},
			wantErr: false,
		},
		{
			name: "Missing hash",
			mapping: &models.ETCMeisaiMapping{
				DTakoRowID: "ROW123456789012345678901",
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
				CreatedBy:  "test_user",
			},
			wantErr: true,
		},
		{
			name: "Empty hash",
			mapping: &models.ETCMeisaiMapping{
				ETCMeisaiHash: "",
				DTakoRowID:    "ROW123456789012345678901",
				CreatedAt:     testTime,
				UpdatedAt:     testTime,
				CreatedBy:     "test_user",
			},
			wantErr: true,
		},
		{
			name: "Missing dtako_row_id",
			mapping: &models.ETCMeisaiMapping{
				ETCMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
				CreatedAt:     testTime,
				UpdatedAt:     testTime,
				CreatedBy:     "test_user",
			},
			wantErr: true,
		},
		{
			name: "Empty dtako_row_id",
			mapping: &models.ETCMeisaiMapping{
				ETCMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
				DTakoRowID:    "",
				CreatedAt:     testTime,
				UpdatedAt:     testTime,
				CreatedBy:     "test_user",
			},
			wantErr: true,
		},
		{
			name: "Missing created_by",
			mapping: &models.ETCMeisaiMapping{
				ETCMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
				DTakoRowID:    "ROW123456789012345678901",
				CreatedAt:     testTime,
				UpdatedAt:     testTime,
			},
			wantErr: true,
		},
		{
			name: "Empty created_by",
			mapping: &models.ETCMeisaiMapping{
				ETCMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
				DTakoRowID:    "ROW123456789012345678901",
				CreatedAt:     testTime,
				UpdatedAt:     testTime,
				CreatedBy:     "",
			},
			wantErr: true,
		},
		{
			name: "Zero created_at",
			mapping: &models.ETCMeisaiMapping{
				ETCMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
				DTakoRowID:    "ROW123456789012345678901",
				CreatedAt:     time.Time{},
				UpdatedAt:     testTime,
				CreatedBy:     "test_user",
			},
			wantErr: true,
		},
		{
			name: "Zero updated_at",
			mapping: &models.ETCMeisaiMapping{
				ETCMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
				DTakoRowID:    "ROW123456789012345678901",
				CreatedAt:     testTime,
				UpdatedAt:     time.Time{},
				CreatedBy:     "test_user",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.mapping.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestETCMeisaiMapping_TableName(t *testing.T) {
	mapping := &models.ETCMeisaiMapping{}
	expected := "etc_meisai_mapping"
	if tableName := mapping.TableName(); tableName != expected {
		t.Errorf("Expected table name %s, got %s", expected, tableName)
	}
}

func TestETCMeisaiMapping_BeforeCreate(t *testing.T) {
	mapping := &models.ETCMeisaiMapping{
		ETCMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
		DTakoRowID:    "ROW123456789012345678901",
		CreatedBy:     "test_user",
	}

	beforeCreate := time.Now()
	mapping.BeforeCreate()
	afterCreate := time.Now()

	// タイムスタンプが設定されているか確認
	if mapping.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set by BeforeCreate")
	}
	if mapping.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set by BeforeCreate")
	}

	// 時間が適切な範囲内か確認
	if mapping.CreatedAt.Before(beforeCreate) || mapping.CreatedAt.After(afterCreate) {
		t.Error("CreatedAt timestamp is out of expected range")
	}
	if mapping.UpdatedAt.Before(beforeCreate) || mapping.UpdatedAt.After(afterCreate) {
		t.Error("UpdatedAt timestamp is out of expected range")
	}

	// CreatedAtとUpdatedAtが同じ時刻に設定されているか確認
	if !mapping.CreatedAt.Equal(mapping.UpdatedAt) {
		t.Error("CreatedAt and UpdatedAt should be equal in BeforeCreate")
	}
}

func TestETCMeisaiMapping_BeforeUpdate(t *testing.T) {
	originalTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	mapping := &models.ETCMeisaiMapping{
		ETCMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
		DTakoRowID:    "ROW123456789012345678901",
		CreatedAt:     originalTime,
		UpdatedAt:     originalTime,
		CreatedBy:     "test_user",
	}

	beforeUpdate := time.Now()
	mapping.BeforeUpdate()
	afterUpdate := time.Now()

	// CreatedAtが変更されていないか確認
	if !mapping.CreatedAt.Equal(originalTime) {
		t.Error("CreatedAt should not be modified by BeforeUpdate")
	}

	// UpdatedAtが更新されているか確認
	if mapping.UpdatedAt.Before(beforeUpdate) || mapping.UpdatedAt.After(afterUpdate) {
		t.Error("UpdatedAt timestamp is out of expected range")
	}

	// UpdatedAtがCreatedAtより新しいか確認
	if !mapping.UpdatedAt.After(mapping.CreatedAt) {
		t.Error("UpdatedAt should be after CreatedAt")
	}
}