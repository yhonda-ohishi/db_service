package mysql

import "time"

// DTakoRows dtako_rowsテーブルのモデル（本番DB）
type DTakoRows struct {
	ID                   string    `gorm:"column:id;primaryKey;size:24" json:"id"`
	OperationNo          string    `gorm:"column:運行NO;size:23" json:"operation_no"`
	ReadDate             time.Time `gorm:"column:読取日;type:date" json:"read_date"`
	OperationDate        time.Time `gorm:"column:運行日;type:date" json:"operation_date"`
	CarCode              int       `gorm:"column:車輌CD" json:"car_code"`
	CarCC                string    `gorm:"column:車輌CC;size:6" json:"car_cc"`
	DriverCode1          *int      `gorm:"column:乗務員CD1" json:"driver_code1"`
	TargetDriverType     int       `gorm:"column:対象乗務員区分" json:"target_driver_type"`
	TargetDriverCode     int       `gorm:"column:対象乗務員CD" json:"target_driver_code"`
	StartWorkDatetime    time.Time `gorm:"column:出社日時;type:datetime" json:"start_work_datetime"`
	EndWorkDatetime      time.Time `gorm:"column:退社日時;type:datetime" json:"end_work_datetime"`
	DepartureDateTime    time.Time `gorm:"column:出庫日時;type:datetime" json:"departure_datetime"`
	ReturnDateTime       time.Time `gorm:"column:帰庫日時;type:datetime" json:"return_datetime"`
	DepartureMeter       float64   `gorm:"column:出庫メーター" json:"departure_meter"`
	ReturnMeter          float64   `gorm:"column:帰庫メーター" json:"return_meter"`
	TotalDistance        float64   `gorm:"column:総走行距離" json:"total_distance"`
	LoadedDistance       *float64  `gorm:"column:実車走行距離" json:"loaded_distance"`
	DestinationCityName  *string   `gorm:"column:行先市町村名;size:40" json:"destination_city_name"`
	DestinationPlaceName *string   `gorm:"column:行先場所名;size:40" json:"destination_place_name"`
	GeneralRoadDriveTime int       `gorm:"column:一般道運転時間" json:"general_road_drive_time"`
	HighwayDriveTime     int       `gorm:"column:高速道運転時間" json:"highway_drive_time"`
	BypassDriveTime      int       `gorm:"column:バイパス運転時間" json:"bypass_drive_time"`
	LoadedDriveTime      int       `gorm:"column:実車走行時間" json:"loaded_drive_time"`
	EmptyDriveTime       int       `gorm:"column:空車走行時間" json:"empty_drive_time"`
	Work1Time            int       `gorm:"column:作業１時間" json:"work1_time"`
	Work2Time            int       `gorm:"column:作業２時間" json:"work2_time"`
	Work3Time            int       `gorm:"column:作業３時間" json:"work3_time"`
	Work4Time            int       `gorm:"column:作業４時間" json:"work4_time"`
	Status1Distance      float64   `gorm:"column:状態１距離" json:"status1_distance"`
	Status1Time          int       `gorm:"column:状態１時間" json:"status1_time"`
}

// TableName テーブル名を指定
func (DTakoRows) TableName() string {
	return "dtako_rows"
}
