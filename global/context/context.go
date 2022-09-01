package context

import (
	"hx/model/common"
	c "context"

	"github.com/gin-gonic/gin"
)

type ContextB interface {
	c.Context
	common.Logger
}

type UserContext interface {
	c.Context
	common.Logger
	Gin() *gin.Context
	Trace() string
	Merchant() *Merchant
}

type Merchant struct {
	ID       string
	Name     string
	Telegram string
}
