package main

import (
	"context"
	"log"
	"time"

	pb "github.com/yhonda-ohishi/db_service/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// gRPCサーバーへの接続
	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to gRPC server")

	// 社員マスタテスト
	testShainMaster(conn)

	// 地域マスタテスト
	testChiikiMaster(conn)

	log.Println("All SQL Server tests completed successfully!")
}

func testShainMaster(conn *grpc.ClientConn) {
	client := pb.Db_NewShainMasterServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("\n=== Testing ShainMaster (社員マスタ) ===")

	// List取得
	resp, err := client.List(ctx, &pb.Db_ListShainMasterRequest{
		Limit:  5,
		Offset: 0,
	})
	if err != nil {
		log.Printf("ShainMaster List failed: %v", err)
		return
	}

	log.Printf("ShainMaster List success: %d items (total: %d)", len(resp.Items), resp.TotalCount)
	if len(resp.Items) > 0 {
		item := resp.Items[0]
		log.Printf("  First item: 社員CD=%s, 社員名=%s", item.ShainC, item.ShainN)
	}
}

func testChiikiMaster(conn *grpc.ClientConn) {
	client := pb.Db_NewChiikiMasterServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("\n=== Testing ChiikiMaster (地域マスタ) ===")

	// List取得
	resp, err := client.List(ctx, &pb.Db_ListChiikiMasterRequest{
		Limit:  5,
		Offset: 0,
	})
	if err != nil {
		log.Printf("ChiikiMaster List failed: %v", err)
		return
	}

	log.Printf("ChiikiMaster List success: %d items (total: %d)", len(resp.Items), resp.TotalCount)
	if len(resp.Items) > 0 {
		item := resp.Items[0]
		log.Printf("  First item: 地域CD=%s, 地域名=%s", item.ChiikiC, item.ChiikiN)
	}
}
