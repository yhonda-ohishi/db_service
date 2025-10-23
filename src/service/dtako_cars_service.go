package service

import (
	"context"

	"github.com/yhonda-ohishi/db_service/src/models/mysql"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DTakoCarsService gRPCサービス実装（本番DB、読み取り専用）
type DTakoCarsService struct {
	proto.UnimplementedDb_DTakoCarsServiceServer
	repo repository.DTakoCarsRepository
}

// NewDTakoCarsService サービスのコンストラクタ
func NewDTakoCarsService(repo repository.DTakoCarsRepository) *DTakoCarsService {
	return &DTakoCarsService{
		repo: repo,
	}
}

// Get 車輌情報取得
func (s *DTakoCarsService) Get(ctx context.Context, req *proto.Db_GetDTakoCarsRequest) (*proto.Db_DTakoCarsResponse, error) {
	car, err := s.repo.GetByID(int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "car not found: %v", err)
	}

	return &proto.Db_DTakoCarsResponse{
		DtakoCars: dtakoCarsModelToProto(car),
	}, nil
}

// List 車輌情報一覧取得
func (s *DTakoCarsService) List(ctx context.Context, req *proto.Db_ListDTakoCarsRequest) (*proto.Db_ListDTakoCarsResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)

	if limit == 0 {
		limit = 100
	}

	cars, totalCount, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list cars: %v", err)
	}

	items := make([]*proto.Db_DTakoCars, len(cars))
	for i, car := range cars {
		items[i] = dtakoCarsModelToProto(car)
	}

	return &proto.Db_ListDTakoCarsResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByCarCode 車輌CDで車輌情報取得
func (s *DTakoCarsService) GetByCarCode(ctx context.Context, req *proto.Db_GetDTakoCarsByCarCodeRequest) (*proto.Db_DTakoCarsResponse, error) {
	car, err := s.repo.GetByCarCode(req.CarCode)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "car not found: %v", err)
	}

	return &proto.Db_DTakoCarsResponse{
		DtakoCars: dtakoCarsModelToProto(car),
	}, nil
}

// dtakoCarsModelToProto ModelからProtoへの変換
func dtakoCarsModelToProto(model *mysql.DTakoCars) *proto.Db_DTakoCars {
	return &proto.Db_DTakoCars{
		Id:                  int32(model.ID),
		CarCode:             model.CarCode,
		CarCc:               model.CarCC,
		CarName:             model.CarName,
		BelongOfficeCode:    int32(model.BelongOfficeCode),
		HighwayCarType:      int32(model.HighwayCarType),
		FerryCarType:        int32(model.FerryCarType),
		EvaluationClassCode: int32(model.EvaluationClassCode),
		IdlingType:          int32(model.IdlingType),
		MaxLoadWeightKg:     int32(model.MaxLoadWeightKg),
		CarClass1:           int32(model.CarClass1),
		CarClass2:           int32(model.CarClass2),
		CarClass3:           int32(model.CarClass3),
		CarClass4:           int32(model.CarClass4),
		CarClass5:           int32(model.CarClass5),
		OperationType:       int32(model.OperationType),
	}
}
