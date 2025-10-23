package service

import (
	"github.com/yhonda-ohishi/db_service/src/models/ichibanboshi"
	pb "github.com/yhonda-ohishi/db_service/src/proto"
)

// convertShainMasterToProto GORMモデルをProtoメッセージに変換
func convertShainMasterToProto(m *ichibanboshi.ShainMaster) *pb.Db_ShainMaster {
	// time.Time から string への変換
	var seinengappi *string
	if m.Seinengappi != nil {
		s := m.Seinengappi.Format("2006-01-02")
		seinengappi = &s
	}
	var nyushaNengappi *string
	if m.NyushaNengappi != nil {
		s := m.NyushaNengappi.Format("2006-01-02")
		nyushaNengappi = &s
	}
	var taishokuNengappi *string
	if m.TaishokuNengappi != nil {
		s := m.TaishokuNengappi.Format("2006-01-02")
		taishokuNengappi = &s
	}
	var jikaiKoshinbi *string
	if m.JikaiKoshinbi != nil {
		s := m.JikaiKoshinbi.Format("2006-01-02")
		jikaiKoshinbi = &s
	}

	return &pb.Db_ShainMaster{
		ShainC:             m.ShainC,
		ShainN:             m.ShainN,
		ShainR:             m.ShainR,
		ShainF:             m.ShainF,
		YubinBango:         m.YubinBango,
		Jusho1:             m.Jusho1,
		Jusho2:             m.Jusho2,
		DenwaBango:         m.DenwaBango,
		KeitaiBango:        m.KeitaiBango,
		ShainK:             m.ShainK,
		Seibetsu:           m.Seibetsu,
		Ketsuekigata:       m.Ketsuekigata,
		Seinengappi:        seinengappi,
		NyushaNengappi:     nyushaNengappi,
		TaishokuNengappi:   taishokuNengappi,
		DaiBunrui1:         m.DaiBunrui1,
		ChuBunrui1:         m.ChuBunrui1,
		ShoBunrui1:         m.ShoBunrui1,
		DaiBunrui2:         m.DaiBunrui2,
		ChuBunrui2:         m.ChuBunrui2,
		ShoBunrui2:         m.ShoBunrui2,
		KodanPlate:         m.KodanPlate,
		UriageMokuhyogaku:  int32(m.UriageMokuhyogaku),
		UntenMenkyoK:       m.UntenMenkyoK,
		MenkyoshoBango:     m.MenkyoshoBango,
		JikaiKoshinbi:      jikaiKoshinbi,
		JishaYoshaK:        m.JishaYoshaK,
		KeisanK:            m.KeisanK,
		ShiharaiRitsu:      m.ShiharaiRitsu,
		HasuK:              m.HasuK,
		BumonC:             m.BumonC,
		UnchinPatternC:     m.UnchinPatternC,
		Kiji1:              m.Kiji1,
		Kiji2:              m.Kiji2,
		Kiji3:              m.Kiji3,
		Kiji4:              m.Kiji4,
		Kiji5:              m.Kiji5,
		Yobi1:              m.Yobi1,
		Yobi2:              m.Yobi2,
		Yobi3:              m.Yobi3,
		Yobi4:              m.Yobi4,
		Yobi5:              m.Yobi5,
		KinmuTaikeiC:       m.KinmuTaikeiC,
	}
}
