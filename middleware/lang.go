package middleware

import (
	"hx/global"
	"hx/util"

	"github.com/gin-gonic/gin"
)

func Lang() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang, _ := c.GetQuery(global.LANGUAGE)
		lang = util.DefaultString(lang, global.Application.DefaultLanguage)
		if len(lang) == 0 {
			lang = global.Application.DefaultLanguage
		}
		c.Set(global.LANGUAGE, lang)
	}
}
