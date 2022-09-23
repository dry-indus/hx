package merchantmod

import "hx/global"

type LoginRequest struct {
	Name     string `json:"name" binding:"required" validate:"required"`
	Password string `json:"password" binding:"required" validate:"required"`
	Code     string `json:"code"`
}

type LoginResponse struct {
	Name     string `json:"name"`
	Telegram string `json:"telegram"`
	// Category 1:餐饮,2:服饰
	Category global.MerchantCategory `json:"category" enums:"1,2"`
}

type LogoutRequest struct {
}

type LogoutResponse struct {
}

type RegisterRequest struct {
	// Name is user account
	Name string `json:"name" binding:"required" validate:"required"`
	// Password is user login password
	Password string `json:"password" binding:"required" validate:"required"`
	// PasswordTwo 二次输入密码，必须和Password 一致
	PasswordTwo string `json:"passwordTwo" binding:"required" validate:"required"`
	Code        string `json:"code"`
	// Telegram 小飞机账号
	Telegram string `json:"telegram" binding:"required" validate:"required"`
	// Category 1:餐饮,2:服饰
	Category global.MerchantCategory `json:"category" enums:"1,2" binding:"required" validate:"required"`
}

type RegisterResponse struct {
	Name     string `json:"name"`
	Telegram string `json:"telegram"`
	// Category 1:餐饮,2:服饰
	Category global.MerchantCategory `json:"category" enums:"1,2"`
}
