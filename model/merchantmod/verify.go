package merchantmod

type SendCodeRequest struct {
	// Name 账号
	Name string `json:"name" binding:"required" validate:"required"`
	// Telegram 用户 id
	TgId int64 `json:"tgId"`
}

type SendCodeResponse struct {
	// VerifyCode 验证码
	VerifyCode string `json:"verifyCode"`
}
