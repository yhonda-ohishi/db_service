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
	proto.UnimplementedDb_DTakoFerryRowsServiceServer
	repo repository.DTakoFerryRowsRepository
}

// NewDTakoFerryRowsService サービスのコンストラクタ
func NewDTakoFerryRowsService(repo repository.DTakoFerryRowsRepository) *DTakoFerryRowsService {
	return &DTakoFerryRowsService{
		repo: repo,
	}
}

// Create フェリー運行データ作成（スタブ実装）
func (s *DTakoFerryRowsService) Create(ctx context.Context, req *proto.Db_CreateDTakoFerryRowsRequest) (*proto.Db_DTakoFerryRowsResponse, error) {
	// 簡略化のため基本実装のみ
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// Get フェリー運行データ取得
func (s *DTakoFerryRowsService) Get(ctx context.Context, req *proto.Db_GetDTakoFerryRowsRequest) (*proto.Db_DTakoFerryRowsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// Update フェリー運行データ更新
func (s *DTakoFerryRowsService) Update(ctx context.Context, req *proto.Db_UpdateDTakoFerryRowsRequest) (*proto.Db_DTakoFerryRowsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// Delete フェリー運行データ削除
func (s *DTakoFerryRowsService) Delete(ctx context.Context, req *proto.Db_DeleteDTakoFerryRowsRequest) (*proto.Db_Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

// List フェリー運行データ一覧取得
func (s *DTakoFerryRowsService) List(ctx context.Context, req *proto.Db_ListDTakoFerryRowsRequest) (*proto.Db_ListDTakoFerryRowsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
