package v1

import (
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
	USER_GROUP_V1 = "api/user/v1"
)

// @title          HaiXian 用户端 API
// @version        1.0
// @termsOfService http://swagger.io/terms/
// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html
// @host           swagger.mik888.com
// @BasePath       /api/user/v1
func Register(router *gin.Engine) {
	redirectU := router.Group("/redirect/user")
	redirectU.GET("/", U(userctr.Land.Redirect))

	user := router.Group(USER_GROUP_V1)
	user.Use(middleware.UAuth.Auth())
	home := user.Group("/home")
	{
		home.POST("/list", U(userctr.Home.List))
		home.POST("/search", U(userctr.Home.Search))
		home.POST("/order/submit", U(userctr.Home.SubmitOrder))
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
	merchant context.Merchant
	trace    string
}

func NewUserContext(c *gin.Context) *UserContext {
	trace := util.DefaultString(c.GetString(global.TRACE), util.UUID().String())
	ctx := &UserContext{
		Context: c,
		Logger: global.DL_LOGGER.WithFields(logrus.Fields{
			"server": "USER",
			"trace":  trace,
		}),
		trace: trace,
	}

	if merchant, _ := merchantser.Merchant.FindByName(ctx, c.GetString(global.MERCHANT)); merchant != nil {
		ctx.merchant = context.Merchant{
			ID:     merchant.ID,
			Name:   merchant.Name,
			TgName: merchant.TgName,
			TgID:   merchant.TgID,
		}
	}

	return ctx
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

func (this *UserContext) Lang() string {
	lang, ok := this.Context.Get(global.LANGUAGE)
	if !ok {
		return global.Application.DefaultLanguage
	}
	return lang.(string)
}

func (this *UserContext) Trace() string {
	if len(this.trace) == 0 {
		this.trace = util.UUID().String()
	}

	return this.trace
}
