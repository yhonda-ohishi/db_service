package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
	proto.RegisterDTakoUriageKeihiServiceServer(grpcServer, dtakoUriageKeihiService)

	etcMeisaiService := service.NewETCMeisaiService(etcMeisaiRepo)
	proto.RegisterETCMeisaiServiceServer(grpcServer, etcMeisaiService)

	dtakoFerryRowsService := service.NewDTakoFerryRowsService(dtakoFerryRowsRepo)
	proto.RegisterDTakoFerryRowsServiceServer(grpcServer, dtakoFerryRowsService)

	// ETC明細マッピングサービスの登録
	etcMeisaiMappingRepo := repository.NewETCMeisaiMappingRepository(db)
	etcMeisaiMappingService := service.NewETCMeisaiMappingService(etcMeisaiMappingRepo)
	proto.RegisterETCMeisaiMappingServiceServer(grpcServer, etcMeisaiMappingService)

	// SQL Serverサービスの登録
	if sqlServerDB != nil {
		// SQL Serverリポジトリの初期化
		untenNippoMeisaiRepo := repository.NewUntenNippoMeisaiRepository(sqlServerDB)
		shainMasterRepo := repository.NewShainMasterRepository(sqlServerDB)
		chiikiMasterRepo := repository.NewChiikiMasterRepository(sqlServerDB)
		chikuMasterRepo := repository.NewChikuMasterRepository(sqlServerDB)

		// SQL Serverサービスの登録
		untenNippoMeisaiService := service.NewUntenNippoMeisaiService(untenNippoMeisaiRepo)
		proto.RegisterUntenNippoMeisaiServiceServer(grpcServer, untenNippoMeisaiService)

		shainMasterService := service.NewShainMasterService(shainMasterRepo)
		proto.RegisterShainMasterServiceServer(grpcServer, shainMasterService)

		chiikiMasterService := service.NewChiikiMasterService(chiikiMasterRepo)
		proto.RegisterChiikiMasterServiceServer(grpcServer, chiikiMasterService)

		chikuMasterService := service.NewChikuMasterService(chikuMasterRepo)
		proto.RegisterChikuMasterServiceServer(grpcServer, chikuMasterService)

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
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	// サーバーの起動
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
