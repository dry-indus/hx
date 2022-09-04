package context

import (
	c "context"
	"hx/model/common"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type ContextB interface {
	c.Context
	common.Logger
}

type UserContext interface {
	ContextB
	Gin() *gin.Context
	Trace() string
	Merchant() *Merchant
	Session() *sessions.Session
}

type MerchantContext interface {
	ContextB
	Gin() *gin.Context
	Trace() string
	Merchant() *Merchant
	Session() *sessions.Session
}

type Merchant struct {
	ID       string
	Name     string
	Telegram string
}
