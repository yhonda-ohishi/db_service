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

// DTakoUriageKeihiService gRPCサービス実装
type DTakoUriageKeihiService struct {
	proto.UnimplementedDb_DTakoUriageKeihiServiceServer
	repo repository.DTakoUriageKeihiRepository
}

// NewDTakoUriageKeihiService サービスのコンストラクタ
func NewDTakoUriageKeihiService(repo repository.DTakoUriageKeihiRepository) *DTakoUriageKeihiService {
	return &DTakoUriageKeihiService{
		repo: repo,
	}
}

// Create 経費精算データ作成
func (s *DTakoUriageKeihiService) Create(ctx context.Context, req *proto.Db_CreateDTakoUriageKeihiRequest) (*proto.Db_DTakoUriageKeihiResponse, error) {
	if req.DtakoUriageKeihi == nil {
		return nil, status.Error(codes.InvalidArgument, "dtako_uriage_keihi is required")
	}

	// ProtoからModelへ変換
	model := protoToModel(req.DtakoUriageKeihi)

	// リポジトリで作成
	if err := s.repo.Create(model); err != nil {
		if err == mysql.ErrDuplicateKey {
			return nil, status.Error(codes.AlreadyExists, "record already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to create record: %v", err)
	}

	// ModelからProtoへ変換して返却
	return &proto.Db_DTakoUriageKeihiResponse{
		DtakoUriageKeihi: modelToProto(model),
	}, nil
}

// Get 経費精算データ取得
func (s *DTakoUriageKeihiService) Get(ctx context.Context, req *proto.Db_GetDTakoUriageKeihiRequest) (*proto.Db_DTakoUriageKeihiResponse, error) {
	// バリデーション
	if req.SrchId == "" || req.Datetime == "" {
		return nil, status.Error(codes.InvalidArgument, "srch_id, datetime, and keihi_c are required")
	}

	// 日時パース
	datetime, err := time.Parse(time.RFC3339, req.Datetime)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid datetime format: %v", err)
	}

	// リポジトリから取得
	model, err := s.repo.GetByCompositeKey(req.SrchId, datetime, req.KeihiC)
	if err != nil {
		if err == mysql.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "record not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get record: %v", err)
	}

	return &proto.Db_DTakoUriageKeihiResponse{
		DtakoUriageKeihi: modelToProto(model),
	}, nil
}

// Update 経費精算データ更新
func (s *DTakoUriageKeihiService) Update(ctx context.Context, req *proto.Db_UpdateDTakoUriageKeihiRequest) (*proto.Db_DTakoUriageKeihiResponse, error) {
	if req.DtakoUriageKeihi == nil {
		return nil, status.Error(codes.InvalidArgument, "dtako_uriage_keihi is required")
	}

	// ProtoからModelへ変換
	model := protoToModel(req.DtakoUriageKeihi)

	// リポジトリで更新
	if err := s.repo.Update(model); err != nil {
		if err == mysql.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "record not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update record: %v", err)
	}

	return &proto.Db_DTakoUriageKeihiResponse{
		DtakoUriageKeihi: modelToProto(model),
	}, nil
}

// Delete 経費精算データ削除
func (s *DTakoUriageKeihiService) Delete(ctx context.Context, req *proto.Db_DeleteDTakoUriageKeihiRequest) (*proto.Db_Empty, error) {
	// バリデーション
	if req.SrchId == "" || req.Datetime == "" {
		return nil, status.Error(codes.InvalidArgument, "srch_id, datetime, and keihi_c are required")
	}

	// 日時パース
	datetime, err := time.Parse(time.RFC3339, req.Datetime)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid datetime format: %v", err)
	}

	// リポジトリから削除
	if err := s.repo.DeleteByCompositeKey(req.SrchId, datetime, req.KeihiC); err != nil {
		if err == mysql.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "record not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete record: %v", err)
	}

	return &proto.Db_Empty{}, nil
}

