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

func TestDTakoUriageKeihiService_Get(t *testing.T) {
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
			SrchId:      "GET_TEST001",
			Datetime:    "2025-09-19T12:00:00Z",
			KeihiC:      1,
			Price:       3000.0,
			Km:          float64Ptr(75.0),
			DtakoRowId:  "DTAKO003",
			DtakoRowIdR: "DTAKO003R",
		},
	}
	_, _ = client.Create(ctx, setupData)

	// テストケース
	testCases := []struct {
		name         string
		request      *proto.GetDTakoUriageKeihiRequest
		wantErr      bool
		expectedCode codes.Code
	}{
		{
			name: "Valid get",
			request: &proto.GetDTakoUriageKeihiRequest{
				SrchId:   "GET_TEST001",
				Datetime: "2025-09-19T12:00:00Z",
				KeihiC:   1,
			},
			wantErr: false,
		},
		{
			name: "Not found",
			request: &proto.GetDTakoUriageKeihiRequest{
				SrchId:   "NONEXISTENT",
				Datetime: "2025-09-19T12:00:00Z",
				KeihiC:   999,
			},
			wantErr:      true,
			expectedCode: codes.NotFound,
		},
		{
			name: "Invalid primary key",
			request: &proto.GetDTakoUriageKeihiRequest{
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
			response, err := client.Get(ctx, tc.request)
			if (err != nil) != tc.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if tc.wantErr && tc.expectedCode != codes.OK {
				st, ok := status.FromError(err)
				if !ok || st.Code() != tc.expectedCode {
					t.Errorf("Expected error code %v, got %v", tc.expectedCode, st.Code())
				}
			}
			if !tc.wantErr && response.DtakoUriageKeihi == nil {
				t.Error("Expected non-nil data for successful get")
			}
		})
	}
}
