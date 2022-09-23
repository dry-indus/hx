package middleware

import (
	"hx/global"
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
		c.Next()
		if session != nil {
			err := session.Save(c.Request, c.Writer)
			if err != nil {
				this.Warningf("failed save session: %s", err)
			}
		}
	}
}

func (this UserAuth) Session(c *gin.Context) *sessions.Session {
	session, err := global.DL_U_SESSION_STORE.Get(c.Request, global.USER_SESSION_KEY)
	if err != nil {
		this.Errorf("failed getting session: %s", err)
		return nil
	}

	{
		sm := util.ValueString(session.Values, global.MERCHANT)
		qm, _ := c.GetQuery(global.MERCHANT)
		merchant := util.DefaultString(sm, qm)
		merchant = util.DefaultString(merchant, global.Application.DefaultMerchantName)
		session.Values[global.MERCHANT] = merchant
	}

	{
		lang, _ := c.GetQuery(global.LANGUAGE_KEY)
		lang = util.DefaultString(lang, global.Application.DefaultLanguage)
		session.Values[global.LANGUAGE_KEY] = lang
	}

	{
		session.Values[global.LastAt] = time.Now().Unix()
	}

	return session
}
