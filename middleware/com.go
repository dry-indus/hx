package middleware

import (
	"hx/global"
	"hx/util"

	"github.com/gin-gonic/gin"
)

func Lang() gin.HandlerFunc {
	return func(c *gin.Context) {
		hlang := c.GetHeader(global.LANGUAGE)
		qlang := c.Query(global.LANGUAGE)
		lang := util.DefaultString(qlang, hlang)
		lang = util.DefaultString(lang, global.Application.DefaultMerchantName)
		c.Set(global.LANGUAGE, lang)
		c.Header(global.LANGUAGE, lang)
		c.Next()
	}
}

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		trace := c.GetHeader(global.TRACE)
		trace = util.DefaultString(trace, util.UUID().String())
		c.Set(global.TRACE, trace)
		c.Header(global.TRACE, trace)
		c.Next()
	}
}
