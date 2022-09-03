package middleware

import (
	"fmt"
	"hx/global"
	"hx/model/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type MerchantAuth struct {
	common.Logger
}

const (
	MERCHANT_SESSION_KEY = "MERCHANT_SESSION"
	MERCHANT_TOKEN_KEY   = "MERCHANT_TOKEN"
)

func NewMerchantAuth() MerchantAuth {
	return MerchantAuth{global.DL_LOGGER.WithFields(logrus.Fields{
		"server": "MERCHANT_AUTH",
	})}
}

func (this MerchantAuth) Auth(redirectPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := this.Session(c)
		if s == nil {
			c.Redirect(http.StatusSeeOther, redirectPath)
			c.Abort()
			return
		}

		token, ok := this.Token(c, s)
		if !ok {
			s.Options.MaxAge = -1
			s.Save(c.Request, c.Writer)
			c.Redirect(http.StatusSeeOther, redirectPath)
			c.Abort()
			return
		}

		c.Set(MERCHANT_TOKEN_KEY, token)
		c.Next()
	}
}

func (this MerchantAuth) Session(c *gin.Context) *sessions.Session {
	s, err := global.DL_M_SESSION_STORE.Get(c.Request, MERCHANT_SESSION_KEY)
	if err != nil {
		this.Warningf("failed getting session: %s", err)
	}
	return s
}

func (this MerchantAuth) Token(c *gin.Context, s *sessions.Session) (string, bool) {
	tokenKV, ok := s.Values[MERCHANT_TOKEN_KEY]
	if !ok {
		this.Warningf("failed getting token! sessionID: %s", s.ID)
		return "", false
	}

	tokenK := fmt.Sprintf("%s", tokenKV)
	if len(tokenK) == 0 {
		this.Warningf("tokenK is empty! sessionID: %s", s.ID)
		return "", false
	}

	token := global.DL_CORE_REDIS.Get(c, tokenK).Val()
	if len(token) == 0 {
		this.Warningf("token is empty! sessionID: %s, tokenKey: %s", s.ID, tokenK)
		return "", false
	}

	return token, true
}
