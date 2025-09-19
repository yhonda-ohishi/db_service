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
	fmt.Println("本番DBサービスは現在無効化されています（プロトコルバッファー定義が必要）")

	fmt.Println("\n=== 全サービステスト完了 ===")
}
