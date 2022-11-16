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
	Lang() string
	Trace() string
}

type UserContext interface {
	ContextB
	Merchant() *Merchant
	Session() *sessions.Session
}

type MerchantContext interface {
	ContextB
	Merchant() *Merchant
	Session() *sessions.Session
}

type Merchant struct {
	ID       primitive.ObjectID      `json:"id"`
	Name     string                  `json:"name"`
	Category global.MerchantCategory `json:"category"`
	TgName   string                  `json:"tgName"`
	TgID     int64                   `json:"tgId"`
}
