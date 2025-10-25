package service

import (
	"context"
	"time"

	"github.com/yhonda-ohishi/db_service/src/models/mysql"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TimeCardDevService gRPCサービス実装（ローカルDB、読み書き可能）
type TimeCardDevService struct {
	proto.UnimplementedDb_TimeCardDevServiceServer
	repo repository.TimeCardDevRepository
}

// NewTimeCardDevService TimeCardDevServiceのコンストラクタ
func NewTimeCardDevService(repo repository.TimeCardDevRepository) *TimeCardDevService {
	return &TimeCardDevService{
		repo: repo,
	}
}

// Create タイムカードデータ作成
func (s *TimeCardDevService) Create(ctx context.Context, req *proto.Db_CreateTimeCardRequest) (*proto.Db_TimeCardResponse, error) {
	if req.TimeCard == nil {
		return nil, status.Errorf(codes.InvalidArgument, "time_card is required")
	}

	// ProtoからModelへの変換
	timeCard, err := protoToTimeCardModel(req.TimeCard)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time_card: %v", err)
	}

	// 作成
	if err := s.repo.Create(timeCard); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create time_card: %v", err)
	}

	return &proto.Db_TimeCardResponse{
		TimeCard: timeCardModelToProto(timeCard),
	}, nil
}

// Get タイムカードデータ取得（複合主キー）
func (s *TimeCardDevService) Get(ctx context.Context, req *proto.Db_GetTimeCardRequest) (*proto.Db_TimeCardResponse, error) {
	// RFC3339形式の文字列をパース
	datetime, err := time.Parse(time.RFC3339, req.Datetime)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid datetime format: %v", err)
	}

	timeCard, err := s.repo.GetByCompositeKey(datetime, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "time_card not found: %v", err)
	}

	return &proto.Db_TimeCardResponse{
		TimeCard: timeCardModelToProto(timeCard),
	}, nil
}

// Update タイムカードデータ更新
func (s *TimeCardDevService) Update(ctx context.Context, req *proto.Db_UpdateTimeCardRequest) (*proto.Db_TimeCardResponse, error) {
	if req.TimeCard == nil {
		return nil, status.Errorf(codes.InvalidArgument, "time_card is required")
	}

	// ProtoからModelへの変換
	timeCard, err := protoToTimeCardModel(req.TimeCard)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid time_card: %v", err)
	}

	// 更新
	if err := s.repo.Update(timeCard); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update time_card: %v", err)
	}

	return &proto.Db_TimeCardResponse{
		TimeCard: timeCardModelToProto(timeCard),
	}, nil
}

// Delete タイムカードデータ削除
func (s *TimeCardDevService) Delete(ctx context.Context, req *proto.Db_DeleteTimeCardRequest) (*proto.Db_Empty, error) {
	// RFC3339形式の文字列をパース
	datetime, err := time.Parse(time.RFC3339, req.Datetime)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid datetime format: %v", err)
	}

	if err := s.repo.Delete(datetime, int(req.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete time_card: %v", err)
	}

	return &proto.Db_Empty{}, nil
}

// List タイムカードデータ一覧取得
func (s *TimeCardDevService) List(ctx context.Context, req *proto.Db_ListTimeCardRequest) (*proto.Db_ListTimeCardResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 100
	}
	offset := int(req.Offset)
	orderBy := ""
	if req.OrderBy != nil {
		orderBy = *req.OrderBy
	}

	timeCards, totalCount, err := s.repo.GetAll(limit, offset, orderBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list time_cards: %v", err)
	}

	items := make([]*proto.Db_TimeCard, len(timeCards))
	for i, timeCard := range timeCards {
		items[i] = timeCardModelToProto(timeCard)
	}

	return &proto.Db_ListTimeCardResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// protoToTimeCardModel ProtoからModelへの変換
func protoToTimeCardModel(pb *proto.Db_TimeCard) (*mysql.TimeCard, error) {
	datetime, err := time.Parse(time.RFC3339, pb.Datetime)
	if err != nil {
		return nil, err
	}
	created, err := time.Parse(time.RFC3339, pb.Created)
	if err != nil {
		return nil, err
	}
	modified, err := time.Parse(time.RFC3339, pb.Modified)
	if err != nil {
		return nil, err
	}

	return &mysql.TimeCard{
		Datetime:    datetime,
		ID:          int(pb.Id),
		MachineIP:   pb.MachineIp,
		State:       pb.State,
		StateDetail: pb.StateDetail,
		Created:     created,
		Modified:    modified,
	}, nil
}
