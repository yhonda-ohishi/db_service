package service

import (
	"context"

	"github.com/yhonda-ohishi/db_service/src/models/mysql"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DriversService gRPCサービス実装（本番DB、読み取り専用）
type DriversService struct {
	proto.UnimplementedDriversServiceServer
	repo repository.DriversRepository
}

// NewDriversService サービスのコンストラクタ
func NewDriversService(repo repository.DriversRepository) *DriversService {
	return &DriversService{
		repo: repo,
	}
}

// Get ドライバー情報取得
func (s *DriversService) Get(ctx context.Context, req *proto.GetDriversRequest) (*proto.DriversResponse, error) {
	driver, err := s.repo.GetByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "driver not found: %v", err)
	}

	return &proto.DriversResponse{
		Drivers: driversModelToProto(driver),
	}, nil
}

// List ドライバー情報一覧取得
func (s *DriversService) List(ctx context.Context, req *proto.ListDriversRequest) (*proto.ListDriversResponse, error) {
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

	drivers, totalCount, err := s.repo.GetAll(limit, offset, orderBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list drivers: %v", err)
	}

	items := make([]*proto.Drivers, len(drivers))
	for i, driver := range drivers {
		items[i] = driversModelToProto(driver)
	}

	return &proto.ListDriversResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByBumon 部門コードでドライバー情報取得
func (s *DriversService) GetByBumon(ctx context.Context, req *proto.GetDriversByBumonRequest) (*proto.ListDriversResponse, error) {
	drivers, err := s.repo.GetByBumon(req.Bumon)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get drivers by bumon: %v", err)
	}

	items := make([]*proto.Drivers, len(drivers))
	for i, driver := range drivers {
		items[i] = driversModelToProto(driver)
	}

	return &proto.ListDriversResponse{
		Items:      items,
		TotalCount: int32(len(drivers)),
	}, nil
}

// driversModelToProto ModelからProtoへの変換
func driversModelToProto(model *mysql.Drivers) *proto.Drivers {
	protoDriver := &proto.Drivers{
		Id:          int32(model.ID),
		Bumon:       model.Bumon,
		KinmuTaikei: int32(model.KinmuTaikei),
	}

	// Optional fields
	if model.Name != nil {
		protoDriver.Name = model.Name
	}
	if model.ShainR != nil {
		protoDriver.ShainR = model.ShainR
	}
	if model.JoinDate != nil {
		date := model.JoinDate.Format("2006-01-02")
		protoDriver.JoinDate = &date
	}
	if model.RetireDate != nil {
		date := model.RetireDate.Format("2006-01-02")
		protoDriver.RetireDate = &date
	}
	if model.Bunrui1 != nil {
		protoDriver.Bunrui1 = model.Bunrui1
	}
	if model.Bunrui2 != nil {
		protoDriver.Bunrui2 = model.Bunrui2
	}
	if model.Kubun != nil {
		val := int32(*model.Kubun)
		protoDriver.Kubun = &val
	}

	return protoDriver
}
