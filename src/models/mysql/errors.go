package mysql

import "errors"

// バリデーションエラー定義
var (
	// 共通エラー
	ErrInvalidPrice   = errors.New("price must be non-negative")
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateKey   = errors.New("duplicate primary key")

	// DTakoUriageKeihi関連
	ErrInvalidSrchID      = errors.New("srch_id cannot be empty")
	ErrInvalidKeihiC      = errors.New("keihi_c must be non-negative")
	ErrInvalidKm          = errors.New("km must be non-negative")
	ErrInvalidDtakoRowID  = errors.New("dtako_row_id cannot be empty")
	ErrInvalidDtakoRowIDR = errors.New("dtako_row_id_r cannot be empty")

	// ETCMeisai関連
	ErrInvalidDateTo     = errors.New("date_to cannot be empty")
	ErrInvalidDateToDate = errors.New("date_to_date cannot be empty")
	ErrInvalidIcFr       = errors.New("ic_fr cannot be empty")
	ErrInvalidIcTo       = errors.New("ic_to cannot be empty")
	ErrInvalidPriceBf    = errors.New("price_bf must be non-negative")
	ErrInvalidDescount   = errors.New("descount must be non-negative")
	ErrInvalidShashu     = errors.New("shashu must be positive")
	ErrInvalidEtcNum     = errors.New("etc_num cannot be empty")
	ErrInvalidHash       = errors.New("hash cannot be empty")
	ErrInvalidCreatedBy  = errors.New("created_by cannot be empty")
	ErrInvalidCreatedAt  = errors.New("created_at cannot be zero")
	ErrInvalidUpdatedAt  = errors.New("updated_at cannot be zero")

	// DTakoFerryRows関連
	ErrInvalidUnkoNo        = errors.New("unko_no cannot be empty")
	ErrInvalidUnkoDate      = errors.New("unko_date cannot be empty")
	ErrInvalidYomitoriDate  = errors.New("yomitori_date cannot be empty")
	ErrInvalidJigyoshoCD    = errors.New("jigyosho_cd must be positive")
	ErrInvalidJigyoshoName  = errors.New("jigyosho_name cannot be empty")
	ErrInvalidSharyoCD      = errors.New("sharyo_cd must be positive")
	ErrInvalidSharyoName    = errors.New("sharyo_name cannot be empty")
	ErrInvalidHyojunRyokin  = errors.New("hyojun_ryokin must be non-negative")
	ErrInvalidKeiyakuRyokin = errors.New("keiyaku_ryokin must be non-negative")
	ErrInvalidMinashiKyori  = errors.New("minashi_kyori must be non-negative")
)
