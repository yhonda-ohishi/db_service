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

// ServiceRegistry holds all db_service gRPC service implementations
type ServiceRegistry struct {
	// ローカルDB用サービス
	ETCMeisaiService        dbproto.ETCMeisaiServiceServer
	DTakoUriageKeihiService dbproto.DTakoUriageKeihiServiceServer
	DTakoFerryRowsService   dbproto.DTakoFerryRowsServiceServer
	ETCMeisaiMappingService dbproto.ETCMeisaiMappingServiceServer

	// 本番DB用サービス（読み取り専用）
	DTakoCarsService         dbproto.DTakoCarsServiceServer
	DTakoEventsService       dbproto.DTakoEventsServiceServer
	DTakoRowsService         dbproto.DTakoRowsServiceServer
	ETCNumService            dbproto.ETCNumServiceServer
	DTakoFerryRowsProdService dbproto.DTakoFerryRowsProdServiceServer
	CarsService              dbproto.CarsServiceServer
	DriversService           dbproto.DriversServiceServer

	// SQL Server (ichibanboshi) 用サービス（読み取り専用）
	UntenNippoMeisaiService dbproto.UntenNippoMeisaiServiceServer
	ShainMasterService      dbproto.ShainMasterServiceServer
	ChiikiMasterService     dbproto.ChiikiMasterServiceServer
	ChikuMasterService      dbproto.ChikuMasterServiceServer
}

// NewServiceRegistry creates a new service registry with all db_service services initialized
// Returns nil if db_service initialization fails
func NewServiceRegistry() *ServiceRegistry {
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

	// Initialize production DB connection (optional)
	prodDB, err := config.NewProdDatabase()
	var dtakoCarsService dbproto.DTakoCarsServiceServer
	var dtakoEventsService dbproto.DTakoEventsServiceServer
	var dtakoRowsService dbproto.DTakoRowsServiceServer
	var etcNumService dbproto.ETCNumServiceServer
	var dtakoFerryRowsProdService dbproto.DTakoFerryRowsProdServiceServer
	var carsService dbproto.CarsServiceServer
	var driversService dbproto.DriversServiceServer

	// Initialize SQL Server (ichibanboshi) connection (optional)
	sqlServerDB, sqlErr := config.NewSQLServerDatabase()
	var untenNippoMeisaiService dbproto.UntenNippoMeisaiServiceServer
	var shainMasterService dbproto.ShainMasterServiceServer
	var chiikiMasterService dbproto.ChiikiMasterServiceServer
	var chikuMasterService dbproto.ChikuMasterServiceServer

	if err == nil && prodDB != nil {
		// Initialize production DB repositories
		dtakoCarsRepo := repository.NewDTakoCarsRepository(prodDB)
		dtakoEventsRepo := repository.NewDTakoEventsRepository(prodDB)
		dtakoRowsRepo := repository.NewDTakoRowsRepository(prodDB)
		etcNumRepo := repository.NewETCNumRepository(prodDB)
		dtakoFerryRowsProdRepo := repository.NewDTakoFerryRowsProdRepository(prodDB)
		carsRepo := repository.NewCarsRepository(prodDB)
		driversRepo := repository.NewDriversRepository(prodDB)

		// Initialize production DB services
		dtakoCarsService = service.NewDTakoCarsService(dtakoCarsRepo)
		dtakoEventsService = service.NewDTakoEventsService(dtakoEventsRepo)
		dtakoRowsService = service.NewDTakoRowsService(dtakoRowsRepo)
		etcNumService = service.NewETCNumService(etcNumRepo)
		dtakoFerryRowsProdService = service.NewDTakoFerryRowsProdService(dtakoFerryRowsProdRepo)
		carsService = service.NewCarsService(carsRepo)
		driversService = service.NewDriversService(driversRepo)

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

		// Production DB services (may be nil if prod DB not available)
		DTakoCarsService:         dtakoCarsService,
		DTakoEventsService:       dtakoEventsService,
		DTakoRowsService:         dtakoRowsService,
		ETCNumService:            etcNumService,
		DTakoFerryRowsProdService: dtakoFerryRowsProdService,
		CarsService:              carsService,
		DriversService:           driversService,

		// SQL Server services (may be nil if SQL Server not available)
		UntenNippoMeisaiService: untenNippoMeisaiService,
		ShainMasterService:      shainMasterService,
		ChiikiMasterService:     chiikiMasterService,
		ChikuMasterService:      chikuMasterService,
	}
}

