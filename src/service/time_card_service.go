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

// TimeCardService gRPCサービス実装（本番DB、読み取り専用）
type TimeCardService struct {
	proto.UnimplementedDb_TimeCardServiceServer
	repo repository.TimeCardRepository
}

// NewTimeCardService TimeCardServiceのコンストラクタ
func NewTimeCardService(repo repository.TimeCardRepository) *TimeCardService {
	return &TimeCardService{
		repo: repo,
	}
}

// Get タイムカードデータ取得（複合主キー）
func (s *TimeCardService) Get(ctx context.Context, req *proto.Db_GetTimeCardRequest) (*proto.Db_TimeCardResponse, error) {
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

// List タイムカードデータ一覧取得
func (s *TimeCardService) List(ctx context.Context, req *proto.Db_ListTimeCardRequest) (*proto.Db_ListTimeCardResponse, error) {
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

// timeCardModelToProto ModelからProtoへの変換
func timeCardModelToProto(model *mysql.TimeCard) *proto.Db_TimeCard {
	return &proto.Db_TimeCard{
		Datetime:    model.Datetime.Format(time.RFC3339),
		Id:          int32(model.ID),
		MachineIp:   model.MachineIP,
		State:       model.State,
		StateDetail: model.StateDetail,
		Created:     model.Created.Format(time.RFC3339),
		Modified:    model.Modified.Format(time.RFC3339),
	}
}
