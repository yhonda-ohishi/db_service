package contract

import (
	"context"
	"testing"
	"time"

	"github.com/yhonda-ohishi/db_service/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestDTakoUriageKeihiService_Update(t *testing.T) {
	// gRPCサーバーへの接続
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewDTakoUriageKeihiServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// テストデータの作成（事前準備）
	setupData := &proto.CreateDTakoUriageKeihiRequest{
		DtakoUriageKeihi: &proto.DTakoUriageKeihi{
			SrchId:      "UPDATE_TEST001",
			Datetime:    "2025-09-19T13:00:00Z",
			KeihiC:      2,
			Price:       4000.0,
			Km:          float64Ptr(100.0),
			DtakoRowId:  "DTAKO004",
			DtakoRowIdR: "DTAKO004R",
		},
	}
	_, _ = client.Create(ctx, setupData)

	// テストケース
	testCases := []struct {
		name         string
		request      *proto.UpdateDTakoUriageKeihiRequest
		wantErr      bool
		expectedCode codes.Code
	}{
		{
			name: "Valid update",
			request: &proto.UpdateDTakoUriageKeihiRequest{
				DtakoUriageKeihi: &proto.DTakoUriageKeihi{
					SrchId:      "UPDATE_TEST001",
					Datetime:    "2025-09-19T13:00:00Z",
					KeihiC:      2,
					Price:       4500.0,            // 更新
					Km:          float64Ptr(110.0), // 更新
					DtakoRowId:  "DTAKO004",
					DtakoRowIdR: "DTAKO004R",
					Manual:      &[]bool{true}[0], // 手動フラグ設定
				},
			},
			wantErr: false,
		},
		{
			name: "Update non-existent record",
			request: &proto.UpdateDTakoUriageKeihiRequest{
				DtakoUriageKeihi: &proto.DTakoUriageKeihi{
					SrchId:      "NONEXISTENT",
					Datetime:    "2025-09-19T14:00:00Z",
					KeihiC:      999,
					Price:       5000.0,
					DtakoRowId:  "DTAKO999",
					DtakoRowIdR: "DTAKO999R",
				},
			},
			wantErr:      true,
			expectedCode: codes.NotFound,
		},
		{
			name: "Invalid update data",
			request: &proto.UpdateDTakoUriageKeihiRequest{
				DtakoUriageKeihi: &proto.DTakoUriageKeihi{
					SrchId:   "UPDATE_TEST001",
					Datetime: "2025-09-19T13:00:00Z",
					KeihiC:   2,
					Price:    -1000.0, // 無効な金額
				},
			},
			wantErr:      true,
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.Update(ctx, tc.request)
			if (err != nil) != tc.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if tc.wantErr && tc.expectedCode != codes.OK {
				st, ok := status.FromError(err)
				if !ok || st.Code() != tc.expectedCode {
					t.Errorf("Expected error code %v, got %v", tc.expectedCode, st.Code())
				}
			}
			if !tc.wantErr {
				// 更新後のデータを確認
				getReq := &proto.GetDTakoUriageKeihiRequest{
					SrchId:   tc.request.DtakoUriageKeihi.SrchId,
					Datetime: tc.request.DtakoUriageKeihi.Datetime,
					KeihiC:   tc.request.DtakoUriageKeihi.KeihiC,
				}
				getResp, _ := client.Get(ctx, getReq)
				if getResp != nil && getResp.DtakoUriageKeihi.Price != tc.request.DtakoUriageKeihi.Price {
					t.Errorf("Price not updated: got %v, want %v",
						getResp.DtakoUriageKeihi.Price,
						tc.request.DtakoUriageKeihi.Price)
				}
			}
		})
	}
}
