package mysql

import "time"

// DTakoEvents dtako_eventsテーブルのモデル（本番DB）
type DTakoEvents struct {
	ID                int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	OperationNo       string    `gorm:"column:運行NO;size:23" json:"operation_no"`
	ReadDate          time.Time `gorm:"column:読取日;type:date" json:"read_date"`
	CarCode           int       `gorm:"column:車輌CD" json:"car_code"`
	CarCC             string    `gorm:"column:車輌CC;size:6" json:"car_cc"`
	TargetDriverType  int       `gorm:"column:対象乗務員区分" json:"target_driver_type"`
	DriverCode1       int       `gorm:"column:乗務員CD1" json:"driver_code1"`
	TargetDriverCode  int       `gorm:"column:対象乗務員CD" json:"target_driver_code"`
	StartDatetime     time.Time `gorm:"column:開始日時;type:datetime" json:"start_datetime"`
	EndDatetime       time.Time `gorm:"column:終了日時;type:datetime" json:"end_datetime"`
	EventCode         *int      `gorm:"column:イベントCD" json:"event_code"`
	EventName         string    `gorm:"column:イベント名;size:20" json:"event_name"`
	StartMileage      float64   `gorm:"column:開始走行距離" json:"start_mileage"`
	EndMileage        float64   `gorm:"column:終了走行距離" json:"end_mileage"`
	SectionTime       int       `gorm:"column:区間時間" json:"section_time"`
	SectionDistance   float64   `gorm:"column:区間距離" json:"section_distance"`
	StartCityCode     *int      `gorm:"column:開始市町村CD" json:"start_city_code"`
	StartCityName     string    `gorm:"column:開始市町村名;size:50" json:"start_city_name"`
	EndCityCode       *int      `gorm:"column:終了市町村CD" json:"end_city_code"`
	EndCityName       string    `gorm:"column:終了市町村名;size:50" json:"end_city_name"`
	StartPlaceCode    *int      `gorm:"column:開始場所CD" json:"start_place_code"`
	StartPlaceName    string    `gorm:"column:開始場所名;size:50" json:"start_place_name"`
	EndPlaceCode      *int      `gorm:"column:終了場所CD" json:"end_place_code"`
	EndPlaceName      string    `gorm:"column:終了場所名;size:50" json:"end_place_name"`
	StartGPSValid     *int      `gorm:"column:開始GPS有効" json:"start_gps_valid"`
	StartGPSLatitude  *int64    `gorm:"column:開始GPS緯度" json:"start_gps_latitude"`
	StartGPSLongitude *int64    `gorm:"column:開始GPS経度" json:"start_gps_longitude"`
	EndGPSValid       *int      `gorm:"column:終了GPS有効" json:"end_gps_valid"`
	EndGPSLatitude    *int64    `gorm:"column:終了GPS緯度" json:"end_gps_latitude"`
	EndGPSLongitude   *int64    `gorm:"column:終了GPS経度" json:"end_gps_longitude"`
}

// TableName テーブル名を指定
func (DTakoEvents) TableName() string {
	return "dtako_events"
}
