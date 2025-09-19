package contract

import (
	"context"
	"testing"
	"time"

	"github.com/yhonda-ohishi/db_service/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestDTakoUriageKeihiService_List(t *testing.T) {
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
	testData := []struct {
		srchId   string
		datetime string
		keihiC   int32
		price    float64
		rowId    string
	}{
		{"LIST_TEST001", "2025-09-19T10:00:00Z", 1, 1000.0, "DTAKO_LIST001"},
		{"LIST_TEST002", "2025-09-19T11:00:00Z", 2, 2000.0, "DTAKO_LIST001"},
		{"LIST_TEST003", "2025-09-19T12:00:00Z", 3, 3000.0, "DTAKO_LIST002"},
		{"LIST_TEST004", "2025-09-20T10:00:00Z", 1, 4000.0, "DTAKO_LIST002"},
		{"LIST_TEST005", "2025-09-21T10:00:00Z", 2, 5000.0, "DTAKO_LIST003"},
	}

	for _, td := range testData {
		req := &proto.CreateDTakoUriageKeihiRequest{
			DtakoUriageKeihi: &proto.DTakoUriageKeihi{
				SrchId:      td.srchId,
				Datetime:    td.datetime,
				KeihiC:      td.keihiC,
				Price:       td.price,
				DtakoRowId:  td.rowId,
				DtakoRowIdR: td.rowId + "R",
			},
		}
		_, _ = client.Create(ctx, req)
	}

	// テストケース
	testCases := []struct {
		name         string
		request      *proto.ListDTakoUriageKeihiRequest
		wantMinCount int32
		wantMaxCount int32
	}{
		{
			name: "List all with limit",
			request: &proto.ListDTakoUriageKeihiRequest{
				Limit:  10,
				Offset: 0,
			},
			wantMinCount: 5,
			wantMaxCount: 10,
		},
		{
			name: "List by dtako_row_id",
			request: &proto.ListDTakoUriageKeihiRequest{
				DtakoRowId: &[]string{"DTAKO_LIST001"}[0],
				Limit:      10,
				Offset:     0,
			},
			wantMinCount: 2,
			wantMaxCount: 2,
		},
		{
			name: "List with date range",
			request: &proto.ListDTakoUriageKeihiRequest{
				StartDate: &[]string{"2025-09-19T00:00:00Z"}[0],
				EndDate:   &[]string{"2025-09-19T23:59:59Z"}[0],
				Limit:     10,
				Offset:    0,
			},
			wantMinCount: 3,
			wantMaxCount: 3,
		},
		{
			name: "List with pagination",
			request: &proto.ListDTakoUriageKeihiRequest{
				Limit:  2,
				Offset: 2,
			},
			wantMinCount: 0,
			wantMaxCount: 2,
		},
		{
			name: "List with high offset",
			request: &proto.ListDTakoUriageKeihiRequest{
				Limit:  10,
				Offset: 100,
			},
			wantMinCount: 0,
			wantMaxCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := client.List(ctx, tc.request)
			if err != nil {
				t.Fatalf("List() error = %v", err)
			}

			itemCount := int32(len(response.Items))
			if itemCount < tc.wantMinCount || itemCount > tc.wantMaxCount {
				t.Errorf("List() returned %d items, want between %d and %d",
					itemCount, tc.wantMinCount, tc.wantMaxCount)
			}

			// ページネーションのテスト
			if tc.request.Limit > 0 && itemCount > tc.request.Limit {
				t.Errorf("List() returned more items than limit: got %d, limit %d",
					itemCount, tc.request.Limit)
			}

			// 総件数の確認
			if response.TotalCount < itemCount {
				t.Errorf("TotalCount (%d) is less than returned items (%d)",
					response.TotalCount, itemCount)
			}
		})
	}
}
