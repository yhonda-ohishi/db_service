package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/yhonda-ohishi/db_service/src/config"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"github.com/yhonda-ohishi/db_service/src/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 設定の検証
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// データベース接続
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := config.CloseDatabase(db); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	// ヘルスチェック
	if err := config.HealthCheck(db); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	// 本番データベース接続（オプション）
	prodDB, err := config.NewProdDatabase()
	if err != nil {
		log.Printf("Production database connection failed: %v (continuing without prod DB)", err)
		prodDB = nil
	} else {
		log.Printf("Production database connected successfully")
		defer func() {
			if err := prodDB.Close(); err != nil {
				log.Printf("Failed to close production database: %v", err)
			}
		}()
	}

	// SQL Serverデータベース接続（オプション）
	sqlServerDB, err := config.NewSQLServerDatabase()
	if err != nil {
		log.Printf("SQL Server database connection failed: %v (continuing without SQL Server)", err)
		sqlServerDB = nil
	} else {
		log.Printf("SQL Server database connected successfully")
		defer func() {
			sqlDB, _ := sqlServerDB.DB.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		}()
	}

	// リポジトリの初期化
	dtakoUriageKeihiRepo := repository.NewDTakoUriageKeihiRepository(db)
	etcMeisaiRepo := repository.NewETCMeisaiRepository(db)
	dtakoFerryRowsRepo := repository.NewDTakoFerryRowsRepository(db)

	// 本番DBリポジトリの初期化（現在無効化）
	// var dtakoRowsRepo repository.DTakoRowsRepository
	// var etcNumRepo repository.ETCNumRepository
	// if prodDB != nil {
	//     dtakoRowsRepo = repository.NewDTakoRowsRepository(prodDB)
	//     etcNumRepo = repository.NewETCNumRepository(prodDB)
	// }

	// gRPCサーバーの作成
	grpcServer := grpc.NewServer()

	// サービスの登録
	dtakoUriageKeihiService := service.NewDTakoUriageKeihiService(dtakoUriageKeihiRepo)
	proto.RegisterDb_DTakoUriageKeihiServiceServer(grpcServer, dtakoUriageKeihiService)

	etcMeisaiService := service.NewETCMeisaiService(etcMeisaiRepo)
	proto.RegisterDb_ETCMeisaiServiceServer(grpcServer, etcMeisaiService)

	dtakoFerryRowsService := service.NewDTakoFerryRowsService(dtakoFerryRowsRepo)
	proto.RegisterDb_DTakoFerryRowsServiceServer(grpcServer, dtakoFerryRowsService)

	// ETC明細マッピングサービスの登録
	etcMeisaiMappingRepo := repository.NewETCMeisaiMappingRepository(db)
	etcMeisaiMappingService := service.NewETCMeisaiMappingService(etcMeisaiMappingRepo)
	proto.RegisterDb_ETCMeisaiMappingServiceServer(grpcServer, etcMeisaiMappingService)

	// SQL Serverサービスの登録
	if sqlServerDB != nil {
		// SQL Serverリポジトリの初期化
		untenNippoMeisaiRepo := repository.NewUntenNippoMeisaiRepository(sqlServerDB)
		shainMasterRepo := repository.NewShainMasterRepository(sqlServerDB)
		chiikiMasterRepo := repository.NewChiikiMasterRepository(sqlServerDB)
		chikuMasterRepo := repository.NewChikuMasterRepository(sqlServerDB)

		// SQL Serverサービスの登録
		untenNippoMeisaiService := service.NewUntenNippoMeisaiService(untenNippoMeisaiRepo)
		proto.RegisterDb_UntenNippoMeisaiServiceServer(grpcServer, untenNippoMeisaiService)

		shainMasterService := service.NewShainMasterService(shainMasterRepo)
		proto.RegisterDb_ShainMasterServiceServer(grpcServer, shainMasterService)

		chiikiMasterService := service.NewChiikiMasterService(chiikiMasterRepo)
		proto.RegisterDb_ChiikiMasterServiceServer(grpcServer, chiikiMasterService)

		chikuMasterService := service.NewChikuMasterService(chikuMasterRepo)
		proto.RegisterDb_ChikuMasterServiceServer(grpcServer, chikuMasterService)

		log.Println("SQL Server services registered: UntenNippoMeisai, ShainMaster, ChiikiMaster, ChikuMaster")
	}

	// 本番DBサービスの登録（現在無効化）
	if prodDB != nil {
		// Note: 本番DBサービスは現在無効化されています（プロトコルバッファー定義が必要）
		log.Println("本番DBサービスは現在無効化されています")
	}

	// リフレクションの登録（開発環境用）
	reflection.Register(grpcServer)

	// リスナーの作成
	listener, err := net.Listen("tcp", cfg.GetGRPCAddress())
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", cfg.GetGRPCAddress(), err)
	}

	log.Printf("gRPC server starting on %s", cfg.GetGRPCAddress())

	// グレースフルシャットダウンの設定
	sigChan := make(chan os.Signal, 1)
	// プラットフォーム固有のシグナルを登録
	signal.Notify(sigChan, getShutdownSignals()...)
	log.Println("Signal handlers registered for graceful shutdown")

	// エラーチャネル
	errChan := make(chan error, 1)

	// サーバーを別ゴルーチンで起動
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			errChan <- err
		}
	}()

	// シャットダウンシグナルまたはエラーを待機
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	}

	// グレースフルシャットダウン開始
	log.Println("Initiating graceful shutdown...")

	// シャットダウン用のコンテキスト（30秒タイムアウト）
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// gRPCサーバーのグレースフルシャットダウン
	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	// タイムアウトまたは完了を待機
	select {
	case <-done:
		log.Println("gRPC server stopped gracefully")
	case <-shutdownCtx.Done():
		log.Println("Graceful shutdown timeout, forcing stop...")
		grpcServer.Stop()
	}

	// データベース接続のクリーンアップ
	log.Println("Cleaning up database connections...")

	// メインデータベースのクローズはdeferで処理される
	// 本番DBとSQL Serverのクローズもdeferで処理される

	log.Println("Shutdown complete")
}
