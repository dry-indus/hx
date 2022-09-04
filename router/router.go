package router

import (
	"fmt"
	"hx/controller/merchantctr"
	"hx/controller/userctr"
	"hx/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	USER_GROUP_V1     = "v1/hx/:merchant"
	MERCHANT_GROUP_V1 = "v1/merchant"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// default allow all origins
	r.Use(cors.Default())

	// swagger
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	hx := r.Group(fmt.Sprintf("/hx/:%s", XX_MERCHANT), BuildUserSession())
	hx.GET("/", U(userctr.Land.Page))

	redirectU := r.Group("/redirect/user")
	redirectU.GET("/", U(userctr.Land.Redirect))
	user := r.Group(USER_GROUP_V1)
	uauth := middleware.NewUserAuth()
	user.Use(uauth.Auth(redirectU.BasePath()))
	initUserGroup(user)

	redirectM := r.Group("/redirect/merchant")
	redirectM.GET("/", M(merchantctr.Land.Redirect))
	merchant := r.Group(MERCHANT_GROUP_V1)
	initMerchantGroup(merchant, redirectM.BasePath())

	return r
}
