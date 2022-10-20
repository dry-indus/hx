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
)

type application struct {
	Port                 string
	DefaultLanguage      string
	DefaultMerchantName  string
	CloseHoken           bool `json:",string"` // 绕开cookie
	VerifyCodeTTLMinutes int  // 验证码有效时间
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
	Endpoint          string //oss-cn-hongkong.aliyuncs.com
	BucketName        string
	ConnectTimeoutSec int64 `json:",string"`
	ReadWriteTimeout  int64 `json:",string"`
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
