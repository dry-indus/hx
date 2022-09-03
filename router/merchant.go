package router

import (
	"hx/controller/merchantctr"
	"hx/global"
	"hx/global/context"
	"hx/middleware"
	"hx/model/common"
	"hx/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func initMerchantGroup(userGroup *gin.RouterGroup) {
	auth := userGroup.Group("auth")
	auth.GET("/", M(merchantctr.Land.Redirect))
	mauth := middleware.NewMerchantAuth()
	auth.Use(mauth.Auth(auth.BasePath()))
	{
		auth.POST("login", M(merchantctr.Auth.Login))
		auth.POST("logout", M(merchantctr.Auth.Logout))
		auth.POST("register", M(merchantctr.Auth.Register))

		// commodity := merchant.Group("commodity")
		// {
		// 	commodity.POST("add", M(userctr.Home.SubmitOrder))
		// 	commodity.POST("modify", M(userctr.Home.SubmitOrder))
		// 	commodity.POST("del", M(userctr.Home.SubmitOrder))
		// 	commodity.POST("publish", M(userctr.Home.SubmitOrder))
		// 	commodity.POST("hide", M(userctr.Home.SubmitOrder))
		// }
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

func MerchantLogger(c *gin.Context) common.Logger {
	trace := util.UUID().String()
	log := global.DL_LOGGER.WithFields(logrus.Fields{
		"server": "MERCHANT",
		"trace":  trace,
	})
	return log
}

func NewMerchantContext(c *gin.Context) *MerchantContext {
	trace := util.UUID().String()
	ctx := &MerchantContext{
		Context: c,
		Logger:  global.DL_LOGGER.WithFields(logrus.Fields{
			"server": "MERCHANT",
			"trace":  trace,
		}),
		trace:   trace,
	}

	val, ok := c.Get("Merchant")
	if !ok {
		panic("not find merchant")
	}

	ctx.merchant, ok = val.(context.Merchant)
	if !ok {
		panic("type is't merchant")
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
