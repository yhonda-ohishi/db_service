package registry

import (
	"fmt"
	"log"

	"github.com/yhonda-ohishi/db_service/src/config"
	dbproto "github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"github.com/yhonda-ohishi/db_service/src/service"
	"google.golang.org/grpc"
)

// RegistryOptions はサービス登録のオプション
type RegistryOptions struct {
	// 除外するサービス名のリスト
	ExcludeServices []string
}

// RegistryOption は関数オプションの型
type RegistryOption func(*RegistryOptions)

// WithExcludeServices 指定したサービスを除外するオプション
func WithExcludeServices(services ...string) RegistryOption {
	return func(o *RegistryOptions) {
		o.ExcludeServices = append(o.ExcludeServices, services...)
	}
}

// isExcluded サービスが除外リストに含まれているかチェック
func (o *RegistryOptions) isExcluded(serviceName string) bool {
	for _, excluded := range o.ExcludeServices {
		if excluded == serviceName {
			return true
		}
	}
	return false
}

// ServiceRegistry holds all db_service gRPC service implementations
type ServiceRegistry struct {
	// ローカルDB用サービス
	ETCMeisaiService        dbproto.Db_ETCMeisaiServiceServer
	DTakoUriageKeihiService dbproto.Db_DTakoUriageKeihiServiceServer
	DTakoFerryRowsService   dbproto.Db_DTakoFerryRowsServiceServer
	ETCMeisaiMappingService dbproto.Db_ETCMeisaiMappingServiceServer
	TimeCardDevService      dbproto.Db_TimeCardDevServiceServer
	TimeCardLogService      dbproto.Db_TimeCardLogServiceServer

	// 本番DB用サービス（読み取り専用）
	DTakoCarsService         dbproto.Db_DTakoCarsServiceServer
	DTakoEventsService       dbproto.Db_DTakoEventsServiceServer
	DTakoRowsService         dbproto.Db_DTakoRowsServiceServer
	ETCNumService            dbproto.Db_ETCNumServiceServer
	DTakoFerryRowsProdService dbproto.Db_DTakoFerryRowsProdServiceServer
	CarsService              dbproto.Db_CarsServiceServer
	DriversService           dbproto.Db_DriversServiceServer
	TimeCardService          dbproto.Db_TimeCardServiceServer

	// SQL Server (ichibanboshi) 用サービス（読み取り専用）
	UntenNippoMeisaiService dbproto.Db_UntenNippoMeisaiServiceServer
	ShainMasterService      dbproto.Db_ShainMasterServiceServer
	ChiikiMasterService     dbproto.Db_ChiikiMasterServiceServer
	ChikuMasterService      dbproto.Db_ChikuMasterServiceServer

	// オプション
	options *RegistryOptions
}

