package service

import (
	"context"

	"github.com/yhonda-ohishi/db_service/src/models"
	"github.com/yhonda-ohishi/db_service/src/proto"
	"github.com/yhonda-ohishi/db_service/src/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DTakoFerryRowsProdService gRPCサービス実装（本番DB、読み取り専用）
type DTakoFerryRowsProdService struct {
	proto.UnimplementedDTakoFerryRowsProdServiceServer
	repo repository.DTakoFerryRowsProdRepository
}

// NewDTakoFerryRowsProdService サービスのコンストラクタ
func NewDTakoFerryRowsProdService(repo repository.DTakoFerryRowsProdRepository) *DTakoFerryRowsProdService {
	return &DTakoFerryRowsProdService{
		repo: repo,
	}
}

// Get フェリー運行データ取得
func (s *DTakoFerryRowsProdService) Get(ctx context.Context, req *proto.GetDTakoFerryRowsProdRequest) (*proto.DTakoFerryRowsProdResponse, error) {
	row, err := s.repo.GetByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "ferry row not found: %v", err)
	}

	return &proto.DTakoFerryRowsProdResponse{
		DtakoFerryRows: dtakoFerryRowsProdModelToProto(row),
	}, nil
}

// List フェリー運行データ一覧取得
func (s *DTakoFerryRowsProdService) List(ctx context.Context, req *proto.ListDTakoFerryRowsProdRequest) (*proto.ListDTakoFerryRowsProdResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)

	if limit == 0 {
		limit = 100
	}

	rows, totalCount, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list ferry rows: %v", err)
	}

	items := make([]*proto.DTakoFerryRowsProd, len(rows))
	for i, row := range rows {
		items[i] = dtakoFerryRowsProdModelToProto(row)
	}

	return &proto.ListDTakoFerryRowsProdResponse{
		Items:      items,
		TotalCount: int32(totalCount),
	}, nil
}

// GetByUnkoNo 運行NOでフェリー運行データ取得
func (s *DTakoFerryRowsProdService) GetByUnkoNo(ctx context.Context, req *proto.GetDTakoFerryRowsProdByUnkoNoRequest) (*proto.ListDTakoFerryRowsProdResponse, error) {
	rows, err := s.repo.GetByUnkoNo(req.UnkoNo)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get ferry rows by unko_no: %v", err)
	}

	items := make([]*proto.DTakoFerryRowsProd, len(rows))
	for i, row := range rows {
		items[i] = dtakoFerryRowsProdModelToProto(row)
	}

	return &proto.ListDTakoFerryRowsProdResponse{
		Items:      items,
		TotalCount: int32(len(rows)),
	}, nil
}

// dtakoFerryRowsProdModelToProto ModelからProtoへの変換
func dtakoFerryRowsProdModelToProto(model *models.DTakoFerryRows) *proto.DTakoFerryRowsProd {
	protoRow := &proto.DTakoFerryRowsProd{
		Id:                int32(model.ID),
		UnkoNo:            model.UnkoNo,
		UnkoDate:          model.UnkoDate.Format("2006-01-02"),
		YomitoriDate:      model.YomitoriDate.Format("2006-01-02"),
		JigyoshoCd:        model.JigyoshoCD,
		JigyoshoName:      model.JigyoshoName,
		SharyoCd:          model.SharyoCD,
		SharyoName:        model.SharyoName,
		JomuinCd1:         model.JomuinCD1,
		JomuinName1:       model.JomuinName1,
		TaishoJomuinKbn:   model.TaishoJomuinKbn,
		KaishiDatetime:    model.KaishiDatetime.Format("2006-01-02T15:04:05Z07:00"),
		ShuryoDatetime:    model.ShuryoDatetime.Format("2006-01-02T15:04:05Z07:00"),
		FerryCompanyCd:    model.FerryCompanyCD,
		FerryCompanyName:  model.FerryCompanyName,
		NoribaCd:          model.NoribaCD,
		NoribaName:        model.NoribaName,
		Bin:               model.Bin,
		OribaCd:           model.OribaCD,
		OribaName:         model.OribaName,
		SeisanKbn:         model.SeisanKbn,
		SeisanKbnName:     model.SeisanKbnName,
		HyojunRyokin:      model.HyojunRyokin,
		KeiyakuRyokin:     model.KeiyakuRyokin,
		KosoShashuKbn:     model.KosoShashuKbn,
		KosoShashuKbnName: model.KosoShashuKbnName,
		MinashiKyori:      model.MinashiKyori,
	}

	// Optional field
	if model.FerrySrch != nil {
		protoRow.FerrySrch = model.FerrySrch
	}

	return protoRow
}
