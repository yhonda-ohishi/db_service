package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yhonda-ohishi/db_service/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// gRPCサーバーに接続
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fmt.Println("=== 既存サービステスト ===")

	// ETCMeisaiサービステスト
	etcClient := proto.NewETCMeisaiServiceClient(conn)
	etcListResp, err := etcClient.List(ctx, &proto.ListETCMeisaiRequest{
		Limit:  5,
		Offset: 0,
	})
	if err != nil {
		log.Printf("ETCMeisai List エラー: %v", err)
	} else {
		fmt.Printf("ETCMeisai総数: %d件\n", etcListResp.TotalCount)
		for i, item := range etcListResp.Items {
			fmt.Printf("  %d. ID:%d, %s -> %s, 料金:%d円\n", i+1, item.Id, item.IcFr, item.IcTo, item.Price)
		}
	}

	fmt.Println("\n=== 本番DBサービステスト ===")

	// DTakoRowsサービステスト（本番DB）
	rowsClient := proto.NewDTakoRowsServiceClient(conn)
	rowsListResp, err := rowsClient.List(ctx, &proto.ListDTakoRowsRequest{
		Limit:  3,
		Offset: 0,
	})
	if err != nil {
		log.Printf("DTakoRows List エラー: %v", err)
	} else {
		fmt.Printf("DTakoRows総数: %d件\n", rowsListResp.TotalCount)
		for i, item := range rowsListResp.Items {
			fmt.Printf("  %d. ID:%s, 運行NO:%s, 車輌CD:%d\n", i+1, item.Id, item.OperationNo, item.CarCode)
		}
	}

	// ETCNumサービステスト（本番DB）
	etcNumClient := proto.NewETCNumServiceClient(conn)
	etcNumListResp, err := etcNumClient.List(ctx, &proto.ListETCNumRequest{
		Limit:  3,
		Offset: 0,
	})
	if err != nil {
		log.Printf("ETCNum List エラー: %v", err)
	} else {
		fmt.Printf("ETCNum総数: %d件\n", etcNumListResp.TotalCount)
		for i, item := range etcNumListResp.Items {
			fmt.Printf("  %d. ETCカード:%s, 車輌ID:%s\n", i+1, item.EtcCardNum, item.CarId)
		}
	}

	fmt.Println("\n=== 全サービステスト完了 ===")
}