// NewServiceRegistry creates a new service registry with all db_service services initialized
// Returns nil if db_service initialization fails
func NewServiceRegistry(opts ...RegistryOption) *ServiceRegistry {
	// デフォルトオプション
	options := &RegistryOptions{}
	for _, opt := range opts {
		opt(options)
	}
	// Load db_service configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: Failed to load db_service config: %v", err)
		return nil
	}

	// Initialize db_service database connection
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize db_service database: %v", err)
		return nil
	}

	// Health check
	if err := config.HealthCheck(db); err != nil {
		log.Printf("Warning: db_service database health check failed: %v", err)
		return nil
	}

	// Initialize local DB repositories
	dtakoUriageKeihiRepo := repository.NewDTakoUriageKeihiRepository(db)
	etcMeisaiRepo := repository.NewETCMeisaiRepository(db)
	dtakoFerryRowsRepo := repository.NewDTakoFerryRowsRepository(db)
	etcMeisaiMappingRepo := repository.NewETCMeisaiMappingRepository(db)
	timeCardDevRepo := repository.NewTimeCardDevRepository(db)
	timeCardLogRepo := repository.NewTimeCardLogRepository(db)

	// Initialize production DB connection (optional)
	prodDB, err := config.NewProdDatabase()
	var dtakoCarsService dbproto.Db_DTakoCarsServiceServer
	var dtakoEventsService dbproto.Db_DTakoEventsServiceServer
	var dtakoRowsService dbproto.Db_DTakoRowsServiceServer
	var etcNumService dbproto.Db_ETCNumServiceServer
	var dtakoFerryRowsProdService dbproto.Db_DTakoFerryRowsProdServiceServer
	var carsService dbproto.Db_CarsServiceServer
	var driversService dbproto.Db_DriversServiceServer
	var timeCardService dbproto.Db_TimeCardServiceServer

	// Initialize SQL Server (ichibanboshi) connection (optional)
	sqlServerDB, sqlErr := config.NewSQLServerDatabase()
	var untenNippoMeisaiService dbproto.Db_UntenNippoMeisaiServiceServer
	var shainMasterService dbproto.Db_ShainMasterServiceServer
	var chiikiMasterService dbproto.Db_ChiikiMasterServiceServer
	var chikuMasterService dbproto.Db_ChikuMasterServiceServer

	if err == nil && prodDB != nil {
		// Initialize production DB repositories
		dtakoCarsRepo := repository.NewDTakoCarsRepository(prodDB)
		dtakoEventsRepo := repository.NewDTakoEventsRepository(prodDB)
		dtakoRowsRepo := repository.NewDTakoRowsRepository(prodDB)
		etcNumRepo := repository.NewETCNumRepository(prodDB)
		dtakoFerryRowsProdRepo := repository.NewDTakoFerryRowsProdRepository(prodDB)
		carsRepo := repository.NewCarsRepository(prodDB)
		driversRepo := repository.NewDriversRepository(prodDB)
		timeCardRepo := repository.NewTimeCardRepository(prodDB)

		// Initialize production DB services
		dtakoCarsService = service.NewDTakoCarsService(dtakoCarsRepo)
		dtakoEventsService = service.NewDTakoEventsService(dtakoEventsRepo)
		dtakoRowsService = service.NewDTakoRowsService(dtakoRowsRepo)
		etcNumService = service.NewETCNumService(etcNumRepo)
		dtakoFerryRowsProdService = service.NewDTakoFerryRowsProdService(dtakoFerryRowsProdRepo)
		carsService = service.NewCarsService(carsRepo)
		driversService = service.NewDriversService(driversRepo)
		timeCardService = service.NewTimeCardService(timeCardRepo)

		log.Println("Production DB services initialized successfully")
	} else {
		log.Printf("Warning: Production DB not available: %v", err)
	}

	if sqlErr == nil && sqlServerDB != nil {
		// Initialize SQL Server repositories
		untenNippoMeisaiRepo := repository.NewUntenNippoMeisaiRepository(sqlServerDB)
		shainMasterRepo := repository.NewShainMasterRepository(sqlServerDB)
		chiikiMasterRepo := repository.NewChiikiMasterRepository(sqlServerDB)
		chikuMasterRepo := repository.NewChikuMasterRepository(sqlServerDB)

		// Initialize SQL Server services
		untenNippoMeisaiService = service.NewUntenNippoMeisaiService(untenNippoMeisaiRepo)
		shainMasterService = service.NewShainMasterService(shainMasterRepo)
		chiikiMasterService = service.NewChiikiMasterService(chiikiMasterRepo)
		chikuMasterService = service.NewChikuMasterService(chikuMasterRepo)

		log.Println("SQL Server (ichibanboshi) services initialized successfully")
	} else {
		log.Printf("Warning: SQL Server not available: %v", sqlErr)
	}

	return &ServiceRegistry{
		// Local DB services
		ETCMeisaiService:        service.NewETCMeisaiService(etcMeisaiRepo),
		DTakoUriageKeihiService: service.NewDTakoUriageKeihiService(dtakoUriageKeihiRepo),
		DTakoFerryRowsService:   service.NewDTakoFerryRowsService(dtakoFerryRowsRepo),
		ETCMeisaiMappingService: service.NewETCMeisaiMappingService(etcMeisaiMappingRepo),
		TimeCardDevService:      service.NewTimeCardDevService(timeCardDevRepo),
		TimeCardLogService:      service.NewTimeCardLogService(timeCardLogRepo),

		// Production DB services (may be nil if prod DB not available)
		DTakoCarsService:         dtakoCarsService,
		DTakoEventsService:       dtakoEventsService,
		DTakoRowsService:         dtakoRowsService,
		ETCNumService:            etcNumService,
		DTakoFerryRowsProdService: dtakoFerryRowsProdService,
		CarsService:              carsService,
		DriversService:           driversService,
		TimeCardService:          timeCardService,

		// SQL Server services (may be nil if SQL Server not available)
		UntenNippoMeisaiService: untenNippoMeisaiService,
		ShainMasterService:      shainMasterService,
		ChiikiMasterService:     chiikiMasterService,
		ChikuMasterService:      chikuMasterService,

		// オプション保存
		options: options,
	}
}

