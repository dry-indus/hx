/**
 * @Author: lzw5399
 * @Date: 2020/8/2 14:01
 * @Desc: user相关api的response model
 */
package response

// user info的api返回值
type UserInfoResponse struct {
	Username string `json:"username"`
}