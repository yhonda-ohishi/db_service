package service

import (
	"context"

	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DTakoFerryRowsService gRPCサービス実装
type DTakoFerryRowsService struct {
	proto.UnimplementedDTakoFerryRowsServiceServer
	repo repository.DTakoFerryRowsRepository
}

// NewDTakoFerryRowsService サービスのコンストラクタ
func NewDTakoFerryRowsService(repo repository.DTakoFerryRowsRepository) *DTakoFerryRowsService {
	return &DTakoFerryRowsService{
		repo: repo,
	}
}

// Create フェリー運行データ作成（スタブ実装）
func (s *DTakoFerryRowsService) Create(ctx context.Context, req *proto.CreateDTakoFerryRowsRequest) (*proto.DTakoFerryRowsResponse, error) {
	// 簡略化のため基本実装のみ
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// Get フェリー運行データ取得
func (s *DTakoFerryRowsService) Get(ctx context.Context, req *proto.GetDTakoFerryRowsRequest) (*proto.DTakoFerryRowsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// Update フェリー運行データ更新
func (s *DTakoFerryRowsService) Update(ctx context.Context, req *proto.UpdateDTakoFerryRowsRequest) (*proto.DTakoFerryRowsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// Delete フェリー運行データ削除
func (s *DTakoFerryRowsService) Delete(ctx context.Context, req *proto.DeleteDTakoFerryRowsRequest) (*proto.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// List フェリー運行データ一覧取得
func (s *DTakoFerryRowsService) List(ctx context.Context, req *proto.ListDTakoFerryRowsRequest) (*proto.ListDTakoFerryRowsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
