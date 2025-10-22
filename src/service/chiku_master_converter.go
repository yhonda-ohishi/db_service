package service

import (
	"github.com/yhonda-ohishi/db_service/src/models/ichibanboshi"
	pb "github.com/yhonda-ohishi/db_service/src/proto"
)

// convertChikuMasterToProto GORMモデルをProtoメッセージに変換
func convertChikuMasterToProto(m *ichibanboshi.ChikuMaster) *pb.ChikuMaster {
	return &pb.ChikuMaster{
		ChikuC:        m.ChikuC,
		ChikuN:        m.ChikuN,
		ChikuR:        m.ChikuR,
		ChikuF:        m.ChikuF,
		ChiikiC:       m.ChiikiC,
		YubinBango:    m.YubinBango,
		Jusho1:        m.Jusho1,
		Jusho2:        m.Jusho2,
		DenwaBango:    m.DenwaBango,
		FaxBango:      m.FAXBango,
		Tantosha:      m.Tantosha,
		Yobi1:         m.Yobi1,
		Yobi2:         m.Yobi2,
		Yobi3:         m.Yobi3,
		Yobi4:         m.Yobi4,
		Yobi5:         m.Yobi5,
		DgrTokuisakiC: m.DGRTokuisakiC,
		DgrTokuisakiH: m.DGRTokuisakiH,
		DgrHinmeiC:    m.DGRHinmeiC,
		DgrHinmeiH:    m.DGRHinmeiH,
	}
}
