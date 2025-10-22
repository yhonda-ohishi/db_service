package service

import (
	"context"

	"github.com/yhonda-ohishi/db_service/src/models/mysql"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ETCNumService gRPCサービス実装（本番DB、読み取り専用）
type ETCNumService struct {
	proto.UnimplementedETCNumServiceServer
	repo repository.ETCNumRepository
}

// NewETCNumService サービスのコンストラクタ
func NewETCNumService(repo repository.ETCNumRepository) *ETCNumService {
	return &ETCNumService{
		repo: repo,
	}
}

// List ETCカード番号一覧取得
func (s *ETCNumService) List(ctx context.Context, req *proto.ListETCNumRequest) (*proto.ListETCNumResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)

	if limit == 0 {
		limit = 100
	}

	etcNums, totalCount, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list etc_num: %v", err)
	}

	items := make([]*proto.ETCNum, len(etcNums))
	for i, etcNum := range etcNums {
		items[i] = etcNumModelToProto(etcNum)
	}

	return &proto.ListETCNumResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByETCCardNum ETCカード番号で取得
func (s *ETCNumService) GetByETCCardNum(ctx context.Context, req *proto.GetETCNumByETCCardNumRequest) (*proto.ListETCNumResponse, error) {
	etcNums, err := s.repo.GetByETCCardNum(req.EtcCardNum)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get etc_num by etc_card_num: %v", err)
	}

	items := make([]*proto.ETCNum, len(etcNums))
	for i, etcNum := range etcNums {
		items[i] = etcNumModelToProto(etcNum)
	}

	return &proto.ListETCNumResponse{
		Items:      items,
		TotalCount: int32(len(etcNums)),
	}, nil
}

// GetByCarID 車輌IDで取得
func (s *ETCNumService) GetByCarID(ctx context.Context, req *proto.GetETCNumByCarIDRequest) (*proto.ListETCNumResponse, error) {
	etcNums, err := s.repo.GetByCarID(req.CarId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get etc_num by car_id: %v", err)
	}

	items := make([]*proto.ETCNum, len(etcNums))
	for i, etcNum := range etcNums {
		items[i] = etcNumModelToProto(etcNum)
	}

	return &proto.ListETCNumResponse{
		Items:      items,
		TotalCount: int32(len(etcNums)),
	}, nil
}

// etcNumModelToProto ModelからProtoへの変換
func etcNumModelToProto(model *mysql.ETCNum) *proto.ETCNum {
	protoETCNum := &proto.ETCNum{
		EtcCardNum: model.ETCCardNum,
		CarId:      model.CarID,
	}

	// Optional fields
	if model.StartDateTime != nil {
		val := model.StartDateTime.Format("2006-01-02T15:04:05Z07:00")
		protoETCNum.StartDateTime = &val
	}
	if model.DueDateTime != nil {
		val := model.DueDateTime.Format("2006-01-02T15:04:05Z07:00")
		protoETCNum.DueDateTime = &val
	}
	if model.ToChange != nil {
		protoETCNum.ToChange = model.ToChange
	}

	return protoETCNum
}
