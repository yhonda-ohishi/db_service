package service

import (
	"context"

	pb "github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShainMasterService 社員マスタサービス
type ShainMasterService struct {
	pb.UnimplementedDb_ShainMasterServiceServer
	repo repository.ShainMasterRepository
}

// NewShainMasterService コンストラクタ
func NewShainMasterService(repo repository.ShainMasterRepository) *ShainMasterService {
	return &ShainMasterService{
		repo: repo,
	}
}

// Get 単一の社員マスタを取得
func (s *ShainMasterService) Get(ctx context.Context, req *pb.Db_GetShainMasterRequest) (*pb.Db_ShainMasterResponse, error) {
	shain, err := s.repo.GetByShainC(req.ShainC)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "社員マスタが見つかりません: %v", err)
	}

	return &pb.Db_ShainMasterResponse{
		ShainMaster: convertShainMasterToProto(shain),
	}, nil
}

// List 社員マスタのリストを取得
func (s *ShainMasterService) List(ctx context.Context, req *pb.Db_ListShainMasterRequest) (*pb.Db_ListShainMasterResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := int(req.Offset)

	orderBy := ""
	if req.OrderBy != nil {
		orderBy = *req.OrderBy
	}

	shainList, totalCount, err := s.repo.GetAll(limit, offset, orderBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "社員マスタの取得に失敗しました: %v", err)
	}

	pbShainList := make([]*pb.Db_ShainMaster, len(shainList))
	for i, shain := range shainList {
		pbShainList[i] = convertShainMasterToProto(shain)
	}

	return &pb.Db_ListShainMasterResponse{
		Items:      pbShainList,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByBumonC 部門Cで社員マスタを取得
func (s *ShainMasterService) GetByBumonC(ctx context.Context, req *pb.Db_GetShainMasterByBumonCRequest) (*pb.Db_ListShainMasterResponse, error) {
	shainList, err := s.repo.GetByBumonC(req.BumonC)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "部門Cでの社員マスタの取得に失敗しました: %v", err)
	}

	pbShainList := make([]*pb.Db_ShainMaster, len(shainList))
	for i, shain := range shainList {
		pbShainList[i] = convertShainMasterToProto(shain)
	}

	return &pb.Db_ListShainMasterResponse{
		Items:      pbShainList,
		TotalCount: int32(len(shainList)),
	}, nil
}
