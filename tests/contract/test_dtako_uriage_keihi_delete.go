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

func TestDTakoUriageKeihiService_Delete(t *testing.T) {
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
			SrchId:       "DELETE_TEST001",
			Datetime:     "2025-09-19T15:00:00Z",
			KeihiC:       3,
			Price:        6000.0,
			Km:           120.0,
			DtakoRowId:   "DTAKO005",
			DtakoRowIdR:  "DTAKO005R",
		},
	}
	_, _ = client.Create(ctx, setupData)

	// テストケース
	testCases := []struct {
		name         string
		request      *proto.DeleteDTakoUriageKeihiRequest
		wantErr      bool
		expectedCode codes.Code
	}{
		{
			name: "Valid delete",
			request: &proto.DeleteDTakoUriageKeihiRequest{
				SrchId:   "DELETE_TEST001",
				Datetime: "2025-09-19T15:00:00Z",
				KeihiC:   3,
			},
			wantErr: false,
		},
		{
			name: "Delete non-existent record",
			request: &proto.DeleteDTakoUriageKeihiRequest{
				SrchId:   "NONEXISTENT",
				Datetime: "2025-09-19T16:00:00Z",
				KeihiC:   999,
			},
			wantErr:      true,
			expectedCode: codes.NotFound,
		},
		{
			name: "Delete with invalid key",
			request: &proto.DeleteDTakoUriageKeihiRequest{
				SrchId:   "",
				Datetime: "",
				KeihiC:   0,
			},
			wantErr:      true,
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.Delete(ctx, tc.request)
			if (err != nil) != tc.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if tc.wantErr && tc.expectedCode != codes.OK {
				st, ok := status.FromError(err)
				if !ok || st.Code() != tc.expectedCode {
					t.Errorf("Expected error code %v, got %v", tc.expectedCode, st.Code())
				}
			}
			if !tc.wantErr {
				// 削除後にデータが存在しないことを確認
				getReq := &proto.GetDTakoUriageKeihiRequest{
					SrchId:   tc.request.SrchId,
					Datetime: tc.request.Datetime,
					KeihiC:   tc.request.KeihiC,
				}
				_, err := client.Get(ctx, getReq)
				if err == nil {
					t.Error("Expected record to be deleted but it still exists")
				}
			}
		})
	}
}