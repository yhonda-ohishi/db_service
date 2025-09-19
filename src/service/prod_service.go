package service

import (
	"context"

	"github.com/yhonda-ohishi/db_service/src/models"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DTakoRowsService 本番DB用DTakoRowsサービス実装
type DTakoRowsService struct {
	proto.UnimplementedDTakoRowsServiceServer
	repo repository.DTakoRowsRepository
}

// NewDTakoRowsService DTakoRowsサービスのコンストラクタ
func NewDTakoRowsService(repo repository.DTakoRowsRepository) *DTakoRowsService {
	return &DTakoRowsService{
		repo: repo,
	}
}

// Get DTakoRowsデータ取得
func (s *DTakoRowsService) Get(ctx context.Context, req *proto.GetDTakoRowsRequest) (*proto.DTakoRowsResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	model, err := s.repo.GetByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get record: %v", err)
	}

	return &proto.DTakoRowsResponse{
		DtakoRows: dtakoRowsModelToProto(model),
	}, nil
}

// List DTakoRowsデータ一覧取得
func (s *DTakoRowsService) List(ctx context.Context, req *proto.ListDTakoRowsRequest) (*proto.ListDTakoRowsResponse, error) {
	models, totalCount, err := s.repo.GetAll(int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list records: %v", err)
	}

	items := make([]*proto.DTakoRowsProd, len(models))
	for i, model := range models {
		items[i] = dtakoRowsModelToProto(model)
	}

	return &proto.ListDTakoRowsResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByOperationNo 運行NO別DTakoRowsデータ取得
func (s *DTakoRowsService) GetByOperationNo(ctx context.Context, req *proto.GetDTakoRowsByOperationNoRequest) (*proto.ListDTakoRowsResponse, error) {
	if req.OperationNo == "" {
		return nil, status.Error(codes.InvalidArgument, "operation_no is required")
	}

	models, err := s.repo.GetByOperationNo(req.OperationNo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get records: %v", err)
	}

	items := make([]*proto.DTakoRowsProd, len(models))
	for i, model := range models {
		items[i] = dtakoRowsModelToProto(model)
	}

	return &proto.ListDTakoRowsResponse{
		Items:      items,
		TotalCount: int32(len(items)),
	}, nil
}

// ETCNumService 本番DB用ETCNumサービス実装
type ETCNumService struct {
	proto.UnimplementedETCNumServiceServer
	repo repository.ETCNumRepository
}

// NewETCNumService ETCNumサービスのコンストラクタ
func NewETCNumService(repo repository.ETCNumRepository) *ETCNumService {
	return &ETCNumService{
		repo: repo,
	}
}

// List ETCNumデータ一覧取得
func (s *ETCNumService) List(ctx context.Context, req *proto.ListETCNumRequest) (*proto.ListETCNumResponse, error) {
	models, totalCount, err := s.repo.GetAll(int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list records: %v", err)
	}

	items := make([]*proto.ETCNumProd, len(models))
	for i, model := range models {
		items[i] = etcNumModelToProto(model)
	}

	return &proto.ListETCNumResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByETCCardNum ETCカード番号別データ取得
func (s *ETCNumService) GetByETCCardNum(ctx context.Context, req *proto.GetETCNumByETCCardNumRequest) (*proto.ListETCNumResponse, error) {
	if req.EtcCardNum == "" {
		return nil, status.Error(codes.InvalidArgument, "etc_card_num is required")
	}

	models, err := s.repo.GetByETCCardNum(req.EtcCardNum)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get records: %v", err)
	}

	items := make([]*proto.ETCNumProd, len(models))
	for i, model := range models {
		items[i] = etcNumModelToProto(model)
	}

	return &proto.ListETCNumResponse{
		Items:      items,
		TotalCount: int32(len(items)),
	}, nil
}

// GetByCarID 車輌ID別データ取得
func (s *ETCNumService) GetByCarID(ctx context.Context, req *proto.GetETCNumByCarIDRequest) (*proto.ListETCNumResponse, error) {
	if req.CarId == "" {
		return nil, status.Error(codes.InvalidArgument, "car_id is required")
	}

	models, err := s.repo.GetByCarID(req.CarId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get records: %v", err)
	}

	items := make([]*proto.ETCNumProd, len(models))
	for i, model := range models {
		items[i] = etcNumModelToProto(model)
	}

	return &proto.ListETCNumResponse{
		Items:      items,
		TotalCount: int32(len(items)),
	}, nil
}

// dtakoRowsModelToProto ModelからProtoへの変換
func dtakoRowsModelToProto(m *models.DTakoRows) *proto.DTakoRowsProd {
	p := &proto.DTakoRowsProd{
		Id:                    m.ID,
		OperationNo:           m.OperationNo,
		ReadDate:              m.ReadDate.Format("2006-01-02"),
		OperationDate:         m.OperationDate.Format("2006-01-02"),
		CarCode:               int32(m.CarCode),
		CarCc:                 m.CarCC,
		TargetDriverType:      int32(m.TargetDriverType),
		TargetDriverCode:      int32(m.TargetDriverCode),
		StartWorkDatetime:     m.StartWorkDatetime.Format("2006-01-02T15:04:05Z07:00"),
		EndWorkDatetime:       m.EndWorkDatetime.Format("2006-01-02T15:04:05Z07:00"),
		DepartureDatetime:     m.DepartureDateTime.Format("2006-01-02T15:04:05Z07:00"),
		ReturnDatetime:        m.ReturnDateTime.Format("2006-01-02T15:04:05Z07:00"),
		DepartureMeter:        m.DepartureMeter,
		ReturnMeter:           m.ReturnMeter,
		TotalDistance:         m.TotalDistance,
		GeneralRoadDriveTime:  int32(m.GeneralRoadDriveTime),
		HighwayDriveTime:      int32(m.HighwayDriveTime),
		BypassDriveTime:       int32(m.BypassDriveTime),
		LoadedDriveTime:       int32(m.LoadedDriveTime),
		EmptyDriveTime:        int32(m.EmptyDriveTime),
		Work1Time:             int32(m.Work1Time),
		Work2Time:             int32(m.Work2Time),
		Work3Time:             int32(m.Work3Time),
		Work4Time:             int32(m.Work4Time),
		Status1Distance:       m.Status1Distance,
		Status1Time:           int32(m.Status1Time),
	}

	if m.DriverCode1 != nil {
		driverCode1 := int32(*m.DriverCode1)
		p.DriverCode1 = &driverCode1
	}
	if m.LoadedDistance != nil {
		p.LoadedDistance = m.LoadedDistance
	}
	if m.DestinationCityName != nil {
		p.DestinationCityName = m.DestinationCityName
	}
	if m.DestinationPlaceName != nil {
		p.DestinationPlaceName = m.DestinationPlaceName
	}

	return p
}

// etcNumModelToProto ModelからProtoへの変換
func etcNumModelToProto(m *models.ETCNum) *proto.ETCNumProd {
	p := &proto.ETCNumProd{
		EtcCardNum: m.ETCCardNum,
		CarId:      m.CarID,
	}

	if m.StartDateTime != nil {
		startDateTime := m.StartDateTime.Format("2006-01-02T15:04:05Z07:00")
		p.StartDateTime = &startDateTime
	}
	if m.DueDateTime != nil {
		dueDateTime := m.DueDateTime.Format("2006-01-02T15:04:05Z07:00")
		p.DueDateTime = &dueDateTime
	}
	if m.ToChange != nil {
		p.ToChange = m.ToChange
	}

	return p
}