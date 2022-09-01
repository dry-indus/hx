/**
 * @Author: lzw5399
 * @Date: 2020/7/30 13:41
 * @Desc: check jwt token
 */
package middleware

import (
	"errors"
	"hx/global"
	"hx/model"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	authConfigs = global.Auth
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		claims, err := parseToken(parts[1])
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 检查令牌信息是否满足要求
		if claims.Issuer != authConfigs.Issuer ||
			claims.Audience != authConfigs.Audience ||
			claims.ExpiresAt < time.Now().Unix() {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		// 将claims存入请求的上下文
		c.Set("currentUser", claims)
		c.Next()
	}
}

// 解析token为
func parseToken(tokenString string) (*model.CustomClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(authConfigs.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
