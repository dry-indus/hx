package router

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

func initMerchantGroup(merchant *gin.RouterGroup, redirect string) {
	mauth := middleware.NewMerchantAuth()
	auth := merchant.Group("/auth")
	{
		auth.POST("/login", M(merchantctr.Auth.Login))
		auth.POST("/register", M(merchantctr.Auth.Register))
		auth.POST("/logout", mauth.Auth(redirect), M(merchantctr.Auth.Logout))

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
	trace    string
	merchant context.Merchant
}

func NewMerchantContext(c *gin.Context) *MerchantContext {
	trace := util.UUID().String()
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

func (this *MerchantContext) Trace() string {
	return this.trace
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
