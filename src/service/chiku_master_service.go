package service

import (
	"context"

	pb "github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ChikuMasterService 地区マスタサービス
type ChikuMasterService struct {
	pb.UnimplementedChikuMasterServiceServer
	repo repository.ChikuMasterRepository
}

// NewChikuMasterService コンストラクタ
func NewChikuMasterService(repo repository.ChikuMasterRepository) *ChikuMasterService {
	return &ChikuMasterService{
		repo: repo,
	}
}

// Get 単一の地区マスタを取得
func (s *ChikuMasterService) Get(ctx context.Context, req *pb.GetChikuMasterRequest) (*pb.ChikuMasterResponse, error) {
	chiku, err := s.repo.GetByChikuC(req.ChikuC)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "地区マスタが見つかりません: %v", err)
	}

	return &pb.ChikuMasterResponse{
		ChikuMaster: convertChikuMasterToProto(chiku),
	}, nil
}

// List 地区マスタのリストを取得
func (s *ChikuMasterService) List(ctx context.Context, req *pb.ListChikuMasterRequest) (*pb.ListChikuMasterResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := int(req.Offset)

	orderBy := ""
	if req.OrderBy != nil {
		orderBy = *req.OrderBy
	}

	chikuList, totalCount, err := s.repo.GetAll(limit, offset, orderBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "地区マスタの取得に失敗しました: %v", err)
	}

	pbChikuList := make([]*pb.ChikuMaster, len(chikuList))
	for i, chiku := range chikuList {
		pbChikuList[i] = convertChikuMasterToProto(chiku)
	}

	return &pb.ListChikuMasterResponse{
		Items:      pbChikuList,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByChiikiC 地域Cで地区マスタを取得
func (s *ChikuMasterService) GetByChiikiC(ctx context.Context, req *pb.GetChikuMasterByChiikiCRequest) (*pb.ListChikuMasterResponse, error) {
	chikuList, err := s.repo.GetByChiikiC(req.ChiikiC)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "地域Cでの地区マスタの取得に失敗しました: %v", err)
	}

	pbChikuList := make([]*pb.ChikuMaster, len(chikuList))
	for i, chiku := range chikuList {
		pbChikuList[i] = convertChikuMasterToProto(chiku)
	}

	return &pb.ListChikuMasterResponse{
		Items:      pbChikuList,
		TotalCount: int32(len(chikuList)),
	}, nil
}
