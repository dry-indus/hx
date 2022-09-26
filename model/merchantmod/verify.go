package merchantmod

type SendCodeRequest struct {
	// Name 账号
	Name string `json:"name" binding:"required" validate:"required"`
	// InvitationCode 邀请码
	InvitationCode string `json:"invitationCode"`
	ChatId         int64  `json:"chatId"`
}

type SendCodeResponse struct {
	// VerifyCode 验证码
	VerifyCode string `json:"verifyCode"`
}
