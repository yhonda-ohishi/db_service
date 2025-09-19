package contract

import (
	"context"
	"testing"
	"time"

	"github.com/yhonda-ohishi/db_service/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestETCMeisaiMappingService_Create(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewETCMeisaiMappingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testCases := []struct {
		name    string
		request *proto.CreateETCMeisaiMappingRequest
		wantErr bool
	}{
		{
			name: "Valid mapping creation",
			request: &proto.CreateETCMeisaiMappingRequest{
				EtcMeisaiMapping: &proto.ETCMeisaiMapping{
					EtcMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
					DtakoRowId:    "ROW123456789012345678901",
					CreatedBy:     "test_user",
					Notes:         &[]string{"Test mapping"}[0],
				},
			},
			wantErr: false,
		},
		{
			name: "Missing hash",
			request: &proto.CreateETCMeisaiMappingRequest{
				EtcMeisaiMapping: &proto.ETCMeisaiMapping{
					DtakoRowId: "ROW123456789012345678901",
					CreatedBy:  "test_user",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing dtako_row_id",
			request: &proto.CreateETCMeisaiMappingRequest{
				EtcMeisaiMapping: &proto.ETCMeisaiMapping{
					EtcMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
					CreatedBy:     "test_user",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing created_by",
			request: &proto.CreateETCMeisaiMappingRequest{
				EtcMeisaiMapping: &proto.ETCMeisaiMapping{
					EtcMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
					DtakoRowId:    "ROW123456789012345678901",
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
				if response == nil || response.EtcMeisaiMapping == nil {
					t.Error("Expected non-nil response for successful creation")
				} else if response.EtcMeisaiMapping.Id == 0 {
					t.Error("Expected auto-generated ID to be non-zero")
				}
			}
		})
	}
}

func TestETCMeisaiMappingService_GetDTakoRowIDByHash(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewETCMeisaiMappingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testCases := []struct {
		name    string
		request *proto.GetDTakoRowIDByHashRequest
		wantErr bool
	}{
		{
			name: "Valid hash lookup",
			request: &proto.GetDTakoRowIDByHashRequest{
				EtcMeisaiHash: "a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456",
			},
			wantErr: false,
		},
		{
			name: "Empty hash",
			request: &proto.GetDTakoRowIDByHashRequest{
				EtcMeisaiHash: "",
			},
			wantErr: true,
		},
		{
			name: "Non-existent hash",
			request: &proto.GetDTakoRowIDByHashRequest{
				EtcMeisaiHash: "nonexistent_hash_123456789012345678901234567890123456789012345",
			},
			wantErr: false, // 空のリストが返される
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := client.GetDTakoRowIDByHash(ctx, tc.request)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetDTakoRowIDByHash() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr && response == nil {
				t.Error("Expected non-nil response for valid hash lookup")
			}
		})
	}
}

func TestETCMeisaiMappingService_List(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := proto.NewETCMeisaiMappingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testCases := []struct {
		name    string
		request *proto.ListETCMeisaiMappingRequest
		wantErr bool
	}{
		{
			name: "Valid list request",
			request: &proto.ListETCMeisaiMappingRequest{
				Limit:  10,
				Offset: 0,
			},
			wantErr: false,
		},
		{
			name: "Filter by hash",
			request: &proto.ListETCMeisaiMappingRequest{
				EtcMeisaiHash: &[]string{"a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456"}[0],
				Limit:         10,
				Offset:        0,
			},
			wantErr: false,
		},
		{
			name: "Filter by dtako_row_id",
			request: &proto.ListETCMeisaiMappingRequest{
				DtakoRowId: &[]string{"ROW123456789012345678901"}[0],
				Limit:      10,
				Offset:     0,
			},
			wantErr: false,
		},
		{
			name: "Invalid limit",
			request: &proto.ListETCMeisaiMappingRequest{
				Limit:  -1,
				Offset: 0,
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := client.List(ctx, tc.request)
			if (err != nil) != tc.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !tc.wantErr {
				if response == nil {
					t.Error("Expected non-nil response for successful list")
				} else if response.TotalCount < 0 {
					t.Error("Expected non-negative total count")
				}
			}
		})
	}
}