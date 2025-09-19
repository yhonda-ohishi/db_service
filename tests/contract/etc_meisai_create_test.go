package contract

import (
	"context"
	"testing"
	"time"

	"github.com/yhonda-ohishi/db_service/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestETCMeisaiService_Create(t *testing.T) {
	// gRPCサーバーへの接続
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewETCMeisaiServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// テストケース
	testCases := []struct {
		name    string
		request *proto.CreateETCMeisaiRequest
		wantErr bool
	}{
		{
			name: "Valid creation",
			request: &proto.CreateETCMeisaiRequest{
				EtcMeisai: &proto.ETCMeisai{
					DateTo:     "2025-09-19T10:00:00Z",
					DateToDate: "2025-09-19",
					IcFr:       "東京IC",
					IcTo:       "横浜IC",
					Price:      1500,
					Shashu:     1,
					EtcNum:     "1234567890123456",
					Hash:       "test_hash_123",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing required fields",
			request: &proto.CreateETCMeisaiRequest{
				EtcMeisai: &proto.ETCMeisai{
					DateTo: "2025-09-19T10:00:00Z",
					// IcFr, IcTo missing
					Price: 1500,
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid price",
			request: &proto.CreateETCMeisaiRequest{
				EtcMeisai: &proto.ETCMeisai{
					DateTo:     "2025-09-19T10:00:00Z",
					DateToDate: "2025-09-19",
					IcFr:       "東京IC",
					IcTo:       "横浜IC",
					Price:      -100, // 無効な料金
					Shashu:     1,
					EtcNum:     "1234567890123456",
					Hash:       "test_hash_invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := client.Create(ctx, tc.request)
			if (err != nil) != tc.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr {
				if response == nil || response.EtcMeisai == nil {
					t.Error("Expected non-nil response for successful creation")
				} else if response.EtcMeisai.Id == 0 {
					t.Error("Expected auto-generated ID to be non-zero")
				}
			}
		})
	}
}
