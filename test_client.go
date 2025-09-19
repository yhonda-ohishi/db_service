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

	// ETCMeisaiクライアントを作成
	etcClient := proto.NewETCMeisaiServiceClient(conn)

	// テストデータの作成
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Create テスト
	createReq := &proto.CreateETCMeisaiRequest{
		EtcMeisai: &proto.ETCMeisai{
			DateTo:     time.Now().Format(time.RFC3339),
			DateToDate: time.Now().Format("2006-01-02"),
			IcFr:       "東京IC",
			IcTo:       "名古屋IC",
			Price:      5000,
			Shashu:     1,
			EtcNum:     "1234-5678-9012-3456",
		},
	}

	resp, err := etcClient.Create(ctx, createReq)
	if err != nil {
		log.Printf("Create failed: %v", err)
	} else {
		fmt.Printf("Created ETC record with ID: %d\n", resp.EtcMeisai.Id)
	}

	// List テスト
	listReq := &proto.ListETCMeisaiRequest{
		Limit:  10,
		Offset: 0,
	}

	listResp, err := etcClient.List(ctx, listReq)
	if err != nil {
		log.Printf("List failed: %v", err)
	} else {
		fmt.Printf("Found %d ETC records\n", listResp.TotalCount)
		for i, item := range listResp.Items {
			fmt.Printf("%d. IC: %s -> %s, Price: %d\n", i+1, item.IcFr, item.IcTo, item.Price)
		}
	}
}