// List 経費精算データ一覧取得
func (s *DTakoUriageKeihiService) List(ctx context.Context, req *proto.Db_ListDTakoUriageKeihiRequest) (*proto.Db_ListDTakoUriageKeihiResponse, error) {
	params := &repository.ListParams{
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}

	// オプションパラメータの設定
	if req.DtakoRowId != nil && *req.DtakoRowId != "" {
		params.DtakoRowID = req.DtakoRowId
	}
	if req.StartDate != nil && *req.StartDate != "" {
		t, err := time.Parse(time.RFC3339, *req.StartDate)
		if err == nil {
			params.StartDate = &t
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := time.Parse(time.RFC3339, *req.EndDate)
		if err == nil {
			params.EndDate = &t
		}
	}

	// リポジトリから取得
	models, totalCount, err := s.repo.List(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list records: %v", err)
	}

	// ModelからProtoへ変換
	items := make([]*proto.Db_DTakoUriageKeihi, len(models))
	for i, model := range models {
		items[i] = modelToProto(model)
	}

	return &proto.Db_ListDTakoUriageKeihiResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// protoToModel ProtoからModelへの変換
func protoToModel(p *proto.Db_DTakoUriageKeihi) *mysql.DTakoUriageKeihi {
	datetime, _ := time.Parse(time.RFC3339, p.Datetime)

	m := &mysql.DTakoUriageKeihi{
		SrchID:      p.SrchId,
		Datetime:    datetime,
		KeihiC:      p.KeihiC,
		Price:       p.Price,
		DtakoRowID:  p.DtakoRowId,
		DtakoRowIDR: p.DtakoRowIdR,
	}

	// オプションフィールド
	if p.Km != nil {
		m.Km = p.Km
	}
	if p.StartSrchId != nil && *p.StartSrchId != "" {
		m.StartSrchID = p.StartSrchId
	}
	if p.StartSrchTime != nil && *p.StartSrchTime != "" {
		t, _ := time.Parse(time.RFC3339, *p.StartSrchTime)
		m.StartSrchTime = &t
	}
	if p.StartSrchPlace != nil {
		m.StartSrchPlace = p.StartSrchPlace
	}
	if p.StartSrchTokui != nil {
		m.StartSrchTokui = p.StartSrchTokui
	}
	if p.EndSrchId != nil {
		m.EndSrchID = p.EndSrchId
	}
	if p.EndSrchTime != nil && *p.EndSrchTime != "" {
		t, _ := time.Parse(time.RFC3339, *p.EndSrchTime)
		m.EndSrchTime = &t
	}
	if p.EndSrchPlace != nil {
		m.EndSrchPlace = p.EndSrchPlace
	}
	if p.Manual != nil {
		m.Manual = p.Manual
	}

	return m
}

// modelToProto ModelからProtoへの変換
func modelToProto(m *mysql.DTakoUriageKeihi) *proto.Db_DTakoUriageKeihi {
	p := &proto.Db_DTakoUriageKeihi{
		SrchId:      m.SrchID,
		Datetime:    m.Datetime.Format(time.RFC3339),
		KeihiC:      m.KeihiC,
		Price:       m.Price,
		DtakoRowId:  m.DtakoRowID,
		DtakoRowIdR: m.DtakoRowIDR,
	}

	// オプションフィールド
	if m.Km != nil {
		p.Km = m.Km
	}
	if m.StartSrchID != nil {
		p.StartSrchId = m.StartSrchID
	}
	if m.StartSrchTime != nil {
		t := m.StartSrchTime.Format(time.RFC3339)
		p.StartSrchTime = &t
	}
	if m.StartSrchPlace != nil {
		p.StartSrchPlace = m.StartSrchPlace
	}
	if m.StartSrchTokui != nil {
		p.StartSrchTokui = m.StartSrchTokui
	}
	if m.EndSrchID != nil {
		p.EndSrchId = m.EndSrchID
	}
	if m.EndSrchTime != nil {
		t := m.EndSrchTime.Format(time.RFC3339)
		p.EndSrchTime = &t
	}
	if m.EndSrchPlace != nil {
		p.EndSrchPlace = m.EndSrchPlace
	}
	if m.Manual != nil {
		p.Manual = m.Manual
	}

	return p
}
