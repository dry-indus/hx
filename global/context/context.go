package context

import (
	c "context"
	"hx/model/common"

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
