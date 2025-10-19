package service

import (
	"context"

	"github.com/yhonda-ohishi/db_service/src/models"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DTakoRowsService gRPCサービス実装（本番DB、読み取り専用）
type DTakoRowsService struct {
	proto.UnimplementedDTakoRowsServiceServer
	repo repository.DTakoRowsRepository
}

// NewDTakoRowsService サービスのコンストラクタ
func NewDTakoRowsService(repo repository.DTakoRowsRepository) *DTakoRowsService {
	return &DTakoRowsService{
		repo: repo,
	}
}

// Get 運行データ取得
func (s *DTakoRowsService) Get(ctx context.Context, req *proto.GetDTakoRowsRequest) (*proto.DTakoRowsResponse, error) {
	row, err := s.repo.GetByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "row not found: %v", err)
	}

	return &proto.DTakoRowsResponse{
		DtakoRows: dtakoRowsModelToProto(row),
	}, nil
}

// List 運行データ一覧取得
func (s *DTakoRowsService) List(ctx context.Context, req *proto.ListDTakoRowsRequest) (*proto.ListDTakoRowsResponse, error) {
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

	rows, totalCount, err := s.repo.GetAll(limit, offset, orderBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list rows: %v", err)
	}

	items := make([]*proto.DTakoRows, len(rows))
	for i, row := range rows {
		items[i] = dtakoRowsModelToProto(row)
	}

	return &proto.ListDTakoRowsResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByOperationNo 運行NOで運行データ取得
func (s *DTakoRowsService) GetByOperationNo(ctx context.Context, req *proto.GetDTakoRowsByOperationNoRequest) (*proto.ListDTakoRowsResponse, error) {
	rows, err := s.repo.GetByOperationNo(req.OperationNo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get rows by operation_no: %v", err)
	}

	items := make([]*proto.DTakoRows, len(rows))
	for i, row := range rows {
		items[i] = dtakoRowsModelToProto(row)
	}

	return &proto.ListDTakoRowsResponse{
		Items:      items,
		TotalCount: int32(len(rows)),
	}, nil
}

// dtakoRowsModelToProto ModelからProtoへの変換
func dtakoRowsModelToProto(model *models.DTakoRows) *proto.DTakoRows {
	protoRow := &proto.DTakoRows{
		Id:                   model.ID,
		OperationNo:          model.OperationNo,
		ReadDate:             model.ReadDate.Format("2006-01-02T15:04:05Z07:00"),
		OperationDate:        model.OperationDate.Format("2006-01-02T15:04:05Z07:00"),
		CarCode:              int32(model.CarCode),
		CarCc:                model.CarCC,
		TargetDriverType:     int32(model.TargetDriverType),
		TargetDriverCode:     int32(model.TargetDriverCode),
		StartWorkDatetime:    model.StartWorkDatetime.Format("2006-01-02T15:04:05Z07:00"),
		EndWorkDatetime:      model.EndWorkDatetime.Format("2006-01-02T15:04:05Z07:00"),
		DepartureDatetime:    model.DepartureDateTime.Format("2006-01-02T15:04:05Z07:00"),
		ReturnDatetime:       model.ReturnDateTime.Format("2006-01-02T15:04:05Z07:00"),
		DepartureMeter:       model.DepartureMeter,
		ReturnMeter:          model.ReturnMeter,
		TotalDistance:        model.TotalDistance,
		GeneralRoadDriveTime: int32(model.GeneralRoadDriveTime),
		HighwayDriveTime:     int32(model.HighwayDriveTime),
		BypassDriveTime:      int32(model.BypassDriveTime),
		LoadedDriveTime:      int32(model.LoadedDriveTime),
		EmptyDriveTime:       int32(model.EmptyDriveTime),
		Work1Time:            int32(model.Work1Time),
		Work2Time:            int32(model.Work2Time),
		Work3Time:            int32(model.Work3Time),
		Work4Time:            int32(model.Work4Time),
		Status1Distance:      model.Status1Distance,
		Status1Time:          int32(model.Status1Time),
	}

	// Optional fields
	if model.DriverCode1 != nil {
		val := int32(*model.DriverCode1)
		protoRow.DriverCode1 = &val
	}
	if model.LoadedDistance != nil {
		protoRow.LoadedDistance = model.LoadedDistance
	}
	if model.DestinationCityName != nil {
		protoRow.DestinationCityName = model.DestinationCityName
	}
	if model.DestinationPlaceName != nil {
		protoRow.DestinationPlaceName = model.DestinationPlaceName
	}

	return protoRow
}
