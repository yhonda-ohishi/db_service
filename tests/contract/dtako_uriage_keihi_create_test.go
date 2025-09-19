package contract

import (
	"context"
	"testing"
	"time"

	"github.com/yhonda-ohishi/db_service/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDTakoUriageKeihiService_Create(t *testing.T) {
	// gRPCサーバーへの接続
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewDTakoUriageKeihiServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// テストデータの作成
	testCases := []struct {
		name    string
		request *proto.CreateDTakoUriageKeihiRequest
		wantErr bool
	}{
		{
			name: "Valid creation",
			request: &proto.CreateDTakoUriageKeihiRequest{
				DtakoUriageKeihi: &proto.DTakoUriageKeihi{
					SrchId:      "TEST001",
					Datetime:    "2025-09-19T10:00:00Z",
					KeihiC:      1,
					Price:       1000.0,
					Km:          float64Ptr(50.5),
					DtakoRowId:  "DTAKO001",
					DtakoRowIdR: "DTAKO001R",
				},
			},
			wantErr: false,
		},
		{
			name: "Duplicate primary key",
			request: &proto.CreateDTakoUriageKeihiRequest{
				DtakoUriageKeihi: &proto.DTakoUriageKeihi{
					SrchId:      "TEST001",
					Datetime:    "2025-09-19T10:00:00Z",
					KeihiC:      1,
					Price:       2000.0,
					DtakoRowId:  "DTAKO002",
					DtakoRowIdR: "DTAKO002R",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing required fields",
			request: &proto.CreateDTakoUriageKeihiRequest{
				DtakoUriageKeihi: &proto.DTakoUriageKeihi{
					SrchId:   "TEST002",
					Datetime: "2025-09-19T11:00:00Z",
					// KeihiC missing
					Price: 1500.0,
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
			if !tc.wantErr && response == nil {
				t.Error("Expected non-nil response for successful creation")
			}
		})
	}
}