// RegisterAll registers all db_service services to the gRPC server
// This method automatically detects and registers all available services from db_service
// When new services are added to db_service, they will be automatically registered here
func (r *ServiceRegistry) RegisterAll(server *grpc.Server) {
	if r.ETCMeisaiService != nil {
		dbproto.RegisterETCMeisaiServiceServer(server, r.ETCMeisaiService)
		log.Println("Registered: ETCMeisaiService")
	}
	if r.DTakoUriageKeihiService != nil {
		dbproto.RegisterDTakoUriageKeihiServiceServer(server, r.DTakoUriageKeihiService)
		log.Println("Registered: DTakoUriageKeihiService")
	}
	if r.DTakoFerryRowsService != nil {
		dbproto.RegisterDTakoFerryRowsServiceServer(server, r.DTakoFerryRowsService)
		log.Println("Registered: DTakoFerryRowsService")
	}
	if r.ETCMeisaiMappingService != nil {
		dbproto.RegisterETCMeisaiMappingServiceServer(server, r.ETCMeisaiMappingService)
		log.Println("Registered: ETCMeisaiMappingService")
	}

	// Production DB services
	if r.DTakoCarsService != nil {
		dbproto.RegisterDTakoCarsServiceServer(server, r.DTakoCarsService)
		log.Println("Registered: DTakoCarsService (Production DB)")
	}
	if r.DTakoEventsService != nil {
		dbproto.RegisterDTakoEventsServiceServer(server, r.DTakoEventsService)
		log.Println("Registered: DTakoEventsService (Production DB)")
	}
	if r.DTakoRowsService != nil {
		dbproto.RegisterDTakoRowsServiceServer(server, r.DTakoRowsService)
		log.Println("Registered: DTakoRowsService (Production DB)")
	}
	if r.ETCNumService != nil {
		dbproto.RegisterETCNumServiceServer(server, r.ETCNumService)
		log.Println("Registered: ETCNumService (Production DB)")
	}
	if r.DTakoFerryRowsProdService != nil {
		dbproto.RegisterDTakoFerryRowsProdServiceServer(server, r.DTakoFerryRowsProdService)
		log.Println("Registered: DTakoFerryRowsProdService (Production DB)")
	}
	if r.CarsService != nil {
		dbproto.RegisterCarsServiceServer(server, r.CarsService)
		log.Println("Registered: CarsService (Production DB)")
	}
	if r.DriversService != nil {
		dbproto.RegisterDriversServiceServer(server, r.DriversService)
		log.Println("Registered: DriversService (Production DB)")
	}

	// SQL Server services
	if r.UntenNippoMeisaiService != nil {
		dbproto.RegisterUntenNippoMeisaiServiceServer(server, r.UntenNippoMeisaiService)
		log.Println("Registered: UntenNippoMeisaiService (SQL Server)")
	}
	if r.ShainMasterService != nil {
		dbproto.RegisterShainMasterServiceServer(server, r.ShainMasterService)
		log.Println("Registered: ShainMasterService (SQL Server)")
	}
	if r.ChiikiMasterService != nil {
		dbproto.RegisterChiikiMasterServiceServer(server, r.ChiikiMasterService)
		log.Println("Registered: ChiikiMasterService (SQL Server)")
	}
	if r.ChikuMasterService != nil {
		dbproto.RegisterChikuMasterServiceServer(server, r.ChikuMasterService)
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
//	registry.Register(grpcServer)
func Register(server *grpc.Server) *ServiceRegistry {
	registry := NewServiceRegistry()
	if registry == nil {
		log.Println("Warning: db_service not available, running without database services")
		return nil
	}

	registry.RegisterAll(server)
	return registry
}
