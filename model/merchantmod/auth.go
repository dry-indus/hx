package merchantmod

import (
	"hx/global"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginRequest struct {
	// Name 商户登录账号
	Name string `json:"name" binding:"required" validate:"required"`
	// Password 商户登录密码
	Password string `json:"password" binding:"required" validate:"required"`
}

type LoginResponse struct {
	// Name 商户登录账号
	Name string `json:"name"`
	// Telegram 小飞机账号
	Telegram string `json:"telegram"`
	// Category 1:餐饮,2:服饰
	Category global.MerchantCategory `json:"category" enums:"1,2"`
}

type LogoutRequest struct {
}

type LogoutResponse struct {
}

type RegisterRequest struct {
	// Name 商户登录账号
	Name string `json:"name" binding:"required" validate:"required"`
	// Password 商户登录密码
	Password string `json:"password" binding:"required" validate:"required"`
	// PasswordTwo 二次输入密码，必须和Password 一致
	PasswordTwo string `json:"passwordTwo" binding:"required" validate:"required"`
	// InvitationCode 邀请码
	InvitationCode string `json:"invitationCode" binding:"required" validate:"required"`
	// Telegram 小飞机账号
	Telegram string `json:"telegram" binding:"required" validate:"required"`
	// VerifyCode 验证码 从小飞机获取
	VerifyCode string `json:"verifyCode" binding:"required" validate:"required"`
	// Category 1:餐饮,2:服饰
	Category global.MerchantCategory `json:"category" enums:"1,2" binding:"required" validate:"required"`
}

type RegisterResponse struct {
	// ID 商户id
	ID primitive.ObjectID `json:"id"`
	// Name 商户登录账号
	Name string `json:"name"`
	// Telegram 小飞机账号
	Telegram string `json:"telegram"`
	// Category 1:餐饮,2:服饰
	Category global.MerchantCategory `json:"category" enums:"1,2"`
}