// RegisterAll registers all db_service services to the gRPC server
// This method automatically detects and registers all available services from db_service
// When new services are added to db_service, they will be automatically registered here
func (r *ServiceRegistry) RegisterAll(server *grpc.Server) {
	if r.ETCMeisaiService != nil {
		dbproto.RegisterDb_ETCMeisaiServiceServer(server, r.ETCMeisaiService)
		log.Println("Registered: ETCMeisaiService")
	}
	if r.DTakoUriageKeihiService != nil {
		dbproto.RegisterDb_DTakoUriageKeihiServiceServer(server, r.DTakoUriageKeihiService)
		log.Println("Registered: DTakoUriageKeihiService")
	}
	if r.DTakoFerryRowsService != nil {
		dbproto.RegisterDb_DTakoFerryRowsServiceServer(server, r.DTakoFerryRowsService)
		log.Println("Registered: DTakoFerryRowsService")
	}
	if r.ETCMeisaiMappingService != nil {
		dbproto.RegisterDb_ETCMeisaiMappingServiceServer(server, r.ETCMeisaiMappingService)
		log.Println("Registered: ETCMeisaiMappingService")
	}
	if r.TimeCardDevService != nil {
		dbproto.RegisterDb_TimeCardDevServiceServer(server, r.TimeCardDevService)
		log.Println("Registered: TimeCardDevService (Local DB)")
	}
	if r.TimeCardLogService != nil {
		dbproto.RegisterDb_TimeCardLogServiceServer(server, r.TimeCardLogService)
		log.Println("Registered: TimeCardLogService (Local DB)")
	}

	// Production DB services
	if r.DTakoCarsService != nil {
		dbproto.RegisterDb_DTakoCarsServiceServer(server, r.DTakoCarsService)
		log.Println("Registered: DTakoCarsService (Production DB)")
	}
	if r.DTakoEventsService != nil && !r.options.isExcluded("DTakoEventsService") {
		dbproto.RegisterDb_DTakoEventsServiceServer(server, r.DTakoEventsService)
		log.Println("Registered: DTakoEventsService (Production DB)")
	}
	if r.DTakoRowsService != nil && !r.options.isExcluded("DTakoRowsService") {
		dbproto.RegisterDb_DTakoRowsServiceServer(server, r.DTakoRowsService)
		log.Println("Registered: DTakoRowsService (Production DB)")
	}
	if r.ETCNumService != nil {
		dbproto.RegisterDb_ETCNumServiceServer(server, r.ETCNumService)
		log.Println("Registered: ETCNumService (Production DB)")
	}
	if r.DTakoFerryRowsProdService != nil {
		dbproto.RegisterDb_DTakoFerryRowsProdServiceServer(server, r.DTakoFerryRowsProdService)
		log.Println("Registered: DTakoFerryRowsProdService (Production DB)")
	}
	if r.CarsService != nil {
		dbproto.RegisterDb_CarsServiceServer(server, r.CarsService)
		log.Println("Registered: CarsService (Production DB)")
	}
	if r.DriversService != nil {
		dbproto.RegisterDb_DriversServiceServer(server, r.DriversService)
		log.Println("Registered: DriversService (Production DB)")
	}
	if r.TimeCardService != nil {
		dbproto.RegisterDb_TimeCardServiceServer(server, r.TimeCardService)
		log.Println("Registered: TimeCardService (Production DB)")
	}

	// SQL Server services
	if r.UntenNippoMeisaiService != nil {
		dbproto.RegisterDb_UntenNippoMeisaiServiceServer(server, r.UntenNippoMeisaiService)
		log.Println("Registered: UntenNippoMeisaiService (SQL Server)")
	}
	if r.ShainMasterService != nil {
		dbproto.RegisterDb_ShainMasterServiceServer(server, r.ShainMasterService)
		log.Println("Registered: ShainMasterService (SQL Server)")
	}
	if r.ChiikiMasterService != nil {
		dbproto.RegisterDb_ChiikiMasterServiceServer(server, r.ChiikiMasterService)
		log.Println("Registered: ChiikiMasterService (SQL Server)")
	}
	if r.ChikuMasterService != nil {
		dbproto.RegisterDb_ChikuMasterServiceServer(server, r.ChikuMasterService)
		log.Println("Registered: ChikuMasterService (SQL Server)")
	}

	fmt.Println("db_service: All services registered successfully")
}

// Register is a convenience function for creating and registering all db_service services
// If db_service initialization fails, it logs a warning but does not fail
//
// Usage in other projects:
//
//	import "github.com/yhonda-ohishi/db_service/src/registry"
//
//	grpcServer := grpc.NewServer()
//	// 全サービスを登録
//	registry.Register(grpcServer)
//	// または特定のサービスを除外
//	registry.Register(grpcServer, registry.WithExcludeServices("DTakoRowsService", "DTakoEventsService"))
func Register(server *grpc.Server, opts ...RegistryOption) *ServiceRegistry {
	registry := NewServiceRegistry(opts...)
	if registry == nil {
		log.Println("Warning: db_service not available, running without database services")
		return nil
	}

	registry.RegisterAll(server)
	return registry
}
