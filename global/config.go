package global

import "go.mongodb.org/mongo-driver/bson/primitive"

var (
	AppName     string
	Application = NewNamespace("Application", &application{}).(*application)
	Logger      = NewNamespace("Logger", &logger{}).(*logger)
	Common      = NewNamespace("Common", &common{}).(*common)
	CoreRedis   = NewNamespace("CoreRedis", &redis{}).(*redis)
	CoreMongo   = NewNamespace("CoreMongo", &mongo{}).(*mongo)
	Auth        = NewNamespace("Auth", &auth{}).(*auth)
)

type application struct {
	Port                    string
	DefaultLanguage         string
	DefaultMerchantName     string
	DefaultMerchantId       primitive.ObjectID
	DefaultMerchantTelegram string
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
