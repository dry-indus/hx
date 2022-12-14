package v1

import (
	"hx/controller/merchantctr"
	"hx/global"
	"hx/global/context"
	"hx/middleware"
	"hx/model/common"
	"hx/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	MERCHANT_GROUP_V1 = "api/merchant/v1"
)

// @title                      HaiXian 商户端 API
// @version                    1.0
// @termsOfService             http://swagger.io/terms/
// @license.name               Apache 2.0
// @license.url                http://www.apache.org/licenses/LICENSE-2.0.html
// @host                       swagger.mik888.com
// @BasePath                   /api/merchant/v1
// @securityDefinitions.apikey Auth
// @in                         header
// @name                       hoken
func Register(router *gin.Engine) {
	redirectM := router.Group("/redirect/merchant")
	redirectM.GET("/", M(merchantctr.Land.Redirect))
	redirect := redirectM.BasePath()

	merchant := router.Group(MERCHANT_GROUP_V1)

	mauth := middleware.NewMerchantAuth()
	auth := merchant.Group("/auth")
	{
		auth.POST("/login", M(merchantctr.Auth.Login))
		auth.POST("/register", M(merchantctr.Auth.Register))
		auth.POST("/logout", mauth.Auth(redirect), M(merchantctr.Auth.Logout))

	}

	verify := merchant.Group("/verify/:sence")
	{
		verify.POST("/code/send", M(merchantctr.Verify.SendCode))
	}

	file := merchant.Group("/file")
	{
		file.POST("/upload", mauth.Auth(redirect), M(merchantctr.File.Upload))
		file.POST("/status", mauth.Auth(redirect), M(merchantctr.File.Status))
	}

	commodity := merchant.Group("/commodity")
	commodity.Use(mauth.Auth(redirect))
	{
		commodity.POST("/list", M(merchantctr.Commodity.List))
		commodity.POST("/add", M(merchantctr.Commodity.Add))
		commodity.POST("/modify", M(merchantctr.Commodity.Modify))
		commodity.POST("/del", M(merchantctr.Commodity.Del))
		commodity.POST("/publish", M(merchantctr.Commodity.Publish))
		commodity.POST("/hide", M(merchantctr.Commodity.Hide))

		tag := commodity.Group("/tag")
		{
			tag.POST("/add", M(merchantctr.Tag.Add))
			tag.POST("/del", M(merchantctr.Tag.Del))
			tag.POST("/stat", M(merchantctr.Tag.Stat))
		}

		sp := commodity.Group("/sp")
		{
			sp.POST("/add", M(merchantctr.SP.Add))
			sp.POST("/modify", M(merchantctr.SP.Modify))
			sp.POST("/del", M(merchantctr.SP.Del))
		}
	}
	// setting := merchant.Group("setting")
	{

	}
}

type MerchantHandlerFunc func(context.MerchantContext)

func M(f MerchantHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(NewMerchantContext(c))
	}
}

type MerchantContext struct {
	*gin.Context
	common.Logger
	merchant context.Merchant
	trace    string
}

func NewMerchantContext(c *gin.Context) *MerchantContext {
	trace := util.DefaultString(c.GetString(global.TRACE), util.UUID().String())
	ctx := &MerchantContext{
		Context: c,
		Logger: global.DL_LOGGER.WithFields(logrus.Fields{
			"server": "MERCHANT",
			"trace":  trace,
		}),
		trace: trace,
	}

	val, _ := c.Get(global.MERCHANT_INFO)
	merchant, ok := val.(context.Merchant)
	if ok {
		ctx.merchant = merchant
	}

	return ctx
}

func (this *MerchantContext) Merchant() *context.Merchant {
	return &this.merchant
}

func (this *MerchantContext) Gin() *gin.Context {
	return this.Context
}

func (this *MerchantContext) Session() *sessions.Session {
	a := middleware.NewMerchantAuth()
	return a.Session(this.Gin())
}

func (this *MerchantContext) Lang() string {
	lang, ok := this.Context.Get(global.LANGUAGE)
	if !ok {
		return global.Application.DefaultLanguage
	}
	return lang.(string)
}

func (this *MerchantContext) Trace() string {
	if len(this.trace) == 0 {
		this.trace = util.UUID().String()
	}

	return this.trace
}
