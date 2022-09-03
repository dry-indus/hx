package middleware

import (
	"hx/global"
	"hx/model/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type UserAuth struct {
	common.Logger
}

const (
	USER_SESSION_KEY = "USER_SESSION"
	USER_TOKEN_KEY   = "USER_TOKEN"
)

func NewUserAuth() UserAuth {
	return UserAuth{global.DL_LOGGER.WithFields(logrus.Fields{
		"server": "USER_AUTH",
	})}
}

func (this UserAuth) Auth(redirectPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := this.Session(c)
		if s == nil {
			c.Redirect(http.StatusSeeOther, redirectPath)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (this UserAuth) Session(c *gin.Context) *sessions.Session {
	s, err := global.DL_U_SESSION_STORE.Get(c.Request, USER_SESSION_KEY)
	if err != nil {
		this.Warningf("failed getting session: %s", err)
	}
	return s
}
