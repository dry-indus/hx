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

func initMerchantGroup(merchant *gin.RouterGroup, redirectPath string) {
	mauth := middleware.NewMerchantAuth()

	auth := merchant.Group("auth")
	{
		auth.POST("login", M(merchantctr.Auth.Login))
		auth.POST("register", M(merchantctr.Auth.Register))
		auth.POST("logout", mauth.Auth(redirectPath), M(merchantctr.Auth.Logout))

		commodity := merchant.Group("commodity")
		commodity.Use(mauth.Auth(redirectPath))
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

func (this *MerchantContext) Session() *sessions.Session {
	a := middleware.NewMerchantAuth()
	return a.Session(this.Gin())
}
