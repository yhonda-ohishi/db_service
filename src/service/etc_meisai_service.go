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

// ETCMeisaiService gRPCサービス実装
type ETCMeisaiService struct {
	proto.UnimplementedETCMeisaiServiceServer
	repo repository.ETCMeisaiRepository
}

// NewETCMeisaiService サービスのコンストラクタ
func NewETCMeisaiService(repo repository.ETCMeisaiRepository) *ETCMeisaiService {
	return &ETCMeisaiService{
		repo: repo,
	}
}

// Create ETC明細データ作成
func (s *ETCMeisaiService) Create(ctx context.Context, req *proto.CreateETCMeisaiRequest) (*proto.ETCMeisaiResponse, error) {
	if req.EtcMeisai == nil {
		return nil, status.Error(codes.InvalidArgument, "etc_meisai is required")
	}

	// ProtoからModelへ変換
	dateTo, _ := time.Parse(time.RFC3339, req.EtcMeisai.DateTo)
	dateToDate, _ := time.Parse("2006-01-02", req.EtcMeisai.DateToDate)

	model := &models.ETCMeisai{
		DateTo:     dateTo,
		DateToDate: dateToDate,
		IcFr:       derefString(req.EtcMeisai.IcFr),
		IcTo:       req.EtcMeisai.IcTo,
		Price:      req.EtcMeisai.Price,
		Shashu:     req.EtcMeisai.Shashu,
		EtcNum:     req.EtcMeisai.EtcNum,
	}

	// オプションフィールド
	if req.EtcMeisai.DateFr != nil && *req.EtcMeisai.DateFr != "" {
		t, _ := time.Parse(time.RFC3339, *req.EtcMeisai.DateFr)
		model.DateFr = &t
	}
	if req.EtcMeisai.PriceBf != nil {
		model.PriceBf = req.EtcMeisai.PriceBf
	}
	if req.EtcMeisai.Descount != nil {
		model.Descount = req.EtcMeisai.Descount
	}
	if req.EtcMeisai.CarIdNum != nil {
		model.CarIDNum = req.EtcMeisai.CarIdNum
	}
	if req.EtcMeisai.Detail != nil {
		model.Detail = req.EtcMeisai.Detail
	}
	// ハッシュが指定されていない場合は自動生成
	if req.EtcMeisai.Hash == "" {
		model.SetHash()
	} else {
		model.Hash = req.EtcMeisai.Hash
	}

	// リポジトリで作成
	if err := s.repo.Create(model); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create record: %v", err)
	}

	// ModelからProtoへ変換して返却
	return &proto.ETCMeisaiResponse{
		EtcMeisai: etcModelToProto(model),
	}, nil
}

// Get ETC明細データ取得
func (s *ETCMeisaiService) Get(ctx context.Context, req *proto.GetETCMeisaiRequest) (*proto.ETCMeisaiResponse, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	model, err := s.repo.GetByID(req.Id)
	if err != nil {
		if err == models.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "record not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get record: %v", err)
	}

	return &proto.ETCMeisaiResponse{
		EtcMeisai: etcModelToProto(model),
	}, nil
}

// Update ETC明細データ更新
func (s *ETCMeisaiService) Update(ctx context.Context, req *proto.UpdateETCMeisaiRequest) (*proto.ETCMeisaiResponse, error) {
	if req.EtcMeisai == nil {
		return nil, status.Error(codes.InvalidArgument, "etc_meisai is required")
	}

	// ProtoからModelへ変換
	model := etcProtoToModel(req.EtcMeisai)

	// リポジトリで更新
	if err := s.repo.Update(model); err != nil {
		if err == models.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "record not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update record: %v", err)
	}

	return &proto.ETCMeisaiResponse{
		EtcMeisai: etcModelToProto(model),
	}, nil
}

// Delete ETC明細データ削除
func (s *ETCMeisaiService) Delete(ctx context.Context, req *proto.DeleteETCMeisaiRequest) (*proto.Empty, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	if err := s.repo.DeleteByID(req.Id); err != nil {
		if err == models.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "record not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete record: %v", err)
	}

	return &proto.Empty{}, nil
}

// List ETC明細データ一覧取得
func (s *ETCMeisaiService) List(ctx context.Context, req *proto.ListETCMeisaiRequest) (*proto.ListETCMeisaiResponse, error) {
	params := &repository.ETCMeisaiListParams{
		Limit:  int(req.Limit),
		Offset: int(req.Offset),
	}

	// オプションパラメータ
	if req.Hash != nil {
		params.Hash = req.Hash
	}
	if req.StartDate != nil && *req.StartDate != "" {
		t, _ := time.Parse(time.RFC3339, *req.StartDate)
		params.StartDate = &t
	}
	if req.EndDate != nil && *req.EndDate != "" {
		t, _ := time.Parse(time.RFC3339, *req.EndDate)
		params.EndDate = &t
	}

	models, totalCount, err := s.repo.List(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list records: %v", err)
	}

	items := make([]*proto.ETCMeisai, len(models))
	for i, model := range models {
		items[i] = etcModelToProto(model)
	}

	return &proto.ListETCMeisaiResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// derefString *stringをstringに変換（nilの場合は空文字列）
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// stringPtr stringを*stringに変換（空文字列の場合はnil）
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// etcProtoToModel ProtoからModelへの変換
func etcProtoToModel(p *proto.ETCMeisai) *models.ETCMeisai {
	dateTo, _ := time.Parse(time.RFC3339, p.DateTo)
	dateToDate, _ := time.Parse("2006-01-02", p.DateToDate)

	m := &models.ETCMeisai{
		ID:         p.Id,
		DateTo:     dateTo,
		DateToDate: dateToDate,
		IcFr:       derefString(p.IcFr),
		IcTo:       p.IcTo,
		Price:      p.Price,
		Shashu:     p.Shashu,
		EtcNum:     p.EtcNum,
	}

	if p.DateFr != nil && *p.DateFr != "" {
		t, _ := time.Parse(time.RFC3339, *p.DateFr)
		m.DateFr = &t
	}
	if p.PriceBf != nil {
		m.PriceBf = p.PriceBf
	}
	if p.Descount != nil {
		m.Descount = p.Descount
	}
	if p.CarIdNum != nil {
		m.CarIDNum = p.CarIdNum
	}
	if p.Detail != nil {
		m.Detail = p.Detail
	}
	m.Hash = p.Hash

	return m
}

// etcModelToProto ModelからProtoへの変換
func etcModelToProto(m *models.ETCMeisai) *proto.ETCMeisai {
	p := &proto.ETCMeisai{
		Id:         m.ID,
		DateTo:     m.DateTo.Format(time.RFC3339),
		DateToDate: m.DateToDate.Format("2006-01-02"),
		IcFr:       stringPtr(m.IcFr),
		IcTo:       m.IcTo,
		Price:      m.Price,
		Shashu:     m.Shashu,
		EtcNum:     m.EtcNum,
	}

	if m.DateFr != nil {
		t := m.DateFr.Format(time.RFC3339)
		p.DateFr = &t
	}
	if m.PriceBf != nil {
		p.PriceBf = m.PriceBf
	}
	if m.Descount != nil {
		p.Descount = m.Descount
	}
	if m.CarIDNum != nil {
		p.CarIdNum = m.CarIDNum
	}
	if m.Detail != nil {
		p.Detail = m.Detail
	}
	p.Hash = m.Hash

	return p
}
