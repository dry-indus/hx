package context

import (
	c "context"
	"hx/global"
	"hx/model/common"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContextB interface {
	c.Context
	Gin() *gin.Context
	common.Logger
}

type UserContext interface {
	ContextB
	Trace() string
	Merchant() *Merchant
	Session() *sessions.Session
}

type MerchantContext interface {
	ContextB
	Trace() string
	Merchant() *Merchant
	Session() *sessions.Session
}

type Merchant struct {
	ID       primitive.ObjectID      `json:"id"`
	Name     string                  `json:"name"`
	Category global.MerchantCategory `json:"category"`
	Telegram string                  `json:"telegram"`
}
