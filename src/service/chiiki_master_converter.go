package service

import (
	"github.com/yhonda-ohishi/db_service/src/models/ichibanboshi"
	pb "github.com/yhonda-ohishi/db_service/src/proto"
)

// convertChiikiMasterToProto GORMモデルをProtoメッセージに変換
func convertChiikiMasterToProto(m *ichibanboshi.ChiikiMaster) *pb.Db_ChiikiMaster {
	return &pb.Db_ChiikiMaster{
		ChiikiC: m.ChiikiC,
		ChiikiN: m.ChiikiN,
		ChiikiR: m.ChiikiR,
		ChiikiF: m.ChiikiF,
	}
}
