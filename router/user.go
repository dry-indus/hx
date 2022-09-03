package router

import (
	"fmt"
	"hx/controller/userctr"
	"hx/global"
	"hx/global/context"
	"hx/middleware"
	"hx/model/common"
	"hx/service/merchantser"
	"hx/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	XX_MERCHANT = "xx_merchant"
)

func initUserGroup(userGroup *gin.RouterGroup) {
	home := userGroup.Group("home")
	{
		home.POST(fmt.Sprintf("list/:%s", XX_MERCHANT), U(userctr.Home.List))
		home.POST("search", U(userctr.Home.Search))
		home.POST("order/info", U(userctr.Home.OrderInfo))
		home.POST("order/submit", U(userctr.Home.SubmitOrder))
	}
}

type UserHandlerFunc func(context.UserContext)

func U(f UserHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(NewUserContext(c))
	}
}

type UserContext struct {
	*gin.Context
	common.Logger
	trace    string
	merchant context.Merchant
}

var defaultMerchant = func() context.Merchant {
	return context.Merchant{
		ID:       "",
		Name:     "default",
		Telegram: "",
	}
}()

func NewUserContext(c *gin.Context) *UserContext {
	trace := util.UUID().String()
	ctx := &UserContext{
		Context: c,
		Logger: global.DL_LOGGER.WithFields(logrus.Fields{
			"server": "USER",
			"trace":  trace,
		}),
		trace:    trace,
		merchant: defaultMerchant,
	}

	merchantName, _ := c.Params.Get(XX_MERCHANT)
	if len(merchantName) == 0 {
		merchantName = c.GetHeader(XX_MERCHANT)
	}

	if len(merchantName) == 0 {
		return ctx
	}

	merchant, err := merchantser.Merchant.FindByName(ctx, merchantName)
	if err != nil {
		ctx.Errorf("merchantser.Merchant.FindByName failed! merchantName: %v, err: %v", merchantName, err)
		return ctx
	}

	ctx.merchant = context.Merchant{
		ID:       merchant.ID,
		Name:     merchant.Name,
		Telegram: merchant.Telegram,
	}

	return ctx
}

func (this *UserContext) Trace() string {
	return this.trace
}

func (this *UserContext) Merchant() *context.Merchant {
	return &this.merchant
}

func (this *UserContext) Gin() *gin.Context {
	return this.Context
}

func (this *UserContext) Session() *sessions.Session {
	a := middleware.NewUserAuth()
	return a.Session(this.Gin())
}
