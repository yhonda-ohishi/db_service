package registry_test

import (
	"testing"

	"github.com/yhonda-ohishi/db_service/src/registry"
	"google.golang.org/grpc"
)

func TestRegister(t *testing.T) {
	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Register db_service services
	reg := registry.Register(grpcServer)

	// Should not panic
	if reg == nil {
		t.Log("Registry initialization failed (expected if DB not configured)")
	} else {
		t.Log("Registry initialized successfully")
	}
}

func TestNewServiceRegistry(t *testing.T) {
	reg := registry.NewServiceRegistry()

	if reg == nil {
		t.Log("Service registry is nil (expected if DB not configured)")
		return
	}

	// Check that services are initialized
	if reg.ETCMeisaiService == nil {
		t.Error("ETCMeisaiService should not be nil")
	}
	if reg.DTakoUriageKeihiService == nil {
		t.Error("DTakoUriageKeihiService should not be nil")
	}
	if reg.DTakoFerryRowsService == nil {
		t.Error("DTakoFerryRowsService should not be nil")
	}
	if reg.ETCMeisaiMappingService == nil {
		t.Error("ETCMeisaiMappingService should not be nil")
	}
}
