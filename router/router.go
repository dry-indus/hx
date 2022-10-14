package router

import (
	mv1 "hx/api/merchant/v1"
	uv1 "hx/api/user/v1"
	"hx/global"
	"strings"

	_ "hx/docs"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	if !strings.EqualFold(global.ENV, "PRO") {
		gin.SetMode(gin.DebugMode)
	}
}

func Run() {
	router := gin.New()
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: global.DL_LOGGER.Writer(),
	}))

	router.Use(gin.Recovery())

	// default allow all origins
	corsR := router.Use(cors.Default())

	uv1.Register(router)
	corsR.GET("/swagger/uv1/*any", ginSwagger.WrapHandler(
		swaggerFiles.NewHandler(),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
		ginSwagger.InstanceName("uv1"),
	))

	mv1.Register(router)
	corsR.GET("/swagger/mv1/*any", ginSwagger.WrapHandler(
		swaggerFiles.NewHandler(),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
		ginSwagger.InstanceName("mv1"),
	))

	port := global.Application.Port
	global.DL_LOGGER.Infof("server listening at port: %v", port)
	_ = router.Run(":" + port)
}
