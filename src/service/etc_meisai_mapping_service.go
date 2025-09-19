package service

import (
	"context"
	"time"

	"github.com/yhonda-ohishi/db_service/src/models"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ETCMeisaiMappingService ETC明細マッピングサービス実装
type ETCMeisaiMappingService struct {
	proto.UnimplementedETCMeisaiMappingServiceServer
	repo repository.ETCMeisaiMappingRepository
}

// NewETCMeisaiMappingService サービスのコンストラクタ
func NewETCMeisaiMappingService(repo repository.ETCMeisaiMappingRepository) *ETCMeisaiMappingService {
	return &ETCMeisaiMappingService{
		repo: repo,
	}
}

// Create マッピング作成
func (s *ETCMeisaiMappingService) Create(ctx context.Context, req *proto.CreateETCMeisaiMappingRequest) (*proto.ETCMeisaiMappingResponse, error) {
	if req.EtcMeisaiMapping == nil {
		return nil, status.Error(codes.InvalidArgument, "etc_meisai_mapping is required")
	}

	// プロトコルバッファーからモデルに変換
	model := etcMeisaiMappingProtoToModel(req.EtcMeisaiMapping)

	// タイムスタンプ設定
	model.BeforeCreate()

	// リポジトリで作成
	if err := s.repo.Create(model); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create mapping: %v", err)
	}

	// レスポンス作成
	return &proto.ETCMeisaiMappingResponse{
		EtcMeisaiMapping: etcMeisaiMappingModelToProto(model),
	}, nil
}

// Get マッピング取得
func (s *ETCMeisaiMappingService) Get(ctx context.Context, req *proto.GetETCMeisaiMappingRequest) (*proto.ETCMeisaiMappingResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// リポジトリから取得
	model, err := s.repo.GetByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "mapping not found: %v", err)
	}

	// レスポンス作成
	return &proto.ETCMeisaiMappingResponse{
		EtcMeisaiMapping: etcMeisaiMappingModelToProto(model),
	}, nil
}

// Update マッピング更新
func (s *ETCMeisaiMappingService) Update(ctx context.Context, req *proto.UpdateETCMeisaiMappingRequest) (*proto.ETCMeisaiMappingResponse, error) {
	if req.EtcMeisaiMapping == nil {
		return nil, status.Error(codes.InvalidArgument, "etc_meisai_mapping is required")
	}
	if req.EtcMeisaiMapping.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// プロトコルバッファーからモデルに変換
	model := etcMeisaiMappingProtoToModel(req.EtcMeisaiMapping)

	// 更新前フック
	model.BeforeUpdate()

	// リポジトリで更新
	if err := s.repo.Update(model); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update mapping: %v", err)
	}

	// レスポンス作成
	return &proto.ETCMeisaiMappingResponse{
		EtcMeisaiMapping: etcMeisaiMappingModelToProto(model),
	}, nil
}

// Delete マッピング削除
func (s *ETCMeisaiMappingService) Delete(ctx context.Context, req *proto.DeleteETCMeisaiMappingRequest) (*proto.Empty, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// リポジトリで削除
	if err := s.repo.DeleteByID(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete mapping: %v", err)
	}

	return &proto.Empty{}, nil
}

// List マッピング一覧取得
func (s *ETCMeisaiMappingService) List(ctx context.Context, req *proto.ListETCMeisaiMappingRequest) (*proto.ListETCMeisaiMappingResponse, error) {
	if req.Limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be positive")
	}
	if req.Offset < 0 {
		return nil, status.Error(codes.InvalidArgument, "offset must be non-negative")
	}

	// パラメータ作成
	params := &repository.ETCMeisaiMappingListParams{
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}

	// オプションパラメータ
	if req.EtcMeisaiHash != nil {
		params.ETCMeisaiHash = req.EtcMeisaiHash
	}
	if req.DtakoRowId != nil {
		params.DTakoRowID = req.DtakoRowId
	}

	// リポジトリから取得
	models, totalCount, err := s.repo.List(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list mappings: %v", err)
	}

	// プロトコルバッファーに変換
	items := make([]*proto.ETCMeisaiMapping, len(models))
	for i, model := range models {
		items[i] = etcMeisaiMappingModelToProto(model)
	}

	return &proto.ListETCMeisaiMappingResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetDTakoRowIDByHash ハッシュからDTakoRowIDを取得
func (s *ETCMeisaiMappingService) GetDTakoRowIDByHash(ctx context.Context, req *proto.GetDTakoRowIDByHashRequest) (*proto.GetDTakoRowIDByHashResponse, error) {
	if req.EtcMeisaiHash == "" {
		return nil, status.Error(codes.InvalidArgument, "etc_meisai_hash is required")
	}

	// リポジトリから取得
	dtakoRowIDs, err := s.repo.GetDTakoRowIDsByHash(req.EtcMeisaiHash)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get dtako_row_ids: %v", err)
	}

	return &proto.GetDTakoRowIDByHashResponse{
		DtakoRowIds: dtakoRowIDs,
	}, nil
}

// etcMeisaiMappingProtoToModel プロトコルバッファーからモデルに変換
func etcMeisaiMappingProtoToModel(p *proto.ETCMeisaiMapping) *models.ETCMeisaiMapping {
	m := &models.ETCMeisaiMapping{
		ID:            p.Id,
		ETCMeisaiHash: p.EtcMeisaiHash,
		DTakoRowID:    p.DtakoRowId,
		CreatedBy:     p.CreatedBy,
	}

	if p.CreatedAt != "" {
		if t, err := time.Parse(time.RFC3339, p.CreatedAt); err == nil {
			m.CreatedAt = t
		}
	}
	if p.UpdatedAt != "" {
		if t, err := time.Parse(time.RFC3339, p.UpdatedAt); err == nil {
			m.UpdatedAt = t
		}
	}
	if p.Notes != nil {
		m.Notes = p.Notes
	}

	return m
}

// etcMeisaiMappingModelToProto モデルからプロトコルバッファーに変換
func etcMeisaiMappingModelToProto(m *models.ETCMeisaiMapping) *proto.ETCMeisaiMapping {
	p := &proto.ETCMeisaiMapping{
		Id:            m.ID,
		EtcMeisaiHash: m.ETCMeisaiHash,
		DtakoRowId:    m.DTakoRowID,
		CreatedAt:     m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     m.UpdatedAt.Format(time.RFC3339),
		CreatedBy:     m.CreatedBy,
	}

	if m.Notes != nil {
		p.Notes = m.Notes
	}

	return p
}