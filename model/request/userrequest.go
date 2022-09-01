/**
 * @Author: lzw5399
 * @Date: 2020/8/2 13:59
 * @Desc: user相关api的request model
 */
package request

// 用户注册api接收的model
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Gender   int    `json:"gender" binding:"gte=1,lte=2"` // male=1，female=2
	Phone    string `json:"phone" binding:"len=11"`
	Email    string `json:"email" binding:"required,email"`
}

// 用户登陆api接收的model
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 更新用户api接收的model
type UserUpdateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
