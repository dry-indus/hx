package middleware

import (
	"fmt"
	"hx/global"
	"hx/global/context"
	"hx/model/common"
	"hx/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type MerchantAuth struct {
	common.Logger
}

func NewMerchantAuth() MerchantAuth {
	return MerchantAuth{global.DL_LOGGER.WithFields(logrus.Fields{
		"server": "MERCHANT_AUTH",
	})}
}

func (this MerchantAuth) Auth(redirectPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := this.Session(c)
		defer this.SaveSession(c, s)

		if s == nil {
			c.Redirect(http.StatusSeeOther, redirectPath)
			c.Abort()
			return
		}

		token, ok := this.Token(c, s)
		if !ok {
			hoken := c.Request.Header.Get(global.HOKEN)
			if len(hoken) == 0 {
				s.Options.MaxAge = -1
				s.Save(c.Request, c.Writer)
				c.Redirect(http.StatusSeeOther, redirectPath)
				c.Abort()
				return
			}
			token = hoken
		}
		c.Set(global.MERCHANT_TOKEN, token)

		merchant := this.Merchant(c, token)
		if merchant.ID.IsZero() {
			s.Options.MaxAge = -1
			s.Save(c.Request, c.Writer)
			c.Redirect(http.StatusSeeOther, redirectPath)
			c.Abort()
			return
		}
		c.Set(global.MERCHANT_INFO, merchant)

		c.Next()
	}
}

func (this MerchantAuth) Session(c *gin.Context) *sessions.Session {
	session, err := global.DL_M_SESSION_STORE.Get(c.Request, global.MERCHANT_SESSION)
	if err != nil {
		this.Warningf("trace: %s, failed getting session: %s", c.GetString(global.TRACE), err)
	}

	return session
}

func (this MerchantAuth) SaveSession(c *gin.Context, s *sessions.Session) {
	if s == nil {
		return
	}

	if merchant, ok := c.Get(global.MERCHANT); ok {
		s.Values[global.MERCHANT] = merchant
	}

	if lang, ok := c.Get(global.LANGUAGE); ok {
		s.Values[global.LANGUAGE] = lang
	}

	if token, ok := c.Get(global.MERCHANT_TOKEN); ok {
		s.Values[global.MERCHANT_TOKEN] = token
	}

	if account, ok := c.Get(global.ACCOUNT); ok {
		s.Values[global.ACCOUNT] = account
	}

	s.Values[global.LastAt] = time.Now().Unix()

	if err := s.Save(c.Request, c.Writer); err != nil {
		this.Warningf("trace: %s, failed save session: %s", c.GetString(global.TRACE), err)
	}
}

func (this MerchantAuth) Token(c *gin.Context, s *sessions.Session) (string, bool) {
	tokenL := util.ValueString(s.Values, global.MERCHANT_TOKEN)
	if len(tokenL) == 0 {
		this.Warningf("trace: %s, failed getting token! sessionID: %s", c.GetString(global.TRACE), s.ID)
		return "", false
	}

	name := util.ValueString(s.Values, global.ACCOUNT)
	if len(name) == 0 {
		this.Warningf("trace: %s, failed getting account! sessionID: %s", c.GetString(global.TRACE), s.ID)
		return "", false
	}

	tokenKey := fmt.Sprintf(global.MERCHANT_TOEKN_KEY_FMT, name)
	tokenR := global.DL_CORE_REDIS.Get(c, tokenKey).Val()
	if len(tokenR) == 0 {
		this.Warningf("trace: %s, tokenR is empty! sessionID: %s, name: %s", c.GetString(global.TRACE), s.ID, name)
		return "", false
	}

	if tokenL != tokenR {
		this.Warningf("trace: %s, does not match! sessionID: %s, name: %s", c.GetString(global.TRACE), s.ID, name)
		return "", false
	}

	return tokenL, true
}

func (this MerchantAuth) Merchant(c *gin.Context, token string) (merchant context.Merchant) {
	infoKey := fmt.Sprintf(global.MERCHANT_INFO_KEY_FMT, token)
	s := global.DL_CORE_REDIS.Get(c, infoKey).Val()
	util.JSON.UnmarshalFromString(s, &merchant)
	return
}
