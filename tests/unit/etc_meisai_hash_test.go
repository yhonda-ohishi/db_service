package unit

import (
	"testing"
	"time"

	"github.com/yhonda-ohishi/db_service/src/models"
)

func TestETCMeisai_GenerateHash(t *testing.T) {
	testTime := time.Date(2025, 9, 19, 10, 0, 0, 0, time.UTC)
	testDate := time.Date(2025, 9, 19, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name     string
		meisai   *models.ETCMeisai
		expected string
	}{
		{
			name: "Basic hash generation",
			meisai: &models.ETCMeisai{
				DateTo:     testTime,
				DateToDate: testDate,
				IcFr:       "東京IC",
				IcTo:       "横浜IC",
				Price:      1500,
				Shashu:     1,
				EtcNum:     "1234567890123456",
			},
		},
		{
			name: "Hash with optional fields",
			meisai: &models.ETCMeisai{
				DateFr:     &testTime,
				DateTo:     testTime,
				DateToDate: testDate,
				IcFr:       "東京IC",
				IcTo:       "横浜IC",
				PriceBf:    &[]int32{2000}[0],
				Descount:   &[]int32{500}[0],
				Price:      1500,
				Shashu:     1,
				CarIDNum:   &[]int32{1234}[0],
				EtcNum:     "1234567890123456",
				Detail:     &[]string{"高速道路"}[0],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash1 := tc.meisai.GenerateHash()
			hash2 := tc.meisai.GenerateHash()

			// 同じデータからは同じハッシュが生成される
			if hash1 != hash2 {
				t.Errorf("Hash generation is not consistent: %s != %s", hash1, hash2)
			}

			// ハッシュが空でない
			if hash1 == "" {
				t.Error("Generated hash should not be empty")
			}

			// SHA256ハッシュの長さは64文字
			if len(hash1) != 64 {
				t.Errorf("Expected hash length 64, got %d", len(hash1))
			}

			// SetHashメソッドのテスト
			tc.meisai.SetHash()
			if tc.meisai.Hash != hash1 {
				t.Errorf("SetHash() didn't set the correct hash: expected %s, got %s", hash1, tc.meisai.Hash)
			}
		})
	}
}

func TestETCMeisai_HashUniqueness(t *testing.T) {
	testTime := time.Date(2025, 9, 19, 10, 0, 0, 0, time.UTC)
	testDate := time.Date(2025, 9, 19, 0, 0, 0, 0, time.UTC)

	meisai1 := &models.ETCMeisai{
		DateTo:     testTime,
		DateToDate: testDate,
		IcFr:       "東京IC",
		IcTo:       "横浜IC",
		Price:      1500,
		Shashu:     1,
		EtcNum:     "1234567890123456",
	}

	meisai2 := &models.ETCMeisai{
		DateTo:     testTime,
		DateToDate: testDate,
		IcFr:       "東京IC",
		IcTo:       "大阪IC", // 異なるIC
		Price:      1500,
		Shashu:     1,
		EtcNum:     "1234567890123456",
	}

	hash1 := meisai1.GenerateHash()
	hash2 := meisai2.GenerateHash()

	// 異なるデータからは異なるハッシュが生成される
	if hash1 == hash2 {
		t.Error("Different data should generate different hashes")
	}
}

func TestETCMeisai_ValidateWithHash(t *testing.T) {
	testTime := time.Date(2025, 9, 19, 10, 0, 0, 0, time.UTC)
	testDate := time.Date(2025, 9, 19, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name    string
		meisai  *models.ETCMeisai
		wantErr bool
	}{
		{
			name: "Valid with hash",
			meisai: &models.ETCMeisai{
				DateTo:     testTime,
				DateToDate: testDate,
				IcFr:       "東京IC",
				IcTo:       "横浜IC",
				Price:      1500,
				Shashu:     1,
				EtcNum:     "1234567890123456",
				Hash:       "valid_hash",
			},
			wantErr: false,
		},
		{
			name: "Missing hash",
			meisai: &models.ETCMeisai{
				DateTo:     testTime,
				DateToDate: testDate,
				IcFr:       "東京IC",
				IcTo:       "横浜IC",
				Price:      1500,
				Shashu:     1,
				EtcNum:     "1234567890123456",
				Hash:       "",
			},
			wantErr: true,
		},
		{
			name: "Invalid price with hash",
			meisai: &models.ETCMeisai{
				DateTo:     testTime,
				DateToDate: testDate,
				IcFr:       "東京IC",
				IcTo:       "横浜IC",
				Price:      -100,
				Shashu:     1,
				EtcNum:     "1234567890123456",
				Hash:       "valid_hash",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.meisai.Validate()
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}