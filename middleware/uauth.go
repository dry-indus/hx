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
		defer this.SaveSession(c, session)

		if session == nil {
			this.Warningf("trace: %s, session is nil", c.GetString(global.TRACE))
			response.InternalServerError(c).Failed("session is nil")
			c.Abort()
			return
		}

		{
			sm := util.ValueString(session.Values, global.MERCHANT)
			hm := c.GetHeader(global.MERCHANT)
			merchant := util.DefaultString(util.DefaultString(hm, sm), global.Application.DefaultMerchantName)
			c.Set(global.MERCHANT, merchant)
		}

		c.Header(global.MERCHANT, c.GetString(global.MERCHANT))
		c.Next()
	}
}

func (this UserAuth) Session(c *gin.Context) *sessions.Session {
	session, err := global.DL_U_SESSION_STORE.Get(c.Request, global.USER_SESSION)
	if err != nil {
		this.Errorf("trace: %s, failed getting session: %s", c.GetString(global.TRACE), err)
		return nil
	}

	return session
}

func (this UserAuth) SaveSession(c *gin.Context, s *sessions.Session) {
	if s == nil {
		return
	}

	if merchant, ok := c.Get(global.MERCHANT); ok {
		s.Values[global.MERCHANT] = merchant
	}

	if lang, ok := c.Get(global.LANGUAGE); ok {
		s.Values[global.LANGUAGE] = lang
	}

	s.Values[global.LastAt] = time.Now().Unix()

	if err := s.Save(c.Request, c.Writer); err != nil {
		this.Warningf("trace: %s, failed save session: %s", c.GetString(global.TRACE), err)
	}
}
