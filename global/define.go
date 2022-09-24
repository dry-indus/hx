package global

const (
	//name-->token
	MERCHANT_TOEKN_KEY_FMT = "MT_%s" // ${name}
	MERCHANT_INFO_KEY_FMT  = "MI_%s" // ${token}
)

const (
	MERCHANT_SESSION = "merchant_session"
	MERCHANT_TOKEN   = "merchant_token"
	MERCHANT_INFO    = "merchant_info"
	ACCOUNT          = "account"
	LANGUAGE         = "language"
	MERCHANT         = "merchant"
	USER_SESSION     = "user_session"
	USER_TOKEN_KEY   = "user_token"
	LastAt           = "LastAt"
)

type MerchantCategory int

const (
	Cater MerchantCategory = 1 // 餐饮
	Dress MerchantCategory = 2 // 服饰
	Sex   MerchantCategory = -1
)

type CommodityStatus int

const (
	Online  CommodityStatus = 0
	Offline CommodityStatus = 1
	Show    CommodityStatus = 2
	Hide    CommodityStatus = 3
)

type ChoiceOpt int

const (
	SingleChoice   ChoiceOpt = 0
	MultipleChoice ChoiceOpt = 1
	MustChoice     ChoiceOpt = 2
)
