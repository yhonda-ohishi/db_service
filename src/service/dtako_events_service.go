package service

import (
	"context"

	"github.com/yhonda-ohishi/db_service/src/models"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DTakoEventsService gRPCサービス実装（本番DB、読み取り専用）
type DTakoEventsService struct {
	proto.UnimplementedDTakoEventsServiceServer
	repo repository.DTakoEventsRepository
}

// NewDTakoEventsService サービスのコンストラクタ
func NewDTakoEventsService(repo repository.DTakoEventsRepository) *DTakoEventsService {
	return &DTakoEventsService{
		repo: repo,
	}
}

// Get イベント情報取得
func (s *DTakoEventsService) Get(ctx context.Context, req *proto.GetDTakoEventsRequest) (*proto.DTakoEventsResponse, error) {
	event, err := s.repo.GetByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "event not found: %v", err)
	}

	return &proto.DTakoEventsResponse{
		DtakoEvents: dtakoEventsModelToProto(event),
	}, nil
}

// List イベント情報一覧取得
func (s *DTakoEventsService) List(ctx context.Context, req *proto.ListDTakoEventsRequest) (*proto.ListDTakoEventsResponse, error) {
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

	events, totalCount, err := s.repo.GetAll(limit, offset, orderBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events: %v", err)
	}

	items := make([]*proto.DTakoEvents, len(events))
	for i, event := range events {
		items[i] = dtakoEventsModelToProto(event)
	}

	return &proto.ListDTakoEventsResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByOperationNo 運行NOでイベント情報取得
func (s *DTakoEventsService) GetByOperationNo(ctx context.Context, req *proto.GetDTakoEventsByOperationNoRequest) (*proto.ListDTakoEventsResponse, error) {
	events, err := s.repo.GetByOperationNo(req.OperationNo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get events by operation_no: %v", err)
	}

	items := make([]*proto.DTakoEvents, len(events))
	for i, event := range events {
		items[i] = dtakoEventsModelToProto(event)
	}

	return &proto.ListDTakoEventsResponse{
		Items:      items,
		TotalCount: int32(len(events)),
	}, nil
}

// dtakoEventsModelToProto ModelからProtoへの変換
func dtakoEventsModelToProto(model *models.DTakoEvents) *proto.DTakoEvents {
	protoEvent := &proto.DTakoEvents{
		Id:              model.ID,
		OperationNo:     model.OperationNo,
		ReadDate:        model.ReadDate.Format("2006-01-02T15:04:05Z07:00"),
		CarCode:         int32(model.CarCode),
		CarCc:           model.CarCC,
		TargetDriverType: int32(model.TargetDriverType),
		DriverCode1:     int32(model.DriverCode1),
		TargetDriverCode: int32(model.TargetDriverCode),
		StartDatetime:   model.StartDatetime.Format("2006-01-02T15:04:05Z07:00"),
		EndDatetime:     model.EndDatetime.Format("2006-01-02T15:04:05Z07:00"),
		EventName:       model.EventName,
		StartMileage:    model.StartMileage,
		EndMileage:      model.EndMileage,
		SectionTime:     int32(model.SectionTime),
		SectionDistance: model.SectionDistance,
		StartCityName:   model.StartCityName,
		EndCityName:     model.EndCityName,
		StartPlaceName:  model.StartPlaceName,
		EndPlaceName:    model.EndPlaceName,
	}

	// Optional fields
	if model.EventCode != nil {
		val := int32(*model.EventCode)
		protoEvent.EventCode = &val
	}
	if model.StartCityCode != nil {
		val := int32(*model.StartCityCode)
		protoEvent.StartCityCode = &val
	}
	if model.EndCityCode != nil {
		val := int32(*model.EndCityCode)
		protoEvent.EndCityCode = &val
	}
	if model.StartPlaceCode != nil {
		val := int32(*model.StartPlaceCode)
		protoEvent.StartPlaceCode = &val
	}
	if model.EndPlaceCode != nil {
		val := int32(*model.EndPlaceCode)
		protoEvent.EndPlaceCode = &val
	}
	if model.StartGPSValid != nil {
		val := int32(*model.StartGPSValid)
		protoEvent.StartGpsValid = &val
	}
	if model.StartGPSLatitude != nil {
		protoEvent.StartGpsLatitude = model.StartGPSLatitude
	}
	if model.StartGPSLongitude != nil {
		protoEvent.StartGpsLongitude = model.StartGPSLongitude
	}
	if model.EndGPSValid != nil {
		val := int32(*model.EndGPSValid)
		protoEvent.EndGpsValid = &val
	}
	if model.EndGPSLatitude != nil {
		protoEvent.EndGpsLatitude = model.EndGPSLatitude
	}
	if model.EndGPSLongitude != nil {
		protoEvent.EndGpsLongitude = model.EndGPSLongitude
	}

	return protoEvent
}
