package v1

import (
	"hx/controller/landingctr"
	"hx/global"
	"hx/global/context"
	"hx/middleware"
	"hx/model/common"
	"hx/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	LANDING_GROUP_V1 = "api/landing/v1"
)

// @title          HaiXian 落地页 API
// @version        1.0
// @termsOfService http://swagger.io/terms/
// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html
// @host           swagger.mik888.com
// @BasePath       /api/landing/v1
func Register(router *gin.Engine) {
	landing := router.Group(LANDING_GROUP_V1)
	landing.Use(middleware.Lang())

	pre := landing.Group("/pre")
	{
		pre.POST("/setting", L(landingctr.Pre.Settting))
	}

	store := landing.Group("/store")
	{
		store.POST("/search", L(landingctr.Store.Search))
	}

	test := landing.Group("/test")
	{
		test.POST("/push", L(landingctr.Store.SearchPush))
	}
}

type UserHandlerFunc func(context.ContextB)

func L(f UserHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(NewLandingContext(c))
	}
}

type LandingContext struct {
	*gin.Context
	common.Logger
	trace string
}

func NewLandingContext(c *gin.Context) *LandingContext {
	trace := util.UUID().String()
	ctx := &LandingContext{
		Context: c,
		Logger: global.DL_LOGGER.WithFields(logrus.Fields{
			"server": "LANDING",
			"trace":  trace,
		}),
		trace: trace,
	}

	return ctx
}

func (this *LandingContext) Gin() *gin.Context {
	return this.Context
}

func (this *LandingContext) Lang() string {
	lang, ok := this.Context.Get(global.LANGUAGE)
	if !ok {
		return global.Application.DefaultLanguage
	}
	return lang.(string)
}

func (this *LandingContext) Trace() string {
	if len(this.trace) == 0 {
		this.trace = util.UUID().String()
	}

	return this.trace
}
