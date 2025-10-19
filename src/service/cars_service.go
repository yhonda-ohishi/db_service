package service

import (
	"context"

	"github.com/yhonda-ohishi/db_service/src/models"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CarsService gRPCサービス実装（本番DB、読み取り専用）
type CarsService struct {
	proto.UnimplementedCarsServiceServer
	repo repository.CarsRepository
}

// NewCarsService サービスのコンストラクタ
func NewCarsService(repo repository.CarsRepository) *CarsService {
	return &CarsService{
		repo: repo,
	}
}

// Get 車両情報取得
func (s *CarsService) Get(ctx context.Context, req *proto.GetCarsRequest) (*proto.CarsResponse, error) {
	car, err := s.repo.GetByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "car not found: %v", err)
	}

	return &proto.CarsResponse{
		Cars: carsModelToProto(car),
	}, nil
}

// List 車両情報一覧取得
func (s *CarsService) List(ctx context.Context, req *proto.ListCarsRequest) (*proto.ListCarsResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)

	if limit == 0 {
		limit = 100
	}

	// order_byパラメータを取得（nilの場合は空文字列）
	orderBy := ""
	if req.OrderBy != nil {
		orderBy = *req.OrderBy
	}

	cars, totalCount, err := s.repo.GetAll(limit, offset, orderBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list cars: %v", err)
	}

	items := make([]*proto.Cars, len(cars))
	for i, car := range cars {
		items[i] = carsModelToProto(car)
	}

	return &proto.ListCarsResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByBumonCodeID 部門コードで車両情報取得
func (s *CarsService) GetByBumonCodeID(ctx context.Context, req *proto.GetCarsByBumonCodeIDRequest) (*proto.ListCarsResponse, error) {
	cars, err := s.repo.GetByBumonCodeID(req.BumonCodeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cars by bumon_code_id: %v", err)
	}

	items := make([]*proto.Cars, len(cars))
	for i, car := range cars {
		items[i] = carsModelToProto(car)
	}

	return &proto.ListCarsResponse{
		Items:      items,
		TotalCount: int32(len(cars)),
	}, nil
}

// carsModelToProto ModelからProtoへの変換
func carsModelToProto(model *models.Cars) *proto.Cars {
	protoCar := &proto.Cars{
		Id:   model.ID,
		Id4:  int32(model.ID4),
		Dai1: int32(model.Dai1),
		Chu1: int32(model.Chu1),
		Sho1: int32(model.Sho1),
		Dai2: int32(model.Dai2),
		Chu2: int32(model.Chu2),
		Sho2: int32(model.Sho2),
	}

	// Optional fields
	if model.Name != nil {
		protoCar.Name = model.Name
	}
	if model.NameR != nil {
		protoCar.NameR = model.NameR
	}
	if model.Shashu != nil {
		protoCar.Shashu = model.Shashu
	}
	if model.Sekisai != nil {
		protoCar.Sekisai = model.Sekisai
	}
	if model.Youseki != nil {
		protoCar.Youseki = model.Youseki
	}
	if model.RegDate != nil {
		protoCar.RegDate = stringPtr(model.RegDate.Format("2006-01-02"))
	}
	if model.NextInspectDate != nil {
		protoCar.NextInspectDate = stringPtr(model.NextInspectDate.Format("2006-01-02"))
	}
	if model.ParchDate != nil {
		protoCar.ParchDate = stringPtr(model.ParchDate.Format("2006-01-02"))
	}
	if model.ScrapDate != nil {
		protoCar.ScrapDate = stringPtr(model.ScrapDate.Format("2006-01-02"))
	}
	if model.BumonCodeID != nil {
		protoCar.BumonCodeId = model.BumonCodeID
	}
	if model.DriverID != nil {
		val := int32(*model.DriverID)
		protoCar.DriverId = &val
	}
	if model.ETC != nil {
		val := int32(*model.ETC)
		protoCar.Etc = &val
	}
	if model.DaiChuSho1 != nil {
		protoCar.Daichusho1 = model.DaiChuSho1
	}
	if model.DaiChuSho2 != nil {
		protoCar.Daichusho2 = model.DaiChuSho2
	}

	return protoCar
}
