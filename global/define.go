package global

import (
	gosonic "github.com/expectedsh/go-sonic/sonic"
)

const (
	//name-->token
	MERCHANT_TOEKN_KEY_FMT = "MT_%s" // ${name}
	MERCHANT_INFO_KEY_FMT  = "MI_%s" // ${token}

	TG_CHAT_INFO_FMT = "TG_CHAT_%v" // ${chatId}
	//name-->code
	VERIFY_CODE_FMT = "VERIFY_CODE_%v_%v_%v" // ${sence} ${name} ${tgId}

	OSS_PROGRESS_HASH_FMT = "OSS_PROGRESS_FMT_%v" // ${taskId} ${fileName}
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
	LastAt           = "lastAt"
	HOKEN            = "hoken"
	TRACE            = "trace"
)

type MerchantCategory int

const (
	All   MerchantCategory = -1 // 所有
	Cater MerchantCategory = 1  // 餐饮
	Dress MerchantCategory = 2  // 服饰
)

type CommodityStatus int

const (
	Online  CommodityStatus = 0
	Offline CommodityStatus = 1
	Show    CommodityStatus = 2
	Hide    CommodityStatus = 3
)

type Sence string

var SenceM = map[string]Sence{}

var (
	RegisterSence = newSence("register")
)

func newSence(s string) Sence {
	SenceM[s] = Sence(s)
	return SenceM[s]
}

func GetSence(s string) Sence {
	return SenceM[s]
}

type SonicBulkPushEvent struct {
	Collection string
	Bucket     string
	Records    []gosonic.IngestBulkRecord
	Lang       string
	Trace      string
}

type SonicSearchEvent struct {
	Collection string
	Bucket     string
	Terms      string
	Limit      int
	Offset     int
	Lang       string
	Result     chan *SonicSearcResult
	Trace      string
}

type SonicSuggestEvent struct {
	Collection string
	Bucket     string
	Word       string
	Limit      int
	Result     chan *SonicSearcResult
	Trace      string
}

type SonicSearcResult struct {
	Results []string
	Err     error
}
