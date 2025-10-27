// Package middleware 定义了Gin框架的中间件
package middleware

import (
	"server/internal/product/user/service"
	"server/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // 导入JWT库，用于解析和验证token
)

// AuthMiddleWare 认证中间件，用于验证JWT token
// 参数 secret: 用于验证JWT签名的密钥
// 返回值: Gin中间件处理函数
func AuthMiddleWare(secret []byte) gin.HandlerFunc {
	// 返回一个中间件处理函数
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		// 检查Authorization头是否为空或不以"Bearer "开头
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Unauthorized(c, response.CodeTokenInvalid, "invalid Authorization")
			c.Abort()
			return
		}
		// 从Authorization头中去除"Bearer "前缀，获取纯token字符串
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// 解析JWT token，将token中的声明信息解析到Claims结构中
		// 参数1: token字符串
		// 参数2: Claims结构体指针，用于存储解析后的声明
		// 参数3: 回调函数，返回用于验证签名的密钥
		token, err := jwt.ParseWithClaims(tokenString, &service.Claims{}, func(token *jwt.Token) (any, error) {
			// 返回密钥用于验证token签名
			return secret, nil
		})
		// 检查token解析是否出错
		if err != nil {
			response.Unauthorized(c, response.CodeTokenInvalid, "invalid Authorization")
			// 终止后续中间件和处理函数的执行
			c.Abort()
			return
		}
		// 将token中的Claims断言为service.Claims类型
		claim, ok := token.Claims.(*service.Claims)

		if ok {
			// 将用户ID存储到Gin上下文中，供后续处理函数使用
			c.Set("userID", claim.UserID)
			c.Set("account", claim.Account)
		}
		// 调用下一个中间件或处理函数
		c.Next()
	}
}
