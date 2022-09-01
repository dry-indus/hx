package router

import (
	"fmt"
	"hx/controller/userctr"
	"hx/global"
	"hx/global/context"
	"hx/model/common"
	"hx/model/db"
	"hx/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	XX_MERCHANT = "xx_merchant"
)

var defaultMerchant = func() context.Merchant {
	return context.Merchant{
		ID:       "",
		Name:     "default",
		Telegram: "",
	}
}()

func initUserGroup(userGroup *gin.RouterGroup) {
	home := userGroup.Group("home")
	{
		home.POST(fmt.Sprintf("list/:%s", XX_MERCHANT), U(userctr.Home.List))
		home.POST("search", U(userctr.Home.Search))
		home.POST("order/info", U(userctr.Home.OrderInfo))
		home.POST("order/submit", U(userctr.Home.SubmitOrder))
	}
}

type HandlerFunc func(context.UserContext)

func U(f HandlerFunc) gin.HandlerFunc {
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

func NewUserContext(c *gin.Context) context.UserContext {
	trace := util.UUID().String()
	log := global.DL_LOGGER.WithFields(logrus.Fields{
		"server": "USER",
		"trace":  trace,
	})

	ctx := &UserContext{
		Context:  c,
		Logger:   log,
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

	merchant, err := db.Merchant.FindOneByName(ctx, merchantName)
	if err != nil {
		ctx.Errorf("db.Merchant.FindOneByName failed! merchantName: %v, err: %v", merchantName, err)
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
