package middleware

import (
	"hx/global"
	"hx/global/response"
	"hx/model/common"
	"hx/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

var UAuth = NewUserAuth()

type UserAuth struct {
	common.Logger
}

func NewUserAuth() UserAuth {
	return UserAuth{global.DL_LOGGER.WithFields(logrus.Fields{
		"server": "USER_AUTH",
	})}
}

func (this UserAuth) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := this.Session(c)
		if session == nil {
			this.Warningf("session is nil")
			response.InternalServerError(c).Failed("session is nil")
			c.Abort()
			return
		}

		c.Header(global.MERCHANT, c.GetString(global.MERCHANT))
		c.Next()
	}
}

func (this UserAuth) Session(c *gin.Context) *sessions.Session {
	session, err := global.DL_U_SESSION_STORE.Get(c.Request, global.USER_SESSION)
	if err != nil {
		this.Errorf("failed getting session: %s", err)
		return nil
	}

	{
		sm := util.ValueString(session.Values, global.MERCHANT)
		km := c.GetHeader(global.MERCHANT)
		merchant := util.DefaultString(util.DefaultString(sm, km), global.Application.DefaultMerchantName)
		session.Values[global.MERCHANT] = merchant
		c.Set(global.MERCHANT, merchant)
	}

	{
		lang, _ := c.GetQuery(global.LANGUAGE)
		lang = util.DefaultString(lang, global.Application.DefaultLanguage)
		session.Values[global.LANGUAGE] = lang
	}

	{
		session.Values[global.LastAt] = time.Now().Unix()
	}

	if err := session.Save(c.Request, c.Writer); err != nil {
		this.Warningf("failed save session: %s", err)
	}

	return session
}
