package service

import (
	"context"

	pb "github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ChiikiMasterService 地域マスタサービス
type ChiikiMasterService struct {
	pb.UnimplementedChiikiMasterServiceServer
	repo repository.ChiikiMasterRepository
}

// NewChiikiMasterService コンストラクタ
func NewChiikiMasterService(repo repository.ChiikiMasterRepository) *ChiikiMasterService {
	return &ChiikiMasterService{
		repo: repo,
	}
}

// Get 単一の地域マスタを取得
func (s *ChiikiMasterService) Get(ctx context.Context, req *pb.GetChiikiMasterRequest) (*pb.ChiikiMasterResponse, error) {
	chiiki, err := s.repo.GetByChiikiC(req.ChiikiC)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "地域マスタが見つかりません: %v", err)
	}

	return &pb.ChiikiMasterResponse{
		ChiikiMaster: convertChiikiMasterToProto(chiiki),
	}, nil
}

// List 地域マスタのリストを取得
func (s *ChiikiMasterService) List(ctx context.Context, req *pb.ListChiikiMasterRequest) (*pb.ListChiikiMasterResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := int(req.Offset)

	orderBy := ""
	if req.OrderBy != nil {
		orderBy = *req.OrderBy
	}

	chiikiList, totalCount, err := s.repo.GetAll(limit, offset, orderBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "地域マスタの取得に失敗しました: %v", err)
	}

	pbChiikiList := make([]*pb.ChiikiMaster, len(chiikiList))
	for i, chiiki := range chiikiList {
		pbChiikiList[i] = convertChiikiMasterToProto(chiiki)
	}

	return &pb.ListChiikiMasterResponse{
		Items:      pbChiikiList,
		TotalCount: int32(totalCount),
	}, nil
}
