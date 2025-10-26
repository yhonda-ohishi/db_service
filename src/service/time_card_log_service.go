package service

import (
	"context"
	"time"

	"github.com/yhonda-ohishi/db_service/src/models/mysql"
	proto "github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TimeCardLogService タイムカードログサービス（ローカルDB、読み書き可能）
type TimeCardLogService struct {
	proto.UnimplementedDb_TimeCardLogServiceServer
	repo repository.TimeCardLogRepository
}

// NewTimeCardLogService TimeCardLogServiceのコンストラクタ
func NewTimeCardLogService(repo repository.TimeCardLogRepository) *TimeCardLogService {
	return &TimeCardLogService{
		repo: repo,
	}
}

// Create タイムカードログ作成
func (s *TimeCardLogService) Create(ctx context.Context, req *proto.Db_CreateTimeCardLogRequest) (*proto.Db_TimeCardLogResponse, error) {
	log, err := protoToTimeCardLogModel(req.Log)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid log: %v", err)
	}

	if err := s.repo.Create(log); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create log: %v", err)
	}

	return &proto.Db_TimeCardLogResponse{
		Log: timeCardLogModelToProto(log),
	}, nil
}

// Get タイムカードログ取得（複合主キー）
func (s *TimeCardLogService) Get(ctx context.Context, req *proto.Db_GetTimeCardLogRequest) (*proto.Db_TimeCardLogResponse, error) {
	log, err := s.repo.GetByCompositeKey(req.Datetime, int(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "log not found: %v", err)
	}

	return &proto.Db_TimeCardLogResponse{
		Log: timeCardLogModelToProto(log),
	}, nil
}

// Update タイムカードログ更新
func (s *TimeCardLogService) Update(ctx context.Context, req *proto.Db_UpdateTimeCardLogRequest) (*proto.Db_TimeCardLogResponse, error) {
	log, err := protoToTimeCardLogModel(req.Log)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid log: %v", err)
	}

	if err := s.repo.Update(log); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update log: %v", err)
	}

	return &proto.Db_TimeCardLogResponse{
		Log: timeCardLogModelToProto(log),
	}, nil
}

// Delete タイムカードログ削除
func (s *TimeCardLogService) Delete(ctx context.Context, req *proto.Db_DeleteTimeCardLogRequest) (*proto.Db_Empty, error) {
	if err := s.repo.Delete(req.Datetime, int(req.Id)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete log: %v", err)
	}

	return &proto.Db_Empty{}, nil
}

// List タイムカードログ一覧取得
func (s *TimeCardLogService) List(ctx context.Context, req *proto.Db_ListTimeCardLogRequest) (*proto.Db_ListTimeCardLogResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)
	orderBy := ""
	if req.OrderBy != nil {
		orderBy = *req.OrderBy
	}

	logs, totalCount, err := s.repo.GetAll(limit, offset, orderBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list logs: %v", err)
	}

	items := make([]*proto.Db_TimeCardLog, len(logs))
	for i, log := range logs {
		items[i] = timeCardLogModelToProto(log)
	}

	return &proto.Db_ListTimeCardLogResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByCardID カードIDでタイムカードログ取得
func (s *TimeCardLogService) GetByCardID(ctx context.Context, req *proto.Db_GetByCardIDRequest) (*proto.Db_ListTimeCardLogResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)

	logs, totalCount, err := s.repo.GetByCardID(req.CardId, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get logs by card_id: %v", err)
	}

	items := make([]*proto.Db_TimeCardLog, len(logs))
	for i, log := range logs {
		items[i] = timeCardLogModelToProto(log)
	}

	return &proto.Db_ListTimeCardLogResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// protoToTimeCardLogModel protoからモデルへの変換
func protoToTimeCardLogModel(pb *proto.Db_TimeCardLog) (*mysql.TimeCardLog, error) {
	created, err := time.Parse(time.RFC3339, pb.Created)
	if err != nil {
		return nil, err
	}

	modified, err := time.Parse(time.RFC3339, pb.Modified)
	if err != nil {
		return nil, err
	}

	return &mysql.TimeCardLog{
		Datetime:    pb.Datetime,
		ID:          int(pb.Id),
		CardID:      pb.CardId,
		MachineIP:   pb.MachineIp,
		State:       pb.State,
		StateDetail: pb.StateDetail,
		Created:     created,
		Modified:    modified,
	}, nil
}

// timeCardLogModelToProto モデルからprotoへの変換
func timeCardLogModelToProto(log *mysql.TimeCardLog) *proto.Db_TimeCardLog {
	return &proto.Db_TimeCardLog{
		Datetime:    log.Datetime,
		Id:          int32(log.ID),
		CardId:      log.CardID,
		MachineIp:   log.MachineIP,
		State:       log.State,
		StateDetail: log.StateDetail,
		Created:     log.Created.Format(time.RFC3339),
		Modified:    log.Modified.Format(time.RFC3339),
	}
}
