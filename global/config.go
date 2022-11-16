package global

var (
	AppName     string
	ENV         string
	Application = NewNamespace("Application", &application{}).(*application)
	Logger      = NewNamespace("Logger", &logger{}).(*logger)
	Common      = NewNamespace("Common", &common{}).(*common)
	CoreRedis   = NewNamespace("CoreRedis", &redis{}).(*redis)
	CoreMongo   = NewNamespace("CoreMongo", &mongo{}).(*mongo)
	Auth        = NewNamespace("Auth", &auth{}).(*auth)
	Telegram    = NewNamespace("Telegram", &telegram{}).(*telegram)
	Oss         = NewNamespace("Oss", &oss{}).(*oss)
	Sonic       = NewNamespace("Sonic", &sonic{}).(*sonic)
	Landing     = NewNamespace("Landing.json", &landing{}).(*landing)
)

type application struct {
	Port                 string
	DefaultLanguage      string
	DefaultMerchantName  string
	CloseHoken           bool `json:",string"` // 绕开cookie
	VerifyCodeTTLMinutes int  // 验证码有效时间
	Domian               string
}

type logger struct {
	LogLevel string
	File     string
}

type common struct {
}

type redis struct {
	Addr     string
	Username string
	Password string
}

type mongo struct {
	Uri string
}

type auth struct {
	Secret        string
	Issuer        string
	Audience      string
	ExpireMinutes int
}

type telegram struct {
	HXBotToken string
	HXBotDebug bool `json:",string"`
}

type oss struct {
	AccessKeyId       string
	AccessKeySecret   string
	Endpoint          string // oss-cn-hongkong.aliyuncs.com
	BucketName        string
	ConnectTimeoutSec int64
	ReadWriteTimeout  int64
	UrlScheme         string // https
	UrlHost           string // oss.hx24h.com
}

type sonic struct {
	Host             string
	Port             int //1491
	Password         string
	ParallelRoutines int
}

type landing struct {
	Entrys []struct {
		// Name 入口名称, eg: 我是商户
		Name string `json:"name"`
		// URL 跳转链接, eg: www.baidu.com
		URL string `json:"url"`
		// BackgroundRPGA 背景色, eg: #F78870
		BackgroundRPGA string `json:"backgroundRPGA"`
		// Show true: 显示，否则不显示
		Show bool `json:",string"`
	} `json:"entrys"`
}

var Namespacem = make(map[string]interface{})
var Namespaces []string

func NewNamespace(namespace string, ptr interface{}) interface{} {
	if v, ok := Namespacem[namespace]; ok {
		return v
	}
	Namespacem[namespace] = ptr
	Namespaces = append(Namespaces, namespace)
	return ptr
}
