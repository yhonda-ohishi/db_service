package service

import (
	"context"

	pb "github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UntenNippoMeisaiService 運転日報明細サービス
type UntenNippoMeisaiService struct {
	pb.UnimplementedDb_UntenNippoMeisaiServiceServer
	repo repository.UntenNippoMeisaiRepository
}

// NewUntenNippoMeisaiService コンストラクタ
func NewUntenNippoMeisaiService(repo repository.UntenNippoMeisaiRepository) *UntenNippoMeisaiService {
	return &UntenNippoMeisaiService{
		repo: repo,
	}
}

// Get 単一の運転日報明細を取得（複合主キー: 日報K, 配車K, 車輌C）
func (s *UntenNippoMeisaiService) Get(ctx context.Context, req *pb.Db_GetUntenNippoMeisaiRequest) (*pb.Db_UntenNippoMeisaiResponse, error) {
	meisai, err := s.repo.GetByNippoK(req.NippoK, req.HaishaK, req.SharyoC)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "運転日報明細が見つかりません: %v", err)
	}

	return &pb.Db_UntenNippoMeisaiResponse{
		UntenNippoMeisai: convertUntenNippoMeisaiToProto(meisai),
	}, nil
}

// List 運転日報明細のリストを取得
func (s *UntenNippoMeisaiService) List(ctx context.Context, req *pb.Db_ListUntenNippoMeisaiRequest) (*pb.Db_ListUntenNippoMeisaiResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := int(req.Offset)

	orderBy := ""
	if req.OrderBy != nil {
		orderBy = *req.OrderBy
	}

	meisaiList, totalCount, err := s.repo.GetAll(limit, offset, orderBy)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "運転日報明細の取得に失敗しました: %v", err)
	}

	pbMeisaiList := make([]*pb.Db_UntenNippoMeisai, len(meisaiList))
	for i, meisai := range meisaiList {
		pbMeisaiList[i] = convertUntenNippoMeisaiToProto(meisai)
	}

	return &pb.Db_ListUntenNippoMeisaiResponse{
		Items:      pbMeisaiList,
		TotalCount: int32(totalCount),
	}, nil
}

// GetBySharyoC 車輌Cで運転日報明細を取得
func (s *UntenNippoMeisaiService) GetBySharyoC(ctx context.Context, req *pb.Db_GetUntenNippoMeisaiBySharyoCRequest) (*pb.Db_ListUntenNippoMeisaiResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}

	meisaiList, err := s.repo.GetBySharyoC(req.SharyoC, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "車輌Cでの運転日報明細の取得に失敗しました: %v", err)
	}

	pbMeisaiList := make([]*pb.Db_UntenNippoMeisai, len(meisaiList))
	for i, meisai := range meisaiList {
		pbMeisaiList[i] = convertUntenNippoMeisaiToProto(meisai)
	}

	return &pb.Db_ListUntenNippoMeisaiResponse{
		Items:      pbMeisaiList,
		TotalCount: int32(len(meisaiList)),
	}, nil
}

// GetByDateRange 日付範囲で運転日報明細を取得
func (s *UntenNippoMeisaiService) GetByDateRange(ctx context.Context, req *pb.Db_GetUntenNippoMeisaiByDateRangeRequest) (*pb.Db_ListUntenNippoMeisaiResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}
	offset := int(req.Offset)

	meisaiList, totalCount, err := s.repo.GetByDateRange(req.StartDate, req.EndDate, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "日付範囲での運転日報明細の取得に失敗しました: %v", err)
	}

	pbMeisaiList := make([]*pb.Db_UntenNippoMeisai, len(meisaiList))
	for i, meisai := range meisaiList {
		pbMeisaiList[i] = convertUntenNippoMeisaiToProto(meisai)
	}

	return &pb.Db_ListUntenNippoMeisaiResponse{
		Items:      pbMeisaiList,
		TotalCount: int32(totalCount),
	}, nil
}
